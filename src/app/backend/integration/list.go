package integration

import (
	integrationApi "github.com/donghoon-khan/kubeportal/src/app/backend/integration/api"
)

type IntegrationGetter interface {
	List() []integrationApi.Integration
}

func (iManager *integrationManager) List() []integrationApi.Integration {
	result := make([]integrationApi.Integration, 0)
	result = append(result, iManager.Metric().List()...)
	return result
}
