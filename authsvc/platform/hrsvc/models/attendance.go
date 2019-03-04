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
	orm.RegisterModel(new(NcAttendance))
}

func GetAttendanceLists(ctx context.Context,siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {

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
			v = " b.name like " + v
		case "2":
			v = " a.staff_id =" + contents[i]
		case "3":
			v = " a.clocktime ='" + contents[i] + "'"
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

	statement = SQL_ATTENDANCE

	if siteid > 1 {
		statement = fmt.Sprintf(statement+
			"where a.siteid=%d %s order by a.attendance_id %s limit ? offset ?", siteid, conditions, sort)
	} else {
		statement = fmt.Sprintf(statement+
			"where 1=1 %s order by a.attendance_id %s limit ? offset ?", conditions, sort)

	}
	cnt, err = o.Raw(statement, num, start).Values(&vs)

	if err != nil {
		logger.Error("不能获考勤列表：%v", err.Error())
		return nil, 0, err
	}

	return vs, cnt, nil
}

func GetAttendanceCount(ctx context.Context,siteid int64, appid int64, stypes, contents []string, order, sort string) (int64, error) {

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
			v = " a.staff_id =" + contents[i]
		case "3":
			v = " a.clocktime ='" + contents[i] + "'"
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
	statement = fmt.Sprintf(SQL_COUNT_ATTENDANCE)

	if siteid > 1 {
		statement = fmt.Sprintf(statement+
			"where (a.siteid=%d) %s order by a.attendance_id %s limit 1 ", siteid, conditions, sort)
	} else {
		statement = fmt.Sprintf(statement+
			"where 1=1 %s order by a.attendance_id %s limit 1 ", conditions, sort)

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

func GetAttendance(ctx context.Context,id int64, siteid int64, appid int64) (interface{}, int64, error) {

	var vs []orm.Params

	o := orm.NewOrm()

	var cnt int64
	var err error

	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return nil, 0, err
	}

	logger.Finest("siteid=%d,appid=%d, id=%d", siteid, appid, id)

	statement := SQL_ATTENDANCE

	if siteid > 1 {
		cnt, err = o.Raw(statement+
			"where a.siteid=? and a.attendance_id =?", siteid, id).Values(&vs)
	} else {
		cnt, err = o.Raw(statement+
			"where a.attendance_id =?", id).Values(&vs)
	}

	if err != nil {
		logger.Error("不能获取考勤信息：%v", err.Error())
		return nil, 0, err
	}

	return vs, cnt, nil
}

func PostAttendance(ctx context.Context,siteid, appid int64, token string, param map[string]interface{}) (id int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v NcAttendance

	if param["id"] != nil {
		v.AttendanceId = utils.Convert2Int64(param["id"])
	}

	if param["siteid"] != nil {
		v.Siteid = utils.Convert2Int64(param["siteid"])
	} else {
		v.Siteid = siteid
	}

	if param["organizationid"] != nil {
		v.OrganizationId = utils.Convert2Int64(param["organizationid"])
	}

	if param["staffid"] != nil {
		v.StaffId = utils.Convert2Int64(param["staffid"])
	}

	if param["clocktime"] != nil {
		v.Clocktime,_ = time.Parse("2006-01-02 15:04:05",param["clocktime"].(string))
	} else {
		v.Clocktime = time.Now()
	}

	if param["clocktimestart"] != nil {
		v.ClockTimeStart,_ = time.Parse("2006-01-02 15:04:05",param["clocktimestart"].(string))
	}
	if param["clocktimeend"] != nil {
		v.ClockTimeEnd,_ = time.Parse("2006-01-02 15:04:05",param["clocktimeend"].(string))
	}

	//通用字段
	if param["remark"] != nil {
		v.Remark = param["remark"].(string)
	}

	if param["operatorid"] != nil {
		v.Operatorid = utils.Convert2Int64(param["operatorid"])
	}

	o := orm.NewOrm()

	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	id, err = o.Insert(&v)
	if err != nil {
		logger.Error("不能插入考勤信息：%v", err.Error())
		return 0, err
	}

	return id, err
}

func PutAttendance(ctx context.Context,siteid, appid, id int64, token string, param map[string]interface{}) (cnt int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v NcAttendance

	if param["id"] != nil {
		v.AttendanceId = utils.Convert2Int64(param["id"])
	}

	if id != v.AttendanceId {
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

	if param["staffid"] != nil {
		v.StaffId = utils.Convert2Int64(param["staffid"])
	}

	if param["clocktime"] != nil {
		v.Clocktime,_ = time.Parse("2006-01-02 15:04:05",param["clocktime"].(string))
	}
	if param["clocktimestart"] != nil {
		v.ClockTimeStart,_ = time.Parse("2006-01-02 15:04:05",param["clocktimestart"].(string))
	}
	if param["clocktimeend"] != nil {
		v.ClockTimeEnd,_ = time.Parse("2006-01-02 15:04:05",param["clocktimeend"].(string))
	}

	//通用字段
	if param["remark"] != nil {
		v.Remark = param["remark"].(string)
	}
	if param["operatorid"] != nil {
		v.Operatorid = utils.Convert2Int64(param["operatorid"])
	}

	//o := orm.NewOrm()

	cnt, err = o.Update(&v)
	if err != nil {
		logger.Error("不能更新考勤信息：%v", err.Error())
		return 0, err
	}

	return cnt, nil
}

func DeleteAttendance(ctx context.Context,id int64, param map[string]interface{}) (num int64, err error) {

	o := orm.NewOrm()
	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	var v NcAttendance

	var ids []interface{}

	if param["id"] != nil {
		ids = param["id"].([]interface{})
	}

	var cnt int64 = 0

	for _, rid := range ids {
		v.AttendanceId = utils.Convert2Int64(rid)
		num, err = o.Delete(&v)
		cnt += num
	}

	return cnt, err
}
