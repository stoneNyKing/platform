package router


import (
	"platform/filesvc/imconf"
	"platform/pfcomm/apis"
)

const(
	BEFE_USER	 = 1
	BEFE_ADMIN	 = 2
)

var befe map[int64]int

func InitBefe() {
	befe = apis.RetrieveAppidBefeFlag(imconf.Config.Sdt,imconf.Config.Sda,imconf.Config.RpcxConfBasepath)
}

