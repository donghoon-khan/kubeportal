package service

import (
	"log"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
)

type Service struct {
	ObjectMeta        api.ObjectMeta    `json:"objectMeta"`
	TypeMeta          api.TypeMeta      `json:"typeMeta"`
	InternalEndpoint  common.Endpoint   `json:"internalEndpoint"`
	ExternalEndpoints []common.Endpoint `json:"externalEndpoints"`
	Selector          map[string]string `json:"selector"`
	Type              v1.ServiceType    `json:"type"`
	ClusterIP         string            `json:"clusterIP"`
}

type ServiceList struct {
	ListMeta api.ListMeta `json:"listMeta"`
	Services []Service    `json:"services"`
	Errors   []error      `json:"errors"`
}

func GetServiceList(kubernetes kubernetes.Interface, nsQuery *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery) (*ServiceList, error) {
	log.Print("Getting list of all services in the cluster")

	channels := &common.ResourceChannels{
		ServiceList: common.GetServiceListChannel(kubernetes, nsQuery, 1),
	}

	return GetServiceListFromChannels(channels, dsQuery)
}

func GetServiceListFromChannels(channels *common.ResourceChannels,
	dsQuery *dataselect.DataSelectQuery) (*ServiceList, error) {
	services := <-channels.ServiceList.List
	err := <-channels.ServiceList.Error
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	return CreateServiceList(services.Items, nonCriticalErrors, dsQuery), nil
}

func CreateServiceList(services []v1.Service, nonCriticalErrors []error, dsQuery *dataselect.DataSelectQuery) *ServiceList {
	serviceList := &ServiceList{
		Services: make([]Service, 0),
		ListMeta: api.ListMeta{TotalItems: len(services)},
		Errors:   nonCriticalErrors,
	}

	serviceCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(services), dsQuery)
	services = fromCells(serviceCells)
	serviceList.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for _, service := range services {
		serviceList.Services = append(serviceList.Services, toService(&service))
	}

	return serviceList
}

func toService(service *v1.Service) Service {
	return Service{
		ObjectMeta:        api.NewObjectMeta(service.ObjectMeta),
		TypeMeta:          api.NewTypeMeta(api.ResourceKindService),
		InternalEndpoint:  common.GetInternalEndpoint(service.Name, service.Namespace, service.Spec.Ports),
		ExternalEndpoints: common.GetExternalEndpoints(service),
		Selector:          service.Spec.Selector,
		ClusterIP:         service.Spec.ClusterIP,
		Type:              service.Spec.Type,
	}
}
