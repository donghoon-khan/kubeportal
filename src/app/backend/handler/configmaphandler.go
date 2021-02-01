package handler

import (
	"net/http"

	"github.com/emicklei/go-restful"

	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/handler/parser"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/configmap"
)

func (apiHandler *APIHandler) installConfigMap(ws *restful.WebService) {
	ws.Route(
		ws.GET("/configmap").
			To(apiHandler.handleGetConfigMapList).
			Writes(configmap.ConfigMapList{}))
	ws.Route(
		ws.GET("/configmap/{namespace}").
			To(apiHandler.handleGetConfigMapList).
			Writes(configmap.ConfigMapList{}))
	ws.Route(
		ws.GET("/configmap/{namespace}/{configmap}").
			To(apiHandler.handleGetConfigMapDetail).
			Writes(configmap.ConfigMapDetail{}))
}

// handleGetConfigMapList godoc
// @Tags ConfigMap
// @Summary Get list of ConfigMap
// @Description Returns a list of ConfigMap from Kubernetes cluster or Namespace
// @Accept  json
// @Produce  json
// @Router /configmap/{namespace} [GET]
// @Param namespace path string false "Namespace"
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
// @Tags ConfigMap
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
