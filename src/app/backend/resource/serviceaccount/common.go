package serviceaccount

import (
	v1 "k8s.io/api/core/v1"

	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
)

type ServiceAccountCell v1.ServiceAccount

func (self ServiceAccountCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
	switch name {
	case dataselect.NameProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Name)
	case dataselect.CreationTimestampProperty:
		return dataselect.StdComparableTime(self.ObjectMeta.CreationTimestamp.Time)
	case dataselect.NamespaceProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Namespace)
	default:
		return nil
	}
}

func toCells(std []v1.ServiceAccount) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = ServiceAccountCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []v1.ServiceAccount {
	std := make([]v1.ServiceAccount, len(cells))
	for i := range std {
		std[i] = v1.ServiceAccount(cells[i].(ServiceAccountCell))
	}
	return std
}
