package args

var Holder = &holder{}

type holder struct {
	port           int
	apiServerHost  string
	kubeConfigFile string
	apiLogLevel    string

	namespace string
}

func (self *holder) GetPort() int {
	return self.port
}

func (self *holder) GetApiServerHost() string {
	return self.apiServerHost
}

func (self *holder) GetKubeConfigFile() string {
	return self.kubeConfigFile
}

func (self *holder) GetApiLogLevel() string {
	return self.apiLogLevel
}

func (self *holder) GetNamespace() string {
	return self.namespace
}
