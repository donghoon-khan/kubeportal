package common

import (
	api "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Condition struct {
	Type               string              `json:"type"`
	Status             api.ConditionStatus `json:"status"`
	LastProbeTime      v1.Time             `json:"lastProbeTime"`
	LastTransitionTime v1.Time             `json:"lastTransitionTime"`
	Reason             string              `json:"reson"`
	Message            string              `json:"message"`
}
