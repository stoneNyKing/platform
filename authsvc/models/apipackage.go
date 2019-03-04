package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"platform/authsvc/imconf"
	"platform/common/utils"
	"time"
)

func init() {
	orm.RegisterModel(new (ApiPackage))
}

func GetApiPackage (id int64) (interface{}, int64, error){

	var vs []orm.Params
	o := orm.NewOrm()

	statement := fmt.Sprintf(SQL_API_PACKAGE,imconf.Config.AuthdbSchema)
	cnt, err := o.Raw(statement + "where a.package_id=?", id).Values(&vs)

	if err != nil {
		logger.Error("不能获取资源信息：%v", err.Error())
		return nil, 0, err
	}

	return vs, cnt, nil
}


func GetApiPackageLists(siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {
	return getApiPackageListCount(1,siteid,appid,stypes,contents,order,sort,num,start)
}
func GetApiPackageCount(siteid int64, appid int64, stypes, contents []string, order, sort string) (int64, error) {
	_,cnt,err := getApiPackageListCount(2,siteid,appid,stypes,contents,order,sort,1,0)

	return cnt,err
}
func getApiPackageListCount (cate int,siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error){

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
			v = " a.name like " + v
		case "2":
			v = " a.price = " +  contents[i]
		case "3":
			v = " a.charge_model = " +  contents[i]
		case "4":
			v = " a.sub_sys_id = " +  contents[i]
		case "5":
			v = " a.status = " +  contents[i]
		case "6":
			v = " a.create_time >= '" + contents[i] + "'"
		case "7":
			v = " a.create_time < '" + contents[i] + "'"
		case "8":
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
		statement = fmt.Sprintf(SQL_API_PACKAGE,imconf.Config.AuthdbSchema)
	} else if cate == 2 {
		statement = fmt.Sprintf(SQL_COUNT_API_PACKAGE,imconf.Config.AuthdbSchema)
	}

	statement = fmt.Sprintf(statement +
		"where 1=1 %s order by a.package_id %s limit ? offset ?", conditions, sort)

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


func PostApiPackage (param map[string]interface{}) (id int64, err error){
	if param == nil {
		return 1, errors.New("no input")
	}

	var v ApiPackage

	if param["id"] != nil {
		v.PackageId = utils.Convert2Int64(param["id"])
	}

	if param["name"] != nil {
		v.Name =  param["name"].(string)
	}
	if param["packagecode"] != nil {
		v.PackageCode =  param["packagecode"].(string)
	}

	if param["price"] != nil {
		v.Price = utils.Convert2Int(param["price"])
	}

	if param["chargemodel"] != nil {
		v.ChargeModel = utils.Convert2Int(param["chargemodel"])
	}

	if param["subsysid"] != nil {
		v.SubSysId = utils.Convert2Int64(param["subsysid"])
	}

	if param["status"] != nil {
		v.Status = utils.Convert2Int16(param["status"])
	}

	if param["createtime"] != nil {
		v.CreateTime,_ = time.Parse("2006-01-02 15:04:05",param["createtime"].(string))
	}else{
		v.CreateTime = time.Now()
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

func PutApiPackage (id int64, param map[string]interface{}) (cnt int64, err error){
	if param == nil {
		return 0, errors.New("no input")
	}

	var v ApiPackage

	if param["id"] != nil {
		v.PackageId = utils.Convert2Int64(param["id"])
	}

	if id != v.PackageId{
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

	if param["name"] != nil {
		v.Name =  param["name"].(string)
	}
	if param["packagecode"] != nil {
		v.PackageCode =  param["packagecode"].(string)
	}

	if param["price"] != nil {
		v.Price = utils.Convert2Int(param["price"])
	}

	if param["chargemodel"] != nil {
		v.ChargeModel = utils.Convert2Int(param["chargemodel"])
	}

	if param["subsysid"] != nil {
		v.SubSysId = utils.Convert2Int64(param["subsysid"])
	}

	if param["status"] != nil {
		v.Status = utils.Convert2Int16(param["status"])
	}

	if param["createtime"] != nil {
		v.CreateTime,_ = time.Parse("2006-01-02 15:04:05",param["createtime"].(string))
	}

	if param["remark"] != nil {
		v.Remark = param["remark"].(string)
	}

	cnt, err = o.Update(&v)

	if err != nil{
		logger.Error("不能更新资源信息：%v", err.Error())
	}

	return cnt, err
}

func DeleteApiPackage (param map[string]interface{}) (num int64, err error){
	if param == nil {
		return 0, errors.New("no input")
	}

	var v ApiPackage

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
		v.PackageId = utils.Convert2Int64(rid)

		num, err = o.Delete(&v)
		cnt += num
	}

	return cnt, err
}
