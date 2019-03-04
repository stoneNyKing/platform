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
	orm.RegisterModel(new(NcGroupStaff))
}

func GetGroupStaffLists(ctx context.Context,siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {

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

		logger.Finest("len=%d,stype[%d]=%s,scontent[%d]=%s", l, i, stypes[i], i, contents[i])

		switch stypes[i] {
		case "1":
			v = " b.name like " + v
		case "2":
			v = " a.status =" + contents[i]
		case "3":
			v = " c.name like " + v
		case "4":
			v = " a.staff_id =" + contents[i]
		case "5":
			v = " a.group_id =" + contents[i]
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

	statement = SQL_GROUP_STAFF

	if siteid > 1 {
		statement = fmt.Sprintf(statement+
			"where a.siteid =%d %s order by a.group_id %s limit ? offset ?", siteid, conditions, sort)
	} else {
		statement = fmt.Sprintf(statement+
			"where 1=1 %s order by a.group_id %s limit ? offset ?", conditions, sort)

	}
	cnt, err = o.Raw(statement, num, start).Values(&vs)

	if err != nil {
		logger.Error("不能获班组员工列表：%v", err.Error())
		return nil, 0, err
	}

	return vs, cnt, nil
}

func GetGroupStaffCount(ctx context.Context,siteid int64, appid int64, stypes, contents []string, order, sort string) (int64, error) {

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
			v = " b.name like " + v
		case "2":
			v = " a.status =" + contents[i]
		case "3":
			v = " c.name like " + v
		case "4":
			v = " a.staff_id =" + contents[i]
		case "5":
			v = " a.group_id =" + contents[i]
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
	statement = fmt.Sprintf(SQL_COUNT_GROUP_STAFF)

	if siteid > 1 {
		statement = fmt.Sprintf(statement+
			"where a.siteid=%d %s order by a.group_id %s limit 1 ", siteid, conditions, sort)
	} else {
		statement = fmt.Sprintf(statement+
			"where 1=1 %s order by a.group_id %s limit 1 ", conditions, sort)

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

func GetGroupStaff(ctx context.Context,id int64, siteid int64, appid int64) (interface{}, int64, error) {

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

	statement := SQL_GROUP_STAFF

	if siteid > 1 {
		cnt, err = o.Raw(statement+
			"where a.siteid =? and a.grp_staff_id =?", siteid, id).Values(&vs)
	} else {
		cnt, err = o.Raw(statement+
			"where a.grp_staff_id =?", id).Values(&vs)
	}

	if err != nil {
		logger.Error("不能获取班组员工信息：%v", err.Error())
		return nil, 0, err
	}

	return vs, cnt, nil
}

func PostGroupStaff(ctx context.Context,siteid, appid int64, token string, param map[string]interface{}) (id int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v NcGroupStaff

	if param["id"] != nil {
		v.GrpStaffId = utils.Convert2Int64(param["id"])
	}

	if param["staffid"] != nil {
		v.StaffId = utils.Convert2Int64(param["staffid"])
	}
	if param["groupid"] != nil {
		v.GroupId = utils.Convert2Int64(param["groupid"])
	}
	if param["position"] != nil {
		v.Position = utils.Convert2Int(param["position"])
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
		logger.Error("不能插入班组员工信息：%v", err.Error())
		return 0, err
	}

	return id, err
}

func PostMultiGroupStaff(ctx context.Context,siteid, appid int64, token string, param map[string]interface{}) (successCnt int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v NcGroupStaff

	if param["id"] != nil {
		v.GrpStaffId = utils.Convert2Int64(param["id"])
	}

	if param["groupid"] != nil {
		v.GroupId = utils.Convert2Int64(param["groupid"])
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

	var vs []NcGroupStaff
	var staffs []interface{}

	if param["staffs"] != nil {
		staffs = param["staffs"].([]interface{})
	}
	var cnt int
	for _, sq := range staffs {
		maps := sq.(map[string]interface{})
		var si NcGroupStaff

		si = v

		if maps["position"] != nil {
			si.Position = utils.Convert2Int(maps["position"])
		}
		if maps["staffid"] != nil {
			si.StaffId = utils.Convert2Int64(maps["staffid"])
		}

		vs = append(vs, si)
		cnt++
	}

	o := orm.NewOrm()

	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	successCnt, err = o.InsertMulti(cnt, &vs)
	if err != nil {
		logger.Error("不能插入班组员工信息：%v", err.Error())
		return 0, err
	}

	return successCnt, err
}

func PutGroupStaff(ctx context.Context,siteid, appid, id int64, token string, param map[string]interface{}) (cnt int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v NcGroupStaff

	if param["id"] != nil {
		v.GrpStaffId = utils.Convert2Int64(param["id"])
	}

	if id != v.GrpStaffId {
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

	if param["staffid"] != nil {
		v.StaffId = utils.Convert2Int64(param["staffid"])
	}
	if param["groupid"] != nil {
		v.GroupId = utils.Convert2Int64(param["groupid"])
	}
	if param["position"] != nil {
		v.Position = utils.Convert2Int(param["position"])
	}

	//通用字段
	if param["remark"] != nil {
		v.Remark = param["remark"].(string)
	}
	if param["status"] != nil {
		v.Status = utils.Convert2Int16(param["status"])
	}

	//o := orm.NewOrm()

	cnt, err = o.Update(&v)
	if err != nil {
		logger.Error("不能更新班组员工信息：%v", err.Error())
		return 0, err
	}

	return cnt, nil
}

func PutMultiGroupStaff(ctx context.Context,siteid, appid, id int64, token string, param map[string]interface{}) (successCnt int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v NcGroupStaff

	if param["groupid"] != nil {
		v.GroupId = utils.Convert2Int64(param["groupid"])
	}

	//通用字段
	if param["remark"] != nil {
		v.Remark = param["remark"].(string)
	}
	if param["status"] != nil {
		v.Status = utils.Convert2Int16(param["status"])
	}
	v.CreateTime = time.Now()

	var vs []NcGroupStaff
	var staffs []interface{}

	if param["staffs"] != nil {
		staffs = param["staffs"].([]interface{})
	}
	var cnt int
	for _, sq := range staffs {
		maps := sq.(map[string]interface{})
		var si NcGroupStaff

		si = v

		if maps["position"] != nil {
			si.Position = utils.Convert2Int(maps["position"])
		}
		if maps["staffid"] != nil {
			si.StaffId = utils.Convert2Int64(maps["staffid"])
		}

		vs = append(vs, si)
		cnt++
	}

	o := orm.NewOrm()

	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	o.Begin()

	_, err = o.QueryTable("nc_group_staff").
		Filter("group_id", v.GroupId).Delete()

	if err != nil {
		o.Rollback()
		logger.Error("不能删除班组员工信息：%v", err.Error())
		return 0, err
	}

	successCnt, err = o.InsertMulti(cnt, &vs)
	if err != nil {
		o.Rollback()
		logger.Error("不能更新班组员工信息：%v", err.Error())
		return 0, err
	}

	o.Commit()

	return successCnt, nil
}

func DeleteGroupStaff(ctx context.Context,id int64, param map[string]interface{}) (num int64, err error) {

	o := orm.NewOrm()
	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	var v NcGroupStaff

	var ids []interface{}

	if param["id"] != nil {
		ids = param["id"].([]interface{})
	}

	var cnt int64 = 0

	for _, rid := range ids {
		v.GrpStaffId = utils.Convert2Int64(rid)
		num, err = o.Delete(&v)
		cnt += num
	}

	return cnt, err
}
