package configmap

import (
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
	api "k8s.io/api/core/v1"
)

type ConfigMapCell api.ConfigMap

func (cmc ConfigMapCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
	switch name {
	case dataselect.NameProperty:
		return dataselect.StdComparableString(cmc.ObjectMeta.Name)
	case dataselect.CreationTimestampProperty:
		return dataselect.StdComparableTime(cmc.ObjectMeta.CreationTimestamp.Time)
	case dataselect.NamespaceProperty:
		return dataselect.StdComparableString(cmc.ObjectMeta.Namespace)
	default:
		return nil
	}
}

func toCells(std []api.ConfigMap) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = ConfigMapCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []api.ConfigMap {
	std := make([]api.ConfigMap, len(cells))
	for i := range std {
		std[i] = api.ConfigMap(cells[i].(ConfigMapCell))
	}
	return std
}
