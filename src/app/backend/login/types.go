package login

import (
	"net/http"

	"golang.org/x/oauth2"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Endpoint     oauth2.Endpoint
	Scopes       []string
}

type LoginManager interface {
	New(config *Config) *http.ServeMux
	issueSession() http.Handler
	profileHandler(w http.ResponseWriter, r *http.Request)
	logoutHandler(w http.ResponseWriter, r *http.Request)
}
