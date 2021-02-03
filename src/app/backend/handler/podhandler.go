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

var podDocsTag = []string{"Pod"}

func (apiHandler *APIHandler) installPod(ws *restful.WebService) {
	ws.Route(
		ws.GET("/pod").
			To(apiHandler.handleGetPodList).
			Writes(pod.PodList{}).
			Doc("List objects of kind Pod").
			Notes("Returns a list of Pod").
			Metadata(restfulspec.KeyOpenAPITags, podDocsTag).
			Returns(200, "OK", pod.PodList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}))
	ws.Route(
		ws.GET("/pod/{namespace}").
			To(apiHandler.handleGetPodListNamespace).
			Writes(pod.PodList{}).
			Doc("List objects of kind Pod in the Namespace").
			Notes("Returns a list of Pod in the Namespace").
			Metadata(restfulspec.KeyOpenAPITags, podDocsTag).
			Param(ws.PathParameter("namespace",
				"Object name and auth scope, such as for teams and projects").DataType("string").Required(true)).
			Returns(200, "OK", pod.PodList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}))
	ws.Route(
		ws.GET("/pod/{namespace}/{name}").
			To(apiHandler.handleGetPodDetail).
			Writes(pod.PodDetail{}).
			Doc("Read the specified Pod").
			Notes("Returns the specified Pod").
			Metadata(restfulspec.KeyOpenAPITags, podDocsTag).
			Param(ws.PathParameter("namespace",
				"Object name and auth scope, such as for teams and projects").DataType("string").Required(true)).
			Param(ws.PathParameter("name", "Name of Pod").DataType("string").Required(true)).
			Returns(200, "OK", pod.PodDetail{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}))
	ws.Route(
		ws.GET("/pod/{namespace}/{name}/container").
			To(apiHandler.handleGetPodContainerList).
			Writes(pod.PodDetail{}).
			Doc("List containers related to a Pod").
			Notes("Returns list of container related to a Pod").
			Metadata(restfulspec.KeyOpenAPITags, podDocsTag).
			Param(ws.PathParameter("namespace",
				"Object name and auth scope, such as for teams and projects").DataType("string").Required(true)).
			Param(ws.PathParameter("name", "Name of Pod").DataType("string").Required(true)).
			Returns(200, "OK", pod.PodDetail{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}))
	ws.Route(
		ws.GET("/pod/{namespace}/{name}/event").
			To(apiHandler.handleGetPodEventList).
			Writes(common.EventList{}).
			Doc("List events related to a Pod").
			Notes("Returns list of event related to a Pod").
			Metadata(restfulspec.KeyOpenAPITags, podDocsTag).
			Param(ws.PathParameter("namespace",
				"Object name and auth scope, such as for teams and projects").DataType("string").Required(true)).
			Param(ws.PathParameter("name", "Name of Pod").DataType("string").Required(true)).
			Returns(200, "OK", common.EventList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}))
	ws.Route(
		ws.GET("/pod/{namespace}/{name}/persistentvolumeclaim").
			To(apiHandler.handleGetPodPersistentVolumeClaimList).
			Writes(persistentvolumeclaim.PersistentVolumeClaimList{}).
			Doc("List PersistentVolumeClaims related to a Pod").
			Notes("Returns list of PersistentVolumeClaim related to a Pod").
			Metadata(restfulspec.KeyOpenAPITags, podDocsTag).
			Param(ws.PathParameter("namespace",
				"Object name and auth scope, such as for teams and projects").DataType("string").Required(true)).
			Param(ws.PathParameter("name", "Name of Pod").DataType("string").Required(true)).
			Returns(200, "OK", persistentvolumeclaim.PersistentVolumeClaimList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}))
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
