package serviceaccount

import (
	"context"
	"log"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ServiceAccountDetail struct {
	ServiceAccount `json:",inline"`
	Errors         []error `json:"errors"`
}

func GetServiceAccountDetail(kubernetes kubernetes.Interface, namespace, name string) (*ServiceAccountDetail, error) {
	log.Printf("Getting details of %s service account in %s namespace", name, namespace)

	raw, err := kubernetes.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return getServiceAccountDetail(raw), nil
}

func getServiceAccountDetail(sa *v1.ServiceAccount) *ServiceAccountDetail {
	return &ServiceAccountDetail{
		ServiceAccount: toServiceAccount(sa),
	}
}
