package kubernetes

import (
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/emicklei/go-restful"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	DefaultQPS                 = 1e6
	DefaultBurst               = 1e6
	DefaultContentType         = "application/vnd.kubernetes.protobuf"
	DefaultCmdConfigName       = "kubernetes"
	JWETokenHeader             = "jweToken"
	DefaultUserAgent           = "kube-portal"
	ImpersonateUserExtraHeader = "Impersonate-Extra-"
)

var Version = "UNKNOWN"

type kubernetesManager struct {
	csrfKey         string
	kubeConfigPath  string
	apiserverHost   string
	inClusterConfig *rest.Config
	//tokenManager authApi.TokenManager
	insecureAPIExtensionsKubernetes apiextensionsclientset.Interface
	//insecurePluginClient pluginclientset.Interface
	insecureKubernetes kubernetes.Interface
	insecureConfig     *rest.Config
}

func (self *kubernetesManager) Kubernetes(req *restful.Request) (kubernetes.Interface, error) {
	if req == nil {
		return nil, errors.NewBadRequest("request can not be nil")
	}
	if self.isSecureModeEnabled(req) {
		return self.secureKubernetes(req)
	}
	return self.InsecureKubernetes(), nil
}

func (self *kubernetesManager) APIExtensionsKubernetes(req *restful.Request) (apiextensionsclientset.Interface, error) {
	if req == nil {
		return nil, errors.NewBadRequest("request can not be nil")
	}

	if self.isSecureModeEnabled(req) {
		return self.secureAPIExtensionsKubernetes(req)
	}

	return self.InsecureAPIExtensionsKubernetes(), nil
}

func (self *kubernetesManager) InsecureKubernetes() kubernetes.Interface {
	return self.insecureKubernetes
}

func (self *kubernetesManager) InsecureAPIExtensionsKubernetes() apiextensionsclientset.Interface {
	return self.insecureAPIExtensionsKubernetes
}

// TODO: Implement check secure mode
func (self *kubernetesManager) isSecureModeEnabled(req *restful.Request) bool {
	return false
}

// TODO: Implement secure kubernetes client
func (self *kubernetesManager) secureKubernetes(req *restful.Request) (kubernetes.Interface, error) {
	return self.InsecureKubernetes(), nil
}

// TODO: Implemenet secure kubernetes extensions client
func (self *kubernetesManager) secureAPIExtensionsKubernetes(req *restful.Request) (apiextensionsclientset.Interface, error) {
	return self.InsecureAPIExtensionsKubernetes(), nil
}
