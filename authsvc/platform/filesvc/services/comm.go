package services

import (
	"platform/filesvc/comm"
	l4g "github.com/libra9z/log4go"

)

var logger l4g.Logger

func InitLogger() {
	logger = comm.Logger
}

