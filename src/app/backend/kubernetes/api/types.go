package api

import (
	"github.com/emicklei/go-restful"
	v1 "k8s.io/api/authorization/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/1.5/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type KubernetesManager interface {
	Kubernetes(req *restful.Request) (kubernetes.Interface, error)
	InsecureKubernetes() kubernetes.Interface
	CanI(req *restful.Request, saar *v1.SelfSubjectAccessReview) bool
	Config(req *restful.Request) (*rest.Config, error)
	ClientCmdConfig(req *restful.Request) (clientcmd.ClientConfig, error)
	CSRFKey() string
	HasAccess(authInfo api.AuthInfo) error
	VerberClient(req *restful.Request) (ResourceVerber, error)
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
