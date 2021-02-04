package handler

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"

	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/handler/parser"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/container"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/persistentvolumeclaim"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/pod"
)

func (apiHandler *APIHandler) installPod(ws *restful.WebService) {
	ws.Route(
		ws.GET("/pod").
			To(apiHandler.handleGetPodList).
			Returns(200, "OK", pod.PodList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List objects of kind Pod").
			Metadata(restfulspec.KeyOpenAPITags, []string{podDocsTag}))
	ws.Route(
		ws.GET("/pod/{namespace}").
			To(apiHandler.handleGetPodListNamespace).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Returns(200, "OK", pod.PodList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List objects of kind Pod in the Namespace").
			Metadata(restfulspec.KeyOpenAPITags, []string{podDocsTag}))
	ws.Route(
		ws.GET("/pod/{namespace}/{name}").
			To(apiHandler.handleGetPodDetail).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Param(ws.PathParameter("name", "Name of Pod").Required(true)).
			Returns(200, "OK", pod.PodDetail{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("Read the specified Pod").
			Metadata(restfulspec.KeyOpenAPITags, []string{podDocsTag}))
	ws.Route(
		ws.GET("/pod/{namespace}/{name}/container").
			To(apiHandler.handleGetPodContainerList).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Param(ws.PathParameter("name", "Name of Pod").Required(true)).
			Returns(200, "OK", pod.PodDetail{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List containers related to a Pod").
			Metadata(restfulspec.KeyOpenAPITags, []string{podDocsTag}))
	ws.Route(
		ws.GET("/pod/{namespace}/{name}/event").
			To(apiHandler.handleGetPodEventList).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Param(ws.PathParameter("name", "Name of Pod").Required(true)).
			Returns(200, "OK", common.EventList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List events related to a Pod").
			Metadata(restfulspec.KeyOpenAPITags, []string{podDocsTag}))
	ws.Route(
		ws.GET("/pod/{namespace}/{name}/persistentvolumeclaim").
			To(apiHandler.handleGetPodPersistentVolumeClaimList).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Param(ws.PathParameter("name", "Name of Pod").Required(true)).
			Returns(200, "OK", persistentvolumeclaim.PersistentVolumeClaimList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List PersistentVolumeClaims related to a Pod").
			Metadata(restfulspec.KeyOpenAPITags, []string{podDocsTag}))
}

func (apiHandler *APIHandler) handleGetPodList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	dataSelect := parser.ParseDataSelectPathParameter(request)
	dataSelect.MetricQuery = dataselect.StandardMetrics
	result, err := pod.GetPodList(k8s, apiHandler.iManager.Metric().Client(), common.NewNamespaceQuery(nil), dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetPodListNamespace(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := parseNamespacePathParameter(request)
	dataSelect := parser.ParseDataSelectPathParameter(request)
	dataSelect.MetricQuery = dataselect.StandardMetrics
	result, err := pod.GetPodList(k8s, apiHandler.iManager.Metric().Client(), namespace, dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetPodDetail(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	podName := request.PathParameter("name")
	result, err := pod.GetPodDetail(k8s, apiHandler.iManager.Metric().Client(), namespace, podName)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetPodContainerList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	podName := request.PathParameter("name")
	result, err := container.GetPodContainers(k8s, namespace, podName)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetPodEventList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	podName := request.PathParameter("name")
	dataSelect := parser.ParseDataSelectPathParameter(request)
	dataSelect.MetricQuery = dataselect.StandardMetrics
	result, err := pod.GetEventsForPod(k8s, dataSelect, namespace, podName)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetPodPersistentVolumeClaimList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	podName := request.PathParameter("name")
	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := persistentvolumeclaim.GetPodPersistentVolumeClaims(k8s, namespace, podName, dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}
