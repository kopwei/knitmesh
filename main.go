package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/kopwei/knitmesh/common"
	"github.com/kopwei/knitmesh/plugin"
	"github.com/kopwei/knitmesh/plugin/listener"
)

var version = "(unreleased version)"

func main() {
	var (
		justVersion bool
		address     string
		nameserver  string
		logLevel    string
	)

	flag.BoolVar(&justVersion, "version", false, "print version and exit")
	flag.StringVar(&logLevel, "log-level", "info", "logging level (debug, info, warning, error)")
	flag.StringVar(&address, "socket", "/run/docker/plugins/weave.sock", "socket on which to listen")
	flag.StringVar(&nameserver, "nameserver", "", "nameserver to provide to containers")

	flag.Parse()

	if justVersion {
		fmt.Printf("weave plugin %s\n", version)
		os.Exit(0)
	}

	common.SetLogLevel(logLevel)

	common.Log.Println("knitmesh plugin", version, "Command line options:", os.Args)

	var d listener.Driver
	d, err := plugin.New(version, nameserver)
	if err != nil {
		common.Log.Fatalf("unable to create driver: %s", err)
	}

	var netlistener net.Listener

	// remove socket from last invocation
	if err := os.Remove(address); err != nil && !os.IsNotExist(err) {
		common.Log.Fatal(err)
	}
	netlistener, err = net.Listen("unix", address)
	if err != nil {
		common.Log.Fatal(err)
	}
	defer netlistener.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	endChan := make(chan error, 1)
	go func() {
		endChan <- listener.Listen(netlistener, d)
	}()

	select {
	case sig := <-sigChan:
		common.Log.Debugf("Caught signal %s; shutting down", sig)
	case err := <-endChan:
		if err != nil {
			common.Log.Errorf("Error from listener: %s", err)
			netlistener.Close()
			os.Exit(1)
		}
	}
}
