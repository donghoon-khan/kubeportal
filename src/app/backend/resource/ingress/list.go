package ingress

import (
	"context"

	networking "k8s.io/api/networking/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
)

type Ingress struct {
	api.ObjectMeta `json:"objectMeta"`
	api.TypeMeta   `json:"typeMeta"`
	Endpoints      []common.Endpoint `json:"endpoints"`
}

type IngressList struct {
	api.ListMeta `json:"listMeta"`
	Items        []Ingress `json:"items"`
	Errors       []error   `json:"errors"`
}

func GetIngressList(kubernetes kubernetes.Interface, namespace *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery) (*IngressList, error) {
	ingressList, err := kubernetes.NetworkingV1().Ingresses(namespace.ToRequestParam()).List(context.TODO(), api.ListEverything)

	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	return toIngressList(ingressList.Items, nonCriticalErrors, dsQuery), nil
}

func getEndpoints(ingress *networking.Ingress) []common.Endpoint {
	endpoints := make([]common.Endpoint, 0)
	if len(ingress.Status.LoadBalancer.Ingress) > 0 {
		for _, status := range ingress.Status.LoadBalancer.Ingress {
			endpoint := common.Endpoint{}
			if status.Hostname != "" {
				endpoint.Host = status.Hostname
			} else if status.IP != "" {
				endpoint.Host = status.IP
			}
			endpoints = append(endpoints, endpoint)
		}
	}
	return endpoints
}

func toIngress(ingress *networking.Ingress) Ingress {
	return Ingress{
		ObjectMeta: api.NewObjectMeta(ingress.ObjectMeta),
		TypeMeta:   api.NewTypeMeta(api.ResourceKindIngress),
		Endpoints:  getEndpoints(ingress),
	}
}

func toIngressList(ingresses []networking.Ingress, nonCriticalErrors []error, dsQuery *dataselect.DataSelectQuery) *IngressList {
	newIngressList := &IngressList{
		ListMeta: api.ListMeta{TotalItems: len(ingresses)},
		Items:    make([]Ingress, 0),
		Errors:   nonCriticalErrors,
	}

	ingresCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(ingresses), dsQuery)
	ingresses = fromCells(ingresCells)
	newIngressList.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for _, ingress := range ingresses {
		newIngressList.Items = append(newIngressList.Items, toIngress(&ingress))
	}

	return newIngressList
}
