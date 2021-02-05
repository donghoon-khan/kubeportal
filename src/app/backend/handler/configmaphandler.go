package handler

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"

	"github.com/donghoon-khan/kubeportal/src/app/backend/docs"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/handler/parser"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/configmap"
)

func (apiHandler *APIHandler) installConfigMap(ws *restful.WebService) {
	ws.Route(
		ws.GET("/configmap").
			To(apiHandler.handleGetConfigMapList).
			Returns(200, "OK", configmap.ConfigMapList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List objects of kind ConfigMap").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.ConfigMapDocsTag}))
	ws.Route(
		ws.GET("/configmap/{namespace}").
			To(apiHandler.handleGetConfigMapListNamespace).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Returns(200, "OK", configmap.ConfigMapList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List objects of kind ConfigMap in the Namespace").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.ConfigMapDocsTag}))
	ws.Route(
		ws.GET("/configmap/{namespace}/{name}").
			To(apiHandler.handleGetConfigMapDetail).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Param(ws.PathParameter("name", "Name of ConfigMap").Required(true)).
			Returns(200, "OK", configmap.ConfigMapDetail{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("Read the specified ConfigMap").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.ConfigMapDocsTag}))
}

func (apiHandler *APIHandler) handleGetConfigMapList(request *restful.Request, response *restful.Response) {
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

func (apiHandler *APIHandler) handleGetConfigMapListNamespace(request *restful.Request, response *restful.Response) {

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

func (apiHandler *APIHandler) handleGetConfigMapDetail(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")
	result, err := configmap.GetConfigMapDetail(k8s, namespace, name)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}
