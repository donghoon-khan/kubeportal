package args

var builder = &holderBuilder{holder: Holder}

type holderBuilder struct {
	holder *holder
}

func (self *holderBuilder) SetPort(port int) *holderBuilder {
	self.holder.port = port
	return self
}

func (self *holderBuilder) SetApiServerHost(apiServerHost string) *holderBuilder {
	self.holder.apiServerHost = apiServerHost
	return self
}

func (self *holderBuilder) SetKubeConfigFile(kubeConfigFile string) *holderBuilder {
	self.holder.kubeConfigFile = kubeConfigFile
	return self
}

func (self *holderBuilder) SetApiLogLevel(apiLogLevel string) *holderBuilder {
	self.holder.apiLogLevel = apiLogLevel
	return self
}

func GetHolderBuilder() *holderBuilder {
	return builder
}
