package clusterrolebinding

import (
	"context"

	rbac "k8s.io/api/rbac/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sClient "k8s.io/client-go/kubernetes"
)

type ClusterRoleBindingDetail struct {
	ClusterRoleBinding `json:",inline"`
	Subjects           []rbac.Subject `json:"subjects,omitempty" protobuf:"bytes,2,rep,name=subjects"`
	RoleRef            rbac.RoleRef   `json:"roleRef" protobuf:"bytes,3,opt,name=roleRef"`
	Errors             []error        `json:"errors"`
}

func GetClusterRoleBindingDetail(kubernetes k8sClient.Interface, name string) (*ClusterRoleBindingDetail, error) {
	rawObject, err := kubernetes.RbacV1().ClusterRoleBindings().Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	cr := toClusterRoleBindingDetail(*rawObject)
	return &cr, nil
}

func toClusterRoleBindingDetail(cr rbac.ClusterRoleBinding) ClusterRoleBindingDetail {
	return ClusterRoleBindingDetail{
		ClusterRoleBinding: toClusterRoleBinding(cr),
		Subjects:           cr.Subjects,
		RoleRef:            cr.RoleRef,
		Errors:             []error{},
	}
}
