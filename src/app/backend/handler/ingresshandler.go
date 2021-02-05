package handler

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"

	"github.com/donghoon-khan/kubeportal/src/app/backend/docs"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/handler/parser"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/ingress"
)

func (apiHandler *APIHandler) installIngress(ws *restful.WebService) {
	ws.Route(
		ws.GET("/ingress").
			To(apiHandler.handleGetIngressList).
			Returns(200, "OK", ingress.IngressList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List objects of kind Ingress").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.IngressDocsTag}))
	ws.Route(
		ws.GET("/ingress/{namespace}").
			To(apiHandler.handleGetIngressListNamespace).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Returns(200, "OK", ingress.IngressList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List objects of kind Ingress in the Namespace").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.IngressDocsTag}))
	ws.Route(
		ws.GET("/ingress/{namespace}/{name}").
			To(apiHandler.handleGetIngressDetail).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Param(ws.PathParameter("name", "Name of Ingress").Required(true)).
			Returns(200, "OK", ingress.IngressDetail{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("Read the specified Ingress").
			Metadata(restfulspec.KeyOpenAPITags, []string{docs.IngressDocsTag}))
}

func (apiHandler *APIHandler) handleGetIngressList(req *restful.Request, res *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(req)
	if err != nil {
		errors.HandleInternalError(res, err)
		return
	}

	dataSelect := parser.ParseDataSelectPathParameter(req)
	result, err := ingress.GetIngressList(k8s, common.NewNamespaceQuery(nil), dataSelect)
	if err != nil {
		errors.HandleInternalError(res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetIngressListNamespace(req *restful.Request, res *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(req)
	if err != nil {
		errors.HandleInternalError(res, err)
		return
	}

	dataSelect := parser.ParseDataSelectPathParameter(req)
	namespace := parseNamespacePathParameter(req)
	result, err := ingress.GetIngressList(k8s, namespace, dataSelect)
	if err != nil {
		errors.HandleInternalError(res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetIngressDetail(req *restful.Request, res *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(req)
	if err != nil {
		errors.HandleInternalError(res, err)
		return
	}

	namespace := req.PathParameter("namespace")
	name := req.PathParameter("name")
	result, err := ingress.GetIngressDetail(k8s, namespace, name)
	if err != nil {
		errors.HandleInternalError(res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, result)
}
