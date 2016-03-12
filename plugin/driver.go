package plugin

import (
	"fmt"
	"sync"

	"github.com/docker/libnetwork/drivers/remote/api"
	"github.com/kopwei/knitmesh/common"
	"github.com/kopwei/knitmesh/plugin/listener"
)

type driver struct {
	version    string
	nameserver string
	sync.RWMutex
	endpoints map[string]struct{}
}

var caps = &api.GetCapabilityResponse{
	Scope: "global",
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

func (driver *driver) GetCapabilities() (*api.GetCapabilityResponse, error) {
	common.Log.Debugf("Get capabilities: responded with %+v", caps)
	return caps, nil
}

func (driver *driver) CreateNetwork(create *api.CreateNetworkRequest) error {
	common.Log.Debugf("Create network request %+v", create)
	common.Log.Infof("Create network %s", create.NetworkID)
	return nil
}

func (driver *driver) DeleteNetwork(delete *api.DeleteNetworkRequest) error {
	common.Log.Debugf("Delete network request: %+v", delete)
	common.Log.Infof("Destroy network %s", delete.NetworkID)
	return nil
}

func (driver *driver) CreateEndpoint(create *api.CreateEndpointRequest) (*api.CreateEndpointResponse, error) {
	common.Log.Debugf("Create endpoint request %+v", create)
	return nil, nil
}

func (driver *driver) DeleteEndpoint(deleteReq *api.DeleteEndpointRequest) error {
	common.Log.Debugf("Delete endpoint request: %+v", deleteReq)
	common.Log.Infof("Delete endpoint %s", deleteReq.EndpointID)
	return nil
}

func (driver *driver) HasEndpoint(endpointID string) bool {
	return false
}

func (driver *driver) EndpointInfo(req *api.EndpointInfoRequest) (*api.EndpointInfoResponse, error) {
	common.Log.Debugf("Endpoint info request: %+v", req)
	common.Log.Infof("Endpoint info %s", req.EndpointID)
	return nil, nil
}

func (driver *driver) JoinEndpoint(j *api.JoinRequest) (*api.JoinResponse, error) {

	return nil, nil
}

func (driver *driver) LeaveEndpoint(leave *api.LeaveRequest) error {
	common.Log.Debugf("Leave request: %+v", leave)
	return nil
}

func (driver *driver) DiscoverNew(disco *api.DiscoveryNotification) error {
	common.Log.Debugf("Dicovery new notification: %+v", disco)
	return nil
}

func (driver *driver) DiscoverDelete(disco *api.DiscoveryNotification) error {
	common.Log.Debugf("Dicovery delete notification: %+v", disco)
	return nil
}
