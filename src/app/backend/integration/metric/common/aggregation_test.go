package common

import (
	"reflect"
	"testing"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	metricApi "github.com/donghoon-khan/kubeportal/src/app/backend/integration/metric/api"
	"k8s.io/apimachinery/pkg/types"
)

func getMetricPromises(metrics []metricApi.Metric) metricApi.MetricPromises {
	metricPromises := metricApi.NewMetricPromises(len(metrics))
	metricPromises.PutMetrics(metrics, nil)
	return metricPromises
}

func TestAggregateMetricPromises(t *testing.T) {
	cases := []struct {
		info         string
		promises     metricApi.MetricPromises
		metricName   string
		aggregations metricApi.AggregationModes
		forceLabel   metricApi.Label
		expected     []metricApi.Metric
	}{
		{
			"should return empty metric when metric name not provided",
			getMetricPromises([]metricApi.Metric{
				{
					DataPoints: []metricApi.DataPoint{
						{X: 0, Y: 5},
						{X: 5, Y: 10},
						{X: 10, Y: 15},
					},
					MetricName: "test-metric",
					Label: metricApi.Label{
						api.ResourceKindPod: []types.UID{"U1"},
					},
				},
			}),
			"",
			metricApi.OnlyDefaultAggregation,
			nil,
			[]metricApi.Metric{
				{
					DataPoints: metricApi.DataPoints{},
					MetricName: "",
					Label:      metricApi.Label{},
					Aggregate:  metricApi.SumAggregation,
				},
			},
		},
		{
			"should override label",
			getMetricPromises([]metricApi.Metric{}),
			"",
			metricApi.OnlyDefaultAggregation,
			metricApi.Label{api.ResourceKindPod: []types.UID{"overridden-uid"}},
			[]metricApi.Metric{
				{
					DataPoints:   metricApi.DataPoints{},
					MetricPoints: []metricApi.MetricPoint{},
					MetricName:   "",
					Label: metricApi.Label{
						api.ResourceKindPod: []types.UID{"overridden-uid"},
					},
					Aggregate: metricApi.SumAggregation,
				},
			},
		},
		{
			"should use default aggregation mode when nothing is provided",
			getMetricPromises([]metricApi.Metric{
				{
					DataPoints: []metricApi.DataPoint{
						{X: 0, Y: 5},
						{X: 5, Y: 10},
						{X: 10, Y: 15},
					},
					MetricName: "test-metric",
					Label: metricApi.Label{
						api.ResourceKindPod: []types.UID{"U1"},
					},
				},
				{
					DataPoints: []metricApi.DataPoint{
						{X: 0, Y: 5},
						{X: 5, Y: 10},
						{X: 10, Y: 15},
					},
					MetricName: "test-metric",
					Label: metricApi.Label{
						api.ResourceKindPod: []types.UID{"U2"},
					},
				},
			}),
			"test-metric",
			nil,
			nil,
			[]metricApi.Metric{
				{
					DataPoints: []metricApi.DataPoint{
						{X: 0, Y: 10},
						{X: 5, Y: 20},
						{X: 10, Y: 30},
					},
					MetricPoints: []metricApi.MetricPoint{},
					MetricName:   "test-metric",
					Label: metricApi.Label{
						api.ResourceKindPod: []types.UID{"U1", "U2"},
					},
					Aggregate: metricApi.SumAggregation,
				},
			},
		},
		{
			"should use sum aggregation mode",
			getMetricPromises([]metricApi.Metric{
				{
					DataPoints: []metricApi.DataPoint{
						{X: 0, Y: 5},
						{X: 5, Y: 10},
						{X: 10, Y: 15},
					},
					MetricName: "test-metric",
					Label: metricApi.Label{
						api.ResourceKindPod: []types.UID{"U1"},
					},
				},
				{
					DataPoints: []metricApi.DataPoint{
						{X: 0, Y: 5},
						{X: 5, Y: 10},
						{X: 10, Y: 15},
					},
					MetricName: "test-metric",
					Label: metricApi.Label{
						api.ResourceKindPod: []types.UID{"U2"},
					},
				},
			}),
			"test-metric",
			metricApi.OnlySumAggregation,
			nil,
			[]metricApi.Metric{
				{
					DataPoints: []metricApi.DataPoint{
						{X: 0, Y: 10},
						{X: 5, Y: 20},
						{X: 10, Y: 30},
					},
					MetricPoints: []metricApi.MetricPoint{},
					MetricName:   "test-metric",
					Label: metricApi.Label{
						api.ResourceKindPod: []types.UID{"U1", "U2"},
					},
					Aggregate: metricApi.SumAggregation,
				},
			},
		},
		{
			"should use min aggregation mode",
			getMetricPromises([]metricApi.Metric{
				{
					DataPoints: []metricApi.DataPoint{
						{X: 0, Y: 5},
						{X: 5, Y: 10},
						{X: 10, Y: 15},
					},
					MetricName: "test-metric",
					Label: metricApi.Label{
						api.ResourceKindPod: []types.UID{"U1"},
					},
				},
				{
					DataPoints: []metricApi.DataPoint{
						{X: 0, Y: 10},
						{X: 5, Y: 15},
						{X: 10, Y: 20},
					},
					MetricName: "test-metric",
					Label: metricApi.Label{
						api.ResourceKindPod: []types.UID{"U2"},
					},
				},
			}),
			"test-metric",
			metricApi.AggregationModes{metricApi.MinAggregation},
			nil,
			[]metricApi.Metric{
				{
					DataPoints: []metricApi.DataPoint{
						{X: 0, Y: 5},
						{X: 5, Y: 10},
						{X: 10, Y: 15},
					},
					MetricPoints: []metricApi.MetricPoint{},
					MetricName:   "test-metric",
					Label: metricApi.Label{
						api.ResourceKindPod: []types.UID{"U1", "U2"},
					},
					Aggregate: metricApi.MinAggregation,
				},
			},
		},
		{
			"should use max aggregation mode",
			getMetricPromises([]metricApi.Metric{
				{
					DataPoints: []metricApi.DataPoint{
						{X: 0, Y: 5},
						{X: 5, Y: 10},
						{X: 10, Y: 15},
					},
					MetricName: "test-metric",
					Label: metricApi.Label{
						api.ResourceKindPod: []types.UID{"U1"},
					},
				},
				{
					DataPoints: []metricApi.DataPoint{
						{X: 0, Y: 10},
						{X: 5, Y: 15},
						{X: 10, Y: 20},
					},
					MetricName: "test-metric",
					Label: metricApi.Label{
						api.ResourceKindPod: []types.UID{"U2"},
					},
				},
			}),
			"test-metric",
			metricApi.AggregationModes{metricApi.MaxAggregation},
			nil,
			[]metricApi.Metric{
				{
					DataPoints: []metricApi.DataPoint{
						{X: 0, Y: 10},
						{X: 5, Y: 15},
						{X: 10, Y: 20},
					},
					MetricPoints: []metricApi.MetricPoint{},
					MetricName:   "test-metric",
					Label: metricApi.Label{
						api.ResourceKindPod: []types.UID{"U1", "U2"},
					},
					Aggregate: metricApi.MaxAggregation,
				},
			},
		},
	}

	for _, c := range cases {
		promises := AggregateMetricPromises(c.promises, c.metricName, c.aggregations,
			c.forceLabel)
		metrics, _ := promises.GetMetrics()
		if !reflect.DeepEqual(metrics, c.expected) {
			t.Errorf("Test Case: %s. Failed to aggregate metrics. Expected: %+v, but got %+v",
				c.info, c.expected, metrics)
		}
	}
}
