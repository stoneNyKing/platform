// ihealth is an service for system
package main

import (
	"github.com/astaxie/beego/orm"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"net/http"
	"os"
	"os/signal"
	"platform/mskit/grace"
	"platform/mskit/rpcx"
	"platform/mskit/sd"
	"platform/mskit/trace"
	tapis "platform/pfcomm/apis"
	"strconv"
	"strings"
	"time"

	"github.com/dchest/captcha"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/strip"
	"github.com/thoas/stats"

	"platform/lib/logstasher"

	"github.com/libra9z/log4go"
	"platform/common/utils"
	"platform/ousvc/apis"
	"platform/ousvc/common"
	"platform/ousvc/config"
	md "platform/ousvc/models"
	"platform/ousvc/rpc"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	mslog "platform/mskit/log"
)

var (
	isChild      bool
	socketOrder  string
	zipkinTracer *zipkin.Tracer
	tracer       opentracing.Tracer
)

func init() {
	GetSettings()
	common.GetLogger()
	apis.InitLogger()
	rpc.InitLogger()
	logstasher.InitLogger(getLogLevel(config.Config.LogLevel), common.Logger)

	//md.InitAppids()
	apis.InitUser()
	apis.InitUsers()
	apis.InitSystem()
	apis.InitWeixin()

	var err error
	if config.Config.DbDriver == "mysql" {
		err = orm.RegisterDataBase("default", "mysql", config.Config.Dsn)
	} else if config.Config.DbDriver == "pgsql" {
		err = orm.RegisterDataBase("default", "postgres", config.Config.Dsn)
	}

	if err != nil {
		panic(err)
	}
}

func getLogLevel(sl string) int {
	var lvl int
	switch config.Config.LogLevel {
	case "DEBUG":
		lvl = int(log4go.DEBUG)
	case "FINEST":
		lvl = int(log4go.FINEST)
	case "INFO":
		lvl = int(log4go.INFO)
	case "TRACE":
		lvl = int(log4go.TRACE)
	case "FINE":
		lvl = int(log4go.FINE)
	case "CRITICAL":
		lvl = int(log4go.CRITICAL)
	case "ERROR":
		lvl = int(log4go.ERROR)
	}

	return lvl
}

var logger = common.Logger

//errorcode: 100000 - 100099
func main() {

	config.Config.ContainerHttp = config.Config.HttpAddress
	config.Config.ContainerHttps = config.Config.HttpsAddress

	common.Logger.Finest("config= %+v", config.Config)

	tracer, zipkinTracer = tapis.CreateTracer(config.Config.RecordAddr, config.Config.ServiceName, mslog.Mslog, config.Config.Debug,
		config.Config.ZipkinUrl, config.Config.AppdashAddr, config.Config.LightstepToken, config.Config.KafkaAddress)

	run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

}

