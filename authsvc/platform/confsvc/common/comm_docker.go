//+build docker
//

package common

import (
	l4g "github.com/libra9z/log4go"
	"platform/confsvc/imconf"
	"platform/common/utils"

	//"fmt"
	"os"
	"path/filepath"
	"time"

)


var Logger l4g.Logger

const(
	PLATVER_V1 	= 10
	PLATVER_V2	= 20		
)


func GetLogger() {

	var err error
	var logf string

	logf = imconf.Config.Logfile

	logf = "/data/logs/confsvc.log"
	dir := filepath.Dir(logf)

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
	}

	f, err := os.Open(logf)
	if err == nil && !os.IsNotExist(err) {
		loggff := logf + "." + utils.GetTimeFormat("20060102150405")
		f.Close()

		os.Rename(logf, loggff)
	}

	initLogger(logf, imconf.Config.LogLevel)
}

func initLogger(filename string, level string) {
	Logger = make(l4g.Logger)

	lvl := l4g.INFO
	switch level {
	case "DEBUG":
		lvl = l4g.DEBUG
	case "FINEST":
		lvl = l4g.FINEST
	case "INFO":
		lvl = l4g.INFO
	case "TRACE":
		lvl = l4g.TRACE
	case "FINE":
		lvl = l4g.FINE
	case "CRITICAL":
		lvl = l4g.CRITICAL
	case "ERROR":
		lvl = l4g.ERROR
	}

//	Logger.AddFilter("stdout", lvl, l4g.NewConsoleLogWriter())

	if _, err := os.Stat(filename); err == nil {
		os.Remove(filename)
	}

	flw := l4g.NewFileLogWriter(filename, true)
	flw.SetRotateSize(imconf.Config.LogMaxSize)
	flw.SetRotateFiles(imconf.Config.LogRotateFiles)

	Logger.AddFilter("logfile", lvl, flw)
	Logger.Info("Current time is : %s\n", time.Now().Format("15:04:05 MST 2006/01/02"))

	return
}

