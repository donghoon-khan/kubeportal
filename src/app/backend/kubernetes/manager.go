package kubernetes

import (
	"github.com/emicklei/go-restful/v3"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	kubernetesapi "github.com/donghoon-khan/kubeportal/src/app/backend/kubernetes/api"
)

// Portal UI default values for kubernetes client configs.
const (
	DefaultQPS                 = 1e6
	DefaultBurst               = 1e6
	DefaultContentType         = "application/vnd.kubernetes.protobuf"
	DefaultCmdConfigName       = "kubernetes"
	JWETokenHeader             = "jweToken"
	DefaultUserAgent           = "kube-portal"
	ImpersonateUserExtraHeader = "Impersonate-Extra-"
)

// Version of this binary
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

func (self *kubernetesManager) buildConfigFromFlags(apiserverHost, kubeConfigPath string) (*rest.Config, error) {
	if len(kubeConfigPath) > 0 || len(apiserverHost) > 0 {
		return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath},
			&clientcmd.ConfigOverrides{ClusterInfo: api.Cluster{Server: apiserverHost}}).ClientConfig()
	}

	if self.isRunningInCluster() {
		return self.inClusterConfig, nil
	}

	return nil, errors.NewInvalid("could not create client config")
}

func (self *kubernetesManager) init() {
	self.initInsecureKubernetes()
}

func (self *kubernetesManager) initConfig(config *rest.Config) {
	config.QPS = DefaultQPS
	config.Burst = DefaultBurst
	config.ContentType = DefaultContentType
	config.UserAgent = DefaultUserAgent + "/" + Version
}

func (self *kubernetesManager) initInsecureKubernetes() {
	self.initInsecureConfig()
	k8sClient, err := kubernetes.NewForConfig(self.insecureConfig)
	if err != nil {
		panic(err)
	}

	apiextensionsclient, err := apiextensionsclientset.NewForConfig(self.insecureConfig)
	if err != nil {
		panic(err)
	}

	//TODO: pluginClient 추가

	self.insecureKubernetes = k8sClient
	self.insecureAPIExtensionsKubernetes = apiextensionsclient
	//TODO: pluginClient mapping
}

func (self *kubernetesManager) initInsecureConfig() {
	config, err := self.buildConfigFromFlags(self.apiserverHost, self.kubeConfigPath)
	if err != nil {
		panic(err)
	}

	self.initConfig(config)
	self.insecureConfig = config
}

// TODO: check is in cluster
func (self *kubernetesManager) isRunningInCluster() bool {
	return false
}

func NewKubernetesManager(kubeConfigPath, apiserverHost string) kubernetesapi.KubernetesManager {
	result := &kubernetesManager{
		kubeConfigPath: kubeConfigPath,
		apiserverHost:  apiserverHost,
	}

	result.init()
	return result
}
