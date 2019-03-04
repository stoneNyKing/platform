//
package rpcx

import (
	"context"
	"encoding/json"
	"errors"
	zipkin "github.com/openzipkin/zipkin-go"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"platform/mskit/trace"
	"strings"
	"platform/mskit/log"
	"time"
)

const (
	JSONRPC_ERR_METHOD_NOT_FOUND = 32601
)
type RpcxServerOptions func(* RpcServer)
type Method func(context.Context,trace.Tracer,int64, int64, string, interface{}) (interface{}, error)

type RpcServer struct {
	Server *server.Server
	logger log.Logger

	Network     string
	ServiceAddr string
	SdType		string
	SdAddress	string
	BasePath	string

	Methods map[string]Method

	//zipkinTracer 		*zipkin.Tracer
	//tracer 		opentracing.Tracer
	tracer 			trace.Tracer
}

var defautlServer *RpcServer

type RpcRequest struct {
	Appid  					int64
	SiteId 					int64
	Id     					int64 //修改某一条记录时的记录标识
	Token  					string
	Req    					string
	WithTracer    			bool
}

type RpcResponse struct {
	Ret string
}

type RpcService interface {
	Services(ctx context.Context, req *RpcRequest, ret *RpcResponse) error
}

type RpcServiceName interface {
	SetServiceName(string)
	GetServiceName() string
}


func RpcRegisterService(servName RpcServiceName, service RpcService, metadata string) {
	if defautlServer != nil && service != nil {
		defautlServer.RegisterService(servName, service, metadata)
	}
}

func RpcRegisterDefaultService(servName RpcServiceName, service RpcService, meta string) {
	if defautlServer != nil {
		defautlServer.RegisterDefaultService(servName, service, meta)
	} else {
		log.Mslog.Log("error","register default services failed.")
	}

}

func RpcRegisterDefaultMethod(methodName string, m Method) {

	if defautlServer != nil {
		defautlServer.RegisterMethod(methodName, m)
	} else {
		log.Mslog.Log("error","register default method failed.")
	}

}
func RegisterMethod(methodName string, m Method) {

	if defautlServer != nil {
		defautlServer.RegisterMethod(methodName, m)
	} else {
		log.Mslog.Log("error","register default method failed.")
	}

}

func RpcGetMethodByName(name string) Method {

	if defautlServer != nil {
		return defautlServer.GetMethodByName(name)
	}

	return nil
}

func RpcGetMethodWithTracer(name string) (Method, trace.Tracer) {

	if defautlServer != nil {
		return defautlServer.GetMethodByName(name),defautlServer.tracer
	}

	return nil,nil
}

func RpcGetMethodWithZipkinTracer(name string) (Method, *zipkin.Tracer) {

	if defautlServer != nil {
		return defautlServer.GetMethodByName(name),defautlServer.tracer.GetZipkinTracer()
	}

	return nil,nil
}

func RpcServe() {

	if defautlServer != nil {
		defautlServer.Serve()
	} else {
		log.Mslog.Log("error","cannot start Rpcx server,default server is nil.")
	}
}

func Serve() {
	RpcServe()
}

func (s *RpcServer) RegisterService(servName RpcServiceName, service RpcService, metadata string) {
	if service != nil {
		err := s.Server.RegisterName(servName.GetServiceName(), service, metadata)
		//s.Server.Register(service,metadata)
		if err != nil {
			s.logger.Log("error", err,"reason","不能注册服务")
		}
	}
}

func (s *RpcServer) RegisterDefaultService(servName RpcServiceName, service RpcService, meta string) {

	if service != nil {
		s.logger.Log("info", "注册服务")
		err := s.Server.RegisterName(servName.GetServiceName(), service, meta)
		//err := s.Server.Register(service,meta)
		if err != nil {
			s.logger.Log("error", err,"reason","不能注册服务")
		}
	} else {
		s.logger.Log("error", "不能注册服务，service为nil")
	}
}

func (s *RpcServer) Serve() error {

	s.logger.Log("rpcx server running on : ", s.ServiceAddr)
	err := s.Server.Serve(s.Network, s.ServiceAddr)

	if err != nil {
		s.logger.Log("cannot run rpcx server: ", err)
		return err
	}
	return nil
}

func (s *RpcServer) RegisterMethod(methodName string, m Method) {

	if methodName == "" {
		return
	}
	s.Methods[methodName] = m
}

func (s *RpcServer) GetMethodByName(name string) Method {

	if name == "" {
		return nil
	}

	if m, ok := s.Methods[name]; ok {
		return m
	}

	return nil
}
func (s *RpcServer) GetMethodWithTracer(name string) (Method,trace.Tracer) {

	if name == "" {
		return nil,nil
	}

	if m, ok := s.Methods[name]; ok {
		return m,s.tracer
	}

	return nil,nil
}


type JSONRpc struct{}

