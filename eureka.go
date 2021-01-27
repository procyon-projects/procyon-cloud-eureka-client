package eureka

type PortWrapper struct {
	Enabled bool
	Port    int
}

type DataCenterName string

const (
	DataCenterNetflix DataCenterName = "Netflix"
	DataCenterAmazon  DataCenterName = "Amazon"
	DataCenterMyOwn   DataCenterName = "MyOwn"
)

type DataCenterInfo struct {
	Name DataCenterName
}

type InstanceStatus string

const (
	InstanceStatusUp           InstanceStatus = "UP"
	InstanceStatusDown         InstanceStatus = "DOWN"
	InstanceStatusStarting     InstanceStatus = "STARTING"
	InstanceStatusOutOfService InstanceStatus = "OUT OF SERVICE"
	InstanceStatusUnknown      InstanceStatus = "UNKNOWN"
)

type Metadata struct {
}

type LeaseInfo struct {
	RenewalIntervalInSecs      int
	DurationInSecs             int
	RegistrationTimestamp      int
	LastRenewalTimestamp       int
	LastRenewalTimestampLegacy int
	EvictionTimestamp          int
	ServiceUpTimestamp         int
}

type ActionType string

const (
	ActionAdded    ActionType = "ADDED"
	ActionModified ActionType = "MODIFIED"
	ActionDeleted  ActionType = "DELETED"
)

type InstanceInfo struct {
	InstanceId                    string
	App                           string
	AppGroupName                  string
	IpAddr                        string
	Port                          PortWrapper
	SecurePort                    PortWrapper
	HomePageUrl                   string
	StatusPageUrl                 string
	HealthCheckUrl                string
	SecureHealthCheckUrl          string
	VipAddress                    string
	SecureVipAddress              string
	CountryId                     int
	DataCenterInfo                DataCenterInfo
	HostName                      string
	Status                        InstanceStatus
	OverriddenStatus              InstanceStatus
	LeaseInfo                     LeaseInfo
	IsCoordinatingDiscoveryServer bool
	Metadata                      Metadata
	LastUpdatedTimestamp          int
	LastDirtyTimestamp            int
	ActionType                    ActionType
	AsgName                       string
}
