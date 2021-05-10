package eureka

import (
	cloud "github.com/procyon-projects/procyon-cloud"
)

type DiscoveryClient struct {
	httpClient HttpClient
}

func newDiscoveryClient(httpClient HttpClient) DiscoveryClient {
	return DiscoveryClient{
		httpClient,
	}
}

func (discoveryClient DiscoveryClient) GetDescription() string {
	return "Procyon Eureka Discovery Client"
}

func (discoveryClient DiscoveryClient) GetServiceInstances(serviceId string) []cloud.ServiceInstance {
	return nil
}

func (discoveryClient DiscoveryClient) GetServices() []string {
	applications, err := discoveryClient.httpClient.GetApplications()

	names := make([]string, 0)

	if err != nil {
		return names
	}

	for _, application := range applications.Applications {
		names = append(names, application.Name)
	}

	return names
}
