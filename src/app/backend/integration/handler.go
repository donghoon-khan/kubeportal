package integration

import (
	"net/http"

	"github.com/donghoon-khan/kubeportal/src/app/backend/integration/api"
	restful "github.com/emicklei/go-restful/v3"
)

type IntegrationHandler struct {
	iManager IntegrationManager
}

func (iHandler IntegrationHandler) Install(ws *restful.WebService) {
	ws.Route(
		ws.GET("/integration/{name}/state").
			To(iHandler.handleGetState).
			Writes(api.IntegrationState{}))
}

// handleGetState godoc
// @Tags integration
// @Summary Return integration state
// @Accept  json
// @Produce  json
// @Router /integration/{name}/state [GET]
// @Success 200 {object} api.IntegrationState
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
