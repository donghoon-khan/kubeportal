package handler

import (
	"net/http"
	"strings"

	"github.com/emicklei/go-restful"

	"github.com/donghoon-khan/kubeportal/src/app/backend/auth"
	authApi "github.com/donghoon-khan/kubeportal/src/app/backend/auth/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/integration"
	k8sApi "github.com/donghoon-khan/kubeportal/src/app/backend/kubernetes/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
)

const (
	RequestLogString  = "[%s] Incoming %s %s %s request from %s: %s"
	ResponseLogString = "[%s] Outcoming response to %s with %d status code"
)

type APIHandler struct {
	iManager integration.IntegrationManager
	kManager k8sApi.KubernetesManager
}

func CreateHttpApiHandler(
	iManager integration.IntegrationManager,
	kManager k8sApi.KubernetesManager,
	authManager authApi.AuthManager) (http.Handler, error) {

	apiHandler := APIHandler{iManager: iManager, kManager: kManager}
	wsContainer := restful.NewContainer()
	wsContainer.EnableContentEncoding(true)

	apiV1Ws := new(restful.WebService)
	InstallFilters(apiV1Ws, kManager)

	apiV1Ws.Path("/api/v1").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	wsContainer.Add(apiV1Ws)

	integrationHandler := integration.NewIntegrationHandler(iManager)
	integrationHandler.Install(apiV1Ws)

	authHandler := auth.NewAuthHandler(authManager)
	authHandler.Install(apiV1Ws)

	apiHandler.installClusterRole(apiV1Ws)
	apiHandler.installClusterRoleBinding(apiV1Ws)
	apiHandler.installConfigMap(apiV1Ws)
	apiHandler.installSecret(apiV1Ws)
	apiHandler.installPod(apiV1Ws)
	apiHandler.installNode(apiV1Ws)

	return wsContainer, nil
}

func parseNamespacePathParameter(request *restful.Request) *common.NamespaceQuery {
	namespace := request.PathParameter("namespace")
	namespaces := strings.Split(namespace, ",")
	var nonEmptyNamespaces []string
	for _, n := range namespaces {
		n = strings.Trim(n, " ")
		if len(n) > 0 {
			nonEmptyNamespaces = append(nonEmptyNamespaces, n)
		}
	}
	return common.NewNamespaceQuery(nonEmptyNamespaces)
}
