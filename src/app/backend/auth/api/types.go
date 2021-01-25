package api

import (
	"time"

	"k8s.io/client-go/tools/clientcmd/api"
)

const (
	EncryptionKeyHolderName     = "kubernetes-dashboard-key-holder"
	CertificateHolderSecretName = "kubernetes-dashboard-certs"
	DefaultTokenTTL             = 900
)

type AuthenticationModes map[AuthenticationMode]bool

type ProtectedResource struct {
	ResourceName      string
	ResourceNamespace string
}

func (self AuthenticationModes) IsEnabled(mode AuthenticationMode) bool {
	_, exists := self[mode]
	return exists
}

func (self AuthenticationModes) Array() []AuthenticationMode {
	modes := []AuthenticationMode{}
	for mode := range self {
		modes = append(modes, mode)
	}

	return modes
}

func (self AuthenticationModes) Add(mode AuthenticationMode) {
	self[mode] = true
}

type AuthenticationMode string

func (self AuthenticationMode) String() string {
	return string(self)
}

const (
	Token AuthenticationMode = "token"
	Basic AuthenticationMode = "basic"
)

type AuthManager interface {
	Login(*LoginSpec) (*AuthResponse, error)
	Refresh(string) (string, error)
	AuthenticationModes() []AuthenticationMode
	AuthenticationSkippable() bool
}

type TokenManager interface {
	Generate(api.AuthInfo) (string, error)
	Decrypt(string) (*api.AuthInfo, error)
	Refresh(string) (string, error)
	SetTokenTTL(time.Duration)
}

type Authenticator interface {
	// GetAuthInfo returns filled AuthInfo structure that can be used for K8S api client creation.
	GetAuthInfo() (api.AuthInfo, error)
}

type LoginSpec struct {
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Token      string `json:"token,omitempty"`
	KubeConfig string `json:"kubeconfig,omitempty"`
}

type AuthResponse struct {
	JWEToken string  `json:"jweToken"`
	Errors   []error `json:"errors" swaggertype:"array,string"`
}

type TokenRefreshSpec struct {
	JWEToken string `json:"jweToken"`
}

type LoginModesResponse struct {
	Modes []AuthenticationMode `json:"modes"`
}

type LoginSkippableResponse struct {
	Skippable bool `json:"skippable"`
}
