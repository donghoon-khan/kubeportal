package handler

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"

	"github.com/donghoon-khan/kubeportal/src/app/backend/docs"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/handler/parser"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/secret"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/serviceaccount"
)

func (apiHandler *APIHandler) installServiceAccount(ws *restful.WebService) {
	ws.Route(
		ws.GET("/serviceaccount").
			To(apiHandler.handleGetServiceAccountList).
			Returns(200, "OK", serviceaccount.ServiceAccountList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List objects of kind ServiceAccount").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.ServiceAccountDocsTag}))
	ws.Route(
		ws.GET("/serviceaccount/{namespace}").
			To(apiHandler.handleGetServiceAccountListNamespace).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Returns(200, "OK", serviceaccount.ServiceAccountList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List objects of kind ServiceAccount in the Namespace").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.ServiceAccountDocsTag}))
	ws.Route(
		ws.GET("/serviceaccount/{namespace}/{name}").
			To(apiHandler.handleGetServiceAccountDetail).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Param(ws.PathParameter("name", "Name of ServiceAccount").Required(true)).
			Returns(200, "OK", serviceaccount.ServiceAccountDetail{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("Read the specified ServiceAccount").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.ServiceAccountDocsTag}))
	ws.Route(
		ws.GET("/serviceaccount/{namespace}/{name}/secret").
			To(apiHandler.handleGetServiceAccountSecrets).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Param(ws.PathParameter("name", "Name of ServiceAccount").Required(true)).
			Returns(200, "OK", secret.SecretList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List Secrets related to a ServiceAccount").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.ServiceAccountDocsTag}))
	ws.Route(
		ws.GET("/serviceaccount/{namespace}/{name}/imagepullsecret").
			To(apiHandler.handleGetServiceAccountImagePullSecrets).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Param(ws.PathParameter("name", "Name of ServiceAccount").Required(true)).
			Returns(200, "OK", secret.SecretList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List imagePullSecrets related to a ServiceAccount").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.ServiceAccountDocsTag}))
}

func (apiHandler *APIHandler) handleGetServiceAccountList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := serviceaccount.GetServiceAccountList(k8s, common.NewNamespaceQuery(nil), dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetServiceAccountListNamespace(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := parseNamespacePathParameter(request)
	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := serviceaccount.GetServiceAccountList(k8s, namespace, dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetServiceAccountDetail(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")
	result, err := serviceaccount.GetServiceAccountDetail(k8s, namespace, name)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetServiceAccountSecrets(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")
	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := serviceaccount.GetServiceAccountSecrets(k8s, namespace, name, dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetServiceAccountImagePullSecrets(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")
	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := serviceaccount.GetServiceAccountImagePullSecrets(k8s, namespace, name, dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}
