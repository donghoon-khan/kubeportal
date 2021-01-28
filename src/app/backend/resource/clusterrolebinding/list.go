package clusterrolebinding

import (
	"log"

	rbac "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
)

type ClusterRoleBindingList struct {
	ListMeta api.ListMeta         `json:"listMeta"`
	Items    []ClusterRoleBinding `json:"items"`
	Errors   []error              `json:"errors" swaggertype:"array,string"`
}

type ClusterRoleBinding struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta   api.TypeMeta   `json:"typeMeta"`
}

func GetClusterRoleBindingList(kubernetes kubernetes.Interface,
	dsQuery *dataselect.DataSelectQuery) (*ClusterRoleBindingList, error) {

	log.Print("Getting list of all clusterRoleBindings in the cluster")
	channels := &common.ResourceChannels{
		ClusterRoleBindingList: common.GetClusterRoleBindingListChannel(kubernetes, 1),
	}
	return GetClusterRoleBindingListFromChannels(channels, dsQuery)
}

func GetClusterRoleBindingListFromChannels(channels *common.ResourceChannels,
	dsQuery *dataselect.DataSelectQuery) (*ClusterRoleBindingList, error) {

	clusterRoleBindings := <-channels.ClusterRoleBindingList.List
	err := <-channels.ClusterRoleBindingList.Error
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}
	clusterRoleBindingList := toClusterRoleBindingList(clusterRoleBindings.Items, nonCriticalErrors, dsQuery)
	return clusterRoleBindingList, nil
}

func toClusterRoleBinding(clusterRoleBinding rbac.ClusterRoleBinding) ClusterRoleBinding {
	return ClusterRoleBinding{
		ObjectMeta: api.NewObjectMeta(clusterRoleBinding.ObjectMeta),
		TypeMeta:   api.NewTypeMeta(api.ResourceKindClusterRoleBinding),
	}
}

func toClusterRoleBindingList(clusterRoleBindings []rbac.ClusterRoleBinding, nonCriticalErrors []error,
	dsQuery *dataselect.DataSelectQuery) *ClusterRoleBindingList {

	result := &ClusterRoleBindingList{
		ListMeta: api.ListMeta{TotalItems: len(clusterRoleBindings)},
		Errors:   nonCriticalErrors,
	}

	items := make([]ClusterRoleBinding, 0)
	for _, item := range clusterRoleBindings {
		items = append(items, toClusterRoleBinding(item))
	}

	clusterRoleBindingCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(items), dsQuery)
	result.ListMeta = api.ListMeta{TotalItems: filteredTotal}
	result.Items = fromCells(clusterRoleBindingCells)
	return result
}
