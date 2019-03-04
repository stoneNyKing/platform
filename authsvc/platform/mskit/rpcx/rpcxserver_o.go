package rpcx

import (
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	metrics "github.com/rcrowley/go-metrics"
	"platform/mskit/log"
	"strings"
	"time"
)

/*
	参数network的定义如下：
	kcp：
	reuseport：
	quic
	default   tcp
*/
func InitRpcServerWithConsul(network, serviceAddr string, consulAddr string, basepath string) {

	defautlServer = NewRpcServerWithConsul(network, serviceAddr, consulAddr, basepath)
	if defautlServer == nil {
		log.Mslog.Log("error","cannot initial rpc server.")
	}
}


func NewRpcServerWithConsul(network, serviceAddr string, consulAddr string, basepath string) *RpcServer {

	s := &RpcServer{
		logger: log.Mslog,
		Server: server.NewServer(),
	}

	if network == "" {
		network = "tcp"
	}

	s.logger.Log("info","开始向consul注册服务...")

	cs := strings.Split(consulAddr, ",")

	s.Network = network
	s.ServiceAddr = serviceAddr
	s.Methods = make(map[string]Method)

	p := &serverplugin.ConsulRegisterPlugin{
		ServiceAddress: network + "@" + serviceAddr,
		ConsulServers:  cs,
		BasePath:       basepath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}

	err := p.Start()
	if err != nil {
		s.logger.Log("error", err)
	}
	s.Server.Plugins.Add(p)

	return s
}
