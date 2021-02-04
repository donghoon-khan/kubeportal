package serviceaccount

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
)

type ServiceAccount struct {
	api.ObjectMeta `json:"objectMeta"`
	api.TypeMeta   `json:"typeMeta"`
}

type ServiceAccountList struct {
	api.ListMeta `json:"listMeta"`
	Items        []ServiceAccount `json:"items"`
	Errors       []error          `json:"errors"`
}

func GetServiceAccountList(kubernetes kubernetes.Interface, namespace *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery) (*ServiceAccountList, error) {
	saList, err := kubernetes.CoreV1().ServiceAccounts(namespace.ToRequestParam()).List(context.TODO(),
		api.ListEverything)

	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	return toServiceAccountList(saList.Items, nonCriticalErrors, dsQuery), nil
}

func toServiceAccount(sa *v1.ServiceAccount) ServiceAccount {
	return ServiceAccount{
		ObjectMeta: api.NewObjectMeta(sa.ObjectMeta),
		TypeMeta:   api.NewTypeMeta(api.ResourceKindServiceAccount),
	}
}

func toServiceAccountList(serviceAccounts []v1.ServiceAccount, nonCriticalErrors []error,
	dsQuery *dataselect.DataSelectQuery) *ServiceAccountList {
	newServiceAccountList := &ServiceAccountList{
		ListMeta: api.ListMeta{TotalItems: len(serviceAccounts)},
		Items:    make([]ServiceAccount, 0),
		Errors:   nonCriticalErrors,
	}

	saCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(serviceAccounts), dsQuery)
	serviceAccounts = fromCells(saCells)

	newServiceAccountList.ListMeta = api.ListMeta{TotalItems: filteredTotal}
	for _, sa := range serviceAccounts {
		newServiceAccountList.Items = append(newServiceAccountList.Items, toServiceAccount(&sa))
	}

	return newServiceAccountList
}
