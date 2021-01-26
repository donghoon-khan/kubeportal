package dataselect

import (
	metricApi "github.com/donghoon-khan/kubeportal/src/app/backend/integration/metric/api"
)

type DataCell interface {
	GetProperty(PropertyName) ComparableValue
}

type MetricDataCell interface {
	DataCell
	GetResourceSelector() *metricApi.ResourceSelector
}

type ComparableValue interface {
	Compare(ComparableValue) int
	Contains(ComparableValue) bool
}

type DataSelector struct {
	GenericDataList            []DataCell
	DataSelectQuery            *DataSelectQuery
	CachedResources            *metricApi.CachedResources
	CumulativeMetricsPromisses metricApi.MetricPromises
	MetricsPromises            metricApi.MetricPromises
}

func (dataSelector DataSelector) Len() int { return len(dataSelector.GenericDataList) }

func (dataSelector DataSelector) Swap(i, j int) {
	dataSelector.GenericDataList[i], dataSelector.GenericDataList[j] =
		dataSelector.GenericDataList[j], dataSelector.GenericDataList[i]
}
