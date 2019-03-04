//
package models

import (
	"context"
	"errors"
	"fmt"
	"platform/common/utils"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"platform/hrsvc/imconf"
	"strconv"
)

func init() {
	orm.RegisterModel(new(NcDepartment))
}


func GetDeptLists(ctx context.Context,siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {
	if num <= 0 {
		num=PAGENUM_MAX
	}

	return getDeptListCount(ctx,1,siteid,appid,stypes,contents,order,sort,num,start)
}
func GetDeptCount(ctx context.Context,siteid int64, appid int64, stypes, contents []string, order, sort string) (int64, error) {
	_,cnt,err := getDeptListCount(ctx,2,siteid,appid,stypes,contents,order,sort,1,0)

	return cnt,err
}

func getDeptListCount(ctx context.Context,cate int,siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {

	var vs []orm.Params

	o := orm.NewOrm()

	var cnt int64
	var err error
	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return nil, 0, err
	}

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
			v = " a.status =" + contents[i]
		case "3":
			v = " a.parent_id =" + contents[i]
		case "4":
			v = " a.level =" + contents[i]
		case "5":
			v = " a.organization_id =" + contents[i]
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
		statement = SQL_DEPARTMENT_LIST
	}else if cate ==2 {
		statement = SQL_COUNT_DEPARTMENT_LIST
	}

	if siteid > 1 {
		statement = fmt.Sprintf(statement+
			"where a.siteid =%d %s order by a.department_id %s limit ? offset ?", siteid, conditions, sort)
	} else {
		statement = fmt.Sprintf(statement+
			"where 1=1 %s order by a.department_id %s limit ? offset ?", conditions, sort)

	}
	cnt, err = o.Raw(statement, num, start).Values(&vs)

	if err != nil {
		logger.Error("不能获取部门列表：%v", err.Error())
		return nil, 0, err
	}

	if cate == 2 {
		if cnt > 0 {
			cnt = utils.Convert2Int64(vs[0]["ucount"])
		}
	}

	return vs, cnt, nil
}

func GetDept(ctx context.Context,id int64, siteid int64, appid int64) (interface{}, int64, error) {

	var vs []orm.Params

	o := orm.NewOrm()

	var cnt int64
	var err error
	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return nil, 0, err
	}

	logger.Finest("siteid=%d,appid=%d,id=%d", siteid, appid, id)

	statement := SQL_DEPARTMENT_LIST

	if siteid > 1 {
		cnt, err = o.Raw(statement+
			"where a.siteid=? and a.department_id =?", siteid, id).Values(&vs)
	} else {
		cnt, err = o.Raw(statement+
			"where a.department_id =?", id).Values(&vs)
	}

	if err != nil {
		logger.Error("不能获取部门信息：%v", err.Error())
		return nil, 0, err
	}

	return vs, cnt, nil
}

func GetDeptTree(ctx context.Context,siteid int64, appid int64, stypes, contents []string, order, sort string) (interface{}, error) {

	var maps []map[string]interface{}

	var err error

	if sort == "" {
		sort = "desc"
	}

	if len(stypes) != len(contents) {
		return nil, errors.New("params number is not match.")
	}

	logger.Finest("siteid=%d,appid=%d,order=%s,sort=%s", siteid, appid, order, sort)

	conditions := ""

	res_id := ""

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
			v = " a.parent_id =" + contents[i]
		case "3":
			v = " a.level =" + contents[i]
		case "4":
			v = ""
			res_id = contents[i]
		case "5":
			v = " a.organization_id =" + contents[i]
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

	logger.Finest("conditions = %s,res_id=%s", conditions, res_id)

	var statement string
	statement = fmt.Sprintf(SQL_DEPARTMENT_LIST)

	if siteid > 1 {
		statement = fmt.Sprintf(statement +
			"where siteid=%d %s order by a.department_id %s ", siteid,conditions, sort)
	}else{
		statement = fmt.Sprintf(statement +
			"where 1=1 %s order by a.department_id %s ", conditions, sort)

	}

	var vs []orm.Params

	o := orm.NewOrm()

	var cnt int64
	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return nil,  err
	}

	cnt, err = o.Raw(statement).Values(&vs)

	if err != nil {
		logger.Error("不能获取列表信息：%v", err.Error())
		return nil, err
	}

	if cnt > 0 {
		for _,v := range vs {
			rid := utils.ConvertToString(v["id"])
			child ,_ ,_:= GetDeptTreeChild(o,siteid,rid)
			v["child"] = child
			maps = append(maps,v)
		}
	}

	return maps, err
}

