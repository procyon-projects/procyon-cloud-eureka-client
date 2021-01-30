package eureka

import (
	core "github.com/procyon-projects/procyon-core"
	"net/url"
	"os"
	"sync"
)

type InstanceInfoProvider interface {
	GetInstanceInfo() *InstanceInfo
}

type DefaultInstanceInfoProvider struct {
	instanceProperties InstanceProperties
	instanceInfo       *InstanceInfo
	instanceInfoMu     sync.RWMutex
	environment        core.Environment
}

func newDefaultInstanceInfoProvider(instanceProperties InstanceProperties, environment core.Environment) *DefaultInstanceInfoProvider {
	return &DefaultInstanceInfoProvider{
		instanceProperties: instanceProperties,
		environment:        environment,
	}
}

func (provider *DefaultInstanceInfoProvider) GetInstanceInfo() *InstanceInfo {
	provider.instanceInfoMu.Lock()
	defer provider.instanceInfoMu.Unlock()
	if provider.instanceInfo != nil {
		return provider.instanceInfo
	}

	port := provider.environment.GetProperty("server.port", "8080").(string)
	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	instanceInfo := &InstanceInfo{
		InstanceId:     provider.instanceProperties.InstanceId,
		App:            provider.instanceProperties.ApplicationName,
		AppGroupName:   provider.instanceProperties.ApplicationGroupName,
		IpAddr:         provider.instanceProperties.IpAddr,
		HomePageUrl:    provider.getUrl(false, hostName, port, provider.instanceProperties.HomePageUrl),
		StatusPageUrl:  provider.getUrl(false, hostName, port, provider.instanceProperties.StatusPageUrl),
		HealthCheckUrl: provider.getUrl(false, hostName, port, provider.instanceProperties.HealthCheckUrl),
		DataCenterInfo: provider.instanceProperties.DataCenterInfo,
		HostName:       provider.instanceProperties.Hostname,
	}

	if provider.instanceProperties.NonSecurePortEnabled {
		instanceInfo.VipAddress = provider.getUrl(false, hostName, port, "")
		instanceInfo.Port = PortWrapper{
			Enabled: true,
			Port:    provider.instanceProperties.NonSecurePort,
		}
	}

	if provider.instanceProperties.SecurePortEnabled {
		instanceInfo.SecureVipAddress = provider.getUrl(true, hostName, port, "")
		instanceInfo.SecurePort = PortWrapper{
			Enabled: true,
			Port:    provider.instanceProperties.SecurePort,
		}
		instanceInfo.SecureHealthCheckUrl = provider.getUrl(true, hostName, port, provider.instanceProperties.HealthCheckUrl)
	}

	provider.instanceInfo = instanceInfo
	return instanceInfo
}

func (provider *DefaultInstanceInfoProvider) getUrl(isSecure bool, hostName string, port string, urlPath string) string {
	scheme := "http"
	if isSecure {
		scheme = "https"
	}

	stringPort := ""
	if port != "80" && port != "443" {
		stringPort = ":" + port
	}
	host := hostName + stringPort

	result := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   urlPath,
	}
	return result.String()
}
