package api

import v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type IntegrationID string

const (
	HeapsterIntegrationID IntegrationID = "heapster"
	SidecarIntegrationID  IntegrationID = "sidecar"
)

type Integration interface {
	HealthCheck() error
	ID() IntegrationID
}

type IntegrationState struct {
	Connected   bool    `json:"connected`
	LastChecked v1.Time `json:"lastChecked"`
	Error       error   `json:"error" swaggertype:"string"`
}
