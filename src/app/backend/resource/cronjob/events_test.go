package cronjob_test

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/cronjob"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
)

func TestGetJobEvents(t *testing.T) {
	cases := []struct {
		namespace, name string
		eventList       *v1.EventList
		expectedActions []string
		expected        *common.EventList
	}{
		{
			namespace,
			name,
			&v1.EventList{
				Items: []v1.Event{{
					Message: eventMessage,
					ObjectMeta: metaV1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
						Labels:    labels,
					}},
				}},
			[]string{"list"},
			&common.EventList{
				ListMeta: api.ListMeta{
					TotalItems: 1,
				},
				Events: []common.Event{{
					TypeMeta: api.TypeMeta{
						Kind: api.ResourceKindEvent,
					},
					ObjectMeta: api.ObjectMeta{
						Name:      name,
						Namespace: namespace,
						Labels:    labels,
					},
					Message: eventMessage,
					Type:    v1.EventTypeNormal,
				}}},
		},
	}

	for _, c := range cases {
		fakeClient := fake.NewSimpleClientset(c.eventList)

		actual, _ := cronjob.GetCronJobEvents(fakeClient, dataselect.NoDataSelect, c.namespace, c.name)

		actions := fakeClient.Actions()
		if len(actions) != len(c.expectedActions) {
			t.Errorf("Unexpected actions: %v, expected %d actions got %d", actions,
				len(c.expectedActions), len(actions))
			continue
		}

		for i, verb := range c.expectedActions {
			if actions[i].GetVerb() != verb {
				t.Errorf("Unexpected action: %+v, expected %s",
					actions[i], verb)
			}
		}

		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("TestGetJobEvents(client,metricClient,%#v, %#v) == \ngot: %#v, \nexpected %#v",
				c.namespace, c.name, actual, c.expected)
		}
	}
}
