package handler

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"

	"github.com/donghoon-khan/kubeportal/src/app/backend/docs"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/handler/parser"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/persistentvolumeclaim"
)

func (apiHandler *APIHandler) installPersistentVolumeClaim(ws *restful.WebService) {
	ws.Route(
		ws.GET("/persistentvolumeclaim").
			To(apiHandler.handleGetPersistentVolumeClaimList).
			Returns(200, "OK", persistentvolumeclaim.PersistentVolumeClaimList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List objects of kind PersistentVolumeClaim").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.PersistentVolumeClaimDocsTag}))
	ws.Route(
		ws.GET("/persistentvolumeclaim/{namespace}").
			To(apiHandler.handleGetPersistentVolumeClaimListNamespace).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Returns(200, "OK", persistentvolumeclaim.PersistentVolumeClaimList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List objects of kind PersistentVolumeClaim in the Namespace").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.PersistentVolumeClaimDocsTag}))
	ws.Route(
		ws.GET("/persistentvolumeclaim/{namespace}/{name}").
			To(apiHandler.handleGetPersistentVolumeClaimDetail).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Param(ws.PathParameter("name", "Name of PersistentVolumeClaim").Required(true)).
			Returns(200, "OK", persistentvolumeclaim.PersistentVolumeClaimDetail{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("Read the specified PersistentVolumeClaim").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.PersistentVolumeClaimDocsTag}))
}

func (apiHandler *APIHandler) handleGetPersistentVolumeClaimList(request *restful.Request, response *restful.Response) {
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

func (apiHandler *APIHandler) handleGetPersistentVolumeClaimListNamespace(request *restful.Request, response *restful.Response) {
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

func (apiHandler *APIHandler) handleGetPersistentVolumeClaimDetail(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")
	result, err := persistentvolumeclaim.GetPersistentVolumeClaimDetail(k8s, namespace, name)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)

}
