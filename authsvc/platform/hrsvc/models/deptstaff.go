//
package models

import (
	"context"
	"errors"
	"fmt"
	"platform/common/utils"
	"time"

	"github.com/astaxie/beego/orm"
	"platform/hrsvc/imconf"
)

func init() {
	orm.RegisterModel(new(NcDepartmentStaff))
}

func GetDeptStaffLists(ctx context.Context,siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {

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
			v = " b.position_id =" + contents[i]
		case "6":
			v = " a.department_id =" + contents[i]
		case "7":
			v = " c.name like " + v
		case "8":
			v = " d.name like " + v
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

	statement = SQL_DEPT_STAFF

	if siteid > 1 {
		statement = fmt.Sprintf(statement+
			"where a.siteid=%d %s order by a.department_id %s limit ? offset ?", siteid, conditions, sort)
	} else {
		statement = fmt.Sprintf(statement+
			"where 1=1 %s order by a.department_id %s limit ? offset ?", conditions, sort)

	}
	cnt, err = o.Raw(statement, num, start).Values(&vs)

	if err != nil {
		logger.Error("不能获取部门列表：%v", err.Error())
		return nil, 0, err
	}

	return vs, cnt, nil
}

func GetDeptStaff(ctx context.Context,id int64, siteid int64, appid int64) (interface{}, int64, error) {

	var vs []orm.Params

	o := orm.NewOrm()

	var cnt int64
	var err error

	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return nil, 0, err
	}

	logger.Finest("siteid=%d,appid=%d,areaid=%d", siteid, appid, id)

	statement := SQL_DEPT_STAFF

	if siteid > 1 {
		cnt, err = o.Raw(statement+
			"where a.siteid=?  and a.department_id =?", siteid, id).Values(&vs)
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

func GetDeptStaffCount(ctx context.Context,siteid int64, appid int64, stypes, contents []string, order, sort string) (int64, error) {

	var vs []orm.Params

	o := orm.NewOrm()

	var cnt int64
	var err error

	if sort == "" {
		sort = "desc"
	}

	logger.Finest("siteid=%d,appid=%d,order=%s,sort=%s", siteid, appid, order, sort)

	if len(stypes) != len(contents) {
		return 0, errors.New("params number is not match.")
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
			v = " b.position_id =" + contents[i]
		case "6":
			v = " a.department_id =" + contents[i]
		case "7":
			v = " c.name like " + v
		case "8":
			v = " d.name like " + v
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
	statement = fmt.Sprintf(SQL_COUNT_DEPT_STAFF)

	if siteid > 1 {
		statement = fmt.Sprintf(statement+
			"where a.siteid=%d %s order by a.department_id %s limit 1 ", siteid, conditions, sort)
	} else {
		statement = fmt.Sprintf(statement+
			"where 1=1 %s order by a.department_id %s limit 1 ", conditions, sort)

	}

	cnt, err = o.Raw(statement).Values(&vs)

	if err != nil {
		logger.Error("不能获取列表数量：%v", err.Error())
		return 0, err
	}

	if cnt > 0 {
		cnt = utils.Convert2Int64(vs[0]["ucount"])
	}

	return cnt, nil
}

func PostDeptStaff(ctx context.Context,siteid, appid int64, token string, param map[string]interface{}) (id int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v NcDepartmentStaff

	if param["id"] != nil {
		v.DepStaffId = utils.Convert2Int64(param["id"])
	}

	if param["siteid"] != nil {
		v.Siteid = utils.Convert2Int64(param["siteid"])
	} else {
		v.Siteid = siteid
	}

	if param["organizationid"] != nil {
		v.OrganizationId = utils.Convert2Int64(param["organizationid"])
	}

	if param["inductiontime"] != nil {
		v.InductionTime,_ = time.Parse("2006-01-02 15:04:05",param["inductiontime"].(string))
	}
	if param["departuretime"] != nil {
		v.DepartureTime,_ = time.Parse("2006-01-02 15:04:05",param["departuretime"].(string))
	}

	if param["staffid"] != nil {
		v.StaffId = utils.Convert2Int64(param["staffid"])
	}
	if param["positionid"] != nil {
		v.PositionId = utils.Convert2Int64(param["positionid"])
	}
	if param["departmentid"] != nil {
		v.DepartmentId = utils.Convert2Int64(param["departmentid"])
	}
	if param["injobflag"] != nil {
		v.InjobFlag = utils.Convert2Int16(param["injobflag"])
	}

	//通用字段
	if param["remark"] != nil {
		v.Remark = param["remark"].(string)
	}

	if param["status"] != nil {
		v.Status = utils.Convert2Int16(param["status"])
	} else {
		v.Status = 1
	}

	v.CreateTime = time.Now()

	o := orm.NewOrm()

	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	id, err = o.Insert(&v)
	if err != nil {
		logger.Error("不能插入员工信息：%v", err.Error())
		return 0, err
	}

	return id, err
}

func PutDeptStaff(ctx context.Context,siteid, appid, id int64, token string, param map[string]interface{}) (cnt int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v NcDepartmentStaff

	if param["id"] != nil {
		v.DepStaffId = utils.Convert2Int64(param["id"])
	}

	if id != v.DepStaffId {
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

	if param["inductiontime"] != nil {
		v.InductionTime,_ = time.Parse("2006-01-02 15:04:05",param["inductiontime"].(string))
	}
	if param["departuretime"] != nil {
		v.DepartureTime,_ = time.Parse("2006-01-02 15:04:05",param["departuretime"].(string))
	}

	if param["injobflag"] != nil {
		v.InjobFlag = utils.Convert2Int16(param["injobflag"])
	}

	if param["staffid"] != nil {
		v.StaffId = utils.Convert2Int64(param["staffid"])
	}
	if param["positionid"] != nil {
		v.PositionId = utils.Convert2Int64(param["positionid"])
	}
	if param["departmentid"] != nil {
		v.DepartmentId = utils.Convert2Int64(param["departmentid"])
	}

	//通用字段
	if param["remark"] != nil {
		v.Remark = param["remark"].(string)
	}
	if param["status"] != nil {
		v.Status = utils.Convert2Int16(param["status"])
	}

	cnt, err = o.Update(&v)
	if err != nil {
		logger.Error("不能更新部门信息：%v", err.Error())
		return 0, err
	}

	return cnt, nil
}

func DeleteDeptStaff(ctx context.Context,id int64, param map[string]interface{}) (num int64, err error) {

	o := orm.NewOrm()
	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	var v NcDepartmentStaff

	var ids []interface{}

	if param["id"] != nil {
		ids = param["id"].([]interface{})
	}

	var cnt int64 = 0

	for _, rid := range ids {
		v.DepStaffId = utils.Convert2Int64(rid)
		num, err = o.Delete(&v)
		cnt += num
	}

	return cnt, err
}
