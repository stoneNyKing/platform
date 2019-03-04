package apis

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	l4g "github.com/libra9z/log4go"
	"platform/ousvc/common"
	"platform/ousvc/config"
)

const (
	PAGENUM_MAX = 1000
)

/*
	用户注册类型
*/
const (
	REGISTER_TYPE_NONE   = 0
	REGISTER_TYPE_PHONE  = 1
	REGISTER_TYPE_NAME   = 2
	REGISTER_TYPE_EMAIL  = 3
	REGISTER_TYPE_IDCARD = 4
	REGISTER_TYPE_IMEI   = 5
	REGISTER_TYPE_SSCARD = 6
)

var logger l4g.Logger

func InitLogger() {
	logger = common.Logger
}

func SetSearchPath(o orm.Ormer, schema string) (err error) {
	if o == nil {
		return errors.New("orm is nil")
	}

	if config.Config.DbDriver == "pgsql" {
		smt := fmt.Sprintf("set search_path to \"%s\";", schema)
		_, err = o.Raw(smt).Exec()
	} else if config.Config.DbDriver == "mysql" {
		return nil
	}

	return err
}

var SQL_USER = "select a.id,a.site_id as siteid,a.name,a.phone,a.email,a.idcard,a.passwd, " +
	"a.type,a.is_active as isactive,a.created,a.updated,a.last_login as lastlogin,a.last_logout as lastlogout, " +
	"a.json,a.rfid,a.weixinid,a.imei,a.nickname,a.image_url as imageurl " +
	"from %s.user a "

var SQL_COUNT_USER = "select count(*) as ucount " +
	"from %s.user a "
