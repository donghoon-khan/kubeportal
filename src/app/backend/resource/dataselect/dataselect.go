package dataselect

import (
	"log"
	"sort"

	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
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
	GenericDataList           []DataCell
	DataSelectQuery           *DataSelectQuery
	CachedResources           *metricApi.CachedResources
	CumulativeMetricsPromises metricApi.MetricPromises
	MetricsPromises           metricApi.MetricPromises
}

func (dataSelector DataSelector) Len() int { return len(dataSelector.GenericDataList) }

func (dataSelector DataSelector) Swap(i, j int) {
	dataSelector.GenericDataList[i], dataSelector.GenericDataList[j] =
		dataSelector.GenericDataList[j], dataSelector.GenericDataList[i]
}

func (dataSelector DataSelector) Less(i, j int) bool {
	for _, sortBy := range dataSelector.DataSelectQuery.SortQuery.SortByList {
		a := dataSelector.GenericDataList[i].GetProperty(sortBy.Property)
		b := dataSelector.GenericDataList[j].GetProperty(sortBy.Property)
		if a == nil || b == nil {
			break
		}
		cmp := a.Compare(b)
		if cmp == 0 {
			continue
		} else {
			return (cmp == -1 && sortBy.Ascending) || (cmp == 1 && !sortBy.Ascending)
		}
	}
	return false
}

func (dataSelector *DataSelector) Sort() *DataSelector {
	sort.Sort(*dataSelector)
	return dataSelector
}

func (dataSelector *DataSelector) Filter() *DataSelector {
	filteredList := []DataCell{}

	for _, c := range dataSelector.GenericDataList {
		matches := true
		for _, filterBy := range dataSelector.DataSelectQuery.FilterQuery.FilterByList {
			v := c.GetProperty(filterBy.Property)
			if v == nil || !v.Contains(filterBy.Value) {
				matches = false
				break
			}
		}
		if matches {
			filteredList = append(filteredList, c)
		}
	}
	dataSelector.GenericDataList = filteredList
	return dataSelector
}

func (dataSelector *DataSelector) getMetrics(metricClient metricApi.MetricClient) ([]metricApi.MetricPromises,
	error) {
	metricPromises := make([]metricApi.MetricPromises, 0)

	if metricClient == nil {
		return metricPromises, errors.NewInternal("No metric client provided. Skipping metrics.")
	}

	metricNames := dataSelector.DataSelectQuery.MetricQuery.MetricNames
	if metricNames == nil {
		return metricPromises, errors.NewInternal("No metrics specified. skipping metrics.")
	}

	selectors := make([]metricApi.ResourceSelector, len(dataSelector.GenericDataList))
	for i, dataCell := range dataSelector.GenericDataList {
		metricDataCell, ok := dataCell.(MetricDataCell)
		if !ok {
			log.Printf("Data cell does not implement MetricDataCell. Skipping. %v", dataCell)
			continue
		}
		selectors[i] = *metricDataCell.GetResourceSelector()
	}

	for _, metricName := range metricNames {
		promises := metricClient.DownloadMetric(selectors, metricName, dataSelector.CachedResources)
		metricPromises = append(metricPromises, promises)
	}

	return metricPromises, nil
}

func (dataSelector *DataSelector) GetMetrics(metricClient metricApi.MetricClient) *DataSelector {
	metricPromisesList, err := dataSelector.getMetrics(metricClient)
	if err != nil {
		log.Print(err)
		return dataSelector
	}

	metricPromises := make(metricApi.MetricPromises, 0)
	for _, promises := range metricPromisesList {
		metricPromises = append(metricPromises, promises...)
	}

	dataSelector.MetricsPromises = metricPromises
	return dataSelector
}

func (dataSelector *DataSelector) GetCumulativeMetrics(metricClient metricApi.MetricClient) *DataSelector {
	metricPromisesList, err := dataSelector.getMetrics(metricClient)
	if err != nil {
		log.Print(err)
		return dataSelector
	}

	metricNames := dataSelector.DataSelectQuery.MetricQuery.MetricNames
	if metricNames == nil {
		log.Print("No metrics specified. Skipping metrics.")
		return dataSelector
	}

	aggregations := dataSelector.DataSelectQuery.MetricQuery.Aggregations
	if aggregations == nil {
		aggregations = metricApi.OnlyDefaultAggregation
	}

	metricPromises := make(metricApi.MetricPromises, 0)
	for i, metricName := range metricNames {
		promises := metricClient.AggregateMetrics(metricPromisesList[i], metricName, aggregations)
		metricPromises = append(metricPromises, promises...)
	}

	dataSelector.CumulativeMetricsPromises = metricPromises
	return dataSelector
}

func (dataSelector *DataSelector) Paginate() *DataSelector {
	pQuery := dataSelector.DataSelectQuery.PaginationQuery
	dataList := dataSelector.GenericDataList
	startIndex, endIndex := pQuery.GetPaginationSettings(len(dataList))

	// Return all items if provided settings do not meet requirements
	if !pQuery.IsValidPagination() {
		return dataSelector
	}
	// Return no items if requested page does not exist
	if !pQuery.IsPageAvailable(len(dataSelector.GenericDataList), startIndex) {
		dataSelector.GenericDataList = []DataCell{}
		return dataSelector
	}

	dataSelector.GenericDataList = dataList[startIndex:endIndex]
	return dataSelector
}

func GenericDataSelect(dataList []DataCell, dsQuery *DataSelectQuery) []DataCell {
	SelectableData := DataSelector{
		GenericDataList: dataList,
		DataSelectQuery: dsQuery,
	}
	return SelectableData.Sort().Paginate().GenericDataList
}

func GenericDataSelectWithFilter(dataList []DataCell, dsQuery *DataSelectQuery) ([]DataCell, int) {
	SelectableData := DataSelector{
		GenericDataList: dataList,
		DataSelectQuery: dsQuery,
	}
	filtered := SelectableData.Filter()
	filteredTotal := len(filtered.GenericDataList)
	processed := filtered.Sort().Paginate()
	return processed.GenericDataList, filteredTotal
}
