package handler

import (
	"errors"
	"net/http"

	"github.com/emicklei/go-restful"
	clientapi "github.com/kubernetes/dashboard/src/app/backend/client/api"
)

const (
	RequestLogString  = "[%s] Incoming %s %s %s request from %s: %s"
	ResponseLogString = "[%s] Outcoming response to %s with %d status code"
)

type APIHandler struct {
	cManager clientapi.ClientManager
}

func CreateHTTPAPIHandler(cManager clientapi.ClientManager) http.Handler {
	apiHandler := APIHandler{cManager: cManager}
	wsContainer := restful.NewContainer()
	wsContainer.EnableContentEncoding(true)

	apiV1Ws := new(restful.WebService)

	apiV1Ws.Path("/api/v1").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	wsContainer.Add(apiV1Ws)

	return wsContainer, nil
}

func (apiHandler *APIHandler) handleDeploy(request *restful.Request, response restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		errors.HandleInternalError(response, err)
	}
}
