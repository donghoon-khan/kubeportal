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

func (self AuthHandler) Install(ws *restful.WebService) {
	ws.Route(
		ws.POST("/login").
			To(self.handleLogin).
			Reads(authApi.LoginSpec{}).
			Writes(authApi.AuthResponse{}))
	ws.Route(
		ws.GET("/login/skippable").
			To(self.handleLoginSkippable).
			Writes(authApi.LoginSkippableResponse{}))
}

// handleLogin godoc
// @Tags Authentication
// @Summary Return JWEToken
// @Accept  json
// @Produce  json
// @Router /login [POST]
// @Param LoginSpec body authApi.LoginSpec true "Information required to authenticate user"
// @Success 200 {object} authApi.AuthResponse
func (self AuthHandler) handleLogin(request *restful.Request, resposne *restful.Response) {
	log.Println("Handle Login")
}

// handleLoginSkippable godoc
// @Tags Authentication
// @Summary Return the authentication skip should be enabled or not
// @Accept  json
// @Produce  json
// @Router /login/skippable [GET]
// @Success 200 {object} authApi.LoginSkippableResponse
func (self *AuthHandler) handleLoginSkippable(request *restful.Request, response *restful.Response) {
	response.WriteHeaderAndEntity(http.StatusOK, authApi.LoginSkippableResponse{Skippable: self.manager.AuthenticationSkippable()})
}

func NewAuthHandler(manager authApi.AuthManager) AuthHandler {
	return AuthHandler{manager: manager}
}
