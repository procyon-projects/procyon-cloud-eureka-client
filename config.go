package eureka

import (
	core "github.com/procyon-projects/procyon-core"
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

	metricPrefix       = "/metrics"
	statusPageUrlPath  = metricPrefix + "/info"
	homePageUrlPath    = "/"
	healthCheckUrlPath = metricPrefix + "/health"
)

type ClientProperties struct {
	RegistryWithEureka bool `json:"registryWithEureka" yaml:"registryWithEureka"`
	FetchRegistry      bool `json:"fetchRegistry" yaml:"fetchRegistry"`
}

func newClientProperties() *ClientProperties {
	return &ClientProperties{
		RegistryWithEureka: true,
		FetchRegistry:      true,
	}
}

func (clientConfiguration *ClientProperties) GetConfigurationPrefix() string {
	return "cloud.eureka.client"
}

type InstanceProperties struct {
	ApplicationName      string         `json:"appName" yaml:"appName"`
	ApplicationGroupName string         `json:"appGroupName" yaml:"appGroupName"`
	DataCenterInfo       DataCenterInfo `json:"dataCenterInfo" yaml:"dataCenterInfo"`
	SecurePort           int            `json:"securePort" yaml:"securePort"`
	NonSecurePort        int            `json:"nonSecurePort" yaml:"nonSecurePort"`
	InstanceId           string         `json:"instanceId" yaml:"instanceId"`
	StatusPageUrl        string         `json:"statusPageUrl" yaml:"statusPageUrl"`
	HomePageUrl          string         `json:"homePageUrl" yaml:"homePageUrl"`
	HealthCheckUrl       string         `json:"healthCheckUrl" yaml:"healthCheckUrl"`
}

func newInstanceProperties(environment core.Environment) *InstanceProperties {
	instanceProperties := &InstanceProperties{
		ApplicationName:      unknown,
		ApplicationGroupName: unknown,
		DataCenterInfo: DataCenterInfo{
			DataCenterMyOwn,
		},
		SecurePort:     securePort,
		NonSecurePort:  nonSecurePort,
		StatusPageUrl:  statusPageUrlPath,
		HomePageUrl:    homePageUrlPath,
		HealthCheckUrl: healthCheckUrlPath,
	}
	instanceProperties.initialize(environment)
	return instanceProperties
}

func (instanceProperties *InstanceProperties) initialize(environment core.Environment) {
	appName := environment.GetProperty("procyon.application.name", unknown).(string)
	instanceProperties.ApplicationName = appName
	instanceProperties.ApplicationGroupName = appName

	// instance id
	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}

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

func (instanceProperties *InstanceProperties) GetConfigurationPrefix() string {
	return "cloud.eureka.instance"
}
