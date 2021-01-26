package dataselect

import (
	metricApi "github.com/donghoon-khan/kubeportal/src/app/backend/integration/metric/api"
)

type DataSelectQuery struct {
	PaginationQuery *PaginationQuery
	SortQuery       *SortQuery
	FilterQuery     *FilterQuery
	MetricQuery     *MetricQuery
}

var NoMetrics = NewMetricQuery(nil, nil)

var StandardMetrics = NewMetricQuery([]string{metricApi.CpuUsage, metricApi.MemoryUsage},
	metricApi.OnlySumAggregation)

type MetricQuery struct {
	MetricNames  []string
	Aggregations metricApi.AggregationModes
}

func NewMetricQuery(metricNames []string, aggregations metricApi.AggregationModes) *MetricQuery {
	return &MetricQuery{
		MetricNames:  metricNames,
		Aggregations: aggregations,
	}
}

type SortQuery struct {
	SortByList []SortBy
}

type SortBy struct {
	Property  PropertyName
	Ascending bool
}

var NoSort = &SortQuery{
	SortByList: []SortBy{},
}

type FilterQuery struct {
	FilterByList []FilterBy
}

type FilterBy struct {
	Property PropertyName
	Value    ComparableValue
}

var NoFilter = &FilterQuery{
	FilterByList: []FilterBy{},
}

var NoDataSelect = NewDataSelectQuery(NoPagination, NoSort, NoFilter, NoMetrics)

var StdMetricsDataSelect = NewDataSelectQuery(NoPagination, NoSort, NoFilter, StandardMetrics)

var DefaultDataSelect = NewDataSelectQuery(DefaultPagination, NoSort, NoFilter, NoMetrics)

var DefaultDataSelectWithMetrics = NewDataSelectQuery(DefaultPagination, NoSort, NoFilter, StandardMetrics)

func NewDataSelectQuery(paginationQuery *PaginationQuery, sortQuery *SortQuery,
	filterQuery *FilterQuery, graphQuery *MetricQuery) *DataSelectQuery {
	return &DataSelectQuery{
		PaginationQuery: paginationQuery,
		SortQuery:       sortQuery,
		FilterQuery:     filterQuery,
		MetricQuery:     graphQuery,
	}
}

func NewSortQuery(sortByListRaw []string) *SortQuery {
	if sortByListRaw == nil || len(sortByListRaw)%2 == 1 {
		return NoSort
	}
	sortByList := []SortBy{}
	for i := 0; i+1 < len(sortByListRaw); i += 2 {
		var ascending bool
		orderOption := sortByListRaw[i]
		if orderOption == "a" {
			ascending = true
		} else if orderOption == "d" {
			ascending = false
		} else {
			return NoSort
		}

		propertyName := sortByListRaw[i+1]
		sortBy := SortBy{
			Property:  PropertyName(propertyName),
			Ascending: ascending,
		}
		sortByList = append(sortByList, sortBy)
	}
	return &SortQuery{
		SortByList: sortByList,
	}
}

func NewFilterQuery(filterByListRaw []string) *FilterQuery {
	if filterByListRaw == nil || len(filterByListRaw)%2 == 1 {
		return NoFilter
	}
	filterByList := []FilterBy{}
	for i := 0; i+1 < len(filterByListRaw); i += 2 {
		propertyName := filterByListRaw[i]
		propertyValue := filterByListRaw[i+1]
		filterBy := FilterBy{
			Property: PropertyName(propertyName),
			Value:    StdComparableString(propertyValue),
		}
		filterByList = append(filterByList, filterBy)
	}
	return &FilterQuery{
		FilterByList: filterByList,
	}
}
