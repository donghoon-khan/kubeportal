package auth

import (
	authApi "github.com/donghoon-khan/kubeportal/src/app/backend/auth/api"
	k8sApi "github.com/donghoon-khan/kubeportal/src/app/backend/kubernetes/api"
)

type authManager struct {
	k8sManager              k8sApi.KubernetesManager
	authenticationModes     authApi.AuthenticationModes
	authenticationSkippable bool
}

/*func (self authManager) Login(spec *authApi.LoginSpec) (*authApi.AuthResponse, error) {
	authenticator, err := self.getAuthenticator(spec)
	if err != nil {
		return nil, err
	}

	authInfo, err := authenticator.GetAuthInfo()
	if err != nil {
		return nil, err
	}

	err = self.healthCheck(authInfo)
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil || len(nonCriticalErrors) > 0 {
		return &authApi.AuthResponse{Errors: nonCriticalErrors}, criticalError
	}

	token, err := self.tokenManager.Generate(authInfo)
	if err != nil {
		return nil, err
	}

	return &authApi.AuthResponse{JWEToken: token, Errors: nonCriticalErrors}, nil
}*/
//TODO
func (self authManager) Login(spec *authApi.LoginSpec) (*authApi.AuthResponse, error) {
	return &authApi.AuthResponse{JWEToken: "JWEToken"}, nil
}

//TODO
func (self authManager) Refresh(jweToken string) (string, error) {
	return "RefreshToken", nil
}

func (self authManager) AuthenticationModes() []authApi.AuthenticationMode {
	return self.authenticationModes.Array()
}

func (self authManager) AuthenticationSkippable() bool {
	return self.authenticationSkippable
}

func NewAuthManager(k8sManager k8sApi.KubernetesManager,
	authenticationModes authApi.AuthenticationModes, authenticationSkippable bool) authApi.AuthManager {
	return &authManager{
		k8sManager:              k8sManager,
		authenticationModes:     authenticationModes,
		authenticationSkippable: authenticationSkippable,
	}
}
