package service

import (
	v1 "k8s.io/api/core/v1"

	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
)

type ServiceCell v1.Service

func (self ServiceCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
	switch name {
	case dataselect.NameProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Name)
	case dataselect.CreationTimestampProperty:
		return dataselect.StdComparableTime(self.ObjectMeta.CreationTimestamp.Time)
	case dataselect.NamespaceProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Namespace)
	case dataselect.TypeProperty:
		return dataselect.StdComparableString(self.Spec.Type)
	default:
		return nil
	}
}

func toCells(std []v1.Service) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = ServiceCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []v1.Service {
	std := make([]v1.Service, len(cells))
	for i := range std {
		std[i] = v1.Service(cells[i].(ServiceCell))
	}
	return std
}
