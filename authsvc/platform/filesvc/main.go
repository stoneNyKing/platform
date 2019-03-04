package main

import (
	"os"
	"os/signal"
	"platform/common/utils"
	"platform/filesvc/comm"
	"platform/filesvc/dbmod"
	"platform/mskit/grace"
	"platform/mskit/sd"
	"strconv"
	l4g "github.com/libra9z/log4go"

	_ "github.com/spf13/viper/remote"
	"platform/filesvc/ss/aliyun"
	"platform/filesvc/models"
	"platform/filesvc/services"
	"platform/filesvc/router"
	"platform/filesvc/imconf"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/astaxie/beego/orm"
)


var logger l4g.Logger

var isChild			bool
var socketOrder		string


func init(){
	GetSettings()
	models.InitLogger()
	services.InitLogger()

	aliyun.InitOs()
	models.InitFilename()
	models.InitDatabase()
	dbmod.InitDbLogger()
	dbmod.InitDatabase()

	var err error
	if imconf.Config.DbDriver == "mysql" {
		err = orm.RegisterDataBase("default", "mysql", imconf.Config.Dsn)
	}else if imconf.Config.DbDriver == "pgsql" {
		err = orm.RegisterDataBase("default", "postgres", imconf.Config.Dsn)
	}

	if err != nil {
		panic(err)
	}

	if imconf.Config.IsAuth {
		router.InitBefe()
	}
}

func main() {

	imconf.Config.ContainerHttp = imconf.Config.HttpAddress
	imconf.Config.ContainerHttps = imconf.Config.HttpsAddress

	comm.Logger.Finest("config= %+v",imconf.Config)

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
