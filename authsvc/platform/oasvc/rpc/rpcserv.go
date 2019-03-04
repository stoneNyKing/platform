package rpc

import (
	"platform/mskit/rpcx"
)


type OaServiceName struct {
	name 		string
}
type OaJSONRpc struct {
	rpcx.JSONRpc
}

func (jr *OaServiceName ) GetServiceName() string {
	return jr.name
}
func (jr *OaServiceName ) SetServiceName(name string)  {
	jr.name = name
}

