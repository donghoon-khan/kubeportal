package service

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/endpoint"
)

func TestGetServiceDetail(t *testing.T) {
	cases := []struct {
		service         *v1.Service
		namespace, name string
		expectedActions []string
		expected        *ServiceDetail
	}{
		{
			service: &v1.Service{ObjectMeta: metaV1.ObjectMeta{
				Name: "svc-1", Namespace: "ns-1", Labels: map[string]string{},
			}},
			namespace: "ns-1", name: "svc-1",
			expectedActions: []string{"get", "list"},
			expected: &ServiceDetail{
				Service: Service{
					ObjectMeta: api.ObjectMeta{
						Name:      "svc-1",
						Namespace: "ns-1",
						Labels:    map[string]string{},
					},
					TypeMeta:          api.TypeMeta{Kind: api.ResourceKindService},
					InternalEndpoint:  common.Endpoint{Host: "svc-1.ns-1"},
					ExternalEndpoints: []common.Endpoint{},
				},
				EndpointList: endpoint.EndpointList{
					Endpoints: []endpoint.Endpoint{},
				},
				Errors: []error{},
			},
		},
		{
			service: &v1.Service{
				ObjectMeta: metaV1.ObjectMeta{
					Name:      "svc-2",
					Namespace: "ns-2",
				},
				Spec: v1.ServiceSpec{
					Selector: map[string]string{"app": "app2"},
				},
			},
			namespace: "ns-2", name: "svc-2",
			expectedActions: []string{"get", "list"},
			expected: &ServiceDetail{
				Service: Service{
					ObjectMeta: api.ObjectMeta{
						Name:      "svc-2",
						Namespace: "ns-2",
					},
					Selector:          map[string]string{"app": "app2"},
					TypeMeta:          api.TypeMeta{Kind: api.ResourceKindService},
					InternalEndpoint:  common.Endpoint{Host: "svc-2.ns-2"},
					ExternalEndpoints: []common.Endpoint{},
				},

				EndpointList: endpoint.EndpointList{
					Endpoints: []endpoint.Endpoint{},
				},
				Errors: []error{},
			},
		},
	}

	for _, c := range cases {
		fakeClient := fake.NewSimpleClientset(c.service)
		actual, _ := GetServiceDetail(fakeClient, c.namespace, c.name)
		actions := fakeClient.Actions()

		if len(actions) != len(c.expectedActions) {
			t.Errorf("Unexpected actions: %v, expected %d actions got %d", actions,
				len(c.expectedActions), len(actions))
			continue
		}

		for i, verb := range c.expectedActions {
			if actions[i].GetVerb() != verb {
				t.Errorf("Unexpected action: %+v, expected %s", actions[i], verb)
			}
		}

		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("GetServiceDetail(client, %#v, %#v) == \ngot %#v, \nexpected %#v", c.namespace,
				c.name, actual, c.expected)
		}
	}
}
