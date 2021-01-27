package clusterrole

import (
	"context"

	rbac "k8s.io/api/rbac/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sClient "k8s.io/client-go/kubernetes"
)

type ClusterRoleDetail struct {
	ClusterRole `json:",inline"`
	Rules       []rbac.PolicyRule `json:"rules"`
	Errors      []error           `json:"errors" swaggertype:"array,string"`
}

func GetClusterRoleDetail(client k8sClient.Interface, name string) (*ClusterRoleDetail, error) {
	rawObject, err := client.RbacV1().ClusterRoles().Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	cr := toClusterRoleDetail(*rawObject)
	return &cr, nil
}

func toClusterRoleDetail(cr rbac.ClusterRole) ClusterRoleDetail {
	return ClusterRoleDetail{
		ClusterRole: toClusterRole(cr),
		Rules:       cr.Rules,
		Errors:      []error{},
	}
}
