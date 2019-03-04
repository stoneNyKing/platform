//+build docker
//

package common

import (
	l4g "github.com/libra9z/log4go"
	"platform/common/utils"
	"platform/ousvc/config"

	"fmt"
	"os"
	"path/filepath"
	"time"
)

var Logger l4g.Logger

func GetLogger() {

	var err error
	var logf string

	//logf = config.Config.Logfile

	logf = "/data/logs/ousvc.log"

	dir := filepath.Dir(logf)

	fmt.Printf("dir=%s\n", dir)

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
	}

	f, err := os.Open(logf)
	if err == nil && !os.IsNotExist(err) {
		fmt.Printf("file exist!\n")
		loggff := logf + "." + utils.GetTimeFormat("[2006-01-02 15:04:05]")
		f.Close()

		os.Rename(logf, loggff)
	}

	initLogger(logf, config.Config.LogLevel)
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
	flw.SetRotateSize(config.Config.LogMaxSize)
	flw.SetRotateFiles(config.Config.LogRotateFiles)

	Logger.AddFilter("logfile", lvl, flw)
	Logger.Info("Current time is : %s\n", time.Now().Format("15:04:05 MST 2006/01/02"))

	return
}
