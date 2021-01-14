package login

import "golang.org/x/oauth2"

type GoogleLogin struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Endpoint     oauth2.Endpoint
	Scopes       []string
}
