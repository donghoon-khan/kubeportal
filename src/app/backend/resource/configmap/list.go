package configmap

import (
	"log"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ConfigMapList struct {
	ListMeta api.ListMeta `json:"listMeta"`
	Items    []ConfigMap  `json:"items"`
	Errors   []error      `json:"errors" swaggertype:"array,string"`
}

type ConfigMap struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta   api.TypeMeta   `json:"typeMeta"`
}

func GetConfigMapList(client kubernetes.Interface, nsQuery *common.NamespaceQuery, dsQuery *dataselect.DataSelectQuery) (*ConfigMapList, error) {
	log.Printf("Getting list config maps in the namespace %s", nsQuery.ToRequestParam())
	channels := &common.ResourceChannels{
		ConfigMapList: common.GetConfigMapListChannel(client, nsQuery, 1),
	}

	return GetConfigMapListFromChannels(channels, dsQuery)
}

func GetConfigMapListFromChannels(channels *common.ResourceChannels, dsQuery *dataselect.DataSelectQuery) (*ConfigMapList, error) {
	configMaps := <-channels.ConfigMapList.List
	err := <-channels.ConfigMapList.Error
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	result := toConfigMapList(configMaps.Items, nonCriticalErrors, dsQuery)

	return result, nil
}

func toConfigMap(meta metaV1.ObjectMeta) ConfigMap {
	return ConfigMap{
		ObjectMeta: api.NewObjectMeta(meta),
		TypeMeta:   api.NewTypeMeta(api.ResourceKindConfigMap),
	}
}

func toConfigMapList(configMaps []v1.ConfigMap, nonCriticalErrors []error, dsQuery *dataselect.DataSelectQuery) *ConfigMapList {
	result := &ConfigMapList{
		Items:    make([]ConfigMap, 0),
		ListMeta: api.ListMeta{TotalItems: len(configMaps)},
		Errors:   nonCriticalErrors,
	}

	configMapCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(configMaps), dsQuery)
	configMaps = fromCells(configMapCells)
	result.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for _, item := range configMaps {
		result.Items = append(result.Items, toConfigMap(item.ObjectMeta))
	}

	return result
}
