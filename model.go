package eureka

type ApplicationsResource struct {
	Applications *Applications `json:"applications,omitempty" xml:"applications,omitempty"`
}

type ApplicationResource struct {
	Application *Application `json:"application,omitempty" xml:"application,omitempty"`
}

type InstanceResource struct {
	InstanceInfo *InstanceInfo `json:"instance,omitempty" xml:"instance,omitempty"`
}

type Error struct {
	Message string `json:"error,omitempty" xml:"error,omitempty"`
}
