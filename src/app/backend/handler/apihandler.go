package handler

import (
	"net/http"

	k8sapi "github.com/donghoon-khan/kubeportal/src/app/backend/kubernetes/api"
	"github.com/emicklei/go-restful"
)

const (
	RequestLogString  = "[%s] Incoming %s %s %s request from %s: %s"
	ResponseLogString = "[%s] Outcoming response to %s with %d status code"
)

type APIHandler struct {
	kManager k8sapi.KubernetesManager
}

func CreateHttpApiHandler(kManager k8sapi.KubernetesManager) (http.Handler, error) {
	//apiHandler := APIHandler{kManager: kManager}
	wsContainer := restful.NewContainer()
	wsContainer.EnableContentEncoding(true)

	apiV1Ws := new(restful.WebService)

	apiV1Ws.Path("/api/v1").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	wsContainer.Add(apiV1Ws)

	/*apiV1Ws.Route(
	apiV1Ws.POST("/appdeployment").
		To(apiHandler.handleDeploy).
		Reads(deployment.AppDeploymentSpec{}).
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
