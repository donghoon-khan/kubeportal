package handler

import (
	"net/http"

	"github.com/emicklei/go-restful"

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
			Writes(pod.PodList{}))
	ws.Route(
		ws.GET("/pod/{namespace}").
			To(apiHandler.handleGetPodList).
			Writes(pod.PodList{}))
	ws.Route(
		ws.GET("/pod/{namespace}/{pod}").
			To(apiHandler.handleGetPodDetail).
			Writes(pod.PodDetail{}))
	ws.Route(
		ws.GET("/pod/{namespace}/{pod}/container").
			To(apiHandler.handleGetPodContainerList).
			Writes(pod.PodDetail{}))
	ws.Route(
		ws.GET("/pod/{namespace}/{pod}/event").
			To(apiHandler.handleGetPodEventList).
			Writes(common.EventList{}))
	ws.Route(
		ws.GET("/pod/{namespace}/{pod}/persistentvolumeclaim").
			To(apiHandler.handleGetPodPersistentVolumeClaimList).
			Writes(persistentvolumeclaim.PersistentVolumeClaimList{}))
}

// handleGetPodList godoc
// @Tags Pod
// @Summary Get list of pod
// @Description Returns a list of pod from Kubernetes cluster or namespace
// @Accept  json
// @Produce  json
// @Router /pod/{namespace} [GET]
// @Param namespace path string false "Name of namespace"
// @Success 200 {object} pod.PodList
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetPodList(request *restful.Request, response *restful.Response) {
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

// handleGetPodDetail godoc
// @Tags Pod
// @Summary Get detail of pod
// @Description Returns a detail of pod
// @Accept  json
// @Produce  json
// @Router /pod/{namespace}/{pod} [GET]
// @Param namespace path string true "Name of namespace"
// @Param pod path string true "Name of pod"
// @Success 200 {object} pod.PodDetail
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetPodDetail(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	podName := request.PathParameter("pod")
	result, err := pod.GetPodDetail(k8s, apiHandler.iManager.Metric().Client(), namespace, podName)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

// handleGetPodContainerList godoc
// @Tags Pod
// @Summary Get containers related to a pod
// @Description Returns a list of container related to a pod in namespace
// @Accept  json
// @Produce  json
// @Router /pod/{namespace}/{pod}/container [GET]
// @Param namespace path string true "Name of namespace"
// @Param pod path string true "Name of pod"
// @Success 200 {object} pod.PodDetail
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetPodContainerList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	podName := request.PathParameter("pod")
	result, err := container.GetPodContainers(k8s, namespace, podName)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

// handleGetPodEventList godoc
// @Tags Pod
// @Summary Get events related to a pod
// @Description Returns list of event related to a pod in namespace
// @Accept  json
// @Produce  json
// @Router /pod/{namespace}/{pod}/event [GET]
// @Param namespace path string true "Name of namespace"
// @Param pod path string true "Name of pod"
// @Success 200 {object} common.EventList
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetPodEventList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	podName := request.PathParameter("pod")
	dataSelect := parser.ParseDataSelectPathParameter(request)
	dataSelect.MetricQuery = dataselect.StandardMetrics
	result, err := pod.GetEventsForPod(k8s, dataSelect, namespace, podName)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

// handleGetPodPersistentVolumeClaimList godoc
// @Tags Pod
// @Summary Get list of persistentvolumeclaim related to a pod
// @Description Returns list of PersistentVolumeClaim related to a pod in namespace
// @Accept  json
// @Produce  json
// @Router /pod/{namespace}/{pod}/persistentvolumeclaim [GET]
// @Param namespace path string true "Name of namespace"
// @Param pod path string true "Name of pod"
// @Success 200 {object} persistentvolumeclaim.PersistentVolumeClaimList
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetPodPersistentVolumeClaimList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	podName := request.PathParameter("pod")
	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := persistentvolumeclaim.GetPodPersistentVolumeClaims(k8s, namespace, podName, dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}
