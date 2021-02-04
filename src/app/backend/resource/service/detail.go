package service

import (
	"context"
	"log"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/endpoint"
)

type ServiceDetail struct {
	Service         `json:",inline"`
	EndpointList    endpoint.EndpointList `json:"endpointList"`
	SessionAffinity v1.ServiceAffinity    `json:"sessionAffinity"`
	Errors          []error               `json:"errors"`
}

func GetServiceDetail(kubernetes kubernetes.Interface, namespace, name string) (*ServiceDetail, error) {
	log.Printf("Getting details of %s service in %s namespace", name, namespace)

	serviceData, err := kubernetes.CoreV1().Services(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	endpointList, err := endpoint.GetServiceEndpoints(kubernetes, namespace, name)
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	service := toServiceDetail(serviceData, *endpointList, nonCriticalErrors)
	return &service, nil
}

func toServiceDetail(service *v1.Service, endpointList endpoint.EndpointList, nonCriticalErrors []error) ServiceDetail {
	return ServiceDetail{
		Service:         toService(service),
		EndpointList:    endpointList,
		SessionAffinity: service.Spec.SessionAffinity,
		Errors:          nonCriticalErrors,
	}
}