func (jr *JSONRpc) Services(ctx context.Context, req *RpcRequest, ret *RpcResponse) error {

	var err error
	if req == nil || ret == nil {
		err = errors.New("input parameter is nil")
		return err
	}

	if req.Req == "" {
		return errors.New("json-rpc request is empty.")
	}

	var vs map[string]interface{}
	err = json.Unmarshal([]byte(req.Req), &vs)

	if err != nil {
		return err
	}

	if vs["jsonrpc"] != nil {
		v := vs["jsonrpc"].(string)
		if v != "2.0" {
			return errors.New("unsupport json-rpc version.")
		}
	}

	var result interface{}
	em := make(map[string]interface{})
	if vs["method"] != nil {
		method := vs["method"].(string)
		log.Mslog.Log("method", method)
		if method != "" {
			var function Method
			var tracer trace.Tracer
			if req.WithTracer {
				function,tracer = RpcGetMethodWithTracer(method)
			}else{
				function = RpcGetMethodByName(method)
			}
			if function != nil {
				result, err = function(ctx,tracer,req.Appid, req.SiteId, req.Token, vs["params"])
			} else {
				log.Mslog.Log("error","没有找对对应的方法。")
			}
		} else {
			em["code"] = JSONRPC_ERR_METHOD_NOT_FOUND
			em["message"] = "该方法不存在或者无效"
			err = errors.New("method not found")
		}
	}

	var rm map[string]interface{}

	rm = make(map[string]interface{})

	rm["jsonrpc"] = "2.0"
	rm["result"] = result
	if err != nil {
		rm["error"] = em
	}
	if vs["id"] != nil {
		rm["id"] = vs["id"]
	}

	r, err := json.Marshal(&rm)

	if err != nil {
		return errors.New("cannot marshal return json.")
	}

	ret.Ret = string(r)

	return nil
}

// v2
func NewRpcxServer(options ...RpcxServerOptions) *RpcServer {

	s := &RpcServer{
		logger: log.Mslog,
		Server: server.NewServer(),
		Methods: make(map[string]Method),
	}

	for _,option := range options {
		option(s)
	}

	s.logger.Log("info","consul registering... ")

	//fmt.Printf("(options=%d)rpcserver=%+v\n",len(options),s)

	cs := strings.Split(s.SdAddress, ",")
	switch s.SdType {
	case "consul":
		p := &serverplugin.ConsulRegisterPlugin{
			ServiceAddress: s.Network + "@" + s.ServiceAddr,
			ConsulServers:  cs,
			BasePath:       s.BasePath,
			Metrics:        metrics.NewRegistry(),
			UpdateInterval: time.Minute,
		}
		err := p.Start()
		if err != nil {
			s.logger.Log("error", err)
		}
		s.Server.Plugins.Add(p)

	case "etcd":
		p := &serverplugin.EtcdRegisterPlugin{
			ServiceAddress: s.Network + "@" + s.ServiceAddr,
			EtcdServers : cs,
			BasePath:       s.BasePath,
			Metrics:        metrics.NewRegistry(),
			UpdateInterval: time.Minute,
		}
		err := p.Start()
		if err != nil {
			s.logger.Log("error", err)
		}
		s.Server.Plugins.Add(p)
	case "zookeeper":
		p := &serverplugin.ZooKeeperRegisterPlugin{
			ServiceAddress: s.Network + "@" + s.ServiceAddr,
			ZooKeeperServers : cs,
			BasePath:       s.BasePath,
			Metrics:        metrics.NewRegistry(),
			UpdateInterval: time.Minute,
		}
		err := p.Start()
		if err != nil {
			s.logger.Log("error", err)
		}
		s.Server.Plugins.Add(p)
	}


	if s.tracer != nil {
		zkp := serverplugin.OpenTracingPlugin{}
		s.Server.Plugins.Add(zkp)
	}
	return s
}

func DefaultRpcServer(options ...RpcxServerOptions) {
	defautlServer = NewRpcxServer(options...)
}

func RpcxBasePathOption( basepath string) RpcxServerOptions {
	return func(c *RpcServer){ c.BasePath = basepath}
}
func RpcxSdTypeOption( sdtype string) RpcxServerOptions {
	return func(c *RpcServer){ c.SdType = sdtype}
}
func RpcxSdAddressOption( sdaddress string) RpcxServerOptions {
	return func(c *RpcServer){ c.SdAddress = sdaddress}
}
func RpcxServiceAddressOption( svraddr string) RpcxServerOptions {
	return func(c *RpcServer){ c.ServiceAddr = svraddr}
}
func RpcxNetworkOption( network string) RpcxServerOptions {
	return func(c *RpcServer){ c.Network = network}
}
func RpcxTracerOption( tracer trace.Tracer) RpcxServerOptions {
	return func(c *RpcServer){ c.tracer = tracer}
}
