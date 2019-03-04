package rpc

import (
	"platform/mskit/rpcx"
)


type ConfServiceName struct {
	name 		string
}
type ConfJSONRpc struct {
	rpcx.JSONRpc
}

func (jr *ConfServiceName ) GetServiceName() string {
	return jr.name
}
func (jr *ConfServiceName ) SetServiceName(name string)  {
	jr.name = name
}

