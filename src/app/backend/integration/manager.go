package integration

import (
	"fmt"

	"github.com/donghoon-khan/kubeportal/src/app/backend/integration/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/integration/metric"
	k8sApi "github.com/donghoon-khan/kubeportal/src/app/backend/kubernetes/api"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IntegrationManager interface {
	IntegrationGetter
	GetState(id api.IntegrationID) (*api.IntegrationState, error)
	Metric() metric.MetricManager
}

type integrationManager struct {
	metric metric.MetricManager
}

func (iManager *integrationManager) Metric() metric.MetricManager {
	return iManager.metric
}

func (iManager *integrationManager) GetState(id api.IntegrationID) (*api.IntegrationState, error) {
	for _, i := range iManager.List() {
		if i.ID() == id {
			return iManager.getState(i), nil
		}
	}
	return nil, fmt.Errorf("Integration with given id %s does not exist", id)
}

func (iManager *integrationManager) getState(integration api.Integration) *api.IntegrationState {
	result := &api.IntegrationState{
		Error: integration.HealthCheck(),
	}
	result.Connected = result.Error == nil
	result.LastChecked = v1.Now()
	return result
}

func NewIntegrationManager(kManager k8sApi.KubernetesManager) IntegrationManager {
	return &integrationManager{
		metric: metric.NewMetricManager(kManager),
	}
}
