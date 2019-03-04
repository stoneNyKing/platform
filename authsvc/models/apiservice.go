package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"platform/common/utils"
	"fmt"
	"platform/authsvc/imconf"
)

func init() {
	orm.RegisterModel(new (ApiService))
}

func GetApiService (id int) (interface{}, int64, error){

	var vs []orm.Params
	o := orm.NewOrm()

	statement := fmt.Sprintf(SQL_API_SERVICE,imconf.Config.AuthdbSchema)
	cnt, err := o.Raw(statement + "where a.service_id=?", id).Values(&vs)

	if err != nil {
		logger.Error("不能获取资源信息：%v", err.Error())
		return nil, 0, err
	}

	return vs, cnt, nil
}

func GetApiServiceLists(siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {
	return getApiServiceListCount(1,siteid,appid,stypes,contents,order,sort,num,start)
}
func GetApiServiceCount(siteid int64, appid int64, stypes, contents []string, order, sort string) (int64, error) {
	_,cnt,err := getApiServiceListCount(2,siteid,appid,stypes,contents,order,sort,1,0)

	return cnt,err
}

func getApiServiceListCount (cate int,siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error){

	var vs []orm.Params

	o := orm.NewOrm()

	var cnt int64
	var err error

	if sort == "" {
		sort = "desc"
	}

	if num <= 0 {
		num = PAGENUM_MAX
	}

	logger.Finest("siteid=%d,appid=%d,order=%s,sort=%s,num=%d,start=%d", siteid, appid, order, sort, num, start)

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
			v = " a.svc_code like " + v
		case "2":
			v = " a.svc_id = " +  contents[i]
		case "3":
			v = " a.route like " + v
		case "4":
			v = " a.web_url like " + v
		case "5":
			v = " a.api_ver " + v
		case "6":
			v = " a.status  = " +  contents[i]
		case "7":
			v = " a.remark like " + v
		}

		if v != "" {
			if i != l-1 {
				conditions = conditions + v + " and "
			} else {
				conditions = conditions + v
			}
		}
	}

	if conditions != "" {
		conditions = " and " + conditions
	}

	logger.Finest("conditions = %s", conditions)
	var statement string

	if cate == 1 {
		statement =fmt.Sprintf(SQL_API_SERVICE,imconf.Config.AuthdbSchema)
	} else if cate == 2 {
		statement =fmt.Sprintf(SQL_COUNT_API_SERVICE,imconf.Config.AuthdbSchema)
	}

	statement = fmt.Sprintf(statement +
		"where 1=1 %s order by a.service_id %s limit ? offset ?", conditions, sort)

	cnt, err = o.Raw(statement, num, start).Values(&vs)

	if err != nil {
		logger.Error("不能获取资源信息列表：%v", err.Error())
		return nil, 0, err
	}
	if cate == 2 {
		if cnt > 0 {
			cnt = utils.Convert2Int64(vs[0]["ucount"])
		}
	}

	return vs, cnt, nil
}



func PostApiService (param map[string]interface{}) (id int64, err error){
	if param == nil {
		return 1, errors.New("no input")
	}

	var v ApiService

	if param["id"] != nil {
		v.ServiceId = utils.Convert2Int64(param["id"])
	}

	if param["svccode"] != nil {
		v.SvcCode =  param["svccode"].(string)
	}

	if param["svcid"] != nil {
		v.SvcId = utils.Convert2Int(param["svcid"])
	}

	if param["route"] != nil {
		v.Route = param["route"].(string)
	}

	if param["weburl"] != nil {
		v.WebUrl = param["weburl"].(string)
	}

	if param["apiver"] != nil {
		v.ApiVer = param["apiver"].(string)
	}

	if param["status"] != nil {
		v.Status = utils.Convert2Int16(param["status"])
	}

	if param["remark"] != nil {
		v.Remark = param["remark"].(string)
	}

	o := orm.NewOrm()
	err = SetSearchPath(o,imconf.Config.AuthdbSchema)
	if err !=nil {
		logger.Error("不能设置search_path :%v",err)
		return 0,err
	}


	id, err = o.Insert(&v)

	if err != nil {
		//o.Rollback()
		logger.Error("不能插入资源信息：%v", err.Error())
	}

	return id, err
}

func PutApiService (id int64, param map[string]interface{}) (cnt int64, err error){
	if param == nil {
		return 0, errors.New("no input")
	}

	var v ApiService

	if param["id"] != nil {
		v.ServiceId = utils.Convert2Int64(param["id"])
	}

	if id != v.ServiceId{
		return 0, errors.New("id is not match.")
	}

	o := orm.NewOrm()
	err = SetSearchPath(o,imconf.Config.AuthdbSchema)
	if err !=nil {
		logger.Error("不能设置search_path :%v",err)
		return 0,err
	}

	err = o.Read(&v)
	if err != nil {
		logger.Error("cannot read record.")
		return 0, err
	}

	if param["svccode"] != nil {
		v.SvcCode =  param["svccode"].(string)
	}

	if param["svcid"] != nil {
		v.SvcId = utils.Convert2Int(param["svcid"])
	}

	if param["route"] != nil {
		v.Route = param["route"].(string)
	}

	if param["weburl"] != nil {
		v.WebUrl = param["weburl"].(string)
	}

	if param["apiver"] != nil {
		v.ApiVer = param["apiver"].(string)
	}

	if param["status"] != nil {
		v.Status = utils.Convert2Int16(param["status"])
	}

	if param["remark"] != nil {
		v.Remark = param["remark"].(string)
	}

	//o := orm.NewOrm()
	cnt, err = o.Update(&v)

	if err != nil{
		logger.Error("不能更新资源信息：%v", err.Error())
	}

	return cnt, err
}

func DeleteApiService (param map[string]interface{}) (num int64, err error){
	if param == nil {
		return 0, errors.New("no input")
	}

	var v ApiService

	var ids []interface{}

	if param["id"] != nil {
		ids = param["id"].([]interface{})
	}

	o := orm.NewOrm()
	err = SetSearchPath(o,imconf.Config.AuthdbSchema)
	if err !=nil {
		logger.Error("不能设置search_path :%v",err)
		return 0,err
	}

	var cnt int64 = 0
	for _, rid := range ids {
		v.ServiceId = utils.Convert2Int64(rid)

		num, err = o.Delete(&v)
		cnt += num
	}

	return cnt, err
}
