package main

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"os"
	"os/signal"
	"platform/common/utils"
	"platform/confsvc/common"
	"platform/confsvc/imconf"
	"platform/confsvc/models"
	"platform/confsvc/router"
	"platform/confsvc/rpc"
	"platform/confsvc/services"
	"platform/mskit/grace"
	"platform/mskit/log"
	"platform/mskit/rpcx"
	"platform/mskit/sd"
	"platform/mskit/trace"
	"platform/pfcomm/apis"
	"strconv"
)

var isChild			bool
var socketOrder		string

func init() {
	GetSettings()
	common.GetLogger()
	services.InitLogger()
	models.InitLogger()
	rpc.InitLogger()

	var err error
	if imconf.Config.DbDriver == "mysql" {
		err = orm.RegisterDataBase("default", "mysql", imconf.Config.Dsn)
	}else if imconf.Config.DbDriver == "pgsql" {
		err = orm.RegisterDataBase("default", "postgres", imconf.Config.Dsn)
	}

	if err != nil {
		panic(err)
	}

	err = orm.RunSyncdb("default", false, true)
	if err != nil {
		fmt.Println(err)
	}

	router.InitBefe()
}


func main() {

	imconf.Config.ContainerHttp = imconf.Config.HttpAddress
	imconf.Config.ContainerHttps = imconf.Config.HttpsAddress

	common.Logger.Finest("config= %+v",imconf.Config)

	run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

}


func run() {

	params := make(map[string]interface{})

	var options []sd.SdOption
	options = append(options,sd.SdTypeOption(imconf.Config.Sdt))
	options = append(options,sd.SdAddressOption(imconf.Config.Sda))
	sdc := sd.NewRegistar(options...)
	go func(){
		if imconf.Config.ServiceConf != "" {
			sdc.RegisterWithConf(nil,"rpcx",imconf.Config.ServiceConf,
				rpcService,
			)
		}else{
			imconf.Config.ContainerRpcx = imconf.Config.RpcxAddr
			rpc.RpcServer(imconf.Config.RpcxNetwork,imconf.Config.ConsulAddress,imconf.Config.ContainerRpcx,imconf.Config.RpcxBasepath)
		}
	}()

	if imconf.Config.HttpEnable {
		go func(){
			var msApp *grace.MicroService
			msApp = grace.NewServer(isChild,socketOrder,imconf.Config.ContainerHttp)
			router.InitRoute(imconf.Config.Prefix,msApp)
			params["host"] = imconf.Config.HttpHost
			params["port"] = imconf.Config.HttpPort
			params["interval"] = imconf.Config.HealthCheckInterval
			params["timeout"] = imconf.Config.HealthCheckTimeout

			if imconf.Config.Env == "prod" {
				if imconf.Config.ServiceConf == "" {
					sdc.Register(msApp,"http",imconf.Config.ServiceName,
						imconf.Config.ContainerHttp,
						httpService,
						params,
					)

				}else{
					sdc.RegisterWithConf(msApp,"http",imconf.Config.ServiceConf,
						httpService,
					)
				}
			} else if imconf.Config.Env == "dev" {
				port := strconv.Itoa(imconf.Config.HttpPort)
				msApp.Serve(imconf.Config.HttpHost, port)
			}
		}()
	}

	if imconf.Config.HttpsEnable {
		go func(){
			var msApp *grace.MicroService
			msApp = grace.NewServer(isChild,socketOrder,imconf.Config.ContainerHttps)
			router.InitRoute(imconf.Config.Prefix,msApp)
			params["host"] = imconf.Config.HttpsHost
			params["port"] = imconf.Config.HttpsPort
			params["interval"] = imconf.Config.HealthCheckInterval
			params["timeout"] = imconf.Config.HealthCheckTimeout

			if imconf.Config.Env == "prod" {
				if imconf.Config.ServiceConf == "" {
					sdc.Register(msApp,"https",imconf.Config.ServiceName,
						imconf.Config.ContainerHttps,
						httpsService,
						params,
					)

				}else{
					sdc.RegisterWithConf(msApp,"https",imconf.Config.ServiceConf,
						httpsService,
					)
				}
			} else if imconf.Config.Env == "dev" {
				port := strconv.Itoa(imconf.Config.HttpsPort)
				msApp.Serve(imconf.Config.HttpsHost, port)
			}
		}()
	}


}

func httpService(msApp *grace.MicroService,param map[string]interface{}) (err error) {

	//user code
	host := utils.ConvertToString(param["host"])
	port := utils.ConvertToString(param["port"])
	if isChild {
		err = msApp.ListenAndServe(host, port)
	}else{
		err = msApp.Serve(host, port)
	}
	return err

}

func httpsService(msApp *grace.MicroService,param map[string]interface{}) (err error) {

	//user code
	var certFile,keyFile,trustfile string
	host := utils.ConvertToString(param["host"])
	port := utils.ConvertToString(param["port"])
	certFile = utils.ConvertToString(param["certfile"])
	keyFile = utils.ConvertToString(param["keyfile"])
	trustfile = utils.ConvertToString(param["trustfile"])

	if trustfile == "" {
		err = msApp.ListenAndServeTLS(certFile,keyFile,host,port)
	}else{
		err = msApp.ListenAndServeMutualTLS(certFile,keyFile,trustfile,host,port)
	}

	return err
}

func rpcService(msApp *grace.MicroService,param map[string]interface{}) error {
	host := utils.ConvertToString(param["host"])
	port := utils.ConvertToString(param["port"])
	consul := utils.ConvertToString(param["consul_address"])
	imconf.Config.RpcxAddr = host +":" + utils.ConvertToString(port)

	imconf.Config.ContainerRpcx = imconf.Config.RpcxAddr

	tracer,_:= apis.CreateTracer(imconf.Config.RecordAddr,imconf.Config.ServiceName,log.Mslog,imconf.Config.Debug,
		imconf.Config.ZipkinUrl,imconf.Config.AppdashAddr,imconf.Config.LightstepToken,imconf.Config.KafkaAddress)

	var options []trace.TraceOption
	options = append(options,trace.WithTracerOption(true))
	options = append(options,trace.OpenTracerOption(tracer))

	mstracer := trace.NewTracer(options...)
	zkoption := rpcx.RpcxTracerOption(mstracer)
	sdtop := rpcx.RpcxSdTypeOption("consul")
	sdop := rpcx.RpcxSdAddressOption(consul)
	rpc.RpcServer(imconf.Config.RpcxNetwork, consul, imconf.Config.ContainerRpcx, imconf.Config.RpcxBasepath,zkoption,sdtop,sdop)

	return nil
}