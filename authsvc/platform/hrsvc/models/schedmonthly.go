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
	orm.RegisterModel(new(NcMonthSchedules))
}

func GetSchedMonthlyLists(ctx context.Context,siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {

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
			v = " a.sched_month ='" + contents[i] + "'"
		case "2":
			v = " a.status =" + contents[i]
		case "3":
			v = " a.staff_id =" + contents[i]
		case "4":
			v = " a.group_id =" + contents[i]
		case "5":
			v = " b.name like " + v
		case "6":
			v = " a.group_id =" + contents[i]
		case "7":
			v = " c.name like " + v
		case "8":
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

	statement = SQL_SCHEDULE_MONTHLY

	if siteid > 1 {
		statement = fmt.Sprintf(statement+
			"where a.siteid =%d %s order by a.group_id %s limit ? offset ?", siteid, conditions, sort)
	} else {
		statement = fmt.Sprintf(statement+
			"where 1=1 %s order by a.group_id %s limit ? offset ?", conditions, sort)

	}
	cnt, err = o.Raw(statement, num, start).Values(&vs)

	if err != nil {
		logger.Error("不能获月排班列表：%v", err.Error())
		return nil, 0, err
	}

	return vs, cnt, nil
}

func GetSchedMonthlyCount(ctx context.Context,siteid int64, appid int64, stypes, contents []string, order, sort string) (int64, error) {

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
			v = " a.sched_month ='" + contents[i] + "'"
		case "2":
			v = " a.status =" + contents[i]
		case "3":
			v = " a.staff_id =" + contents[i]
		case "4":
			v = " a.group_id =" + contents[i]
		case "5":
			v = " b.name like " + v
		case "6":
			v = " a.group_id =" + contents[i]
		case "7":
			v = " c.name like " + v
		case "8":
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
	statement = fmt.Sprintf(SQL_COUNT_SCHEDULE_MONTHLY)

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

func GetSchedMonthly(ctx context.Context,id int64, siteid int64, appid int64) (interface{}, int64, error) {

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

	statement := SQL_SCHEDULE_MONTHLY

	if siteid > 1 {
		cnt, err = o.Raw(statement+
			"where a.siteid=? and a.month_sched_id =?", siteid, id).Values(&vs)
	} else {
		cnt, err = o.Raw(statement+
			"where a.month_sched_id =?", id).Values(&vs)
	}

	if err != nil {
		logger.Error("不能获取月排班信息：%v", err.Error())
		return nil, 0, err
	}

	return vs, cnt, nil
}

func PostSchedMonthly(ctx context.Context,siteid, appid int64, token string, param map[string]interface{}) (id int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v NcMonthSchedules

	if param["id"] != nil {
		v.MonthSchedId = utils.Convert2Int64(param["id"])
	}

	if param["siteid"] != nil {
		v.Siteid = utils.Convert2Int64(param["siteid"])
	} else {
		v.Siteid = siteid
	}

	if param["organizationid"] != nil {
		v.OrganizationId = utils.Convert2Int64(param["organizationid"])
	}

	if param["schedmonth"] != nil {
		v.SchedMonth,_ = time.Parse("2006-01-02 15:04:05",param["schedmonth"].(string))
	}

	if param["groupid"] != nil {
		v.GroupId = utils.Convert2Int64(param["groupid"])
	}
	if param["staffid"] != nil {
		v.StaffId = utils.Convert2Int64(param["staffid"])
	}
	if param["schedplanid"] != nil {
		v.SchedPlanId = utils.Convert2Int64(param["schedplanid"])
	}

	if param["day1"] != nil {
		v.Day1 = utils.Convert2Int64(param["day1"])
	}
	if param["day2"] != nil {
		v.Day2 = utils.Convert2Int64(param["day2"])
	}
	if param["day3"] != nil {
		v.Day3 = utils.Convert2Int64(param["day3"])
	}
	if param["day4"] != nil {
		v.Day4 = utils.Convert2Int64(param["day4"])
	}
	if param["day5"] != nil {
		v.Day5 = utils.Convert2Int64(param["day5"])
	}
	if param["day6"] != nil {
		v.Day6 = utils.Convert2Int64(param["day6"])
	}
	if param["day7"] != nil {
		v.Day7 = utils.Convert2Int64(param["day7"])
	}
	if param["day8"] != nil {
		v.Day8 = utils.Convert2Int64(param["day8"])
	}
	if param["day9"] != nil {
		v.Day9 = utils.Convert2Int64(param["day9"])
	}
	if param["day10"] != nil {
		v.Day10 = utils.Convert2Int64(param["day10"])
	}
	if param["day11"] != nil {
		v.Day11 = utils.Convert2Int64(param["day11"])
	}
	if param["day12"] != nil {
		v.Day12 = utils.Convert2Int64(param["day12"])
	}
	if param["day13"] != nil {
		v.Day13 = utils.Convert2Int64(param["day13"])
	}
	if param["day14"] != nil {
		v.Day14 = utils.Convert2Int64(param["day14"])
	}
	if param["day15"] != nil {
		v.Day15 = utils.Convert2Int64(param["day15"])
	}
	if param["day16"] != nil {
		v.Day16 = utils.Convert2Int64(param["day16"])
	}
	if param["day17"] != nil {
		v.Day17 = utils.Convert2Int64(param["day17"])
	}
	if param["day18"] != nil {
		v.Day18 = utils.Convert2Int64(param["day18"])
	}
	if param["day19"] != nil {
		v.Day19 = utils.Convert2Int64(param["day19"])
	}
	if param["day20"] != nil {
		v.Day20 = utils.Convert2Int64(param["day20"])
	}
	if param["day21"] != nil {
		v.Day21 = utils.Convert2Int64(param["day21"])
	}
	if param["day22"] != nil {
		v.Day22 = utils.Convert2Int64(param["day22"])
	}
	if param["day23"] != nil {
		v.Day23 = utils.Convert2Int64(param["day23"])
	}
	if param["day24"] != nil {
		v.Day24 = utils.Convert2Int64(param["day24"])
	}
	if param["day25"] != nil {
		v.Day25 = utils.Convert2Int64(param["day25"])
	}
	if param["day26"] != nil {
		v.Day26 = utils.Convert2Int64(param["day26"])
	}
	if param["day27"] != nil {
		v.Day27 = utils.Convert2Int64(param["day27"])
	}
	if param["day28"] != nil {
		v.Day28 = utils.Convert2Int64(param["day28"])
	}
	if param["day29"] != nil {
		v.Day29 = utils.Convert2Int64(param["day29"])
	}
	if param["day30"] != nil {
		v.Day30 = utils.Convert2Int64(param["day30"])
	}
	if param["day31"] != nil {
		v.Day31 = utils.Convert2Int64(param["day31"])
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
		logger.Error("不能插入月排班信息：%v", err.Error())
		return 0, err
	}

	return id, err
}

func PutSchedMonthly(ctx context.Context,siteid, appid, id int64, token string, param map[string]interface{}) (cnt int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v NcMonthSchedules

	if param["id"] != nil {
		v.MonthSchedId = utils.Convert2Int64(param["id"])
	}

	if id != v.MonthSchedId {
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

	if param["schedmonth"] != nil {
		v.SchedMonth,_ = time.Parse("2006-01-02 15:04:05",param["schedmonth"].(string))
	}

	if param["groupid"] != nil {
		v.GroupId = utils.Convert2Int64(param["groupid"])
	}
	if param["staffid"] != nil {
		v.StaffId = utils.Convert2Int64(param["staffid"])
	}
	if param["schedplanid"] != nil {
		v.SchedPlanId = utils.Convert2Int64(param["schedplanid"])
	}

	if param["day1"] != nil {
		v.Day1 = utils.Convert2Int64(param["day1"])
	}
	if param["day2"] != nil {
		v.Day2 = utils.Convert2Int64(param["day2"])
	}
	if param["day3"] != nil {
		v.Day3 = utils.Convert2Int64(param["day3"])
	}
	if param["day4"] != nil {
		v.Day4 = utils.Convert2Int64(param["day4"])
	}
	if param["day5"] != nil {
		v.Day5 = utils.Convert2Int64(param["day5"])
	}
	if param["day6"] != nil {
		v.Day6 = utils.Convert2Int64(param["day6"])
	}
	if param["day7"] != nil {
		v.Day7 = utils.Convert2Int64(param["day7"])
	}
	if param["day8"] != nil {
		v.Day8 = utils.Convert2Int64(param["day8"])
	}
	if param["day9"] != nil {
		v.Day9 = utils.Convert2Int64(param["day9"])
	}
	if param["day10"] != nil {
		v.Day10 = utils.Convert2Int64(param["day10"])
	}
	if param["day11"] != nil {
		v.Day11 = utils.Convert2Int64(param["day11"])
	}
	if param["day12"] != nil {
		v.Day12 = utils.Convert2Int64(param["day12"])
	}
	if param["day13"] != nil {
		v.Day13 = utils.Convert2Int64(param["day13"])
	}
	if param["day14"] != nil {
		v.Day14 = utils.Convert2Int64(param["day14"])
	}
	if param["day15"] != nil {
		v.Day15 = utils.Convert2Int64(param["day15"])
	}
	if param["day16"] != nil {
		v.Day16 = utils.Convert2Int64(param["day16"])
	}
	if param["day17"] != nil {
		v.Day17 = utils.Convert2Int64(param["day17"])
	}
	if param["day18"] != nil {
		v.Day18 = utils.Convert2Int64(param["day18"])
	}
	if param["day19"] != nil {
		v.Day19 = utils.Convert2Int64(param["day19"])
	}
	if param["day20"] != nil {
		v.Day20 = utils.Convert2Int64(param["day20"])
	}
	if param["day21"] != nil {
		v.Day21 = utils.Convert2Int64(param["day21"])
	}
	if param["day22"] != nil {
		v.Day22 = utils.Convert2Int64(param["day22"])
	}
	if param["day23"] != nil {
		v.Day23 = utils.Convert2Int64(param["day23"])
	}
	if param["day24"] != nil {
		v.Day24 = utils.Convert2Int64(param["day24"])
	}
	if param["day25"] != nil {
		v.Day25 = utils.Convert2Int64(param["day25"])
	}
	if param["day26"] != nil {
		v.Day26 = utils.Convert2Int64(param["day26"])
	}
	if param["day27"] != nil {
		v.Day27 = utils.Convert2Int64(param["day27"])
	}
	if param["day28"] != nil {
		v.Day28 = utils.Convert2Int64(param["day28"])
	}
	if param["day29"] != nil {
		v.Day29 = utils.Convert2Int64(param["day29"])
	}
	if param["day30"] != nil {
		v.Day30 = utils.Convert2Int64(param["day30"])
	}
	if param["day31"] != nil {
		v.Day31 = utils.Convert2Int64(param["day31"])
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
		logger.Error("不能更新月排班信息：%v", err.Error())
		return 0, err
	}

	return cnt, nil
}

func DeleteSchedMonthly(ctx context.Context,id int64, param map[string]interface{}) (num int64, err error) {

	o := orm.NewOrm()
	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	var v NcMonthSchedules

	var ids []interface{}

	if param["id"] != nil {
		ids = param["id"].([]interface{})
	}

	var cnt int64 = 0

	for _, rid := range ids {
		v.MonthSchedId = utils.Convert2Int64(rid)
		num, err = o.Delete(&v)
		cnt += num
	}

	return cnt, err
}
