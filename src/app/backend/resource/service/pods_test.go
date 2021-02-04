package service

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	metricApi "github.com/donghoon-khan/kubeportal/src/app/backend/integration/metric/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/pod"
)

func TestGetServicePods(t *testing.T) {
	cases := []struct {
		namespace, name string
		service         *v1.Service
		podList         *v1.PodList
		expected        *pod.PodList
	}{
		{
			"ns-1",
			"svc-1",
			&v1.Service{ObjectMeta: metaV1.ObjectMeta{
				Name: "svc-1", Namespace: "ns-1", Labels: map[string]string{"app": "test"},
			}, Spec: v1.ServiceSpec{Selector: map[string]string{}}},
			&v1.PodList{Items: []v1.Pod{
				{ObjectMeta: metaV1.ObjectMeta{
					Name:      "pod-1",
					Namespace: "ns-1",
					UID:       "test-uid",
				}},
			}},
			&pod.PodList{
				ListMeta:          api.ListMeta{TotalItems: 1},
				CumulativeMetrics: make([]metricApi.Metric, 0),
				Pods: []pod.Pod{
					{
						ObjectMeta: api.ObjectMeta{
							Name:      "pod-1",
							UID:       "test-uid",
							Namespace: "ns-1"},
						TypeMeta: api.TypeMeta{Kind: api.ResourceKindPod},
						Status:   string(v1.PodUnknown),
						Warnings: []common.Event{},
					},
				},
				Errors: []error{},
			},
		},
	}
	for _, c := range cases {
		fakeClient := fake.NewSimpleClientset(c.service, c.podList)

		actual, _ := GetServicePods(fakeClient, nil, c.namespace, c.name, dataselect.NoDataSelect)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("GetServicePods == \ngot %#v, \nexpected %#v", actual, c.expected)
		}

	}
}
