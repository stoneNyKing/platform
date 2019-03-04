package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	l4g "github.com/libra9z/log4go"
	"platform/confsvc/common"
	"platform/confsvc/imconf"
)

var logger l4g.Logger

func InitLogger() {
	logger = common.Logger
}

const(
	
	PAGENUM_MAX = 1000
)



var	SQL_DOMAIN_CONF = "select a.domain_id as id,a.domain,a.domaingrp,a.name,a.keyid,a.value,a.siteid,a.action " +
	"from %s.sys_domain_conf a "
var	SQL_COUNT_DOMAIN_CONF = "select count(*) as ucount " +
	"from %s.sys_domain_conf a "

var	SQL_APPID = "select a.appid as id,a.appid,a.appkey,a.remark,a.json,a.status,a.created,a.updated,a.befe_flag as befeflag " +
	"from %s.appid a "
var	SQL_COUNT_APPID = "select count(*) as ucount " +
	"from %s.appid a "




func SetSearchPath( o orm.Ormer, schema string ) (err error) {
	if o == nil {
		return errors.New("orm is nil")
	}

	if imconf.Config.DbDriver == "pgsql" {
		smt := fmt.Sprintf("set search_path to \"%s\";",schema)
		_,err  = o.Raw(smt).Exec()
	}else if imconf.Config.DbDriver == "mysql" {
		return nil
	}

	return err
}
