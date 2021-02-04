package endpoint

import (
	"log"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
)

type Endpoint struct {
	ObjectMeta api.ObjectMeta    `json:"objectMeta"`
	TypeMeta   api.TypeMeta      `json:"typeMeta"`
	Host       string            `json:"host"`
	NodeName   *string           `json:"nodeName"`
	Ready      bool              `json:"ready"`
	Ports      []v1.EndpointPort `json:"ports"`
}

func GetServiceEndpoints(kubernetes kubernetes.Interface, namespace, name string) (*EndpointList, error) {
	endpointList := &EndpointList{
		Endpoints: make([]Endpoint, 0),
		ListMeta:  api.ListMeta{TotalItems: 0},
	}

	serviceEndpoints, err := GetEndpoints(kubernetes, namespace, name)
	if err != nil {
		return endpointList, err
	}

	endpointList = toEndpointList(serviceEndpoints)
	log.Printf("Found %d endpoints related to %s service in %s namespace", len(endpointList.Endpoints), name, namespace)
	return endpointList, nil
}

func GetEndpoints(kubernetes kubernetes.Interface, namespace, name string) ([]v1.Endpoints, error) {
	fieldSelector, err := fields.ParseSelector("metadata.name" + "=" + name)
	if err != nil {
		return nil, err
	}

	channels := &common.ResourceChannels{
		EndpointList: common.GetEndpointListChannelWithOptions(kubernetes,
			common.NewSameNamespaceQuery(namespace),
			metaV1.ListOptions{
				LabelSelector: labels.Everything().String(),
				FieldSelector: fieldSelector.String(),
			},
			1),
	}

	endpointList := <-channels.EndpointList.List
	if err := <-channels.EndpointList.Error; err != nil {
		return nil, err
	}

	return endpointList.Items, nil
}

func toEndpoint(address v1.EndpointAddress, ports []v1.EndpointPort, ready bool) *Endpoint {
	return &Endpoint{
		TypeMeta: api.NewTypeMeta(api.ResourceKindEndpoint),
		Host:     address.IP,
		Ports:    ports,
		Ready:    ready,
		NodeName: address.NodeName,
	}
}
