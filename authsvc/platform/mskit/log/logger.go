package log

import (
	"github.com/go-kit/kit/log"
	"os"
)

type Logger log.Logger

var Mslog Logger

func init() {
	//logger = kitlog.NewLogfmtLogger(os.Stdout)
	// Logging domain.
	Mslog = log.NewLogfmtLogger(os.Stdout)
	Mslog = log.With(Mslog, "ts", log.DefaultTimestampUTC)
	Mslog = log.With(Mslog, "caller", log.DefaultCaller)
}