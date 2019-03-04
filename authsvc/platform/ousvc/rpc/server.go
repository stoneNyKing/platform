package rpc

import (
	"platform/mskit/rpcx"
)

func RpcServer(network, consul, rpcxaddr, path string,op... rpcx.RpcxServerOptions) {

	var serviceName OusvcServiceName
	serviceName.SetServiceName("OusvcJSONRpc")

	logger.Finest("rpcxnetwork=%s,rpcxaddress=%s,consul-address=%s,servicename=%s", network,
		rpcxaddr, consul, serviceName.GetServiceName())

	var options []rpcx.RpcxServerOptions

	options = append(options,rpcx.RpcxBasePathOption(path))
	options = append(options,rpcx.RpcxNetworkOption(network))
	options = append(options,rpcx.RpcxServiceAddressOption(rpcxaddr))

	options = append(options,op...)

	rpcx.DefaultRpcServer(options...)

	rpcx.RpcRegisterService(&serviceName, new(OusvcJSONRpc), "")

	//注册处理的handle 函数
	rpcx.RpcRegisterDefaultMethod("CheckUser", CheckUser)
	rpcx.RpcRegisterDefaultMethod("AddUser", AddUser)
	rpcx.RpcRegisterDefaultMethod("ReadOrCreateUser", ReadOrCreateUser)
	rpcx.RpcRegisterDefaultMethod("UpdateUser", UpdateUser)
	rpcx.RpcRegisterDefaultMethod("GetUserList", GetUserList)
	rpcx.RpcRegisterDefaultMethod("DeleteUser", DeleteUser)

	rpcx.RpcServe()
}
