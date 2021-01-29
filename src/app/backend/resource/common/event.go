package common

import (
	api "github.com/donghoon-khan/kubeportal/src/app/backend/api"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EventList struct {
	ListMeta api.ListMeta `json:"listMeta"`
	Events   []Event      `json:"events"`
	Errors   []error      `json:"errors" swaggertype:"array,string"`
}

type Event struct {
	ObjectMeta      api.ObjectMeta `json:"objectMeta"`
	TypeMeta        api.TypeMeta   `json:"typeMeta"`
	Message         string         `json:"message"`
	SourceComponent string         `json:"sourceComponent"`
	SourceHost      string         `json:"sourceHost"`
	SubObject       string         `json:"object"`
	Count           int32          `json:"count"`
	FirstSeen       v1.Time        `json:"firstSeen"`
	LastSeen        v1.Time        `json:"lastSeen"`
	Reason          string         `json:"reason"`
	Type            string         `json:"type"`
}
