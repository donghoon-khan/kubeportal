package api

import (
	"strings"

	"github.com/donghoon-khan/kubeportal/src/app/backend/args"
)

func ToAuthenticationModes(modes []string) AuthenticationModes {
	result := AuthenticationModes{}
	modesMap := map[string]bool{}

	for _, mode := range []AuthenticationMode{Token, Basic} {
		modesMap[mode.String()] = true
	}

	for _, mode := range modes {
		if _, exists := modesMap[mode]; exists {
			result.Add(AuthenticationMode(mode))
		}
	}

	return result
}

func ShouldRejectRequest(url string) bool {
	for _, protectedResource := range protectedResources {
		if strings.Contains(url, protectedResource.ResourceName) && strings.Contains(url, protectedResource.ResourceNamespace) {
			return true
		}
	}
	return false
}

var protectedResources = []ProtectedResource{
	{EncryptionKeyHolderName, args.Holder.GetNamespace()},
	{CertificateHolderSecretName, args.Holder.GetNamespace()},
}
