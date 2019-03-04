package router


import (
	"platform/authsvc/imconf"
	"platform/pfcomm/apis"
)

const(
	BEFE_USER	 = 1
	BEFE_ADMIN	 = 2
)

var befe map[int64]int

func InitBefe() {
	befe = apis.GetAppidBefeFlag(imconf.Config.RpcxConfBasepath,imconf.Config.ConsulAddress)
}

