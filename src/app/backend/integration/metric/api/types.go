package api

import (
	"fmt"
	"time"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	integrationapi "github.com/donghoon-khan/kubeportal/src/app/backend/integration/api"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type MetricClient interface {
	DownloadMetric(selectors []ResourceSelector, metricName string,
		cachedResources *CachedResources) MetricPromises
	DownloadMetrics(selectors []ResourceSelector, metricNames []string,
		cachedResources *CachedResources) MetricPromises
	AggregateMetrics(metrics MetricPromises, metricName string,
		aggregations AggregationModes) MetricPromises

	integrationapi.Integration
}

type CachedResources struct {
	Pods []v1.Pod
}

var NoResourceCache = &CachedResources{}

type AggregationMode string

const (
	SumAggregation     = "sum"
	MaxAggregation     = "max"
	MinAggregation     = "min"
	DefaultAggregation = SumAggregation
)

type AggregationModes []AggregationMode

var OnlySumAggregation = AggregationModes{SumAggregation}
var OnlyDefaultAggregation = AggregationModes{DefaultAggregation}

var AggregatingFunctions = map[AggregationMode]func([]int64) int64{
	SumAggregation: SumAggregate,
	MaxAggregation: MaxAggregate,
	MinAggregation: MinAggregate,
}

var DerivedResources = map[api.ResourceKind]api.ResourceKind{
	api.ResourceKindDeployment:            api.ResourceKindPod,
	api.ResourceKindReplicaSet:            api.ResourceKindPod,
	api.ResourceKindReplicationController: api.ResourceKindPod,
	api.ResourceKindStatefulSet:           api.ResourceKindPod,
	api.ResourceKindJob:                   api.ResourceKindPod,
}

type ResourceSelector struct {
	Namespace    string
	ResourceType api.ResourceKind
	ResourceName string
	Selector     map[string]string
	UID          types.UID
}

const (
	CpuUsage    = "cpu/usage_rate"
	MemoryUsage = "memory/usage"
)

type DataPoints []DataPoint

type DataPoint struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

type MetricPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     uint64    `json:"value"`
}

type Label map[api.ResourceKind][]types.UID

func (self Label) AddMetricLabel(other Label) Label {
	if other == nil {
		return self
	}

	uniqueMap := map[types.UID]bool{}
	for _, v := range self {
		for _, t := range v {
			uniqueMap[t] = true
		}
	}

	for k, v := range other {
		for _, t := range v {
			if _, exists := uniqueMap[t]; !exists {
				self[k] = append(self[k], t)
			}
		}
	}
	return self
}

type Metric struct {
	DataPoints   `json:"dataPoints"`
	MetricPoints []MetricPoint `json:"metricPoints"`
	MetricName   string        `json:"metricName"`
	Label        `json:"-"`
	Aggregate    AggregationMode `json:"aggregation,omitempty"`
}

type SidecarMetric struct {
	DataPoints   `json:"dataPoints"`
	MetricPoints []MetricPoint `json:"metricPoints"`
	MetricName   string        `json:"metricName"`
	UIDs         []string      `json:"uids"`
}

type SidecarMetricResultList struct {
	Items []SidecarMetric `json:"items"`
}

type MetricResultList struct {
	Items []Metric `json:"items"`
}

func (metric *SidecarMetric) AddMetricPoint(item MetricPoint) []MetricPoint {
	metric.MetricPoints = append(metric.MetricPoints, item)
	return metric.MetricPoints
}

func (metric *Metric) AddMetricPoint(item MetricPoint) []MetricPoint {
	metric.MetricPoints = append(metric.MetricPoints, item)
	return metric.MetricPoints
}

func (self Metric) String() string {
	return "{\nDataPoints: " + fmt.Sprintf("%v", self.DataPoints) +
		"\nMetricPoints: " + fmt.Sprintf("%v", self.MetricPoints) +
		"\nMetricName: " + self.MetricName +
		"\nLabel: " + fmt.Sprintf("%v", self.Label) +
		"\nAggregate: " + fmt.Sprintf("%v", self.Aggregate)
}

type MetricPromise struct {
	Metric chan *Metric
	Error  chan error
}

func (self MetricPromise) GetMetric() (*Metric, error) {
	err := <-self.Error
	if err != nil {
		return nil, err
	}
	return <-self.Metric, nil
}

func NewMetricPromise() MetricPromise {
	return MetricPromise{
		Metric: make(chan *Metric, 1),
		Error:  make(chan error, 1),
	}
}

type MetricPromises []MetricPromise

func (self MetricPromises) GetMetrics() ([]Metric, error) {
	result := make([]Metric, 0)

	for _, metricPromise := range self {
		metric, err := metricPromise.GetMetric()
		if err != nil {
			continue
		}

		if metric == nil {
			continue
		}

		result = append(result, *metric)
	}

	return result, nil
}

func (self MetricPromises) PutMetrics(metrics []Metric, err error) {
	for i, metricPromise := range self {
		if err != nil {
			metricPromise.Metric <- nil
		} else {
			metricPromise.Metric <- &metrics[i]
		}
		metricPromise.Error <- err
	}
}

func NewMetricPromises(length int) MetricPromises {
	result := make(MetricPromises, length)
	for i := 0; i < length; i++ {
		result[i] = NewMetricPromise()
	}
	return result
}

func SumAggregate(values []int64) int64 {
	result := int64(0)
	for _, e := range values {
		result += e
	}
	return result
}

func MaxAggregate(values []int64) int64 {
	result := values[0]
	for _, e := range values {
		if e > result {
			result = e
		}
	}
	return result
}

func MinAggregate(values []int64) int64 {
	result := values[0]
	for _, e := range values {
		if e < result {
			result = e
		}
	}
	return result
}
