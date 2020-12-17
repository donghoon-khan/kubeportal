package args

var Holder = &holder{}

type holder struct {
	port          int
	apiServerHost string
}

func (self *holder) GetPort() int {
	return self.port
}

func (self *holder) GetApiServerHost() string {
	return self.apiServerHost
}
