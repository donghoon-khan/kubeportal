package handler

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"

	"github.com/donghoon-khan/kubeportal/src/app/backend/docs"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/handler/parser"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/clusterrole"
)

func (apiHandler *APIHandler) installClusterRole(ws *restful.WebService) {
	ws.Route(
		ws.GET("/clusterrole").
			To(apiHandler.handleGetClusterRoleList).
			Returns(200, "OK", clusterrole.ClusterRoleList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List objects of kind ClusterRole").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.ClusterRoleDocsTag}))
	ws.Route(
		ws.GET("/clusterrole/{name}").
			To(apiHandler.handleGetClusterRoleDetail).
			Param(ws.PathParameter("name", "Name of ClusterRole").Required(true)).
			Returns(200, "OK", clusterrole.ClusterRoleDetail{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("Read the specified ClusterRole").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.ClusterRoleDocsTag}))
}

func (apiHandler *APIHandler) handleGetClusterRoleList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := clusterrole.GetClusterRoleList(k8s, dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetClusterRoleDetail(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
	}

	name := request.PathParameter("name")
	result, err := clusterrole.GetClusterRoleDetail(k8s, name)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}
