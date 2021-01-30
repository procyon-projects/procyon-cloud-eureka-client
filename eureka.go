package eureka

type PortWrapper struct {
	Enabled bool `json:"@enabled" xml:",chardata"`
	Port    int  `json:"$" xml:"enabled,attr"`
}

type DataCenterName string

const (
	DataCenterNetflix DataCenterName = "Netflix"
	DataCenterAmazon  DataCenterName = "Amazon"
	DataCenterMyOwn   DataCenterName = "MyOwn"
)

type DataCenterInfo struct {
	Name DataCenterName `json:"name" xml:"name"`
}

type InstanceStatus string

const (
	InstanceStatusUp           InstanceStatus = "UP"
	InstanceStatusDown         InstanceStatus = "DOWN"
	InstanceStatusStarting     InstanceStatus = "STARTING"
	InstanceStatusOutOfService InstanceStatus = "OUT_OF_SERVICE"
	InstanceStatusUnknown      InstanceStatus = "UNKNOWN"
)

type Metadata struct {
}

type LeaseInfo struct {
	RenewalIntervalInSecs      int `json:"renewalIntervalInSecs,omitempty" xml:"renewalIntervalInSecs,omitempty"`
	DurationInSecs             int `json:"durationInSecs,omitempty" xml:"durationInSecs,omitempty"`
	RegistrationTimestamp      int `json:"registrationTimestamp,omitempty" xml:"registrationTimestamp,omitempty"`
	LastRenewalTimestamp       int `json:"lastRenewalTimestamp,omitempty" xml:"lastRenewalTimestamp,omitempty"`
	LastRenewalTimestampLegacy int `json:"renewalTimestamp,omitempty" xml:"renewalTimestamp,omitempty"`
	EvictionTimestamp          int `json:"evictionTimestamp,omitempty" xml:"evictionTimestamp,omitempty"`
	ServiceUpTimestamp         int `json:"serviceUpTimestamp,omitempty" xml:"serviceUpTimestamp,omitempty"`
}

type ActionType string

const (
	ActionAdded    ActionType = "ADDED"
	ActionModified ActionType = "MODIFIED"
	ActionDeleted  ActionType = "DELETED"
)

type InstanceInfo struct {
	InstanceId                    string         `json:"instanceId,omitempty" xml:"instanceId,omitempty"`
	App                           string         `json:"app" xml:"app"`
	AppGroupName                  string         `json:"appGroupName" xml:"appGroupName"`
	IpAddr                        string         `json:"ipAddr" xml:"ipAddr"`
	Port                          PortWrapper    `json:"port" xml:"port"`
	SecurePort                    PortWrapper    `json:"securePort" xml:"securePort"`
	HomePageUrl                   string         `json:"homePageUrl" xml:"homePageUrl"`
	StatusPageUrl                 string         `json:"statusPageUrl" xml:"statusPageUrl"`
	HealthCheckUrl                string         `json:"healthCheckUrl" xml:"healthCheckUrl"`
	SecureHealthCheckUrl          string         `json:"secureHealthCheckUrl" xml:"secureHealthCheckUrl"`
	VipAddress                    string         `json:"vipAddress" xml:"vipAddress"`
	SecureVipAddress              string         `json:"secureVipAddress" xml:"secureVipAddress"`
	CountryId                     int            `json:"countryId" xml:"countryId"`
	DataCenterInfo                DataCenterInfo `json:"dataCenterInfo" xml:"dataCenterInfo"`
	HostName                      string         `json:"hostName" xml:"hostName"`
	Status                        InstanceStatus `json:"status" xml:"status"`
	OverriddenStatus              InstanceStatus `json:"overriddenstatus" xml:"overriddenstatus"`
	LeaseInfo                     LeaseInfo      `json:"leaseInfo" xml:"leaseInfo"`
	IsCoordinatingDiscoveryServer bool           `json:"isCoordinatingDiscoveryServer" xml:"isCoordinatingDiscoveryServer"`
	Metadata                      Metadata       `json:"metadata" xml:"metadata"`
	LastUpdatedTimestamp          int            `json:"lastUpdatedTimestamp" xml:"lastUpdatedTimestamp"`
	LastDirtyTimestamp            int            `json:"lastDirtyTimestamp" xml:"lastDirtyTimestamp"`
	ActionType                    ActionType     `json:"actionType" xml:"actionType"`
	AsgName                       string         `json:"asgName" xml:"asgName"`
}
