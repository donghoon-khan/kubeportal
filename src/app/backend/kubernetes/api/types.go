package api

import (
	"github.com/emicklei/go-restful"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

const (
	CsrfTokenSecretName = "kube-portal-csrf"
	CsrfTokenSecretData = "csrf"
)

type KubernetesManager interface {
	Kubernetes(req *restful.Request) (kubernetes.Interface, error)
	InsecureKubernetes() kubernetes.Interface

	APIExtensionsKubernetes(req *restful.Request) (apiextensionsclientset.Interface, error)
	InsecureAPIExtensionsKubernetes() apiextensionsclientset.Interface

	//PluginKubernetes(req *restful.Request) (pluginclientset.Interface, error)
	//InsecurePluginKubernetes() pluginclientset.Interface

	//CanI(req *restful.Request, saar *v1.SelfSubjectAccessReview) bool
	//Config(req *restful.Request) (*rest.Config, error)
	//ClientCmdConfig(req *restful.Request) (clientcmd.ClientConfig, error)
	//CSRFKey() string
	//HasAccess(authInfo api.AuthInfo) error
	//VerberClient(req *restful.Request) (ResourceVerber, error)
	//SetTokenManager(manager authApi.TokenManager)
}

type ResourceVerber interface {
	Put(kind string, namespaceSet bool, namespace string, name string,
		object *runtime.Unknown) error
	Get(kind string, namespaceSet bool, namespace string, name string) (runtime.Object, error)
	Delete(kind string, namespaceSet bool, namespace string, name string) error
}

type CanIResponse struct {
	Allowed bool `json:"allowed"`
}

type CsrfTokenManager interface {
	Token() string
}
