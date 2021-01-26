package handler

import (
	"net/http"

	"github.com/donghoon-khan/kubeportal/src/app/backend/auth"
	authApi "github.com/donghoon-khan/kubeportal/src/app/backend/auth/api"
	k8sApi "github.com/donghoon-khan/kubeportal/src/app/backend/kubernetes/api"
	"github.com/emicklei/go-restful"
)

const (
	RequestLogString  = "[%s] Incoming %s %s %s request from %s: %s"
	ResponseLogString = "[%s] Outcoming response to %s with %d status code"
)

type APIHandler struct {
	kManager k8sApi.KubernetesManager
}

func CreateHttpApiHandler(kManager k8sApi.KubernetesManager,
	authManager authApi.AuthManager) (http.Handler, error) {
	//apiHandler := APIHandler{kManager: kManager}
	wsContainer := restful.NewContainer()
	wsContainer.EnableContentEncoding(true)

	apiV1Ws := new(restful.WebService)

	apiV1Ws.Path("/api/v1").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	wsContainer.Add(apiV1Ws)

	authHandler := auth.NewAuthHandler(authManager)
	authHandler.Install(apiV1Ws)

	//apiV1Ws.Route(
	//apiV1Ws.GET("/namespace").
	//To(apiHandler.handleGetNamespaces).
	//Writes(namespace.NamespaceList{}))

	/*apiV1Ws.Route(
	apiV1Ws.POST("/appdeployment").
		To(apiHandler.handleDeploy).
		Reads(deployment.AppDepl	oymentSpec{}).
		Writes(deployment.AppDeploymentSpec{}))*/

	return wsContainer, nil
}

func (apiHandler *APIHandler) handleDeploy(request *restful.Request, response *restful.Response) {
	/*k8sClient, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}*/

	/*appDeploymentSpec := new(deployment.AppDeploymentSpec)
	if err := request.ReadEntity(appDeploymentSpec); err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	if err := deployment.DeployApp(appDeploymentSpec, k8sClient); err != nil {
		errors.HandleInternalError(response, err)
		return
	}*/

	//response.WriteHeaderAndEntity(http.StatusCreated, appDeploymentSpec)
}

func (apiHandler *APIHandler) handleGetNamespaces(request *restful.Request, response *restful.Response) {
	/*k8sClient, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	dataSelect := parser.ParseDataSelectPathParameter(request)*/
}
