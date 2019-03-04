package config

import (
	"fmt"
	"sync"
	//"platform/caller/base"
	"github.com/spf13/viper"
	"platform/common/utils"
	"strconv"
)

func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

type ImConf struct {
	Dsn string

	HttpHost  string
	HttpPort  int
	HttpsHost string
	HttpsPort int

	Env string // run mode, "dev" or "prod"

	DbAddr   string
	DbPort   int
	DbUser   string
	DbPasswd string

	Database      string
	DbDriver      string
	ObjectsSchema string

	SmsUrl      string
	ServiceName string

	SessionKey       string
	SessionStoreIP   string
	SessionStorePort string

	Logfile        string
	LogLevel       string
	LogConfig      string
	LogMaxSize     int
	LogRotateFiles int
	//Redis
	RedisAddr   string
	RedisDb     int
	RedisPool   int
	RedisEnable bool
	RedisPort   int

	ContainerHttps string
	ContainerHttp  string

	HttpEnable  bool
	HttpsEnable bool

	//consul
	HttpAddress  string
	HttpsAddress string

	ConsulToken   string
	ConsulAddress string

	Sdt   				string
	Sda 				string

	ServiceConf         string
	HealthCheckInterval string
	HealthCheckTimeout  string

	ContainerAddr string
	ContainerRpcx string

	Prefix         string
	ApiProxyPrefix string

	//rpcx
	RpcxOamBasepath    string
	RpcxSmsBasepath    string
	RpcxConfBasepath   string
	RpcxSysmgrBasepath string
	RpcxBasepath       string
	RpcxAddr           string
	RpcxNetwork        string

	KafkaAddress string
	ZipkinUrl    string

	//监控地址
	Debug          bool
	DebugAddr      string
	RecordAddr     string
	AppdashAddr    string
	LightstepToken string
}

var (
	Config *ImConf
)

var once sync.Once

func (config *ImConf) ReadConf() {

	//修改为viper配置
	config.GetViperConfiguration()
}

func (config *ImConf) GetViperConfiguration() {

	config.DbDriver = viper.GetString("self.system.dbdriver")
	if config.DbDriver == "" {
		config.DbDriver = "pgsql"
	}

	//oadb dsn
	config.DbPort = viper.GetInt(config.DbDriver + ".oadb.port")

	config.DbUser = viper.GetString(config.DbDriver + ".oadb.user")
	config.DbPasswd = viper.GetString(config.DbDriver + ".oadb.passwd")
	config.Database = viper.GetString(config.DbDriver + ".oadb.database")
	config.DbAddr = viper.GetString(config.DbDriver + ".oadb.host")
	config.ObjectsSchema = viper.GetString(config.DbDriver + ".oadb.schema")

	//oadb dsn
	port := ""
	port = strconv.Itoa(viper.GetInt(config.DbDriver + ".oadb.port"))
	if config.DbDriver == "mysql" {
		config.Dsn = config.DbUser + ":" + config.DbPasswd + "@tcp(" +
			config.DbAddr + ":" + port + ")/" + config.Database
	} else if config.DbDriver == "pgsql" {
		config.Dsn = fmt.Sprintf("postgres://%s:%s@%s/%s?connect_timeout=10&sslmode=disable",
			config.DbUser, config.DbPasswd,
			config.DbAddr+":"+port, config.Database)
	}

	//session
	config.SessionKey = viper.GetString("self.session.key")
	config.SessionStoreIP = viper.GetString("self.session.storeip")
	config.SessionStorePort = viper.GetString("self.session.storeport")
	config.SmsUrl = viper.GetString("self.system.sms_url")

	config.Logfile = viper.GetString("self.system.logfile")
	config.LogConfig = viper.GetString("self.system.logconfig")
	config.LogLevel = viper.GetString("self.system.loglevel")
	config.LogMaxSize = viper.GetInt("self.system.logmaxsize")
	if config.LogMaxSize == 0 {
		config.LogMaxSize = 104857600
	}
	config.LogRotateFiles = viper.GetInt("self.system.log_rotate_files")
	if config.LogRotateFiles == 0 {
		config.LogRotateFiles = 10
	}

	config.HealthCheckInterval = viper.GetString("consul.check.interval")
	config.HealthCheckTimeout = viper.GetString("consul.check.timeout")

	config.ServiceName = viper.GetString("self.system.service_name")
	if config.ServiceName == "" {
		config.ServiceName = "ms-9060-oasvc"
	}

	config.RedisPort = viper.GetInt("redis.queue.port")
	config.RedisAddr = viper.GetString("redis.queue.host") + ":" + strconv.Itoa(config.RedisPort)
	config.RedisPool = viper.GetInt("redis.queue.pool")
	config.RedisDb = viper.GetInt("redis.queue.database")
	config.RedisEnable = viper.GetBool("redis.queue.enable")

	config.Env = viper.GetString("self.system.env")
	if config.Env == "" {
		config.Env = "prod"
	}

	config.HttpEnable = viper.GetBool("self.http_enable")
	config.HttpsEnable = viper.GetBool("self.https_enable")
	config.Prefix = viper.GetString("self.prefix")

	if config.HttpEnable {
		config.HttpHost = viper.GetString("self.http.host")
		config.HttpPort = viper.GetInt("self.http.port")
		config.HttpAddress = config.HttpHost + ":" + utils.ConvertToString(config.HttpPort)
	}
	if config.HttpsEnable {
		config.HttpsHost = viper.GetString("self.https.host")
		config.HttpsPort = viper.GetInt("self.https.port")
		config.HttpsAddress = config.HttpsHost + ":" + utils.ConvertToString(config.HttpsPort)
	}

	config.RpcxOamBasepath = viper.GetString("rpcx.server.oam_base_path")
	config.RpcxSmsBasepath = viper.GetString("rpcx.server.sms_base_path")
	config.RpcxConfBasepath = viper.GetString("rpcx.server.conf_base_path")
	config.RpcxSysmgrBasepath = viper.GetString("rpcx.server.sysmgr_base_path")

	config.RpcxAddr = utils.Hostname2IPv4(viper.GetString("self.rpcx.address"))
	config.RpcxNetwork = viper.GetString("self.rpcx.network")
	config.RpcxBasepath = viper.GetString("self.rpcx.base_path")

	config.KafkaAddress = viper.GetString("self.zipkin.kafka_address")
	config.ZipkinUrl = viper.GetString("self.zipkin.zipkin_url")
}