func run() {

	md.InitAppids()
	params := make(map[string]interface{})

	var options []sd.SdOption
	options = append(options, sd.SdTypeOption(config.Config.Sdt))
	options = append(options, sd.SdAddressOption(config.Config.Sda))
	sdc := sd.NewRegistar(options...)

	go func() {
		if config.Config.ServiceConf != "" {
			sdc.RegisterWithConf(nil, "rpcx", config.Config.ServiceConf,
				rpcService,
			)
		} else {
			config.Config.ContainerRpcx = config.Config.RpcxAddr
			rpc.RpcServer(config.Config.RpcxNetwork, config.Config.ConsulAddress, config.Config.ContainerRpcx, config.Config.RpcxBasepath)
		}
	}()

	if config.Config.HttpEnable {
		go func() {
			var msApp *grace.MicroService
			msApp = grace.NewServer(isChild, socketOrder, config.Config.ContainerHttp)
			params["host"] = config.Config.HttpHost
			params["port"] = config.Config.HttpPort
			params["interval"] = config.Config.HealthCheckInterval
			params["timeout"] = config.Config.HealthCheckTimeout

			if config.Config.Env == "prod" {
				if config.Config.ServiceConf == "" {
					sdc.Register(msApp, "http", config.Config.ServiceName,
						config.Config.ContainerHttp,
						httpService,
						params,
					)

				} else {
					sdc.RegisterWithConf(msApp, "http", config.Config.ServiceConf,
						httpService,
					)
				}
			} else if config.Config.Env == "dev" {
				s := &http.Server{
					Addr:           config.Config.HttpHost + ":" + strconv.Itoa(config.Config.HttpPort),
					Handler:        PublicHandler(config.Config.Prefix),
					ReadTimeout:    30 * time.Second,
					WriteTimeout:   30 * time.Second,
					MaxHeaderBytes: 1 << 20,
				}
				logger.Info("Listening...main[%s:%d]", config.Config.HttpHost, config.Config.HttpPort)
				logger.Error(s.ListenAndServe())
			}
		}()
	}

	if config.Config.HttpsEnable {
		go func() {
			var msApp *grace.MicroService
			msApp = grace.NewServer(isChild, socketOrder, config.Config.ContainerHttps)
			params["host"] = config.Config.HttpsHost
			params["port"] = config.Config.HttpsPort
			params["interval"] = config.Config.HealthCheckInterval
			params["timeout"] = config.Config.HealthCheckTimeout

			if config.Config.Env == "prod" {
				if config.Config.ServiceConf == "" {
					sdc.Register(msApp, "https", config.Config.ServiceName,
						config.Config.ContainerHttps,
						httpsService,
						params,
					)

				} else {
					sdc.RegisterWithConf(msApp, "https", config.Config.ServiceConf,
						httpsService,
					)
				}
			} else if config.Config.Env == "dev" {
				s := &http.Server{
					Addr:           config.Config.HttpsHost + ":" + strconv.Itoa(config.Config.HttpsPort),
					Handler:        PublicHandler(config.Config.Prefix),
					ReadTimeout:    30 * time.Second,
					WriteTimeout:   30 * time.Second,
					MaxHeaderBytes: 1 << 20,
				}
				logger.Info("Listening...main[%s:%d]", config.Config.HttpsHost, config.Config.HttpsPort)
				logger.Error(s.ListenAndServe())
			}
		}()
	}

}

func httpService(msApp *grace.MicroService, param map[string]interface{}) (err error) {

	//user code
	host := utils.ConvertToString(param["host"])
	port := utils.ConvertToString(param["port"])

	//msApp.Server.Handler = PublicHandler(config.Config.Prefix)
	msApp.Server.ReadTimeout = 30 * time.Second
	msApp.Server.WriteTimeout = 30 * time.Second
	msApp.Server.MaxHeaderBytes = 1 << 20

	var hcHandler *HttpService
	// create global zipkin http server middleware
	if zipkinTracer != nil {
		serverMiddleware := zipkinhttp.NewServerMiddleware(
			zipkinTracer, zipkinhttp.TagResponseSize(true),
		)
		hcHandler = &HttpService{handler: serverMiddleware(PublicHandler(config.Config.Prefix))}
	} else {
		hcHandler = &HttpService{handler: PublicHandler(config.Config.Prefix)}
	}
	msApp.Server.Handler = hcHandler

	if isChild {
		err = msApp.ListenAndServe(host, port)
	} else {
		err = msApp.Serve(host, port)
	}

	return err

}

