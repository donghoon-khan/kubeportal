package handler

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"

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
			Writes(node.NodeList{}))
	ws.Route(
		ws.GET("/node/{node}").
			To(apiHandler.handleGetNodeDetail).
			Writes(node.NodeDetail{}))
	ws.Route(
		ws.GET("/node/{node}/event").
			To(apiHandler.handleGetNodeEventList).
			Writes(common.EventList{}))
	ws.Route(
		ws.GET("/node/{node}/pod").
			To(apiHandler.handleGetNodePods).
			Writes(pod.PodList{}))
}

// handleGetNodeList godoc
// @Tags Node
// @Summary Get list of node
// @Description Returns a list of node
// @Accept  json
// @Produce  json
// @Router /node [GET]
// @Success 200 {object} node.NodeList
// @Failure 401 {string} string "Unauthorized"
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

// handleGetNodeDetail godoc
// @Tags Node
// @Summary Get detail of node
// @Description Returns a detail of node
// @Accept  json
// @Produce  json
// @Router /node/{node} [GET]
// @Param node path string true "Name of node"
// @Success 200 {object} node.NodeDetail
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetNodeDetail(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	nodeName := request.PathParameter("node")
	dataSelect := parser.ParseDataSelectPathParameter(request)
	dataSelect.MetricQuery = dataselect.StandardMetrics
	result, err := node.GetNodeDetail(k8s, apiHandler.iManager.Metric().Client(), nodeName, dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

// handleGetNodeEventList godoc
// @Tags Node
// @Summary Get events related to a node
// @Description Returns list of event related to node
// @Accept  json
// @Produce  json
// @Router /node/{node}/event [GET]
// @Param node path string true "Name of node"
// @Success 200 {object} common.EventList
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetNodeEventList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	nodeName := request.PathParameter("node")
	dataSelect := parser.ParseDataSelectPathParameter(request)
	dataSelect.MetricQuery = dataselect.StandardMetrics
	result, err := event.GetNodeEvents(k8s, dataSelect, nodeName)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

// handleGetNodePods godoc
// @Tags Node
// @Summary Get pods related to a node
// @Description Returns list of pod related to node
// @Accept  json
// @Produce  json
// @Router /node/{node}/pod [GET]
// @Param node path string true "Name of node"
// @Success 200 {object} pod.PodList
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetNodePods(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	nodeName := request.PathParameter("node")
	dataSelect := parser.ParseDataSelectPathParameter(request)
	dataSelect.MetricQuery = dataselect.StandardMetrics
	result, err := node.GetNodePods(k8s, apiHandler.iManager.Metric().Client(), dataSelect, nodeName)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}
