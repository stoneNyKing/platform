package models


import (
	"github.com/astaxie/beego/orm"
	"platform/confsvc/imconf"
	"fmt"
	"errors"
	"platform/common/utils"
)


func init(){
	orm.RegisterModel(new(SysDomainConf))
}


func GetDomainConfLists(siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {
	return getDomainConfListCount(1,siteid,appid,stypes,contents,order,sort,num,start)
}
func GetDomainConfCount(siteid int64, appid int64, stypes, contents []string, order, sort string) (int64, error) {
	_,cnt,err := getDomainConfListCount(2,siteid,appid,stypes,contents,order,sort,1,0)

	return cnt,err
}

func getDomainConfListCount(cate int,siteid int64,appid int64,stypes []string, contents []string,order,sort string,num,start int64)(interface{},int64,error) {

	var vs []orm.Params
	
	var cnt int64
	var err error
	
	if sort == "" {
		sort = "asc"
	}
	
	if num <= 0 {
		num=PAGENUM_MAX
	}

	logger.Finest("siteid=%d,appid=%d,stypes=%v,contents=%v",siteid,appid,stypes,contents)
	
	if len(stypes) != len(contents) {
		return nil,0,errors.New("params number is not match.")
	} 
	
	conditions := ""

	l := len(stypes)
	var v string
	for i := 0;i < l;i++ {
		if contents[i]!= "" {
			v = "'%" + contents[i] + "%'"
		}
		
		logger.Finest("len=%d,stype[%d]=%s,scontent[%d]=%s",l,i,stypes[i],i,contents[i])
		
		switch stypes[i] {
		case "1":
			v = " a.domain='" + contents[i]+"'"
		case "2":
			v = " a.domaingrp=" + contents[i]
		case "3":	
			v = " a.keyid=" + contents[i]
		}
		
		if v!= "" {
			if i != l-1 {
				conditions = conditions + v +" and "
			}else{
				conditions = conditions + v 
			}
		}
	}
	
	logger.Finest("conditions = %s",conditions)
	
	if conditions != "" {
		conditions = " and " + conditions
	}	
	
	logger.Finest("siteid=%d,appid=%d,order=%s,sort=%s,num=%d,start=%d",siteid,appid,order,sort,num,start)

	var sqlstr,statement string

	if cate == 1 {
		sqlstr = SQL_DOMAIN_CONF
	} else {
		sqlstr = SQL_COUNT_DOMAIN_CONF
		num = 1
		start = 0
	}

	statement = fmt.Sprintf(sqlstr +
						"where (siteid=0 or siteid=?) %s order by a.domain_id %s limit ? offset ?",
						imconf.Config.SysconfdbSchema,conditions,sort)

	o := orm.NewOrm()

	err = SetSearchPath(o,imconf.Config.SysconfdbSchema)
	if err !=nil {
		logger.Error("不能设置search_path :%v",err)
		return nil,0,err
	}

	cnt,err = o.Raw(statement,siteid,num,start).Values(&vs)
			
	if err != nil {
		logger.Error("不能获取domain conf配置信息列表：%v",err.Error())
		return nil,0,err
	}

	if cate == 2 {
		if cnt > 0 {
			cnt = utils.Convert2Int64(vs[0]["ucount"])
		}
	}
	return vs,cnt,nil
}

func GetDomainConf(id int64,siteid int64,appid int64)(interface{},int64,error) {

	var vs []orm.Params
	
	var cnt int64
	var err error
	
	logger.Finest("siteid=%d,appid=%d,areaid=%d",siteid,appid,id)
	
	statement := fmt.Sprintf(SQL_DOMAIN_CONF +
					"where (siteid=0 or siteid=?) and a.domain_id=? ",
					imconf.Config.SysconfdbSchema)

	o := orm.NewOrm()

	err = SetSearchPath(o,imconf.Config.SysconfdbSchema)
	if err !=nil {
		logger.Error("不能设置search_path :%v",err)
		return nil,0,err
	}

	cnt,err = o.Raw(statement,siteid,id).Values(&vs)

	if err != nil {
		logger.Error("不能获取domain conf配置信息列表：%v",err.Error())
		return nil,0,err
	}

	return vs,cnt,nil
}



func PostDomainConf(appid,siteid int64,token string, param map[string]interface{} )( id int64, err error) {
	
	if param == nil {
		return 0,errors.New("no input")
	}
	
	var v SysDomainConf
	
	if param["domain"] != nil {
		v.Domain = param["domain"].(string)
	}
	if param["value"] != nil {
		v.Value = param["value"].(string)
	}

	if param["name"] != nil {
		v.Name = param["name"].(string)
	}

	if param["action"] != nil {
		v.Action = param["action"].(string)
	}

	if param["siteid"] != nil {
		v.Siteid = utils.Convert2Int64(param["siteid"])
	}else{
		v.Siteid = siteid
	}

	if param["domaingrp"] != nil {
		v.Domaingrp = utils.Convert2Int(param["domaingrp"])
	}
	if param["keyid"] != nil {
		v.Keyid = utils.Convert2Int64(param["keyid"])
	}

	o := orm.NewOrm()

	err = SetSearchPath(o,imconf.Config.SysconfdbSchema)
	if err !=nil {
		logger.Error("不能设置search_path :%v",err)
		return 0,err
	}

	ed,err := o.Insert(&v)
	if err.Error != nil {
		logger.Error("不能插入domain conf信息：%v",err.Error())
		return 0,err
	}
			
	return ed, nil
}

func PutDomainConf( id int64, param map[string]interface{} )( cnt int64,err error) {
	
	if param == nil {
		return 0,errors.New("no input")
	}

	var v SysDomainConf
	
	if param["id"]!= nil {
		v.DomainId = utils.Convert2Int64(param["id"])
	}
	
	if id != v.DomainId {
		return 0,errors.New("id is not match.")
	}

	o := orm.NewOrm()

	err = SetSearchPath(o,imconf.Config.SysconfdbSchema)
	if err !=nil {
		logger.Error("不能设置search_path :%v",err)
		return 0,err
	}

	err = o.Read(&v)
	if err != nil {
		logger.Error("cannot read record.")
		return 0, err
	}

	if param["domain"] != nil {
		v.Domain = param["domain"].(string)
	}
	if param["value"] != nil {
		v.Value = param["value"].(string)
	}
	if param["name"] != nil {
		v.Name = param["name"].(string)
	}
	if param["action"] != nil {
		v.Action = param["action"].(string)
	}

	if param["siteid"] != nil {
		v.Siteid = utils.Convert2Int64(param["siteid"])
	}
	if param["domaingrp"] != nil {
		v.Domaingrp = utils.Convert2Int(param["domaingrp"])
	}
	if param["keyid"] != nil {
		v.Keyid = utils.Convert2Int64(param["keyid"])
	}


	cnt, err = o.Update(&v)
	if err != nil {
		logger.Error("不能更新domain conf信息：%v", err.Error())
		return 0, err
	}
	
	return cnt,nil
}	

func DeleteDomainConf(id int64,param map[string]interface{}) (num int64, err error) {

	var v SysDomainConf

	var ids []interface{}
	
	if param["id"] != nil {
		ids = param["id"].([]interface{})
	}
	
	var cnt int64 = 0

	o := orm.NewOrm()
	err = SetSearchPath(o,imconf.Config.SysconfdbSchema)
	if err !=nil {
		logger.Error("不能设置search_path :%v",err)
		return 0,err
	}

	for  _,rid := range ids {
		v.DomainId = utils.Convert2Int64(rid)
		c,_ := o.Delete(&v)
		cnt += c
	}


	return cnt, err
}

