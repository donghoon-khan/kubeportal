package clusterrolebinding

import "github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"

type ClusterRoleBindingCell ClusterRoleBinding

func (crbc ClusterRoleBindingCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
	switch name {
	case dataselect.NameProperty:
		return dataselect.StdComparableString(crbc.ObjectMeta.Name)
	case dataselect.CreationTimestampProperty:
		return dataselect.StdComparableTime(crbc.ObjectMeta.CreationTimestamp.Time)
	case dataselect.NamespaceProperty:
		return dataselect.StdComparableString(crbc.ObjectMeta.Namespace)
	default:
		return nil
	}
}

func toCells(std []ClusterRoleBinding) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = ClusterRoleBindingCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []ClusterRoleBinding {
	std := make([]ClusterRoleBinding, len(cells))
	for i := range std {
		std[i] = ClusterRoleBinding(cells[i].(ClusterRoleBindingCell))
	}
	return std
}
