package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"platform/common/utils"
	"fmt"
	"platform/authsvc/imconf"
	"time"
)

func init() {
	orm.RegisterModel(new (ApiPackageService))
}

func GetApiPkgService (id int) (interface{}, int64, error){

	var vs []orm.Params
	o := orm.NewOrm()

	statement := fmt.Sprintf(SQL_API_PKG_SERVICE,imconf.Config.AuthdbSchema,imconf.Config.AuthdbSchema)
	cnt, err := o.Raw(statement + "where a.pkg_service_id=?", id).Values(&vs)

	if err != nil {
		logger.Error("不能获取资源信息：%v", err.Error())
		return nil, 0, err
	}

	return vs, cnt, nil
}


func GetApiPkgServiceLists(siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {
	return getApiPkgServiceListCount(1,siteid,appid,stypes,contents,order,sort,num,start)
}
func GetApiPkgServiceCount(siteid int64, appid int64, stypes, contents []string, order, sort string) (int64, error) {
	_,cnt,err := getApiPkgServiceListCount(2,siteid,appid,stypes,contents,order,sort,1,0)

	return cnt,err
}
func getApiPkgServiceListCount (cate int,siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error){

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
			v = " a.service_id = " +  contents[i]
		case "2":
			v = " a.package_id = " +  contents[i]
		case "3":
			v = " a.usage_counts = " +  contents[i]
		case "4":
			v = " a.status = " +  contents[i]
		case "5":
			v = " a.create_time >= '" + contents[i] + "'"
		case "6":
			v = " a.create_time  < '" + contents[i] + "'"
		case "7":
			v = " b.route= '"+ contents[i] + "'"
		case "8":
			v = " b.svc_id = "+ contents[i]
		case "9":
			v = " b.svc_code = '"+ contents[i] + "'"
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
		statement = fmt.Sprintf(SQL_API_PKG_SERVICE,imconf.Config.AuthdbSchema,imconf.Config.AuthdbSchema)
	} else if cate == 2 {
		statement = fmt.Sprintf(SQL_COUNT_API_PKG_SERVICE,imconf.Config.AuthdbSchema,imconf.Config.AuthdbSchema)
	}

	statement = fmt.Sprintf(statement +
		"where 1=1 %s order by a.pkg_service_id %s limit ? offset ?", conditions, sort)

	cnt, err = o.Raw(statement, num,start).Values(&vs)

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

func PostApiPkgService (param map[string]interface{}) (id int64, err error){
	if param == nil {
		return 1, errors.New("no input")
	}

	var v ApiPackageService

	if param["id"] != nil {
		v.PkgServiceId = utils.Convert2Int64(param["id"])
	}

	if param["serviceid"] != nil {
		v.ServiceId =  utils.Convert2Int(param["serviceid"])
	}

	if param["packageid"] != nil {
		v.PackageId = utils.Convert2Int(param["packageid"])
	}

	if param["totalcounts"] != nil {
		v.TotalCounts = utils.Convert2Int(param["totalcounts"])
	}
	if param["dailycounts"] != nil {
		v.DailyCounts = utils.Convert2Int(param["dailycounts"])
	}

	if param["status"] != nil {
		v.Status = utils.Convert2Int16(param["status"])
	}

	if param["createtime"] != nil {
		v.CreateTime,_ = time.Parse("2006-01-02 15:04:05",param["createtime"].(string))
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

func PutApiPkgService (id int64, param map[string]interface{}) (cnt int64, err error){
	if param == nil {
		return 0, errors.New("no input")
	}

	var v ApiPackageService

	if param["id"] != nil {
		v.PkgServiceId = utils.Convert2Int64(param["id"])
	}

	if id != v.PkgServiceId{
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

	if param["serviceid"] != nil {
		v.ServiceId =  utils.Convert2Int(param["serviceid"])
	}

	if param["packageid"] != nil {
		v.PackageId = utils.Convert2Int(param["packageid"])
	}

	if param["totalcounts"] != nil {
		v.TotalCounts = utils.Convert2Int(param["totalcounts"])
	}
	if param["dailycounts"] != nil {
		v.DailyCounts = utils.Convert2Int(param["dailycounts"])
	}

	if param["status"] != nil {
		v.Status = utils.Convert2Int16(param["status"])
	}

	if param["createtime"] != nil {
		v.CreateTime,_ = time.Parse("2006-01-02 15:04:05",param["createtime"].(string))
	}

	cnt, err = o.Update(&v)

	if err != nil{
		logger.Error("不能更新资源信息：%v", err.Error())
	}

	return cnt, err
}

func DeleteApiPkgService (param map[string]interface{}) (num int64, err error){
	if param == nil {
		return 0, errors.New("no input")
	}

	var v ApiPackageService

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
		v.PkgServiceId = utils.Convert2Int64(rid)

		num, err = o.Delete(&v)
		cnt += num
	}

	return cnt, err
}
