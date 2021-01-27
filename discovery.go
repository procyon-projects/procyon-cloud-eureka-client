package eureka

import (
	cloud "github.com/procyon-projects/procyon-cloud"
)

type DiscoveryClient struct {
}

func (discoveryClient DiscoveryClient) GetDescription() string {
	return "Procyon Eureka Discovery Client"
}

func (discoveryClient DiscoveryClient) GetServiceInstances(serviceId string) []cloud.ServiceInstance {
	return nil
}

func (discoveryClient DiscoveryClient) GetServices() []string {
	return nil
}
