package handler

import (
	"net/http"

	"github.com/emicklei/go-restful"

	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/handler/parser"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/persistentvolumeclaim"
)

func (apiHandler *APIHandler) installPersistentVolumeClaim(ws *restful.WebService) {
	ws.Route(
		ws.GET("/persistentvolumeclaim").
			To(apiHandler.handleGetPersistentVolumeClaimList).
			Writes(persistentvolumeclaim.PersistentVolumeClaimList{}))
	ws.Route(
		ws.GET("/persistentvolumeclaim/{namespace}").
			To(apiHandler.handleGetPersistentVolumeClaimList).
			Writes(persistentvolumeclaim.PersistentVolumeClaimList{}))
	ws.Route(
		ws.GET("/persistentvolumeclaim/{namespace}/{persistentvolumeclaim}").
			To(apiHandler.handleGetPersistentVolumeClaimDetail).
			Writes(persistentvolumeclaim.PersistentVolumeClaimDetail{}))
}

// handleGetPersistentVolumeClaimList godoc
// @Tags PersistentVolumeClaim
// @Summary Get list of persistentvolumeclaim
// @Description Returns a list of persistentvolumeclaim from Kubernetes cluster or namespace
// @Accept  json
// @Produce  json
// @Router /persistenvolumeclaim/{namespace} [GET]
// @Param namespace path string false "Name of namespace"
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
// @Tags PersistentVolumeClaim
// @Summary Get detail of persistentvolumeclaim
// @Description Returns a detail of persistentvolumeclaim
// @Accept  json
// @Produce  json
// @Router /persistenvolumeclaim/{namespace}/{persistentvolumeclaim} [GET]
// @Param namespace path string true "Name of namespace"
// @Param persistentvolumeclaim path string true "Name of persistentvolumeclaim"
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
