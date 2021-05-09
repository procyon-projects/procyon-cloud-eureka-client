package eureka

import cloud "github.com/procyon-projects/procyon-cloud"

type ServiceRegistry struct {
}

func (serviceRegistry ServiceRegistry) Register(instance cloud.ServiceInstance) {

}

func (serviceRegistry ServiceRegistry) Deregister(instance cloud.ServiceInstance) {

}

func (serviceRegistry ServiceRegistry) SetStatus(instance cloud.ServiceInstance, status string) {

}

func (serviceRegistry ServiceRegistry) GetStatus(instance cloud.ServiceInstance) interface{} {
	return nil
}
