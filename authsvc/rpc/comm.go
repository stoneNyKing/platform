package rpc

import (
	"platform/authsvc/common"
	l4g "github.com/libra9z/log4go"

)

var logger l4g.Logger

func InitLogger() {
	logger = common.Logger
}
