
package models

import (
	"errors"
	"platform/filesvc/imconf"
	"platform/common/utils"
	"platform/filesvc/dbmod"
	borm "github.com/astaxie/beego/orm"
	"fmt"
)

func init(){
	borm.RegisterModel(new(dbmod.FileConf))
}

func GetFileConfLists(siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {
	if num <= 0 {
		num=PAGENUM_MAX
	}

	return getFileConfListCount(1,siteid,appid,stypes,contents,order,sort,num,start)
}
func GetFileConfCount(siteid int64, appid int64, stypes, contents []string, order, sort string) (int64, error) {
	_,cnt,err := getFileConfListCount(2,siteid,appid,stypes,contents,order,sort,1,0)

	return cnt,err
}


func getFileConfListCount(cate int,siteid int64, appid int64, stypes, contents []string, order, sort string,num, start int64) (interface{}, int64, error) {


	var vs [] borm.Params

	var cnt int64
	var err error

	if sort == "" {
		sort = "desc"
	}

	logger.Finest("siteid=%d,appid=%d,order=%s,sort=%s", siteid, appid, order, sort)

	if len(stypes) != len(contents) {
		return nil,0, errors.New("params number is not match.")
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
			v = " a.status= " + contents[i]
		case "2":
			v = " a.create_time >  '" + contents[i] +"'"
		case "3":
			v = " a.create_time <=  '" + contents[i] +"'"
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


	o := borm.NewOrm()

	var sqlstr,statement string
	if cate == 1 {
		sqlstr = fmt.Sprintf(SQL_FileConf)
	}else if cate == 2 {
		sqlstr = fmt.Sprintf(SQL_COUNT_FileConf)
	}
	if siteid > 1 {
		statement = fmt.Sprintf(sqlstr +
			"where siteid=%d %s order by a.file_conf_id %s limit ? offset ? ", siteid,conditions, sort)
	}else{
		statement = fmt.Sprintf(sqlstr +
			"where 1=1 %s order by a.file_conf_id %s limit ? offset ? ", conditions, sort)
	}

	cnt, err = o.Raw(statement,num,start).Values(&vs)

	if err != nil {
		logger.Error("不能获取列表数量：%v", err.Error())
		return nil, 0, err
	}

	if cate == 2 {
		if cnt > 0 {
			cnt = utils.Convert2Int64(vs[0]["ucount"])
		}
	}


	return vs,cnt, nil
}


func GetFileConf(id int64, siteid int64, appid int64) (interface{}, int64, error) {

	var vs []borm.Params
	var err error
	var cnt int64

	logger.Finest("siteid=%d,appid=%d,id=%d", siteid, appid, id)

	statement := fmt.Sprintf(SQL_FileConf)

	o := borm.NewOrm()
	if siteid > 1 {
		cnt, err = o.Raw(statement +
			"where siteid=%d and a.file_conf_id =?", siteid, id).Values(&vs)
	}else{
		cnt, err = o.Raw(statement +
			"where a.file_conf_id =?", id).Values(&vs)
	}

	if err != nil {
		logger.Error("不能获取信息：%v", err.Error())
		return nil, 0, err
	}

	return vs, cnt, nil
}


func PostFileConf(appid,siteid int64,token string, param map[string]interface{}) (id int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v dbmod.FileConf


	if param["id"] != nil {
		v.FileConfId = utils.Convert2Int64(param["id"])
	}

	if param["siteid"] != nil {
		v.Siteid = utils.Convert2Int64(param["siteid"])
	}

	if param["status"] != nil {
		v.Status = utils.Convert2Int(param["status"])
	}

	if param["filecategory"] != nil {
		v.FileCategory = utils.Convert2Int(param["filecategory"])
	}

	if param["filetype"] != nil {
		v.FileType = utils.Convert2Int(param["filetype"])
	}
	if param["storagetype"] != nil {
		v.StorageType = utils.Convert2Int(param["storagetype"])
	}

	if param["createtime"] != nil {
		v.CreateTime = utils.ConvertToString(param["createtime"])
	}else{
		v.CreateTime = utils.GetTimeFormat("2006-01-02 15:04:05")
	}

	if param["remark"] != nil {
		v.Remark = utils.ConvertToString(param["remark"])
	}

	if param["filepath"] != nil {
		v.FilePath = utils.ConvertToString(param["filepath"])
	}
	if param["fileprefix"] != nil {
		v.FilePrefix = utils.ConvertToString(param["fileprefix"])
	}
	if param["template"] != nil {
		v.Template = utils.ConvertToString(param["template"])
	}

	o := borm.NewOrm()
	err = BSetSearchPath(o,imconf.Config.FiledbSchema)
	if err !=nil {
		logger.Error("不能设置search_path :%v",err)
		return 0,err
	}


	id,err = o.Insert(&v)

	if err != nil {
		logger.Error("不能插入信息：%v", err.Error())
		return 0, err
	}

	return id, nil
}


func PutFileConf(id int64,appid,siteid int64,token string, param map[string]interface{}) (cnt int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v dbmod.FileConf

	if param["id"] != nil {
		v.FileConfId = utils.Convert2Int64(param["id"])
	}

	if id != v.FileConfId {
		return 0, errors.New("id is not match.")
	}

	err = SetSearchPath(imconf.Config.FiledbSchema)
	if err !=nil {
		logger.Error("不能设置search_path :%v",err)
		return 0,err
	}

	_,err = orm.Get(&v)
	if err != nil {
		logger.Error("cannot read record.")
		return 0, err
	}


	if param["siteid"] != nil {
		v.Siteid = utils.Convert2Int64(param["siteid"])
	}

	if param["status"] != nil {
		v.Status = utils.Convert2Int(param["status"])
	}

	if param["filecategory"] != nil {
		v.FileCategory = utils.Convert2Int(param["filecategory"])
	}

	if param["filetype"] != nil {
		v.FileType = utils.Convert2Int(param["filetype"])
	}
	if param["storagetype"] != nil {
		v.StorageType = utils.Convert2Int(param["storagetype"])
	}

	if param["createtime"] != nil {
		v.CreateTime = utils.ConvertToString(param["createtime"])
	}

	if param["remark"] != nil {
		v.Remark = utils.ConvertToString(param["remark"])
	}

	if param["filepath"] != nil {
		v.FilePath = utils.ConvertToString(param["filepath"])
	}
	if param["fileprefix"] != nil {
		v.FilePrefix = utils.ConvertToString(param["fileprefix"])
	}
	if param["template"] != nil {
		v.Template = utils.ConvertToString(param["template"])
	}

	o := borm.NewOrm()
	err = BSetSearchPath(o,imconf.Config.FiledbSchema)
	if err !=nil {
		logger.Error("不能设置search_path :%v",err)
		return 0,err
	}


	cnt, err = o.Update(&v)
	if err != nil {
		logger.Error("不能更新记录：%v", err.Error())
		return 0, err
	}

	return cnt, nil
}


func DeleteFileConf(id int64,appid,siteid int64,token string, param map[string]interface{}) (num int64, err error) {

	o := borm.NewOrm()
	err = BSetSearchPath(o,imconf.Config.FiledbSchema)
	if err !=nil {
		logger.Error("不能设置search_path :%v",err)
		return 0,err
	}
	var v dbmod.FileConf

	var ids []interface{}

	if param["id"] != nil {
		ids = param["id"].([]interface{})
	}

	var cnt int64 = 0

	for _, rid := range ids {
		v.FileConfId = utils.Convert2Int64(rid)
		num, err = o.Delete(&v)
		cnt += num
	}

	return cnt, err
}

