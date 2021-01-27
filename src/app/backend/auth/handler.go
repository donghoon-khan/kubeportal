package auth

import (
	"log"
	"net/http"

	authApi "github.com/donghoon-khan/kubeportal/src/app/backend/auth/api"
	"github.com/emicklei/go-restful"
)

type AuthHandler struct {
	manager authApi.AuthManager
}

func (authHandler AuthHandler) Install(ws *restful.WebService) {
	ws.Route(
		ws.POST("/login").
			To(authHandler.handleLogin).
			Reads(authApi.LoginSpec{}).
			Writes(authApi.AuthResponse{}))
	ws.Route(
		ws.GET("/login/skippable").
			To(authHandler.handleLoginSkippable).
			Writes(authApi.LoginSkippableResponse{}))
}

// handleLogin godoc
// @Tags Authentication
// @Summary Login
// @Description Returns a JWEToken by login spec
// @Accept  json
// @Produce  json
// @Router /login [POST]
// @Param LoginSpec body authApi.LoginSpec true "Information required to authenticate user"
// @Success 200 {object} authApi.AuthResponse
func (authHandler AuthHandler) handleLogin(request *restful.Request, resposne *restful.Response) {
	log.Println("Handle Login")
}

// handleLoginSkippable godoc
// @Tags Authentication
// @Summary Login skippable
// @Description Returns a authentication skip should be enabled or not
// @Accept  json
// @Produce  json
// @Router /login/skippable [GET]
// @Success 200 {object} authApi.LoginSkippableResponse
func (authHandler *AuthHandler) handleLoginSkippable(request *restful.Request, response *restful.Response) {
	response.WriteHeaderAndEntity(http.StatusOK,
		authApi.LoginSkippableResponse{Skippable: authHandler.manager.AuthenticationSkippable()})
}

func NewAuthHandler(manager authApi.AuthManager) AuthHandler {
	return AuthHandler{manager: manager}
}
