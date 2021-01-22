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
		ws.GET("/login/skippable").
			To(self.handleLoginSkippable).
			Writes(authApi.LoginSkippableResponse{}))
}

func (self AuthHandler) handleLogin(request *restful.Request, resposne *restful.Response) {
	log.Println("Handle Login")
}

func (self *AuthHandler) handleLoginSkippable(request *restful.Request, response *restful.Response) {
	response.WriteHeaderAndEntity(http.StatusOK, authApi.LoginSkippableResponse{Skippable: self.manager.AuthenticationSkippable()})
}

func NewAuthHandler(manager authApi.AuthManager) AuthHandler {
	return AuthHandler{manager: manager}
}
