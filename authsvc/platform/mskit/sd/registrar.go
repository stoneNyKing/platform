package sd

import (
	"platform/common/utils"
	"platform/mskit/grace"
)

type Registar interface {
	Register(app *grace.MicroService, schema, name string, prefix string,  callback ServiceCallback, params map[string]interface{})
	RegisterWithConf(app *grace.MicroService, schema string, fname string, callbacks ...ServiceCallback)
}

type serviceDiscovery struct {
	SdType 		string
	SdAddress 	string
	SdToken  	string
}

type SdOption func(*serviceDiscovery)

func NewRegistar(options ...SdOption) Registar {
	s := &serviceDiscovery{}
	for _,option := range options {
		option(s)
	}

	return s
}

func SdTypeOption(sdtype string) SdOption {
	return func(s *serviceDiscovery) { s.SdType = sdtype }
}
func SdAddressOption(sdaddress string) SdOption {
	return func(s *serviceDiscovery) { s.SdAddress = sdaddress }
}
func SdTokenOption(sdtoken string) SdOption {
	return func(s *serviceDiscovery) { s.SdToken = sdtoken }
}

func (s *serviceDiscovery) Register(app *grace.MicroService, schema, name string, prefix string,  callback ServiceCallback, params map[string]interface{}) {
	if params == nil {
		logger.Log("error","no input parameters")
		return
	}
	addr :=""
	if params["host"] != nil && params["port"] != nil {
		addr = utils.ConvertToString(params["host"])+":"+utils.ConvertToString(params["port"])
	}else{
		logger.Log("error","host or port not set.")
		return
	}
	switch s.SdType {
	case "consul":
		Register(app,schema,name,prefix,addr,s.SdAddress,s.SdToken,callback,params)
	}
}

func (s *serviceDiscovery) RegisterWithConf(app *grace.MicroService, schema string, fname string, callbacks ...ServiceCallback){
	switch s.SdType {
	case "consul":
		RegisterWithConf(app,schema,fname,s.SdAddress,s.SdToken,callbacks...)
	}
}