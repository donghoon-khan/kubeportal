package handler

import (
	"net/http"
	"strings"

	"github.com/donghoon-khan/kubeportal/src/app/backend/auth"
	authApi "github.com/donghoon-khan/kubeportal/src/app/backend/auth/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/handler/parser"
	"github.com/donghoon-khan/kubeportal/src/app/backend/integration"
	k8sApi "github.com/donghoon-khan/kubeportal/src/app/backend/kubernetes/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/clusterrole"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/clusterrolebinding"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/configmap"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/persistentvolumeclaim"
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
		apiV1Ws.GET("/clusterrole/{clusterrole}").
			To(apiHandler.handleGetClusterRoleDetail).
			Writes(clusterrole.ClusterRoleDetail{}))

	/* ClusterRoleBinding */
	apiV1Ws.Route(
		apiV1Ws.GET("/clusterrolebinding").
			To(apiHandler.handleGetClusterRoleBindingList).
			Writes(clusterrolebinding.ClusterRoleBindingList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/clusterrolebinding/{clusterrolebinding}").
			To(apiHandler.handleGetClusterRoleBindingDetail).
			Writes(clusterrolebinding.ClusterRoleBindingDetail{}))

	/* ConfigMap */
	apiV1Ws.Route(
		apiV1Ws.GET("/configmap").
			To(apiHandler.handleGetAllConfigMapList).
			Writes(configmap.ConfigMapList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/configmap/{namespace}").
			To(apiHandler.handleGetConfigMapList).
			Writes(configmap.ConfigMapList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/configmap/{namespace}/{configmap}").
			To(apiHandler.handleGetConfigMapDetail).
			Writes(configmap.ConfigMapDetail{}))

	/* PersistentVolumeClaim */
	apiV1Ws.Route(
		apiV1Ws.GET("/persistentvolumeclaim").
			To(apiHandler.handleGetAllPersistentVolumeClaimList).
			Writes(persistentvolumeclaim.PersistentVolumeClaimList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/persistentvolumeclaim/{namespace}").
			To(apiHandler.handleGetPersistentVolumeClaimList).
			Writes(persistentvolumeclaim.PersistentVolumeClaimList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/persistentvolumeclaim/{namespace}/{persistentvolumeclaim}").
			To(apiHandler.handleGetPersistentVolumeClaimDetail).
			Writes(persistentvolumeclaim.PersistentVolumeClaimDetail{}))

	return wsContainer, nil
}

// handleGetClusterRoleList godoc
// @Tags Kubernetes
// @Summary Get list of ClusterRole
// @Description Returns a list of ClusterRole
// @Accept  json
// @Produce  json
// @Router /clusterrole [GET]
// @Success 200 {object} clusterrole.ClusterRoleList
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetClusterRoleList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := clusterrole.GetClusterRoleList(k8s, dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

// handleGetClusterRoleDetail godoc
// @Tags Kubernetes
// @Summary Get detail of ClusterRole
// @Description Returns a detail of ClusterRole
// @Accept  json
// @Produce  json
// @Router /clusterrole/{clusterrole} [GET]
// @Param clusterrole path string true "Name of ClusterRole"
// @Success 200 {object} clusterrole.ClusterRoleDetail
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetClusterRoleDetail(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
	}

	crName := request.PathParameter("clusterrole")
	result, err := clusterrole.GetClusterRoleDetail(k8s, crName)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

// handleGetClusterRoleBindingList godoc
// @Tags Kubernetes
// @Summary Get list of ClusterRoleBinding
// @Description Returns a list of ClusterRoleBinding
// @Accept  json
// @Produce  json
// @Router /clusterrolebinding [GET]
// @Success 200 {object} clusterrolebinding.ClusterRoleBindingList
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetClusterRoleBindingList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := clusterrolebinding.GetClusterRoleBindingList(k8s, dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

// handleGetClusterRoleBindingDetail godoc
// @Tags Kubernetes
// @Summary Get detail of ClusterRoleBinding
// @Description Returns a detail of ClusterRoleBinding
// @Accept  json
// @Produce  json
// @Router /clusterrolebinding/{clusterrolebinding} [GET]
// @Param clusterrolebinding path string true "Name of ClusterRoleBinding"
// @Success 200 {object} clusterrolebinding.ClusterRoleBindingDetail
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetClusterRoleBindingDetail(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	crbName := request.PathParameter("clusterrolebinding")
	result, err := clusterrolebinding.GetClusterRoleBindingDetail(k8s, crbName)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

// handleGetAllConfigMapList godoc
// @Tags Kubernetes
// @Summary Get list of ConfigMap from cluster
// @Description Returns a list of ConfigMap from kubernetes cluster
// @Accept  json
// @Produce  json
// @Router /configmap [GET]
// @Success 200 {object} configmap.ConfigMapList
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetAllConfigMapList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := configmap.GetConfigMapList(k8s, common.NewNamespaceQuery(nil), dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

// handleGetConfigMapList godoc
// @Tags Kubernetes
// @Summary Get list of ConfigMap from Namespace
// @Description Returns a list of ConfigMap from Namespace
// @Accept  json
// @Produce  json
// @Router /configmap/{namespace} [GET]
// @Param namespace path string true "Namespace"
// @Success 200 {object} configmap.ConfigMapList
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetConfigMapList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := parseNamespacePathParameter(request)
	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := configmap.GetConfigMapList(k8s, namespace, dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

// handleGetConfiMapDetail godoc
// @Tags Kubernetes
// @Summary Get detail of ConfigMap
// @Description Returns a detail of ConfigMap
// @Accept  json
// @Produce  json
// @Router /configmap/{namespace}/{configmap} [GET]
// @Param namespace path string true "Namespace"
// @Param configmap path string true "Name of ConfigMap"
// @Success 200 {object} configmap.ConfigMapDetail
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetConfigMapDetail(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	configmapName := request.PathParameter("configmap")
	result, err := configmap.GetConfigMapDetail(k8s, namespace, configmapName)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

// handleGetAllPersistentVolumeClaimList godoc
// @Tags Kubernetes
// @Summary Get list of PersistentVolumeClaim from cluster
// @Description Returns a list of PersistentVolumeClaim from kubernetes cluster
// @Accept  json
// @Produce  json
// @Router /persistenvolumeclaim [GET]
// @Success 200 {object} persistentvolumeclaim.PersistentVolumeClaimList
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetAllPersistentVolumeClaimList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := persistentvolumeclaim.GetPersistentVolumeClaimList(k8s, common.NewNamespaceQuery(nil), dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

// handleGetPersistentVolumeClaimList godoc
// @Tags Kubernetes
// @Summary Get list of PersistentVolumeClaim from namespace
// @Description Returns a list of PersistentVolumeClaim from kubernetes namespace
// @Accept  json
// @Produce  json
// @Router /persistenvolumeclaim/{namespace} [GET]
// @Param namespace path string true "Namespace"
// @Success 200 {object} persistentvolumeclaim.PersistentVolumeClaimList
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetPersistentVolumeClaimList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := parseNamespacePathParameter(request)
	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := persistentvolumeclaim.GetPersistentVolumeClaimList(k8s, namespace, dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

// handleGetPersistentVolumeClaimDetail godoc
// @Tags Kubernetes
// @Summary Get detail of PersistentVolumeClaim from namespace
// @Description Returns a detail of PersistentVolumeClaim from kubernetes namespace
// @Accept  json
// @Produce  json
// @Router /persistenvolumeclaim/{namespace}/{persistentvolumeclaim} [GET]
// @Param namespace path string true "Namespace"
// @Param persistentvolumeclaim path string true "Name of PersistentVolumeClaim"
// @Success 200 {object} persistentvolumeclaim.PersistentVolumeClaimDetail
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetPersistentVolumeClaimDetail(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	pvcName := request.PathParameter("persistentvolumeclaim")
	result, err := persistentvolumeclaim.GetPersistentVolumeClaimDetail(k8s, namespace, pvcName)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)

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
