package models

import (
	"platform/authsvc/common"
	l4g "github.com/libra9z/log4go"

	"github.com/astaxie/beego/orm"
	"fmt"
	"errors"
	"platform/authsvc/imconf"
)

var logger l4g.Logger

func InitLogger() {
	logger = common.Logger
}


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


const(
	PAGENUM_MAX=1000
)

var	SQL_APP_LICENSE =  "select a.license_id as id, a.siteid as siteid,a.api_key as apikey,a.org_code as orgcode,a.status,a.userid,a.organization_id as organizationid "+
		"a.create_time as createtime,a.remark,a.license "+
		"from %s.sec_site_info a "
var	SQL_COUNT_APP_LICENSE =  "select count(*) as ucount "+
		"from %s.sec_site_info a "

var	SQL_API_LICENSE_COUNTS = "select a.pkg_service_id as id,a.service_id as serviceid,a.package_id as packageid,a.daily_counts as dailycounts,a.total_counts as totalcounts,"+
		"b.svc_code as svccode,b.svc_id as svcid,b.route as path "+
		"from %s.api_package_service a "+
		"left join %s.api_service b on a.service_id=b.service_id "


var	SQL_API_PACKAGE =  "select a.package_id as id,a.name,a.price,a.charge_model as chargemodel,a.sub_sys_id as subsysid,a.status,a.create_time as createtime,a.remark,a.package_code as packagecode "+
		"from %s.api_package a "

var	SQL_COUNT_API_PACKAGE =  "select count(*) as ucount "+
		"from %s.api_package a "

var	SQL_API_PKG_SERVICE = "select a.pkg_service_id as id,a.service_id as serviceid,a.package_id as packageid,a.total_counts as totalcounts,a.daily_counts as dailycounts,a.status,a.create_time as createtime,b.svc_code as svccode,b.svc_id as svcid,b.route as path "+
		"from %s.api_package_service a "+
		"left join %s.api_service b on a.service_id=b.service_id "
var	SQL_COUNT_API_PKG_SERVICE = "select count(*) as ucount "+
		"from %s.api_package_service a "+
		"left join %s.api_service b on a.service_id=b.service_id "

var	SQL_API_SERVICE =  "select a.service_id as id,a.svc_code as svccode,a.svc_id as svcid,a.route,a.web_url as weburl,a.api_ver as apiver,a.status,a.remark "+
		"from %s.api_service a "

var	SQL_COUNT_API_SERVICE =  "select count(*) as ucount "+
		"from %s.api_service a "

