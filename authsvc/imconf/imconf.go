// imconf project imconf.go
package imconf

import (
	//"flag"
	"fmt"
	"strconv"
	"sync"
	//修改为viper
	"platform/common/utils"

	"github.com/spf13/viper"
)

var once sync.Once

var (
	Appid int = 1
	Modid int = 103

	Mysqlpwd string = "InterMa140"
)

var SecIsused bool = true
var HasRedis bool = true

type ImConf struct {
	Dsn         string

	RedisAddr   string
	RedisDb     int
	RedisPool   int
	RedisEnable bool
	RedisPort   int

	//log level
	LogLevel string
	Logfile  string
	LogMaxSize		int
	LogRotateFiles 	int

	//数据库名
	AuthdbName   	string
	AuthdbSchema   	string

	DbDriver     string

	//consul
	ConsulToken   string
	ConsulAddress string
	Sdt 		  string
	Sda 		  string

	ServiceConf				string
	HealthCheckInterval		string
	HealthCheckTimeout		string

	//api
	UsercheckUrl  string
	IsAuth        bool

	Env         string
	ServiceName string

	//rpcx
	RpcxOamBasepath 	string
	RpcxBasepath 		string
	RpcxConfBasepath 	string
	RpcxAddr		 	string
	RpcxNetwork		 	string

	//consul
	HttpAddress		       	string
	HttpsAddress    	   	string

	//oa ou URL
	OasvcUrl 		string
	OusvcUrl 		string

	Prefix 				string

	//rest api
	HttpEnable 		bool
	HttpsEnable 	bool

	HttpHost   			string
	HttpPort   			int
	HttpsHost   		string
	HttpsPort   		int

	ContainerHttps			string
	ContainerHttp			string
	ContainerRpcx			string

	KafkaAddress          string
	ZipkinUrl		      string

	//监控地址
	Debug 					bool
	DebugAddr			string
	RecordAddr			string
	HttpAddr			string
	AppdashAddr			string
	LightstepToken		string
}

var (
	Config *ImConf
)


func (config *ImConf) ReadConf() {
	//修改为viper配置
	config.GetViperConfiguration()
}

func (config *ImConf) GetViperConfiguration() {
	config.LogLevel = "INFO"

	config.DbDriver = viper.GetString("self.system.dbdriver")
	if config.DbDriver == "" {
		config.DbDriver = "pgsql"
	}

	port := ""
	port = strconv.Itoa(viper.GetInt(config.DbDriver+".authdb.port"))
	//authdb dsn
	if config.DbDriver == "mysql" {
		config.Dsn = viper.GetString(config.DbDriver+".authdb.user") + ":" + viper.GetString(config.DbDriver+".authdb.passwd") + "@tcp(" +
			viper.GetString(config.DbDriver+".authdb.host") + ":" + port + ")/" + viper.GetString(config.DbDriver+".authdb.database")
	}else if config.DbDriver == "pgsql" {
		config.Dsn = fmt.Sprintf("postgres://%s:%s@%s/%s?connect_timeout=10&sslmode=disable",
			viper.GetString(config.DbDriver+".authdb.user"), viper.GetString(config.DbDriver+".authdb.passwd"),
			viper.GetString(config.DbDriver+".authdb.host") +":"+port, viper.GetString(config.DbDriver+".authdb.database"))
	}

	config.AuthdbName = viper.GetString(config.DbDriver+".authdb.database")
	config.AuthdbSchema = viper.GetString(config.DbDriver+".authdb.schema")

	config.RedisPort = viper.GetInt("redis.queue.port")
	config.RedisAddr = viper.GetString("redis.queue.host") + ":" + strconv.Itoa(config.RedisPort)
	config.RedisPool = viper.GetInt("redis.queue.pool")
	config.RedisDb = viper.GetInt("redis.queue.database")
	config.RedisEnable = viper.GetBool("redis.queue.enable")
	HasRedis = config.RedisEnable

	config.RpcxOamBasepath = viper.GetString("rpcx.server.oam_base_path")
	config.RpcxConfBasepath = viper.GetString("rpcx.server.conf_base_path")
	config.RpcxBasepath = viper.GetString("rpcx.server.auth_base_path")
	if config.RpcxBasepath == "" {
		config.RpcxBasepath = viper.GetString("rpcx.server.authsvc_base_path")
	}

	config.RpcxAddr = utils.Hostname2IPv4(viper.GetString("self.rpcx.address"))
	config.RpcxNetwork = viper.GetString("self.rpcx.network")

	config.LogLevel = viper.GetString("self.system.loglevel")
	config.Logfile = viper.GetString("self.system.logfile")
	config.LogMaxSize = viper.GetInt("self.system.logmaxsize")
	if config.LogMaxSize == 0 {
		config.LogMaxSize = 104857600
	}
	config.LogRotateFiles = viper.GetInt("self.system.log_rotate_files")
	if config.LogRotateFiles == 0 {
		config.LogRotateFiles = 10
	}

	config.UsercheckUrl = viper.GetString("self.system.usercheck_url")
	config.IsAuth = viper.GetBool("self.system.auth")

	config.OasvcUrl = viper.GetString("self.system.oasvc_service_url")
	config.OusvcUrl = viper.GetString("self.system.ousvc_service_url")

	config.Env = viper.GetString("self.system.env")
	if config.Env == "" {
		config.Env = "prod"
	}
	config.ServiceName = viper.GetString("self.system.service_name")
	if config.ServiceName == "" {
		config.ServiceName = "ms-" + strconv.Itoa(Appid) + "-" + strconv.Itoa(Modid)
	}

	config.HealthCheckInterval = viper.GetString("consul.check.interval")
	config.HealthCheckTimeout = viper.GetString("consul.check.timeout")

	config.HttpEnable = viper.GetBool("self.http_enable")
	config.HttpsEnable = viper.GetBool("self.https_enable")
	config.Prefix = viper.GetString("self.prefix")

	if config.HttpEnable {
		config.HttpHost = viper.GetString("self.http.host")
		config.HttpPort = viper.GetInt("self.http.port")
		config.HttpAddress = config.HttpHost+":"+utils.ConvertToString(config.HttpPort)
	}
	if config.HttpsEnable {
		config.HttpsHost = viper.GetString("self.https.host")
		config.HttpsPort = viper.GetInt("self.https.port")
		config.HttpsAddress = config.HttpsHost+":"+utils.ConvertToString(config.HttpsPort)
	}

	config.KafkaAddress = viper.GetString("self.zipkin.kafka_address")
	config.ZipkinUrl = viper.GetString("self.zipkin.zipkin_url")

}
