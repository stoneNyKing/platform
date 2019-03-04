// imconf project imconf.go
package imconf

import (
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
	Modid int = 114

	Mysqlpwd string = "InterMa140"
)

var SecIsused bool = true
var HasRedis bool = true

type ImConf struct {
	Dsn          string
	ObjectsdbDsn string

	Httpservice string
	RedisAddr   string
	RedisDb     int
	RedisPool   int
	RedisEnable bool
	RedisPort   int

	//log level
	LogLevel       string
	Logfile        string
	LogMaxSize     int
	LogRotateFiles int

	//数据库名
	ObjectsdbName   string
	ObjectsdbSchema string
	StaffdbName     string
	StaffdbSchema   string
	DbDriver        string
	DefaultSchema   string

	//consul
	ConsulToken   string
	ConsulAddress string
	Sdt			  string
	Sda			  string

	//api
	UsercheckUrl  string
	IsAuth        bool

	Env           string
	ServiceName   string
	ServiceUrl    string

	ServiceConf         string
	HealthCheckInterval string
	HealthCheckTimeout  string

	//rpcx
	RpcxOamBasepath string
	RpcxConfBasepath string
	RpcxOaBasepath string

	RpcxBasepath string
	RpcxAddr     string
	RpcxNetwork  string

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

	ZipkinUrl       	string
	KafkaAddress		string

	//监控地址
	Debug          bool
	DebugAddr      string
	RecordAddr     string
	HttpAddr       string
	AppdashAddr    string
	LightstepToken string
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

	//staffdb dsn
	config.DbDriver = viper.GetString("self.system.dbdriver")
	if config.DbDriver == "" {
		config.DbDriver = "pgsql"
	}

	port := ""
	port = strconv.Itoa(viper.GetInt(config.DbDriver + ".staffdb.port"))
	//staffdb dsn
	if config.DbDriver == "mysql" {
		config.Dsn = viper.GetString(config.DbDriver+".staffdb.user") + ":" + viper.GetString(config.DbDriver+".staffdb.passwd") + "@tcp(" +
			viper.GetString(config.DbDriver+".staffdb.host") + ":" + port + ")/" + viper.GetString(config.DbDriver+".staffdb.database")
	} else if config.DbDriver == "pgsql" {
		config.Dsn = fmt.Sprintf("postgres://%s:%s@%s/%s?connect_timeout=10&sslmode=disable",
			viper.GetString(config.DbDriver+".staffdb.user"), viper.GetString(config.DbDriver+".staffdb.passwd"),
			viper.GetString(config.DbDriver+".staffdb.host")+":"+port, viper.GetString(config.DbDriver+".staffdb.database"))
	}

	config.StaffdbName = viper.GetString(config.DbDriver + ".staffdb.database")
	config.StaffdbSchema = viper.GetString(config.DbDriver + ".staffdb.schema")

	config.DefaultSchema = config.StaffdbSchema

	//oadb dsn
	port = strconv.Itoa(viper.GetInt(config.DbDriver + ".oadb.port"))
	//oadb dsn
	if config.DbDriver == "mysql" {
		config.ObjectsdbDsn = viper.GetString(config.DbDriver+".oadb.user") + ":" + viper.GetString(config.DbDriver+".oadb.passwd") + "@tcp(" +
			viper.GetString(config.DbDriver+".oadb.host") + ":" + port + ")/" + viper.GetString(config.DbDriver+".oadb.database")
	} else if config.DbDriver == "pgsql" {
		config.ObjectsdbDsn = fmt.Sprintf("postgres://%s:%s@%s/%s?connect_timeout=10&sslmode=disable",
			viper.GetString(config.DbDriver+".oadb.user"), viper.GetString(config.DbDriver+".oadb.passwd"),
			viper.GetString(config.DbDriver+".oadb.host")+":"+port, viper.GetString(config.DbDriver+".oadb.database"))
	}

	config.ObjectsdbName = viper.GetString(config.DbDriver + ".oadb.database")
	config.ObjectsdbSchema = viper.GetString(config.DbDriver + ".oadb.schema")

	config.RedisPort = viper.GetInt("redis.queue.port")
	config.RedisAddr = utils.Hostname2IPv4(viper.GetString("redis.queue.host")) + ":" + strconv.Itoa(config.RedisPort)
	config.RedisPool = viper.GetInt("redis.queue.pool")
	config.RedisDb = viper.GetInt("redis.queue.database")
	config.RedisEnable = viper.GetBool("redis.queue.enable")
	HasRedis = config.RedisEnable

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

	config.ServiceUrl = viper.GetString("self.system.service_url")

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

	config.RpcxOamBasepath = viper.GetString("rpcx.server.oam_base_path")
	config.RpcxOaBasepath = viper.GetString("rpcx.server.oasvc_base_path")
	config.RpcxConfBasepath = viper.GetString("rpcx.server.conf_base_path")
	config.RpcxBasepath = viper.GetString("rpcx.server.hrsvc_base_path")

	config.RpcxAddr = utils.Hostname2IPv4(viper.GetString("self.rpcx.address"))
	config.RpcxNetwork = viper.GetString("self.rpcx.network")

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

	config.ZipkinUrl = viper.GetString("self.zipkin.zipkin_url")
	config.KafkaAddress = viper.GetString("self.zipkin.kafka_address")

}
