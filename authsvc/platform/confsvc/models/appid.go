package models


import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"platform/common/utils"
	"platform/confsvc/imconf"
	"time"
)


func init(){
	orm.RegisterModel(new(Appid))
}


func GetAppidLists(siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {
	return getAppidListCount(1,siteid,appid,stypes,contents,order,sort,num,start)
}
func GetAppidCount(siteid int64, appid int64, stypes, contents []string, order, sort string) (int64, error) {
	_,cnt,err := getAppidListCount(2,siteid,appid,stypes,contents,order,sort,1,0)

	return cnt,err
}

func getAppidListCount(cate int,siteid int64,appid int64,stypes []string, contents []string,order,sort string,num,start int64)(interface{},int64,error) {

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
		
		switch stypes[i] {
		case "1":
			v = " a.appid='" + contents[i]+"'"
		case "2":
			v = " a.appkey=" + contents[i]
		case "3":	
			v = " a.status=" + contents[i]
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
		sqlstr = SQL_APPID
	} else {
		sqlstr = SQL_COUNT_APPID
		num = 1
		start = 0
	}

	statement = fmt.Sprintf(sqlstr +
						"where 1=1 %s order by a.appid %s limit ? offset ?",
						imconf.Config.SysconfdbSchema,conditions,sort)

	o := orm.NewOrm()

	cnt,err = o.Raw(statement,num,start).Values(&vs)
			
	if err != nil {
		logger.Error("不能获取appid配置信息列表：%v",err.Error())
		return nil,0,err
	}

	return vs,cnt,nil
}

func GetAppid(id int64,siteid int64,appid int64)(interface{},int64,error) {

	var vs []orm.Params
	
	var cnt int64
	var err error
	
	logger.Finest("siteid=%d,appid=%d,areaid=%d",siteid,appid,id)
	
	statement := fmt.Sprintf(SQL_APPID +
					"where 1=1 and a.appid=? ",
					imconf.Config.SysconfdbSchema)

	o := orm.NewOrm()
	cnt,err = o.Raw(statement,id).Values(&vs)

	if err != nil {
		logger.Error("不能获取domain conf配置信息列表：%v",err.Error())
		return nil,0,err
	}

	return vs,cnt,nil
}



func PostAppid(appid,siteid int64,token string, param map[string]interface{} )( id int64, err error) {
	
	if param == nil {
		return 0,errors.New("no input")
	}
	
	var v Appid

	if param["id"]!= nil {
		v.Appid = utils.Convert2Int64(param["id"])
	}
	if param["appkey"] != nil {
		v.Appkey = utils.ConvertToString(param["appkey"])
	}
	if param["remark"] != nil {
		v.Remark = utils.ConvertToString(param["remark"])
	}
	if param["json"] != nil {
		v.Json = utils.ConvertToString(param["json"])
	}
	if param["updated"] != nil {
		v.Updated = time.Now()
	}
	if param["siteid"] != nil {
		v.SiteId = utils.Convert2Int64(param["siteid"])
	}
	if param["status"] != nil {
		v.Status = utils.Convert2Int(param["status"])
	}

	o := orm.NewOrm()

	err = SetSearchPath(o,imconf.Config.SysconfdbSchema)
	if err !=nil {
		logger.Error("不能设置search_path :%v",err)
		return 0,err
	}

	ed,err := o.Insert(&v)
	if err.Error != nil {
	    logger.Error("不能插入appid信息：%v",err.Error())
	    return 0,err
	}
	
			
	return ed, nil
}

func PutAppid( id int64, param map[string]interface{} )( cnt int64,err error) {
	
	if param == nil {
		return 0,errors.New("no input")
	}

	var v Appid
	
	if param["id"]!= nil {
		v.Appid = utils.Convert2Int64(param["id"])
	}
	
	if id != v.Appid {
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

	if param["appkey"] != nil {
		v.Appkey = utils.ConvertToString(param["appkey"])
	}
	if param["remark"] != nil {
		v.Remark = utils.ConvertToString(param["remark"])
	}
	if param["json"] != nil {
		v.Json = utils.ConvertToString(param["json"])
	}
	if param["updated"] != nil {
		v.Updated = time.Now()
	}
	if param["siteid"] != nil {
		v.SiteId = utils.Convert2Int64(param["siteid"])
	}
	if param["status"] != nil {
		v.Status = utils.Convert2Int(param["status"])
	}

	cnt, err = o.Update(&v)
	if err != nil {
		logger.Error("不能更新appid信息：%v", err.Error())
		return 0, err
	}
	
	return cnt,nil
}	

func DeleteAppid(id int64,param map[string]interface{}) (num int64, err error) {

	var v Appid

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
		v.Appid = utils.Convert2Int64(rid)
		c,_ := o.Delete(&v)
		cnt += c
	}


	return cnt, err
}

