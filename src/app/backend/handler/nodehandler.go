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
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/event"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/node"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/pod"
)

func (apiHandler *APIHandler) installNode(ws *restful.WebService) {
	ws.Route(
		ws.GET("/node").
			To(apiHandler.handleGetNodeList).
			Returns(200, "OK", node.NodeList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List objects of kind Node").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.NodeDocsTag}))
	ws.Route(
		ws.GET("/node/{name}").
			To(apiHandler.handleGetNodeDetail).
			Param(ws.PathParameter("name", "Name of Node").Required(true)).
			Returns(200, "OK", node.NodeDetail{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("Read the specified Node").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.NodeDocsTag}))
	ws.Route(
		ws.GET("/node/{name}/event").
			To(apiHandler.handleGetNodeEvents).
			Param(ws.PathParameter("name", "Name of Node").Required(true)).
			Returns(200, "OK", common.EventList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List events related to a Node").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.NodeDocsTag}))
	ws.Route(
		ws.GET("/node/{name}/pod").
			To(apiHandler.handleGetNodePods).
			Param(ws.PathParameter("name", "Name of Node").Required(true)).
			Returns(200, "OK", pod.PodList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("list Pods related to a Node").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.NodeDocsTag}))
}

func (apiHandler *APIHandler) handleGetNodeList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	dataSelect := parser.ParseDataSelectPathParameter(request)
	dataSelect.MetricQuery = dataselect.StandardMetrics
	result, err := node.GetNodeList(k8s, dataSelect, apiHandler.iManager.Metric().Client())
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetNodeDetail(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	name := request.PathParameter("name")
	dataSelect := parser.ParseDataSelectPathParameter(request)
	dataSelect.MetricQuery = dataselect.StandardMetrics
	result, err := node.GetNodeDetail(k8s, apiHandler.iManager.Metric().Client(), name, dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetNodeEvents(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	nodeName := request.PathParameter("name")
	dataSelect := parser.ParseDataSelectPathParameter(request)
	dataSelect.MetricQuery = dataselect.StandardMetrics
	result, err := event.GetNodeEvents(k8s, dataSelect, nodeName)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetNodePods(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	nodeName := request.PathParameter("name")
	dataSelect := parser.ParseDataSelectPathParameter(request)
	dataSelect.MetricQuery = dataselect.StandardMetrics
	result, err := node.GetNodePods(k8s, apiHandler.iManager.Metric().Client(), dataSelect, nodeName)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}
