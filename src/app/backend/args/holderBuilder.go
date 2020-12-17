package args

var builder = &holderBuilder{holder: Holder}

type holderBuilder struct {
	holder *holder
}

func (self *holderBuilder) SetPort(port int) *holderBuilder {
	self.holder.port = port
	return self
}

func (self *holderBuilder) setApiServerHost(apiServerHost string) *holderBuilder {
	self.holder.apiServerHost = apiServerHost
	return self
}

func GetHolderBuilder() *holderBuilder {
	return builder
}
