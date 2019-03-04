package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"platform/confsvc/imconf"

)

func init() {
	orm.RegisterModel(new(SysSiteConf))
}

func GetSiteDNSConf(siteid int64, appid int64, stypes []string, contents []string, sitedns string) (interface{}, int64, error) {
	var maps []orm.Params

	var cnt int64
	var err error

	logger.Finest("siteid=%d,appid=%d,stypes=%v,contents=%v", siteid, appid, stypes, contents)

	if len(stypes) != len(contents) {
		return nil, 0, errors.New("params number is not match.")
	}

	conditions := ""

	l := len(stypes)
	var v string
	for i := 0; i < l; i++ {
		if contents[i] != "" {
			v = "'%" + contents[i] + "%'"
		}

		switch stypes[i] {
		case "1":
			v = " a.status=" + contents[i]
		}

		if v != "" {
			if i != l-1 {
				conditions = conditions + v + " and "
			} else {
				conditions = conditions + v
			}
		}
	}

	logger.Finest("conditions = %s", conditions)

	if conditions != "" {
		conditions = " and " + conditions
	}

	statement := fmt.Sprintf("select a.site_dns,a.skin_style,a.status,a.root_area_id as rootareaid "+
		"from %s.sys_site_conf a "+
		"where 1=1 %s and site_dns=?",
		imconf.Config.SysconfdbSchema, conditions)

	o := orm.NewOrm()

	err = SetSearchPath(o,imconf.Config.SysconfdbSchema)
	if err !=nil {
		logger.Error("不能设置search_path :%v",err)
		return nil,0,err
	}

	cnt,err = o.Raw(statement,siteid).Values(&maps)

	if err != nil {
		logger.Error("不能获取domain conf配置信息列表：%v",err.Error())
		return nil,0,err
	}


	return maps, cnt, nil
}
