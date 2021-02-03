package handler

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"

	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/handler/parser"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/clusterrole"
)

var clusterRoleDocsTag = []string{"ClusterRole"}

func (apiHandler *APIHandler) installClusterRole(ws *restful.WebService) {
	ws.Route(
		ws.GET("/clusterrole").
			To(apiHandler.handleGetClusterRoleList).
			Writes(clusterrole.ClusterRoleList{}).
			Doc("List objects of kind ClusterRole").
			Notes("Returns a list of ClusterRole").
			Metadata(restfulspec.KeyOpenAPITags, clusterRoleDocsTag).
			Returns(200, "OK", clusterrole.ClusterRoleList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}))
	ws.Route(
		ws.GET("/clusterrole/{clusterrole}").
			To(apiHandler.handleGetClusterRoleDetail).
			Writes(clusterrole.ClusterRoleDetail{}).
			Doc("Read the specified ClusterRole").
			Notes("Returns the specified ClusterRole").
			Metadata(restfulspec.KeyOpenAPITags, clusterRoleDocsTag).
			Param(ws.PathParameter("name", "Name of ClusterRole").DataType("string").Required(true)).
			Returns(200, "OK", clusterrole.ClusterRoleDetail{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}))
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

	crName := request.PathParameter("clusterrole")
	result, err := clusterrole.GetClusterRoleDetail(k8s, crName)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}
