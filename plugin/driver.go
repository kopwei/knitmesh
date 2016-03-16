package plugin

import (
	"fmt"
	"sync"

	"github.com/docker/go-plugins-helpers/network"
	"github.com/kopwei/goovs"
	"github.com/kopwei/knitmesh/common"
	"github.com/kopwei/knitmesh/plugin/listener"
)

const (
	generalOpt = "com.docker.network.generic"
)

const (
	optName = "com.github.kopwei.knitmesh.name"
	//optName = "com.github.kopwei.knitmesh.name"
)

type driver struct {
	version    string
	nameserver string
	sync.RWMutex
	endpoints   map[string]struct{}
	ovsdbClient goovs.OvsClient
	networks    map[string]*networkInfo
}

type networkInfo struct {
	bridgeName string
}

var caps = &network.CapabilitiesResponse{
	Scope: network.LocalScope,
}

// New ist used to intialize the plugin object
func New(version string, nameserver string) (listener.Driver, error) {
	//client, err := docker.NewClient("unix:///var/run/docker.sock")
	//if err != nil {
	//	return nil, errorf("could not connect to docker: %s", err)
	//}

	driver := &driver{
		nameserver: nameserver,
		version:    version,
		endpoints:  make(map[string]struct{}),
		networks:   make(map[string]*networkInfo),
	}
	var err error
	driver.ovsdbClient, err = goovs.GetOVSClient("unix", "")
	if err != nil {
		return nil, err
	}
	/*
		_, err = NewWatcher(client, driver)
		if err != nil {
			return nil, err
		}
	*/
	return driver, nil
}

func errorf(format string, a ...interface{}) error {
	common.Log.Errorf(format, a...)
	return fmt.Errorf(format, a...)
}

func getBrNameOption(request *network.CreateNetworkRequest) (string, error) {
	options, ok := request.Options[generalOpt]
	if !ok {
		return "", fmt.Errorf("Generic optionaa doesn't exists in options")
	}
	optionMap := options.(map[string]interface{})
	brname, ok := optionMap[optName]
	if !ok {
		return "", fmt.Errorf("Bridge name is not specified in option")
	}
	return brname.(string), nil
}

func (driver *driver) GetCapabilities() (*network.CapabilitiesResponse, error) {
	common.Log.Debugf("Get capabilities: responded with %+v", caps)
	return caps, nil
}

func (driver *driver) CreateNetwork(create *network.CreateNetworkRequest) error {
	common.Log.Debugf("Create network request %+v", create)
	name, err := getBrNameOption(create)
	if err != nil {
		return err
	}
	err = driver.ovsdbClient.CreateBridge(name)
	if err != nil {
		common.Log.Debugf("Failed to create ovs bridge %s due to %s", name, err.Error())
	}
	netInfo := &networkInfo{bridgeName: name}
	driver.networks[create.NetworkID] = netInfo
	//driver.ovsdbClient.CreateBridge(create.Options[""])
	common.Log.Infof("Create network %s", create.NetworkID)
	return nil
}

func (driver *driver) DeleteNetwork(deletereq *network.DeleteNetworkRequest) error {
	common.Log.Debugf("Delete network request: %+v", deletereq)
	netInfo, ok := driver.networks[deletereq.NetworkID]
	if !ok {
		return fmt.Errorf("Failed to delete network due to invalid network id %s", deletereq.NetworkID)
	}
	err := driver.ovsdbClient.DeleteBridge(netInfo.bridgeName)
	if err != nil {
		return fmt.Errorf("Failed to delete ovs bridge %s due to %s", netInfo.bridgeName, err.Error())
	}
	delete(driver.networks, deletereq.NetworkID)
	common.Log.Infof("Destroy network %s", deletereq.NetworkID)
	return nil
}

func (driver *driver) CreateEndpoint(create *network.CreateEndpointRequest) (*network.CreateEndpointResponse, error) {
	common.Log.Debugf("Create endpoint request %+v", create)
	return nil, nil
}

func (driver *driver) DeleteEndpoint(deleteReq *network.DeleteEndpointRequest) error {
	common.Log.Debugf("Delete endpoint request: %+v", deleteReq)
	common.Log.Infof("Delete endpoint %s", deleteReq.EndpointID)
	return nil
}

func (driver *driver) HasEndpoint(endpointID string) bool {
	return false
}

func (driver *driver) EndpointInfo(req *network.InfoRequest) (*network.InfoResponse, error) {
	common.Log.Debugf("Endpoint info request: %+v", req)
	common.Log.Infof("Endpoint info %s", req.EndpointID)
	return nil, nil
}

func (driver *driver) Join(j *network.JoinRequest) (*network.JoinResponse, error) {

	return nil, nil
}

func (driver *driver) Leave(leave *network.LeaveRequest) error {
	common.Log.Debugf("Leave request: %+v", leave)
	return nil
}

func (driver *driver) DiscoverNew(disco *network.DiscoveryNotification) error {
	common.Log.Debugf("Dicovery new notification: %+v", disco)
	return nil
}

func (driver *driver) DiscoverDelete(disco *network.DiscoveryNotification) error {
	common.Log.Debugf("Dicovery delete notification: %+v", disco)
	return nil
}
