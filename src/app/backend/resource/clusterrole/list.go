package clusterrole

import (
	"log"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes"
)

type ClusterRoleList struct {
	ListMeta api.ListMeta  `json:"listMeta"`
	Items    []ClusterRole `json:"items"`
	Errors   []error       `json:"errors"`
}

type ClusterRole struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta   api.TypeMeta   `json:"typeMeta"`
}

func GetClusterRoleList(kubernetes kubernetes.Interface, dsQuery *dataselect.DataSelectQuery) (*ClusterRoleList, error) {
	log.Println("Getting list of RBAC roles")
	channels := &common.ResourceChannels{
		ClusterRoleList: common.GetClusterRoleListChannel(kubernetes, 1),
	}
	return GetClusterRoleListFromChannels(channels, dsQuery)
}

func GetClusterRoleListFromChannels(channels *common.ResourceChannels, dsQuery *dataselect.DataSelectQuery) (*ClusterRoleList, error) {
	clusterRoles := <-channels.ClusterRoleList.List
	err := <-channels.ClusterRoleList.Error
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	result := toClusterRoleLists(clusterRoles.Items, nonCriticalErrors, dsQuery)
	return result, nil
}

func toClusterRole(role rbac.ClusterRole) ClusterRole {
	return ClusterRole{
		ObjectMeta: api.NewObjectMeta(role.ObjectMeta),
		TypeMeta:   api.NewTypeMeta(api.ResourceKindClusterRole),
	}
}

func toClusterRoleLists(clusterRoles []rbac.ClusterRole, nonCriticalErrors []error,
	dsQuery *dataselect.DataSelectQuery) *ClusterRoleList {
	result := &ClusterRoleList{
		ListMeta: api.ListMeta{TotalItems: len(clusterRoles)},
		Errors:   nonCriticalErrors,
	}

	items := make([]ClusterRole, 0)
	for _, item := range clusterRoles {
		items = append(items, toClusterRole(item))
	}

	roleCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(items), dsQuery)
	result.ListMeta = api.ListMeta{TotalItems: filteredTotal}
	result.Items = fromCells(roleCells)
	return result
}
