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

	Env      string // run mode, "dev" or "prod"
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
	LogMaxSize     int
	LogRotateFiles int

	//Redis
	RedisAddr   string
	RedisDb     int
	RedisPool   int
	RedisEnable bool
	RedisPort   int

	//consul
	HttpAddress  string
	HttpsAddress string

	//oa ou URL
	OasvcUrl string
	OusvcUrl string

	Prefix string

	//rest api
	HttpEnable  bool
	HttpsEnable bool

	HttpHost  string
	HttpPort  int
	HttpsHost string
	HttpsPort int

	ContainerHttps string
	ContainerHttp  string
	ContainerRpcx  string

	//consul
	ConsulToken   string
	ConsulAddress string
	Sdt           string
	Sda           string

	ServiceConf         string
	HealthCheckInterval string
	HealthCheckTimeout  string

	ApiProxyPrefix string

	//weixin相关
	WxAppId          string
	WxOriId          string
	WxToken          string
	WxAppSecret      string
	WxScope          string
	WxRedirectUrl    string
	WxAesKeyEncode   string
	WxMyAppid        int
	WxMainPage       string
	WxMainHost       string
	WxMainPageNoAuth string
	WxServiceUrl     string

	//rpcx
	RpcxOamBasepath  string
	RpcxSmsBasepath  string
	RpcxConfBasepath string
	RpcxBasepath     string
	RpcxAddr         string
	RpcxNetwork      string

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

	//oudb dsn
	config.DbPort = viper.GetInt(config.DbDriver + ".oudb.port")

	config.DbUser = viper.GetString(config.DbDriver + ".oudb.user")
	config.DbPasswd = viper.GetString(config.DbDriver + ".oudb.passwd")
	config.Database = viper.GetString(config.DbDriver + ".oudb.database")
	config.DbAddr = viper.GetString(config.DbDriver + ".oudb.host")
	config.ObjectsSchema = viper.GetString(config.DbDriver + ".oudb.schema")

	//oadb dsn
	port := ""
	port = strconv.Itoa(viper.GetInt(config.DbDriver + ".oudb.port"))
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
		config.ServiceName = "ms-9058-ousvc"
	}

	//微信相关
	config.WxAppId = viper.GetString("weixin.appid")
	config.WxOriId = viper.GetString("weixin.oriid")
	config.WxToken = viper.GetString("weixin.token")
	config.WxAppSecret = viper.GetString("weixin.app_secret")
	config.WxRedirectUrl = viper.GetString("weixin.redirect_url")
	config.WxScope = viper.GetString("weixin.scope")
	config.WxAesKeyEncode = viper.GetString("weixin.aes_key_encode")
	config.WxMyAppid = viper.GetInt("weixin.myappid")
	config.WxMainHost = viper.GetString("weixin.main_host")
	config.WxMainPage = viper.GetString("weixin.main_page")
	config.WxMainPageNoAuth = viper.GetString("weixin.main_page_noauth")
	config.WxServiceUrl = viper.GetString("weixin.service_url")

	config.RedisPort = viper.GetInt("redis.queue.port")
	config.RedisAddr = viper.GetString("redis.queue.host") + ":" + strconv.Itoa(config.RedisPort)
	config.RedisPool = viper.GetInt("redis.queue.pool")
	config.RedisDb = viper.GetInt("redis.queue.database")
	config.RedisEnable = viper.GetBool("redis.queue.enable")

	config.Env = viper.GetString("self.system.env")
	if config.Env == "" {
		config.Env = "prod"
	}

	config.RpcxOamBasepath = viper.GetString("rpcx.server.oam_base_path")
	config.RpcxSmsBasepath = viper.GetString("rpcx.server.sms_base_path")
	config.RpcxConfBasepath = viper.GetString("rpcx.server.conf_base_path")

	config.RpcxAddr = utils.Hostname2IPv4(viper.GetString("self.rpcx.address"))
	config.RpcxNetwork = viper.GetString("self.rpcx.network")
	config.RpcxBasepath = viper.GetString("self.rpcx.base_path")

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

	config.KafkaAddress = viper.GetString("self.zipkin.kafka_address")
	config.ZipkinUrl = viper.GetString("self.zipkin.zipkin_url")
}
