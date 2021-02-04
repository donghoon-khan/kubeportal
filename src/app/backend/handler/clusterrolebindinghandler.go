package handler

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"

	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/handler/parser"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/clusterrolebinding"
)

func (apiHandler *APIHandler) installClusterRoleBinding(ws *restful.WebService) {
	ws.Route(
		ws.GET("/clusterrolebinding").
			To(apiHandler.handleGetClusterRoleBindingList).
			Returns(200, "OK", clusterrolebinding.ClusterRoleBindingList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List objects of kind ClusterRoleBinding").
			Metadata(restfulspec.KeyOpenAPITags, []string{clusterRoleBindingDocsTag}))
	ws.Route(
		ws.GET("/clusterrolebinding/{name}").
			To(apiHandler.handleGetClusterRoleBindingDetail).
			Param(ws.PathParameter("name", "Name of ClusterRoleBinding").DataType("string").Required(true)).
			Returns(200, "OK", clusterrolebinding.ClusterRoleBindingDetail{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("Read the specified ClusterRoleBinding").
			Metadata(restfulspec.KeyOpenAPITags, []string{clusterRoleBindingDocsTag}))
}

func (apiHandler *APIHandler) handleGetClusterRoleBindingList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := clusterrolebinding.GetClusterRoleBindingList(k8s, dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetClusterRoleBindingDetail(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	name := request.PathParameter("name")
	result, err := clusterrolebinding.GetClusterRoleBindingDetail(k8s, name)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}
