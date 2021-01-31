package eureka

type ApplicationsResponse struct {
	Applications *Applications `json:"applications,omitempty" xml:"applications,omitempty"`
}

type ApplicationResponse struct {
	Applications *Application `json:"application,omitempty" xml:"application,omitempty"`
}

type InstanceInfoResponse struct {
	InstanceInfo *InstanceInfo `json:"instance,omitempty" xml:"instance,omitempty"`
}
