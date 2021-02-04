package ingress

import (
	networking "k8s.io/api/networking/v1"

	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
)

type IngressCell networking.Ingress

func (self IngressCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
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

func toCells(std []networking.Ingress) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = IngressCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []networking.Ingress {
	std := make([]networking.Ingress, len(cells))
	for i := range std {
		std[i] = networking.Ingress(cells[i].(IngressCell))
	}
	return std
}
