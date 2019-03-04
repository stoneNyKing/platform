package rpcx

import (
	"context"
	"fmt"
	"github.com/smallnest/rpcx/client"
	"platform/mskit/trace"
	"strings"
)

func RpcCallWithConsul(basepath, consuladdr, serviceName, methodName string, selectMode int, req *RpcRequest, ret *RpcResponse) error {

	ss := strings.Split(consuladdr, ";")

	d := client.NewConsulDiscovery(basepath, serviceName, ss, nil)

	if selectMode < 0 {
		selectMode = int(client.RandomSelect)
	}

	client := client.NewXClient(serviceName, client.Failtry, client.SelectMode(selectMode), d, client.DefaultOption)
	defer client.Close()

	serviceMethod := methodName
	err := client.Call(context.Background(), serviceMethod, req, ret)
	if err != nil {
		fmt.Printf("error for %s: %v \n", serviceMethod, err)
	} else {
		fmt.Printf("%s: call success.\n", serviceMethod)
	}

	return nil
}

func RpcxCall(ctx context.Context,tracer trace.Tracer,
			sdtype,sdaddr string,
			basepath, serviceName,service, methodName string,
			failMode client.FailMode,selectMode client.SelectMode,
			req *RpcRequest) (ret *RpcResponse,err error) {


	var options []ClientOption

	options = append(options,BasePathOption(basepath))
	options = append(options,SdAddressOption(sdaddr))
	options = append(options,SdTypeOption(sdtype))
	options = append(options,FailModeOption(failMode))
	options = append(options,SelectModeOption(selectMode))
	options = append(options,MethodOption(methodName))
	options = append(options,ServiceOption(service))
	options = append(options,ServiceNameOption(serviceName))

	options = append(options,RpcxClientOpenTracing(tracer))

	resp := RpcResponse{}
	c := NewClient(&resp,options...)
	pc := c.GetClientPool()
	c.client = pc.Get().(client.XClient)
	defer c.Close()

	r,err := c.Endpoint()(ctx,req)
	if r != nil {
		return  r.(*RpcResponse),err
	}
	return nil,err
}