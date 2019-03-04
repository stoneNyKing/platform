package rpc

import (
	"platform/mskit/rpcx"
)


type AuthServiceName struct {
	name 		string
}
type AuthJSONRpc struct {
	rpcx.JSONRpc
}

func (jr *AuthServiceName ) GetServiceName() string {
	return jr.name
}
func (jr *AuthServiceName ) SetServiceName(name string)  {
	jr.name = name
}

