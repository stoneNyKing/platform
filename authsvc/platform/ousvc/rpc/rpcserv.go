package rpc

import (
	"platform/mskit/rpcx"
)

type OusvcServiceName struct {
	name string
}
type OusvcJSONRpc struct {
	rpcx.JSONRpc
}

func (jr *OusvcServiceName) GetServiceName() string {
	return jr.name
}
func (jr *OusvcServiceName) SetServiceName(name string) {
	jr.name = name
}
