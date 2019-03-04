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
	orm.RegisterModel(new(NcWeekSchedules))
}

func GetSchedWeeklyLists(ctx context.Context,siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {

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
			v = " c.name like " + v
		case "2":
			v = " a.status =" + contents[i]
		case "3":
			v = " b.name like " + v
		case "4":
			v = " a.group_id =" + contents[i]
		case "5":
			v = " a.staff_id =" + contents[i]
		case "6":
			v = " a.sched_type =" + contents[i]
		case "7":
			v = " a.sched_plan_id =" + contents[i]
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

	statement = SQL_WEEKLY_SCHEDULE

	if siteid > 1 {
		statement = fmt.Sprintf(statement+
			"where a.siteid=%d %s order by a.group_id %s limit ? offset ?", siteid, conditions, sort)
	} else {
		statement = fmt.Sprintf(statement+
			"where 1=1 %s order by a.group_id %s limit ? offset ?", conditions, sort)

	}
	cnt, err = o.Raw(statement, num, start).Values(&vs)

	if err != nil {
		logger.Error("不能获取周排班列表：%v", err.Error())
		return nil, 0, err
	}

	return vs, cnt, nil
}

func GetSchedWeeklyCount(ctx context.Context,siteid int64, appid int64, stypes, contents []string, order, sort string) (int64, error) {

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
			v = " c.name like " + v
		case "2":
			v = " a.status =" + contents[i]
		case "3":
			v = " b.name like " + v
		case "4":
			v = " a.group_id =" + contents[i]
		case "5":
			v = " a.staff_id =" + contents[i]
		case "6":
			v = " a.sched_type =" + contents[i]
		case "7":
			v = " a.sched_plan_id =" + contents[i]
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
	statement = fmt.Sprintf(SQL_COUNT_WEEKLY_SCHEDULE)

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

func GetSchedWeekly(ctx context.Context,id int64, siteid int64, appid int64) (interface{}, int64, error) {

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

	statement := SQL_WEEKLY_SCHEDULE

	if siteid > 1 {
		cnt, err = o.Raw(statement+
			"where a.siteid=? and a.week_sched_id =?", siteid, id).Values(&vs)
	} else {
		cnt, err = o.Raw(statement+
			"where a.week_sched_id =?", id).Values(&vs)
	}

	if err != nil {
		logger.Error("不能获取周排班信息：%v", err.Error())
		return nil, 0, err
	}

	return vs, cnt, nil
}

func PostSchedWeekly(ctx context.Context,siteid, appid int64, token string, param map[string]interface{}) (id int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v NcWeekSchedules

	if param["id"] != nil {
		v.WeekSchedId = utils.Convert2Int64(param["id"])
	}

	if param["siteid"] != nil {
		v.Siteid = utils.Convert2Int64(param["siteid"])
	} else {
		v.Siteid = siteid
	}

	if param["organizationid"] != nil {
		v.OrganizationId = utils.Convert2Int64(param["organizationid"])
	}

	if param["groupid"] != nil {
		v.GroupId = utils.Convert2Int64(param["groupid"])
	}
	if param["staffid"] != nil {
		v.StaffId = utils.Convert2Int64(param["staffid"])
	}
	if param["schedtype"] != nil {
		v.SchedType = utils.Convert2Int16(param["schedtype"])
	}

	if param["monday"] != nil {
		v.Monday = utils.Convert2Int64(param["monday"])
	}
	if param["tuesday"] != nil {
		v.Tuesday = utils.Convert2Int64(param["tuesday"])
	}
	if param["wednesday"] != nil {
		v.Wednesday = utils.Convert2Int64(param["wednesday"])
	}
	if param["thursday"] != nil {
		v.Thursday = utils.Convert2Int64(param["thursday"])
	}
	if param["friday"] != nil {
		v.Friday = utils.Convert2Int64(param["friday"])
	}
	if param["saturday"] != nil {
		v.Saturday = utils.Convert2Int64(param["saturday"])
	}
	if param["sunday"] != nil {
		v.Sunday = utils.Convert2Int64(param["sunday"])
	}

	if param["starttime"] != nil {
		v.WeekStartDate,_ = time.Parse("2006-01-02 15:04:05",param["starttime"].(string))
	}
	if param["endtime"] != nil {
		v.WeekEndDate,_ = time.Parse("2006-01-02 15:04:05",param["endtime"].(string))
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

	id, err = o.Insert(&v)
	if err != nil {
		logger.Error("不能插入周排班信息：%v", err.Error())
		return 0, err
	}

	return id, err
}

func PutSchedWeekly(ctx context.Context,siteid, appid, id int64, token string, param map[string]interface{}) (cnt int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v NcWeekSchedules

	if param["id"] != nil {
		v.WeekSchedId = utils.Convert2Int64(param["id"])
	}

	if id != v.WeekSchedId {
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

	if param["groupid"] != nil {
		v.GroupId = utils.Convert2Int64(param["groupid"])
	}
	if param["staffid"] != nil {
		v.StaffId = utils.Convert2Int64(param["staffid"])
	}
	if param["schedtype"] != nil {
		v.SchedType = utils.Convert2Int16(param["schedtype"])
	}

	if param["monday"] != nil {
		v.Monday = utils.Convert2Int64(param["monday"])
	}
	if param["tuesday"] != nil {
		v.Tuesday = utils.Convert2Int64(param["tuesday"])
	}
	if param["wednesday"] != nil {
		v.Wednesday = utils.Convert2Int64(param["wednesday"])
	}
	if param["thursday"] != nil {
		v.Thursday = utils.Convert2Int64(param["thursday"])
	}
	if param["friday"] != nil {
		v.Friday = utils.Convert2Int64(param["friday"])
	}
	if param["saturday"] != nil {
		v.Saturday = utils.Convert2Int64(param["saturday"])
	}
	if param["sunday"] != nil {
		v.Sunday = utils.Convert2Int64(param["sunday"])
	}

	if param["starttime"] != nil {
		v.WeekStartDate,_ = time.Parse("2006-01-02 15:04:05",param["starttime"].(string))
	}
	if param["endtime"] != nil {
		v.WeekEndDate,_ = time.Parse("2006-01-02 15:04:05",param["endtime"].(string))
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
		logger.Error("不能更新周排班信息：%v", err.Error())
		return 0, err
	}

	return cnt, nil
}

func PutMultiSchedWeekly(ctx context.Context,siteid, appid, id int64, token string, param map[string]interface{}) (cnt int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v NcWeekSchedules

	if param["siteid"] != nil {
		v.Siteid = utils.Convert2Int64(param["siteid"])
	}
	if param["organizationid"] != nil {
		v.OrganizationId = utils.Convert2Int64(param["organizationid"])
	}

	if param["groupid"] != nil {
		v.GroupId = utils.Convert2Int64(param["groupid"])
	}
	if param["schedtype"] != nil {
		v.SchedType = utils.Convert2Int16(param["schedtype"])
	}

	if param["starttime"] != nil {
		v.WeekStartDate,_ = time.Parse("2006-01-02 15:04:05",param["starttime"].(string))
	}
	if param["endtime"] != nil {
		v.WeekEndDate,_ = time.Parse("2006-01-02 15:04:05",param["endtime"].(string))
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

	var qas []interface{}

	if param["schedules"] != nil {
		qas = param["schedules"].([]interface{})
	}

	o := orm.NewOrm()
	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	for _, sq := range qas {
		maps := sq.(map[string]interface{})

		var qa NcWeekSchedules
		qa = v

		if maps["id"] != nil {
			qa.WeekSchedId = utils.Convert2Int64(maps["id"])
		}
		if maps["staffid"] != nil {
			qa.StaffId = utils.Convert2Int64(maps["staffid"])
		}
		if maps["monday"] != nil {
			qa.Monday = utils.Convert2Int64(maps["monday"])
		}
		if maps["tuesday"] != nil {
			qa.Tuesday = utils.Convert2Int64(maps["tuesday"])
		}
		if maps["wednesday"] != nil {
			qa.Wednesday = utils.Convert2Int64(maps["wednesday"])
		}
		if maps["thursday"] != nil {
			qa.Thursday = utils.Convert2Int64(maps["thursday"])
		}
		if maps["friday"] != nil {
			qa.Friday = utils.Convert2Int64(maps["friday"])
		}
		if maps["saturday"] != nil {
			qa.Saturday = utils.Convert2Int64(maps["saturday"])
		}
		if maps["sunday"] != nil {
			qa.Sunday = utils.Convert2Int64(maps["sunday"])
		}

		cnt, err = o.Update(&qa, "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday")
		if err != nil {
			logger.Error("不能更新周排班信息：%v", err.Error())
			return 0, err
		}

		cnt++

	}

	return cnt, nil
}

func DeleteSchedWeekly(ctx context.Context,id int64, param map[string]interface{}) (num int64, err error) {

	o := orm.NewOrm()
	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	var v NcWeekSchedules

	var ids []interface{}

	if param["id"] != nil {
		ids = param["id"].([]interface{})
	}

	var cnt int64 = 0

	for _, rid := range ids {
		v.WeekSchedId = utils.Convert2Int64(rid)
		num, err = o.Delete(&v)
		cnt += num
	}

	return cnt, err
}
