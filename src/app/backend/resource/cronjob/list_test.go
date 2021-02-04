package cronjob_test

import (
	"reflect"
	"testing"

	batch "k8s.io/api/batch/v1beta1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	metricApi "github.com/donghoon-khan/kubeportal/src/app/backend/integration/metric/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/cronjob"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
)

func TestGetCronJobListFromChannels(t *testing.T) {
	cases := []struct {
		raw           batch.CronJobList
		rawError      error
		expected      *cronjob.CronJobList
		expectedError error
	}{
		{
			batch.CronJobList{},
			nil,
			&cronjob.CronJobList{
				ListMeta:          api.ListMeta{},
				CumulativeMetrics: make([]metricApi.Metric, 0),
				Status:            common.ResourceStatus{},
				Items:             []cronjob.CronJob{},
				Errors:            []error{},
			},
			nil,
		},
		{
			batch.CronJobList{},
			customError,
			nil,
			customError,
		},
		{
			batch.CronJobList{
				Items: []batch.CronJob{
					{
						ObjectMeta: metaV1.ObjectMeta{
							Name:      name,
							Namespace: namespace,
							Labels:    labels,
						},
					},
					{
						ObjectMeta: metaV1.ObjectMeta{
							Name:      name,
							Namespace: namespace,
							Labels:    labels,
						},
					},
				},
			},
			nil,
			&cronjob.CronJobList{
				ListMeta:          api.ListMeta{TotalItems: 2},
				CumulativeMetrics: make([]metricApi.Metric, 0),
				Status:            common.ResourceStatus{Failed: 2},
				Items: []cronjob.CronJob{{
					ObjectMeta: api.ObjectMeta{
						Name:      name,
						Namespace: namespace,
						Labels:    labels,
					},
					TypeMeta: api.TypeMeta{Kind: api.ResourceKindCronJob},
				}, {
					ObjectMeta: api.ObjectMeta{
						Name:      name,
						Namespace: namespace,
						Labels:    labels,
					},
					TypeMeta: api.TypeMeta{Kind: api.ResourceKindCronJob},
				}},
				Errors: []error{},
			},
			nil,
		},
	}

	for _, c := range cases {
		channels := &common.ResourceChannels{
			CronJobList: common.CronJobListChannel{
				List:  make(chan *batch.CronJobList, 1),
				Error: make(chan error, 1),
			},
		}

		channels.CronJobList.Error <- c.rawError
		channels.CronJobList.List <- &c.raw

		actual, err := cronjob.GetCronJobListFromChannels(channels, dataselect.NoDataSelect, nil)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("GetCronJobListFromChannels() ==\n %#v\nExpected: %#v", actual, c.expected)
		}
		if !reflect.DeepEqual(err, c.expectedError) {
			t.Errorf("GetCronJobListFromChannels() ==\n %#v\nExpected: %#v", err, c.expectedError)
		}
	}
}
