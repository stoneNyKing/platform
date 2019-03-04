package rpc

import (
	l4g "github.com/libra9z/log4go"
	"platform/ousvc/common"
)

var logger l4g.Logger

func InitLogger() {
	logger = common.Logger
}
