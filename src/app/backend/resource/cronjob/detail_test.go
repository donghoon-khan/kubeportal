package cronjob_test

import (
	"reflect"
	"testing"

	batch "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/cronjob"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
)

func TestGetJobDetail(t *testing.T) {
	cases := []struct {
		namespace, name string
		expectedActions []string
		raw             *batch.CronJob
		expected        *cronjob.CronJobDetail
	}{
		{
			namespace,
			name,
			[]string{"get"},
			&batch.CronJob{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
					Labels:    labels,
				},
				Spec: batch.CronJobSpec{
					Suspend: &suspend,
				},
			},
			&cronjob.CronJobDetail{
				CronJob: cronjob.CronJob{
					ObjectMeta: api.ObjectMeta{
						Name:      name,
						Namespace: namespace,
						Labels:    labels,
					},
					TypeMeta: api.TypeMeta{Kind: api.ResourceKindCronJob},
					Suspend:  &suspend,
				},
			},
		},
	}

	for _, c := range cases {
		fakeClient := fake.NewSimpleClientset(c.raw)
		dataselect.DefaultDataSelectWithMetrics.MetricQuery = dataselect.NoMetrics
		actual, _ := cronjob.GetCronJobDetail(fakeClient, c.namespace, c.name)

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
			t.Errorf("GetCronJobDetail() got:\n%#v,\nexpected:\n%#v", actual, c.expected)
		}
	}
}
