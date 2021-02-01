package secret

import (
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
	api "k8s.io/api/core/v1"
)

type SecretCell api.Secret

func (secretCell SecretCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
	switch name {
	case dataselect.NameProperty:
		return dataselect.StdComparableString(secretCell.ObjectMeta.Name)
	case dataselect.CreationTimestampProperty:
		return dataselect.StdComparableTime(secretCell.ObjectMeta.CreationTimestamp.Time)
	case dataselect.NamespaceProperty:
		return dataselect.StdComparableString(secretCell.ObjectMeta.Namespace)
	default:
		return nil
	}
}

func toCells(std []api.Secret) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = SecretCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []api.Secret {
	std := make([]api.Secret, len(cells))
	for i := range std {
		std[i] = api.Secret(cells[i].(SecretCell))
	}
	return std
}
