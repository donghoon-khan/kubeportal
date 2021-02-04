package handler

import (
	"strings"

	"github.com/emicklei/go-restful/v3"

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
	authManager authApi.AuthManager) (*restful.Container, error) {

	apiHandler := APIHandler{iManager: iManager, kManager: kManager}
	wsContainer := restful.NewContainer()
	wsContainer.EnableContentEncoding(true)

	k8sWs := new(restful.WebService)
	InstallFilters(k8sWs, kManager)

	k8sWs.Path("/api/v1/kubernetes").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Param(
			k8sWs.QueryParameter("itemPerPage",
				"The number of items per page can be configured adding a query parameter named itemsPerPage"+
					" `e.g. itemPerPage=10`").
				DataType("int")).
		Param(k8sWs.HeaderParameter("token", "token")).
		Param(
			k8sWs.QueryParameter("page", "The number of page `e.g. page=1`").DataType("int")).
		Param(
			k8sWs.QueryParameter("sortBy",
				"be used to sort the result list {ascending or descending},"+
					"{name or creationTimestamp or namespace or status or type}"+
					" `e.g. sortBy=asc,name`").
				DataType("Collection of string(csv)")).
		Param(
			k8sWs.QueryParameter("filterBy",
				"be used to get filterd result list "+
					"{name or creationTimestamp or namespace or statusor type}"+
					" `e.g. filterBy=namespace,kube-system`").
				DataType("Collection of string(csv)"))

	apiHandler.installClusterRole(k8sWs)
	apiHandler.installClusterRoleBinding(k8sWs)
	apiHandler.installConfigMap(k8sWs)
	apiHandler.installCronJob(k8sWs)
	apiHandler.installIngress(k8sWs)
	apiHandler.installSecret(k8sWs)
	apiHandler.installPersistentVolumeClaim(k8sWs)
	apiHandler.installPod(k8sWs)
	apiHandler.installNode(k8sWs)
	apiHandler.installService(k8sWs)
	wsContainer.Add(k8sWs)

	integrationHandler := integration.NewIntegrationHandler(iManager)
	integrationWs := new(restful.WebService)
	integrationWs.Path("/api/v1/integration").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	integrationHandler.Install(integrationWs)
	wsContainer.Add(integrationWs)

	authHandler := auth.NewAuthHandler(authManager)
	authWs := new(restful.WebService)
	authWs.Path("/api/v1/authentication").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	authHandler.Install(authWs)
	wsContainer.Add(authWs)

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
