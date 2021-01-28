package handler

import (
	"net/http"

	"github.com/donghoon-khan/kubeportal/src/app/backend/auth"
	authApi "github.com/donghoon-khan/kubeportal/src/app/backend/auth/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/handler/parser"
	"github.com/donghoon-khan/kubeportal/src/app/backend/integration"
	k8sApi "github.com/donghoon-khan/kubeportal/src/app/backend/kubernetes/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/clusterrole"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/clusterrolebinding"
	"github.com/emicklei/go-restful"
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

	/* ClusterRole */
	apiV1Ws.Route(
		apiV1Ws.GET("/clusterrole").
			To(apiHandler.handleGetClusterRoleList).
			Writes(clusterrole.ClusterRoleList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/clusterrole/{name}").
			To(apiHandler.handleGetClusterRoleDetail).
			Writes(clusterrole.ClusterRoleDetail{}))

	/* ClusterRoleBinding */
	apiV1Ws.Route(
		apiV1Ws.GET("/clusterrolebinding").
			To(apiHandler.handleGetClusterRoleBindingList).
			Writes(clusterrolebinding.ClusterRoleBindingList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/clusterrolebinding/{name}").
			To(apiHandler.handleGetClusterRoleBindingDetail).
			Writes(clusterrolebinding.ClusterRoleBindingDetail{}))

	return wsContainer, nil
}

// handleGetClusterRoleList godoc
// @Tags Kubernetes
// @Summary List of ClusterRole
// @Description Returns the list of ClusterRole from kubernetes cluster
// @Accept  json
// @Produce  json
// @Router /clusterrole [GET]
// @Success 200 {object} clusterrole.ClusterRoleList
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetClusterRoleList(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := clusterrole.GetClusterRoleList(k8sClient, dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

// handleGetClusterRoleDetail godoc
// @Tags Kubernetes
// @Summary Detail of ClusterRole
// @Description Returns the Detail of ClusterRole from kubernetes cluster
// @Accept  json
// @Produce  json
// @Router /clusterrole/{name} [GET]
// @Param name path string true "Name of ClusterRole"
// @Success 200 {object} clusterrole.ClusterRoleDetail
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetClusterRoleDetail(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
	}

	name := request.PathParameter("name")
	result, err := clusterrole.GetClusterRoleDetail(k8sClient, name)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

// handleGetClusterRoleBindingList godoc
// @Tags Kubernetes
// @Summary List of ClusterRoleBinding
// @Description Returns the list of ClusterRoleBinding from kubernetes cluster
// @Accept  json
// @Produce  json
// @Router /clusterrolebinding [GET]
// @Success 200 {object} clusterrolebinding.ClusterRoleBindingList
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetClusterRoleBindingList(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := clusterrolebinding.GetClusterRoleBindingList(k8sClient, dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

// handleGetClusterRoleBindingDetail godoc
// @Tags Kubernetes
// @Summary Detail of ClusterRoleBinding
// @Description Returns the Detail of ClusterRoleBinding from kubernetes cluster
// @Accept  json
// @Produce  json
// @Router /clusterrolebinding/{name} [GET]
// @Param name path string true "Name of ClusterRoleBinding"
// @Success 200 {object} clusterrolebinding.ClusterRoleBindingDetail
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetClusterRoleBindingDetail(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	name := request.PathParameter("name")
	result, err := clusterrolebinding.GetClusterRoleBindingDetail(k8sClient, name)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}
