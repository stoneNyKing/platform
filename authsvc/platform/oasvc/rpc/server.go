package rpc

import (
	"platform/mskit/rpcx"
)

func RpcServer(network, consul, rpcxaddr, path string,op... rpcx.RpcxServerOptions) {

	var serviceName OaServiceName
	serviceName.SetServiceName("OaJSONRpc")

	logger.Finest("rpcxnetwork=%s,rpcxaddress=%s,consul-address=%s,servicename=%s", network,
		rpcxaddr, consul, serviceName.GetServiceName())

	//rpcx.InitRpcServerWithConsul(network, rpcxaddr, consul, path)
	var options []rpcx.RpcxServerOptions

	options = append(options,rpcx.RpcxBasePathOption(path))
	options = append(options,rpcx.RpcxNetworkOption(network))
	options = append(options,rpcx.RpcxServiceAddressOption(rpcxaddr))
	options = append(options,op...)

	rpcx.DefaultRpcServer(options...)

	rpcx.RpcRegisterService(&serviceName, new(OaJSONRpc), "")

	//注册处理的handle 函数
	rpcx.RegisterMethod("GetAdminList", GetAdminList)
	rpcx.RegisterMethod("AddAdmin", AddAdmin)
	rpcx.RegisterMethod("UpdateAdmin", UpdateAdmin)
	rpcx.RegisterMethod("ReadOrCreateAdmin", ReadOrCreateAdmin)
	rpcx.RegisterMethod("DeleteAdmin", DeleteAdmin)

	rpcx.Serve()
}
