package admins

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"platform/oasvc/config"
)

const (
	PAGENUM_MAX = 1000

	JOBNUMBER_PREFIX = "jobnumber_site_"
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

var SQL_ADMIN = "select a.id, a.site_id as siteid,a.role_id as roleid,a.name,a.phone,a.email,a.created,a.job_number as jobnumber,a.organization_id as organizationid, " +
	"a.updated,a.description,a.effective_time as effectivetime,a.expire_time as expiretime,a.passwd, a.type, a.image_url as imageurl,a.gender,a.realname, " +
	"b.name as rolename,b.startpage,c.resource_id as rootresid, " +
	"d.name as orgname " +
	"from %s.admin a " +
	"left join %s.role b on a.role_id=b.id " +
	"left join %s.role_resource c on c.role_id = b.id " +
	"left join %s.site d on a.site_id = d.id "

var SQL_COUNT_ADMIN = "select count(*) as ucount " +
	"from %s.admin a " +
	"left join %s.role b on a.role_id=b.id " +
	"left join %s.role_resource c on c.role_id = b.id " +
	"left join %s.site d on a.site_id = d.id "
