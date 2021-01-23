package common

import (
	"sort"

	metricApi "github.com/donghoon-khan/kubeportal/src/app/backend/integration/metric/api"
)

type SortableInt64 []int64

func (s SortableInt64) Len() int           { return len(s) }
func (s SortableInt64) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SortableInt64) Less(i, j int) bool { return s[i] < s[j] }

func AggregateData(metricList []metricApi.Metric, metricName string,
	aggregationName metricApi.AggregationMode) metricApi.Metric {

	_, isAggregateAvailable := metricApi.AggregatingFunctions[aggregationName]
	if !isAggregateAvailable {
		aggregationName = metricApi.DefaultAggregation
	}

	aggrMap, newLabel := AggregatingMapFromDataList(metricList, metricName)
	Xs := SortableInt64{}
	for k := range aggrMap {
		Xs = append(Xs, k)
	}
	newDataPoints := []metricApi.DataPoint{}
	sort.Sort(Xs)
	for _, x := range Xs {
		y := metricApi.AggregatingFunctions[aggregationName](aggrMap[x])
		newDataPoints = append(newDataPoints, metricApi.DataPoint{X: x, Y: y})
	}

	metricPoints := []metricApi.MetricPoint{}
	if len(metricList) == 1 {
		metricPoints = metricList[0].MetricPoints
	}

	return metricApi.Metric{
		DataPoints:   newDataPoints,
		MetricPoints: metricPoints,
		MetricName:   metricName,
		Label:        newLabel,
		Aggregate:    aggregationName,
	}
}

func AggregatingMapFromDataList(metricList []metricApi.Metric, metricName string) (
	map[int64][]int64, metricApi.Label) {
	newLabel := metricApi.Label{}

	aggrMap := make(map[int64][]int64, 0)
	for _, data := range metricList {
		if data.MetricName != metricName {
			continue
		}
		newLabel = newLabel.AddMetricLabel(data.Label)
		for _, dataPoint := range data.DataPoints {
			_, isXPresent := aggrMap[dataPoint.X]
			if !isXPresent {
				aggrMap[dataPoint.X] = []int64{}
			}
			aggrMap[dataPoint.X] = append(aggrMap[dataPoint.X], dataPoint.Y)
		}

	}
	return aggrMap, newLabel
}

func AggregateMetricPromises(metricPromises metricApi.MetricPromises, metricName string,
	aggregations metricApi.AggregationModes, forceLabel metricApi.Label) metricApi.MetricPromises {
	if aggregations == nil || len(aggregations) == 0 {
		aggregations = metricApi.OnlyDefaultAggregation
	}
	result := metricApi.NewMetricPromises(len(aggregations))
	go func() {
		metricList, err := metricPromises.GetMetrics()
		if err != nil {
			result.PutMetrics(metricList, err)
			return
		}
		aggrResult := []metricApi.Metric{}
		for _, aggregation := range aggregations {
			aggregated := AggregateData(metricList, metricName, aggregation)
			if forceLabel != nil {
				aggregated.Label = forceLabel
			}
			aggrResult = append(aggrResult, aggregated)
		}
		result.PutMetrics(aggrResult, nil)
	}()
	return result
}
