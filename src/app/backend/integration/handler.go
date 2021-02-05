package integration

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"

	"github.com/donghoon-khan/kubeportal/src/app/backend/docs"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/integration/api"
)

type IntegrationHandler struct {
	iManager IntegrationManager
}

func (iHandler IntegrationHandler) Install(ws *restful.WebService) {
	ws.Route(
		ws.GET("/integration/{name}/state").
			To(iHandler.handleGetState).
			Writes(api.IntegrationState{}).
			Doc("Get state of integration").
			Metadata(restfulspec.KeyOpenAPITags, docs.IntegrationDocsTag).
			Param(ws.PathParameter("name", "Name of integration").Required(true)).
			Returns(200, "OK", api.IntegrationState{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}))
}

func (iHandler IntegrationHandler) handleGetState(request *restful.Request, response *restful.Response) {
	iName := request.PathParameter("name")
	state, err := iHandler.iManager.GetState(api.IntegrationID(iName))
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error()+"\n")
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, state)
}

func NewIntegrationHandler(iManager IntegrationManager) IntegrationHandler {
	return IntegrationHandler{iManager: iManager}
}
