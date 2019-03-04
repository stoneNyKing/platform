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
	orm.RegisterModel(new(NcSchedPlan))
}

func GetSchedPlanLists(ctx context.Context,siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {

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
			v = " a.group_id =" + contents[i]
		case "4":
			v = " b.name like " + v
		case "5":
			v = " a.plan_type =" + contents[i]
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

	statement = SQL_SCHED_PLAN

	if siteid > 1 {
		statement = fmt.Sprintf(statement+
			"where a.siteid=%d %s order by a.sched_plan_id %s limit ? offset ?", siteid, conditions, sort)
	} else {
		statement = fmt.Sprintf(statement+
			"where 1=1 %s order by a.sched_plan_id %s limit ? offset ?", conditions, sort)

	}
	cnt, err = o.Raw(statement, num, start).Values(&vs)

	if err != nil {
		logger.Error("不能获获取列表：%v", err.Error())
		return nil, 0, err
	}

	return vs, cnt, nil
}

func GetSchedPlanCount(ctx context.Context,siteid int64, appid int64, stypes, contents []string, order, sort string) (int64, error) {

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
			v = " a.group_id =" + contents[i]
		case "4":
			v = " b.name like " + v
		case "5":
			v = " a.plan_type =" + contents[i]
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
	statement = fmt.Sprintf(SQL_COUNT_SCHED_PLAN)

	if siteid > 1 {
		statement = fmt.Sprintf(statement+
			"where a.siteid=%d %s order by a.sched_plan_id %s limit 1 ", siteid, conditions, sort)
	} else {
		statement = fmt.Sprintf(statement+
			"where 1=1 %s order by a.sched_plan_id %s limit 1 ", conditions, sort)

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

func GetSchedPlan(ctx context.Context,id int64, siteid int64, appid int64) (interface{}, int64, error) {

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

	statement := SQL_SCHED_PLAN

	if siteid > 1 {
		cnt, err = o.Raw(statement+
			"where a.siteid=? and a.sched_plan_id =?", siteid, id).Values(&vs)
	} else {
		cnt, err = o.Raw(statement+
			"where a.sched_plan_id =?", id).Values(&vs)
	}

	if err != nil {
		logger.Error("不能获取排班方案信息：%v", err.Error())
		return nil, 0, err
	}

	return vs, cnt, nil
}

func PostSchedPlan(ctx context.Context,siteid, appid int64, token string, param map[string]interface{}) (id int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v NcSchedPlan

	if param["id"] != nil {
		v.SchedPlanId = utils.Convert2Int64(param["id"])
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
	if param["groupid"] != nil {
		v.GroupId = utils.Convert2Int64(param["groupid"])
	}
	if param["plantype"] != nil {
		v.PlanType = utils.Convert2Int16(param["plantype"])
	}

	//通用字段
	if param["remark"] != nil {
		v.Remark = param["remark"].(string)
	}

	if param["operatorid"] != nil {
		v.Operatorid = utils.Convert2Int64(param["operatorid"])
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

	o.Begin()

	id, err = o.Insert(&v)
	if err != nil {
		o.Rollback()
		logger.Error("不能插入信息：%v", err.Error())
		return 0, err
	}

	var gs []NcGroupStaff

	_, err = o.QueryTable("nc_group_staff").Filter("group_id", v.GroupId).Filter("status", 1).All(&gs, "GroupId", "StaffId")
	if err != nil {
		o.Rollback()
		logger.Error("不能获取班组员工信息：%v", err.Error())
		return 0, err
	}

	cnt := 0
	if len(gs) > 0 {
		if v.PlanType == PLAN_TYPE_WEEK {
			var ws []NcWeekSchedules
			for _, val := range gs {
				var w NcWeekSchedules
				w.GroupId = val.GroupId
				w.Siteid = v.Siteid
				w.OrganizationId = v.OrganizationId
				w.Status = 1
				w.CreateTime = time.Now()
				w.SchedPlanId = id
				w.StaffId = val.StaffId
				cnt++
				ws = append(ws, w)
			}

			_, err = o.InsertMulti(cnt, ws)
			if err != nil {
				o.Rollback()
				logger.Error("不能插入排班信息：%v", err.Error())
				return 0, err
			}

		} else if v.PlanType == PLAN_TYPE_MONTH {
			var ms []NcMonthSchedules
			for _, val := range gs {
				var m NcMonthSchedules
				m.GroupId = val.GroupId
				m.Siteid = v.Siteid
				m.OrganizationId = v.OrganizationId
				m.Status = 1
				m.CreateTime = v.CreateTime
				m.SchedPlanId = id
				m.StaffId = val.StaffId
				cnt++
				ms = append(ms, m)
			}

			_, err = o.InsertMulti(cnt, ms)
			if err != nil {
				o.Rollback()
				logger.Error("不能插入排班信息：%v", err.Error())
				return 0, err
			}
		}

	}

	o.Commit()

	return id, err
}

func PutSchedPlan(ctx context.Context,siteid, appid, id int64, token string, param map[string]interface{}) (cnt int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v NcSchedPlan

	if param["id"] != nil {
		v.SchedPlanId = utils.Convert2Int64(param["id"])
	}

	if id != v.SchedPlanId {
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
	if param["groupid"] != nil {
		v.GroupId = utils.Convert2Int64(param["groupid"])
	}
	if param["plantype"] != nil {
		v.PlanType = utils.Convert2Int16(param["plantype"])
	}

	//通用字段
	if param["remark"] != nil {
		v.Remark = param["remark"].(string)
	}
	if param["status"] != nil {
		v.Status = utils.Convert2Int16(param["status"])
	}
	if param["operatorid"] != nil {
		v.Operatorid = utils.Convert2Int64(param["operatorid"])
	}

	//o := orm.NewOrm()

	cnt, err = o.Update(&v)
	if err != nil {
		logger.Error("不能更新信息：%v", err.Error())
		return 0, err
	}

	return cnt, nil
}

func DeleteSchedPlan(ctx context.Context,id int64, param map[string]interface{}) (num int64, err error) {

	o := orm.NewOrm()
	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	var v NcSchedPlan

	var ids []interface{}

	if param["id"] != nil {
		ids = param["id"].([]interface{})
	}

	var cnt int64 = 0

	for _, rid := range ids {
		v.SchedPlanId = utils.Convert2Int64(rid)
		o.Begin()
		_, err = o.Raw("delete from nc_week_schedules where sched_plan_id=? ", v.SchedPlanId).Exec()
		if err != nil {
			o.Rollback()
			continue
		}
		_, err = o.Raw("delete from nc_month_schedules where sched_plan_id=? ", v.SchedPlanId).Exec()
		if err != nil {
			o.Rollback()
			continue
		}
		num, err = o.Delete(&v)
		if err != nil {
			o.Rollback()
			continue
		}
		o.Commit()
		cnt += num
	}

	return cnt, err
}
