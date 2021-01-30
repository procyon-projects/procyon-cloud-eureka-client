package eureka

import (
	"net/url"
	"strconv"
)

type ServiceInstance struct {
	instanceInfo *InstanceInfo
}

func newServiceInstance(instanceInfo *InstanceInfo) ServiceInstance {
	if instanceInfo == nil {
		panic("Instance info required")
	}
	return ServiceInstance{
		instanceInfo,
	}
}

func (serviceInstance ServiceInstance) GetInstanceId() string {
	return serviceInstance.instanceInfo.InstanceId
}

func (serviceInstance ServiceInstance) GetServiceId() string {
	return serviceInstance.instanceInfo.AppName
}

func (serviceInstance ServiceInstance) GetURL() url.URL {
	return url.URL{
		Scheme: serviceInstance.GetScheme(),
		Host:   serviceInstance.GetHost() + ":" + strconv.Itoa(serviceInstance.GetPort()),
	}
}

func (serviceInstance ServiceInstance) GetScheme() string {
	if serviceInstance.IsSecure() {
		return "https"
	}
	return "http"
}

func (serviceInstance ServiceInstance) GetHost() string {
	return serviceInstance.instanceInfo.HostName
}

func (serviceInstance ServiceInstance) GetPort() int {
	if serviceInstance.IsSecure() {
		return serviceInstance.instanceInfo.SecurePort.Port
	}
	return serviceInstance.instanceInfo.Port.Port
}

func (serviceInstance ServiceInstance) IsSecure() bool {
	if serviceInstance.instanceInfo.SecurePort != nil {
		return true
	}
	return false
}

func (serviceInstance ServiceInstance) GetMetadata() map[string]string {
	return serviceInstance.instanceInfo.Metadata
}
