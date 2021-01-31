package eureka

import (
	core "github.com/procyon-projects/procyon-core"
	"net"
	"os"
	"strconv"
)

const (
	unknown       = "unknown"
	DefaultPrefix = "/eureka"
	DefaultUrl    = "http://localhost:8761" + DefaultPrefix + "/"
	DefaultZone   = "defaultZone"

	securePort    = 443
	nonSecurePort = 80

	watchPrefix        = "/watch"
	statusPageUrlPath  = watchPrefix + "/info"
	homePageUrlPath    = "/"
	healthCheckUrlPath = watchPrefix + "/health"
)

type ClientProperties struct {
	RegistryWithEureka bool `json:"registryWithEureka,omitempty" yaml:"registryWithEureka,omitempty"`
	FetchRegistry      bool `json:"fetchRegistry,omitempty" yaml:"fetchRegistry,omitempty"`
}

func newClientProperties() *ClientProperties {
	return &ClientProperties{
		RegistryWithEureka: true,
		FetchRegistry:      true,
	}
}

func (clientConfiguration *ClientProperties) GetConfigurationPrefix() string {
	return "procyon.cloud.eureka.client"
}

type InstanceProperties struct {
	ApplicationName      string         `json:"appName,omitempty" yaml:"appName,omitempty"`
	ApplicationGroupName string         `json:"appGroupName,omitempty" yaml:"appGroupName,omitempty"`
	IpAddr               string         `json:"ipAddr,omitempty" yaml:"ipAddr,omitempty"`
	DataCenterInfo       DataCenterInfo `json:"dataCenterInfo,omitempty" yaml:"dataCenterInfo,omitempty"`
	SecurePort           int            `json:"securePort,omitempty" yaml:"securePort,omitempty"`
	NonSecurePort        int            `json:"nonSecurePort,omitempty" yaml:"nonSecurePort,omitempty"`
	NonSecurePortEnabled bool           `json:"nonSecurePortEnabled,omitempty" yaml:"nonSecurePortEnabled,omitempty"`
	SecurePortEnabled    bool           `json:"securePortEnabled,omitempty" yaml:"securePortEnabled,omitempty"`
	InstanceId           string         `json:"instanceId,omitempty" yaml:"instanceId,omitempty"`
	StatusPageUrl        string         `json:"statusPageUrl,omitempty" yaml:"statusPageUrl,omitempty"`
	HomePageUrl          string         `json:"homePageUrl,omitempty" yaml:"homePageUrl,omitempty"`
	HealthCheckUrl       string         `json:"healthCheckUrl,omitempty" yaml:"healthCheckUrl,omitempty"`
	Hostname             string         `json:"hostname,omitempty" yaml:"hostname,omitempty"`
}

func newInstanceProperties(environment core.Environment) *InstanceProperties {
	instanceProperties := &InstanceProperties{
		ApplicationName:      unknown,
		ApplicationGroupName: unknown,
		DataCenterInfo: DataCenterInfo{
			DataCenterMyOwn,
			"com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
		},
		SecurePort:           securePort,
		NonSecurePort:        nonSecurePort,
		NonSecurePortEnabled: true,
		SecurePortEnabled:    false,
		StatusPageUrl:        statusPageUrlPath,
		HomePageUrl:          homePageUrlPath,
		HealthCheckUrl:       healthCheckUrlPath,
	}
	instanceProperties.initialize(environment)
	return instanceProperties
}

func (instanceProperties *InstanceProperties) initialize(environment core.Environment) {
	appName := environment.GetProperty("procyon.application.name", unknown).(string)
	instanceProperties.ApplicationName = appName
	instanceProperties.ApplicationGroupName = appName

	// instance id, hostname and ipAddr
	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	instanceProperties.Hostname = hostName
	instanceProperties.IpAddr = instanceProperties.getFirstNonLoopbackIpAddr(hostName)

	namePart := instanceProperties.combineParts(hostName, appName, ":")
	port := environment.GetProperty("server.port", "8080").(string)
	instanceProperties.InstanceId = instanceProperties.combineParts(namePart, port, ":")

	// non secure port
	var parsedPort int64
	parsedPort, err = strconv.ParseInt(port, 10, 32)
	instanceProperties.NonSecurePort = int(parsedPort)
}

func (instanceProperties *InstanceProperties) combineParts(firstPart, secondPart, separator string) string {
	combined := ""
	if firstPart != "" && secondPart != "" {
		combined = firstPart + separator + secondPart
	} else if firstPart != "" {
		combined = firstPart
	} else if secondPart != "" {
		combined = secondPart
	}
	return combined
}

func (instanceProperties *InstanceProperties) getFirstNonLoopbackIpAddr(hostName string) string {
	addresses, err := net.LookupIP(hostName)
	if err != nil {
		return unknown
	} else {
		for _, ip := range addresses {
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}
			return ip.String()
		}
	}
	return unknown
}

func (instanceProperties *InstanceProperties) GetConfigurationPrefix() string {
	return "procyon.cloud.eureka.instance"
}
