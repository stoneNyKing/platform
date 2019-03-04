package models

import (
	l4g "github.com/libra9z/log4go"
	"platform/filesvc/comm"

	"fmt"
	"github.com/go-xorm/xorm"
	"platform/filesvc/imconf"
	borm "github.com/astaxie/beego/orm"
	"errors"
)

var logger l4g.Logger


const(
	PAGENUM_MAX		= 1000
	FILE_CATEGORY_ALIOSS  = 2
)


var (
	orm    *xorm.Engine
	tables []interface{}
)

const (
	SQL_FileConf = "select a.file_conf_id as id,a.siteid,a.file_category as filecategory,a.file_type as filetype,a.storage_type as storagetype,a.status,a.create_time as createtime,a.remark,a.file_path as filepath,a.file_prefix as fileprefix,a.template "+
		"from file_conf a "
	SQL_COUNT_FileConf = "select count(*) as ucount "+
		"from file_conf a "
	SQL_FileInfo = "select a.fileid as id,a.siteid,a.filekey,a.name,a.orig_url as origurl,a.local_file as localfile,a.storage_type as storagetype,a.status,a.create_time as createtime,a.remark,a.hash,a.redirect,a.location,a.file_no as fileno,a.file_owner as fileowner,a.file_size as filesize "+
		"from file_info a "
	SQL_COUNT_FileInfo = "select count(*) as ucount "+
		"from file_info a "
)


func InitLogger() {
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

func Truncate() {
	orm.Exec("truncate table file_conf")
	orm.Exec("truncate table file_info")
	orm.Exec("truncate table file_storage")
}

func SetSearchPath( schema string ) (err error) {

	if imconf.Config.DbDriver == "pgsql" {
		smt := fmt.Sprintf("set search_path to \"%s\";",schema)
		_,err  = orm.Exec(smt)
	}else if imconf.Config.DbDriver == "mysql" {
		return nil
	}

	return err
}

func BSetSearchPath( o borm.Ormer, schema string ) (err error) {

	if o == nil {
		return errors.New("orm is nil")
	}

	if imconf.Config.DbDriver == "pgsql" {
		smt := fmt.Sprintf("set search_path to \"%s\";", schema)
		_, err = o.Raw(smt).Exec()
	} else if imconf.Config.DbDriver == "mysql" {
		return nil
	}


	return err
}

