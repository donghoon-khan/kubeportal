package handler

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"

	"github.com/donghoon-khan/kubeportal/src/app/backend/docs"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/handler/parser"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/pod"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/service"
)

func (apiHandler *APIHandler) installService(ws *restful.WebService) {
	ws.Route(
		ws.GET("/service").
			To(apiHandler.handleGetServiceList).
			Returns(200, "OK", service.ServiceList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List objects of kind Service").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.ServiceDocsTag}))
	ws.Route(
		ws.GET("/service/{namespace}").
			To(apiHandler.handleGetServiceListNamespace).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Returns(200, "OK", service.ServiceList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List objects of kind Service in the Namespace").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.ServiceDocsTag}))
	ws.Route(
		ws.GET("/service/{namespace}/{name}").
			To(apiHandler.handleGetServiceDetail).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Param(ws.PathParameter("name", "Name of Service").Required(true)).
			Returns(200, "OK", service.ServiceDetail{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("Read the specified Service").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.ServiceDocsTag}))
	ws.Route(
		ws.GET("/service/{namespace}/{name}/event").
			To(apiHandler.handleGetServiceEvents).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Param(ws.PathParameter("name", "Name of Service").Required(true)).
			Returns(200, "OK", common.EventList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List events related to a Service").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.ServiceDocsTag}))
	ws.Route(
		ws.GET("/service/{namespace}/{name}/pod").
			To(apiHandler.handleGetServicePods).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Param(ws.PathParameter("name", "Name of Service").Required(true)).
			Returns(200, "OK", pod.PodList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List Pods related to a Service").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.ServiceDocsTag}))
}

func (apiHandler *APIHandler) handleGetServiceList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := service.GetServiceList(k8s, common.NewNamespaceQuery(nil), dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetServiceListNamespace(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := parseNamespacePathParameter(request)
	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := service.GetServiceList(k8s, namespace, dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetServiceDetail(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("service")
	result, err := service.GetServiceDetail(k8s, namespace, name)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetServiceEvents(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")
	dataSelect := parser.ParseDataSelectPathParameter(request)
	dataSelect.MetricQuery = dataselect.StandardMetrics
	result, err := service.GetServiceEvents(k8s, dataSelect, namespace, name)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetServicePods(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")
	dataSelect := parser.ParseDataSelectPathParameter(request)
	dataSelect.MetricQuery = dataselect.StandardMetrics
	result, err := service.GetServicePods(k8s, apiHandler.iManager.Metric().Client(), namespace, name, dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}
