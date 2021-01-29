package pod

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetPodEvents(t *testing.T) {
	cases := []struct {
		namespace, podName string
		eventList          *v1.EventList
		podList            *v1.PodList
		expected           *common.EventList
	}{
		{
			"ns-1", "pod-1",
			&v1.EventList{Items: []v1.Event{
				{
					Message: "test-message",
					ObjectMeta: metaV1.ObjectMeta{
						Name: "ev-1", Namespace: "ns-1",
						Labels: map[string]string{"app": "test"},
					},
					InvolvedObject: v1.ObjectReference{UID: "test-uid"}},
			}},
			&v1.PodList{Items: []v1.Pod{
				{ObjectMeta: metaV1.ObjectMeta{
					Name:      "pod-1",
					Namespace: "ns-1",
					UID:       "test-uid",
				}},
			}},
			&common.EventList{
				ListMeta: api.ListMeta{TotalItems: 1},
				Events: []common.Event{{
					TypeMeta: api.TypeMeta{Kind: api.ResourceKindEvent},
					ObjectMeta: api.ObjectMeta{Name: "ev-1", Namespace: "ns-1",
						Labels: map[string]string{"app": "test"}},
					Message: "test-message",
					Type:    v1.EventTypeNormal,
				}},
				Errors: []error{},
			},
		},
	}

	for _, c := range cases {
		fakeClient := fake.NewSimpleClientset(c.podList, c.eventList)

		actual, _ := GetEventsForPod(fakeClient, dataselect.NoDataSelect, c.namespace, c.podName)

		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("GetEventsForPods == \ngot %#v, \nexpected %#v", actual,
				c.expected)
		}
	}
}
