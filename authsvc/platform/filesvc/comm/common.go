package comm

import (
		"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	//"flag"

	l4g "github.com/libra9z/log4go"
	"github.com/spf13/viper"

	"github.com/go-xorm/xorm"
	"platform/filesvc/imconf"
	"fmt"
	"strconv"
)

var once sync.Once

const (
	STORAGE_TYPE_FILESYSTEM = 0 //"fs"
	STORAGE_TYPE_DATABASE   = 1 //"db"
	STORAGE_TYPE_ALIOSS     = 2 //"oss"
)

const (
	FILE_TYPE_REGULAR = 0 //"regular"
	FILE_TYPE_IMAGE   = 1 //"image"
	FILE_TYPE_AUDIO   = 2 //"audio"
	FILE_TYPE_VIDEO   = 3 //"video"
)

var Logger l4g.Logger


type Size struct {
	Width  int
	Height int
}

var apptypes map[string]string
var atlock sync.Mutex

var (

	//consul
	Address       string
	ConsulToken   string
	ConsulAddress string
	Prefix 			string
	FilePath		string
)

func InitApptype() {

	apptypes = make(map[string]string)

	maps := viper.GetStringMapString("apps")
	for k, v := range maps {

		s := strings.Split(k, "_")
		if len(s) >= 2 {
			RegisterAppType(s[1], v)
		}
	}
}


func GetFileType(fn string) int {
	var ftype int

	ext := filepath.Ext(fn)

	if ext == "" {
		ftype = FILE_TYPE_REGULAR
	}

	ext = strings.ToLower(ext)

	switch ext {
	case ".mp3", ".mp2", ".ogg", ".amr", ".aac", ".flac", ".wav", ".mid":
		ftype = FILE_TYPE_AUDIO
	case ".jpg", ".jpeg", ".png", ".tif", ".tiff", ".bmp", ".gif":
		ftype = FILE_TYPE_IMAGE
	default:
		ftype = FILE_TYPE_REGULAR
	}
	return ftype
}

func GetKeyString(key string) string {
	return viper.GetString(key)
}

func InitLogger(filename string, level string) {
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

	Logger.AddFilter("stdout", lvl, l4g.NewConsoleLogWriter())

	if _, err := os.Stat(filename); err == nil {
		os.Remove(filename)
	}

	flw := l4g.NewFileLogWriter(filename, true)
	flw.SetRotateSize(imconf.Config.LogMaxSize)
	flw.SetRotateFiles(imconf.Config.LogRotateFiles)

	Logger.AddFilter("logfile", lvl, flw)
	Logger.Info("Current time is : %s\n", time.Now().Format("15:04:05 MST 2006/01/02"))

	return
}

func RegisterAppType(appid string, apptype string) {
	atlock.Lock()
	defer atlock.Unlock()

	apptypes[appid] = apptype
}
func GetAppType(appid string) string {
	atlock.Lock()
	defer atlock.Unlock()

	v, _ := apptypes[appid]

	return v
}

func SetEngine(driver,host string, port int, user string, passwd string, database,schema string) (orm *xorm.Engine,err error) {

	if driver == "mysql" {
		orm, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
			user, passwd, host+":"+strconv.Itoa(port), database))
		if err != nil {
			Logger.Error("models.init(fail to conntect database): %v", err)
			return nil,err
		}

	} else if driver == "pgsql" {
		orm, err = xorm.NewEngine("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?connect_timeout=10&sslmode=disable",
			user, passwd, host+":"+strconv.Itoa(port), database))
		if err != nil {
			Logger.Error("models.init(fail to conntect database): ) %v", err)
			return nil,err
		}
		smt := fmt.Sprintf("SET SEARCH_PATH TO \"%s\"", schema)
		orm.Exec(smt)
	}

	return orm,nil
}