func httpsService(msApp *grace.MicroService, param map[string]interface{}) error {

	//user code
	var certFile, keyFile string
	host := utils.ConvertToString(param["host"])
	port := utils.ConvertToString(param["port"])
	certFile = utils.ConvertToString(param["certfile"])
	keyFile = utils.ConvertToString(param["keyfile"])
	trustfile := utils.ConvertToString(param["trustfile"])

	logger.Fine("certfile=%s,keyfile=%s,trustfile=%s", certFile, keyFile, trustfile)

	msApp.Server.ReadTimeout = 30 * time.Second
	msApp.Server.WriteTimeout = 30 * time.Second
	msApp.Server.MaxHeaderBytes = 1 << 20

	var hcHandler *HttpService
	// create global zipkin http server middleware
	if zipkinTracer != nil {
		serverMiddleware := zipkinhttp.NewServerMiddleware(
			zipkinTracer, zipkinhttp.TagResponseSize(true),
		)
		hcHandler = &HttpService{handler: serverMiddleware(PublicHandler(config.Config.Prefix))}
	} else {
		hcHandler = &HttpService{handler: PublicHandler(config.Config.Prefix)}
	}
	msApp.Server.Handler = hcHandler

	var err error
	if trustfile == "" {
		err = msApp.ListenAndServeTLS(certFile, keyFile, host, port)
	} else {
		err = msApp.ListenAndServeMutualTLS(certFile, keyFile, trustfile, host, port)
	}

	return err
}

func rpcService(msApp *grace.MicroService, param map[string]interface{}) error {
	host := utils.ConvertToString(param["host"])
	port := utils.ConvertToString(param["port"])
	consul := utils.ConvertToString(param["consul_address"])
	config.Config.RpcxAddr = host + ":" + utils.ConvertToString(port)

	config.Config.ContainerRpcx = config.Config.RpcxAddr
	var options []trace.TraceOption
	options = append(options, trace.WithTracerOption(true))
	options = append(options, trace.OpenTracerOption(tracer))

	mstracer := trace.NewTracer(options...)
	zkoption := rpcx.RpcxTracerOption(mstracer)

	sdtop := rpcx.RpcxSdTypeOption("consul")
	sdop := rpcx.RpcxSdAddressOption(consul)
	rpc.RpcServer(config.Config.RpcxNetwork, consul, config.Config.ContainerRpcx, config.Config.RpcxBasepath, zkoption, sdtop, sdop)

	return nil
}

func checkrequest(c martini.Context, req *http.Request, r render.Render) {
	req.ParseForm()
	for k, v := range req.Form {
		str := k + strings.Join(v, "")
		if strings.Contains(str, "<") || strings.Contains(str, ">") {
			r.JSON(200, map[string]interface{}{"Ret": 100001, "Msg": "参数含有非法字符"})
			return
		}
	}
}

func PublicHandler(prefix string) *martini.ClassicMartini {
	middleware := stats.New()

	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(logstasher.LoggerByLevel())
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PUT", "OPTIONS", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "accept", "x-requested-with", "Content-Type", "Content-Range", "Content-Disposition", "Content-Description"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           30 * time.Minute,
	}))
	m.Use(func(c martini.Context, w http.ResponseWriter, r *http.Request) {
		beginning, recorder := middleware.Begin(w)
		c.Next()
		middleware.End(beginning, stats.WithRecorder(recorder))
	})

	//m.Any(prefix+"/health", checkrequest, ping)

	//userapi
	m.Any(prefix+"/user/.*", checkrequest, strip.Prefix(prefix+"/user"), apis.UserHander().ServeHTTP)
	m.Any(prefix+"/system/.*", checkrequest, strip.Prefix(prefix+"/system"), apis.SystemHander().ServeHTTP)
	m.Any(prefix+"/token/.*", checkrequest, strip.Prefix(prefix+"/token"), apis.TokenHander().ServeHTTP)
	m.Get(prefix+"/captcha/.*", checkrequest, captcha.Server(80, 40).ServeHTTP)
	m.Any(prefix+"/qr/.*", checkrequest, strip.Prefix(prefix+"/qr"), apis.QrHander().ServeHTTP)
	m.Any(prefix+"/weixin/.*", checkrequest, strip.Prefix(prefix+"/weixin"), apis.WeixinHander().ServeHTTP)

	return m
}