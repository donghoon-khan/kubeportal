package metric

import (
	"fmt"
	"log"
	"time"

	integrationApi "github.com/donghoon-khan/kubeportal/src/app/backend/integration/api"
	metricApi "github.com/donghoon-khan/kubeportal/src/app/backend/integration/metric/api"
	k8sApi "github.com/donghoon-khan/kubeportal/src/app/backend/kubernetes/api"
	"k8s.io/apimachinery/pkg/util/wait"
)

type MetricManager interface {
	AddClient(metricApi.MetricClient) MetricManager
	Client() metricApi.MetricClient
	Enable(integrationApi.IntegrationID) error
	EnableWithRetry(id integrationApi.IntegrationID, period time.Duration)
	List() []integrationApi.Integration
	//ConfigureSidecar(host string) MetricManager
	//ConfigureHeapster(host string) MetricManager
}

type metricManager struct {
	manager k8sApi.KubernetesManager
	clients map[integrationApi.IntegrationID]metricApi.MetricClient
	active  metricApi.MetricClient
}

func (mManager *metricManager) AddClient(client metricApi.MetricClient) MetricManager {
	if client != nil {
		mManager.clients[client.ID()] = client
	}
	return mManager
}

func (mManager *metricManager) Client() metricApi.MetricClient {
	return mManager.active
}

func (mManager *metricManager) Enable(id integrationApi.IntegrationID) error {
	metricClient, exists := mManager.clients[id]
	if !exists {
		return fmt.Errorf("No metric client found for integration id: %s", id)
	}

	err := metricClient.HealthCheck()
	if err != nil {
		return fmt.Errorf("Health check failed: %s", err.Error())
	}

	mManager.active = metricClient
	return nil
}

func (mManager *metricManager) EnableWithRetry(id integrationApi.IntegrationID, period time.Duration) {
	go wait.Forever(func() {
		metricClient, exists := mManager.clients[id]
		if !exists {
			log.Printf("Metric client with given id %s does not exist.", id)
			return
		}

		err := metricClient.HealthCheck()
		if err != nil {
			mManager.active = nil
			log.Printf("Metric client health check failed: %s. Retrying in %d seconds.", err, period)
			return
		}

		if mManager.active == nil {
			log.Printf("Successful request to %s", id)
			mManager.active = metricClient
		}
	}, period*time.Second)
}

func (mManager *metricManager) List() []integrationApi.Integration {
	result := make([]integrationApi.Integration, 0)
	for _, c := range mManager.clients {
		result = append(result, c.(integrationApi.Integration))
	}

	return result
}

func NewMetricManager(manager k8sApi.KubernetesManager) MetricManager {
	return &metricManager{
		manager: manager,
		clients: make(map[integrationApi.IntegrationID]metricApi.MetricClient),
	}
}
