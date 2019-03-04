package rpc

import (
	"platform/mskit/rpcx"
)


func RpcServer(network,consul,rpcxaddr,path string,op... rpcx.RpcxServerOptions) {

	var serviceName AuthServiceName
	serviceName.SetServiceName("AuthJSONRpc")

	logger.Finest("rpcxnetwork=%s,rpcxaddress=%s,consul-address=%s,servicename=%s",network,
		rpcxaddr,consul,serviceName.GetServiceName())
	var options []rpcx.RpcxServerOptions

	options = append(options,rpcx.RpcxBasePathOption(path))
	options = append(options,rpcx.RpcxNetworkOption(network))
	options = append(options,rpcx.RpcxServiceAddressOption(rpcxaddr))
	options = append(options,op...)

	rpcx.DefaultRpcServer(options...)

	rpcx.RpcRegisterService(&serviceName,new(AuthJSONRpc),"")

	//注册处理的handle 函数
	rpcx.RpcRegisterDefaultMethod("CheckAuth",CheckAuth)

	rpcx.RpcRegisterDefaultMethod("AddLicense",AddLicense)
	rpcx.RpcRegisterDefaultMethod("DeleteLicense",DeleteLicense)
	rpcx.RpcRegisterDefaultMethod("UpdateLicense",UpdateLicense)
	rpcx.RpcRegisterDefaultMethod("GetLicense",GetLicense)
	rpcx.RpcRegisterDefaultMethod("GetLicenseCounts",GetLicenseCounts)
	rpcx.RpcRegisterDefaultMethod("GetPkgServices",GetPkgServices)

	rpcx.RpcServe()
}

