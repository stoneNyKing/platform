package dbmod

import (
	"fmt"
	"github.com/go-xorm/xorm"
	l4g "github.com/libra9z/log4go"
	"github.com/spf13/viper"
	"platform/common/utils"
	"platform/filesvc/comm"
	"platform/filesvc/imconf"
)

type Options struct {
	Category 		int
}

var logger l4g.Logger

var Option *Options

var (
	orm    *xorm.Engine
	tables []interface{}
)


func init() {
	Option = &Options{}
	tables = append(tables, new(FileConf), new(FileInfo), new(FileStorage))
}

func InitDbLogger() {
	logger = comm.Logger
}

func InitDatabase() {
	NewEngine(imconf.Config.DbDriver,imconf.Config.DbAddr, imconf.Config.DbPort, imconf.Config.DbUser, imconf.Config.DbPasswd, imconf.Config.FiledbName, imconf.Config.FiledbSchema)
}


func NewEngine(driver,host string, port int, user string, passwd string, database,schema string) (err error) {
	if orm,err = comm.SetEngine(driver,host, port, user, passwd, database,schema); err != nil {
		return err
	}
	if err = orm.Sync(tables...); err != nil {
		return fmt.Errorf("sync database struct error: %v\n", err)
	}
	return nil
}

type URLTemplater interface {
	GetURLTemplate(*Options) string
}


func (c *Options) GetSitePath(site string,cate, stype int) string {

	fc := &FileConf{Siteid:utils.Convert2Int64(site),FileType:stype,FileCategory:cate}
	_,err :=orm.Get(fc)

	if err != nil {
		logger.Error("不能获取文件配置信息：site=%s,ftype=%d,err=%v",site,stype,err)
		return ""
	}

	return fc.FilePath
}

func (c *Options) GetSitePrefix(site string,cate, stype int) string {
	fc := &FileConf{Siteid:utils.Convert2Int64(site),FileType:stype,FileCategory:cate}
	_,err :=orm.Get(fc)

	if err != nil {
		logger.Error("不能获取文件配置信息：site=%s,ftype=%d,err=%v",site,stype,err)
		return ""
	}

	return fc.FilePrefix
}

func (c *Options) GetSiteStorage(site string,cate, stype int) (bool,int) {
	fc := &FileConf{Siteid:utils.Convert2Int64(site),FileType:stype,FileCategory:cate}
	b,err :=orm.Get(fc)

	if err != nil {
		logger.Error("不能获取文件配置信息：site=%s,ftype=%d,err=%v",site,stype,err)
		return false,0
	}

	return b,fc.StorageType
}

func (c *Options) GetSiteURLTemplate(site string, cate, stype int) string {
	fc := &FileConf{Siteid:utils.Convert2Int64(site),FileType:stype,FileCategory:cate}
	_,err :=orm.Get(fc)

	if err != nil {
		logger.Error("不能获取文件配置信息：site=%s,ftype=%d,err=%v",site,stype,err)
		return ""
	}

	return fc.Template
}

func (c *Options) GetServiceId() string {
	return "107"
}

func (c *Options) GetString(key string) string {
	return viper.GetString(key)
}

func (c *Options) GetCategory() int {
	return c.Category
}
func (c *Options) SetCategory(cate int) {
	c.Category = cate
}

