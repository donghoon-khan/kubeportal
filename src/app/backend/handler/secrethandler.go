package handler

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"

	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/handler/parser"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/secret"
)

func (apiHandler *APIHandler) installSecret(ws *restful.WebService) {
	ws.Route(
		ws.GET("/secret").
			To(apiHandler.handleGetSecretList).
			Returns(200, "OK", secret.SecretList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List objects of kind Secret").
			Metadata(restfulspec.KeyOpenAPITags, []string{secretDocsTag}))
	ws.Route(
		ws.GET("/secret/{namespace}").
			To(apiHandler.handleGetSecretListNamespace).
			Param(ws.PathParameter("namespace",
				"Object name and auth scope, such as for teams and projects").DataType("string").Required(true)).
			Returns(200, "OK", secret.SecretList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List objects of kind Secret in the Namespace").
			Metadata(restfulspec.KeyOpenAPITags, []string{secretDocsTag}))
	ws.Route(
		ws.GET("/secret/{namespace}/{name}").
			To(apiHandler.handleGetSecretDetail).
			Param(ws.PathParameter("namespace",
				"Object name and auth scope, such as for teams and projects").DataType("string").Required(true)).
			Param(ws.PathParameter("name", "Name of Secret").DataType("string").Required(true)).
			Returns(200, "OK", secret.SecretDetail{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("Read the specified Secret").
			Metadata(restfulspec.KeyOpenAPITags, []string{secretDocsTag}))
}

func (apiHandler *APIHandler) handleGetSecretList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := secret.GetSecretList(k8s, common.NewNamespaceQuery(nil), dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetSecretListNamespace(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	dataSelect := parser.ParseDataSelectPathParameter(request)
	namespace := parseNamespacePathParameter(request)
	result, err := secret.GetSecretList(k8s, namespace, dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetSecretDetail(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")
	result, err := secret.GetSecretDetail(k8s, namespace, name)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}
