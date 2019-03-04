package rpc

import (
	"platform/mskit/rpcx"
)


type HrServiceName struct {
	name 		string
}
type HrJSONRpc struct {
	rpcx.JSONRpc
}

func (jr *HrServiceName ) GetServiceName() string {
	return jr.name
}
func (jr *HrServiceName ) SetServiceName(name string)  {
	jr.name = name
}

