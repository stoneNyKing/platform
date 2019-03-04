package logstasher

import (
	"net/http"
	"time"

	"github.com/libra9z/log4go"
	"github.com/go-martini/martini"
)

var log log4go.Logger
var loglevel log4go.Level

func InitLogger(lvl int,logger log4go.Logger) {
	log = logger
	loglevel = log4go.Level(lvl)
}

// Logger returns a middleware handler prints the request in a Logstash-JSON compatiable format
func Logger() martini.Handler {
	if log == nil {
		log = log4go.Global
	}
	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		start := time.Now()
		rw := res.(martini.ResponseWriter)
		c.Next()
		if req.URL.String() != "/health" {
			log.Fine("IP:%s\tHost:%s\tUserAgent:%s\tMethod:%s\tURL:%s\tStatus:%d\tSize:%d\tDuration:%fs\tHead:%v\tParams:%v", req.RemoteAddr, req.Host, req.UserAgent(), req.Method, req.URL.String(), rw.Status(), rw.Size(), time.Since(start).Seconds(), req.Header, map[string][]string(req.Form))
		}
	}
}


func LoggerByLevel() martini.Handler {
	if log == nil {
		log = log4go.Global
	}
	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		start := time.Now()
		rw := res.(martini.ResponseWriter)
		c.Next()
		log.Logf(log4go.Level(loglevel),"IP:%s\tHost:%s\tUserAgent:%s\tMethod:%s\tURL:%s\tStatus:%d\tSize:%d\tDuration:%fs\tHead:%v\tParams:%v", req.RemoteAddr, req.Host, req.UserAgent(), req.Method, req.URL.String(), rw.Status(), rw.Size(), time.Since(start).Seconds(), req.Header, map[string][]string(req.Form))
	}
}
