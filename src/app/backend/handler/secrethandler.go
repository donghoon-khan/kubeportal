package handler

import (
	"net/http"

	"github.com/emicklei/go-restful"

	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/handler/parser"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/secret"
)

func (apiHandler *APIHandler) installSecret(ws *restful.WebService) {
	ws.Route(
		ws.GET("/secret").
			To(apiHandler.handleGetSecretList).
			Writes(secret.SecretList{}))
	ws.Route(
		ws.GET("/secret/{namespace}").
			To(apiHandler.handleGetSecretList).
			Writes(secret.SecretList{}))
	ws.Route(
		ws.GET("/secret/{namespace}/{secret}").
			To(apiHandler.handleGetSecretDetail).
			Writes(secret.SecretDetail{}))
}

// handleGetSecretList godoc
// @Tags Secret
// @Summary Get list of Secret
// @Description Returns a list of Secret from Kubernetes cluster or Namespace
// @Accept  json
// @Produce  json
// @Router /secret/{namespace} [GET]
// @Param namespace path string false "Namespace"
// @Success 200 {object} secret.SecretList
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetSecretList(request *restful.Request, response *restful.Response) {
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

// handleGetSecretDetail godoc
// @Tags Secret
// @Summary Get detail of Secret
// @Description Returns a detail of Secret
// @Accept  json
// @Produce  json
// @Router /secret/{namespace}/{secret} [GET]
// @Param namespace path string true "Namespace"
// @Param secret path string true "Name of Secret"
// @Success 200 {object} secret.SecretDetail
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetSecretDetail(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	secretName := request.PathParameter("secret")
	result, err := secret.GetSecretDetail(k8s, namespace, secretName)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}