func GetDeptTreeChild(o orm.Ormer,siteid int64, resids string) ([]map[string]interface{}, int64, error) {
	var maps []map[string]interface{}

	var vs []orm.Params

	if o == nil {
		o = orm.NewOrm()
	}

	var cnt int64
	var err error
	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return nil, 0, err
	}

	statement := fmt.Sprintf(SQL_DEPARTMENT_TREE+
		"where a.parent_id=%s and a.siteid =%d order by a.department_id desc ; ", resids, siteid)

	//logger.Finest("statement=%s", statement)

	cnt, err = o.Raw(statement).Values(&vs)

	if err != nil {
		logger.Error("不能获取部门信息列表：%v", err.Error())
		return nil, 0, err
	}

	var res_ids string = ""
	if cnt > 0 {
		for _, v := range vs {
			m := make(map[string]interface{})

			for key, value := range v {
				m[key] = value
			}
			if v["id"] != nil {
				res_ids = v["id"].(string)
				if res_ids != "" {
					pas, count, err := GetDeptTreeChild(o,siteid, res_ids)

					if err == nil && count > 0 {
						m["child"] = pas
					}
				} else {
					m["child"] = "[]"
				}
			}
			maps = append(maps, m)
		}

	}

	return maps, cnt, nil
}

func PostDept(ctx context.Context,siteid, appid int64, token string, param map[string]interface{}) (id int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v NcDepartment

	if param["id"] != nil {
		v.DepartmentId = utils.Convert2Int64(param["id"])
	}

	if param["siteid"] != nil {
		v.Siteid = utils.Convert2Int64(param["siteid"])
	} else {
		v.Siteid = siteid
	}

	if param["organizationid"] != nil {
		v.OrganizationId = utils.Convert2Int64(param["organizationid"])
	}

	if param["name"] != nil {
		v.Name = param["name"].(string)
	}

	if param["parentid"] != nil {
		v.ParentId = utils.Convert2Int64(param["parentid"])
	}
	if param["level"] != nil {
		v.Level = utils.Convert2Int16(param["level"])
	}

	//通用字段
	if param["remark"] != nil {
		v.Remark = param["remark"].(string)
	}

	if param["operatorid"] != nil {
		v.Creatorid = utils.Convert2Int64(param["operatorid"])
		v.Updaterid = v.Creatorid
	}

	if param["status"] != nil {
		v.Status = utils.Convert2Int16(param["status"])
	} else {
		v.Status = 1
	}

	v.CreateTime,_ = time.Parse("2006-01-02 15:04:05",utils.GetTimeFormat("2006-01-02 15:04:05"))
	v.UpdateTime = v.CreateTime

	o := orm.NewOrm()

	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	id, err = o.Insert(&v)
	if err != nil {
		logger.Error("不能插入员工信息：%v", err.Error())
		if strings.Contains(err.Error(),"1062") {
			err = errors.New("科室名称重复。")
		}
		return 0, err
	}

	return id, err
}

func PutDept(ctx context.Context,siteid, appid, id int64, token string, param map[string]interface{}) (cnt int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v NcDepartment

	if param["id"] != nil {
		v.DepartmentId = utils.Convert2Int64(param["id"])
	}

	if id != v.DepartmentId {
		return 0, errors.New("id is not match.")
	}

	o := orm.NewOrm()
	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	err = o.Read(&v)
	if err != nil {
		logger.Error("cannot read record.")
		return 0, err
	}

	if param["siteid"] != nil {
		v.Siteid = utils.Convert2Int64(param["siteid"])
	}
	if param["organizationid"] != nil {
		v.OrganizationId = utils.Convert2Int64(param["organizationid"])
	}

	if param["name"] != nil {
		v.Name = param["name"].(string)
	}

	if param["parentid"] != nil {
		v.ParentId = utils.Convert2Int64(param["parentid"])
	}
	if param["level"] != nil {
		v.Level = utils.Convert2Int16(param["level"])
	}

	//通用字段
	if param["remark"] != nil {
		v.Remark = param["remark"].(string)
	}
	if param["status"] != nil {
		v.Status = utils.Convert2Int16(param["status"])
	}

	if param["operatorid"] != nil {
		v.Updaterid = utils.Convert2Int64(param["operatorid"])
	}

	v.CreateTime,_ = time.Parse("2006-01-02 15:04:05",v.CreateTime.Format("2006-01-02 15:04:05"))
	v.UpdateTime,_ = time.Parse("2006-01-02 15:04:05",utils.GetTimeFormat("2006-01-02 15:04:05"))

	cnt, err = o.Update(&v)
	if err != nil {
		logger.Error("不能更新部门信息：%v", err.Error())
		return 0, err
	}

	return cnt, nil
}

func DeleteDept(ctx context.Context,siteid, id int64, param map[string]interface{}) (num int64, err error) {

	o := orm.NewOrm()
	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	var v NcDepartment

	var ids []interface{}

	if param["id"] != nil {
		ids = param["id"].([]interface{})
	}

	var cnt int64 = 0

	for _, rid := range ids {
		// 需要判断在当前科室下有没有子科室
		_, childs, _ := GetDeptTreeChild(o,siteid, strconv.FormatInt(utils.Convert2Int64(rid), 10))
		if childs > 0 {
			err = errors.New("该科室下有子科室不能删除")
			continue
		}
		v.DepartmentId = utils.Convert2Int64(rid)
		num, err = o.Delete(&v)
		if err != nil {
			if strings.Contains(err.Error(),"1451") {
				err = errors.New("该科室已被关联,不能删除。")
			}
			return 0, err
		}
		cnt += num
	}

	return cnt, err
}
