//
package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/smallnest/rpcx/client"
	"platform/common/utils"
	"platform/mskit/trace"
	"platform/pfcomm/apis"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"platform/hrsvc/imconf"
	md "platform/pfcomm/models"
)

func init() {
	orm.RegisterModel(new(NcStaff))
	orm.RegisterModel(new(md.Role))
}


func GetStaffLists(ctx context.Context,siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {
	if num <= 0 {
		num=PAGENUM_MAX
	}

	return getStaffListCount(ctx,1,siteid,appid,stypes,contents,order,sort,num,start)
}
func GetStaffCount(ctx context.Context,siteid int64, appid int64, stypes, contents []string, order, sort string) (int64, error) {
	_,cnt,err := getStaffListCount(ctx,2,siteid,appid,stypes,contents,order,sort,1,0)

	return cnt,err
}

/*
	搜索条件:
	stype = 1: 姓名
	stype = 2: 员工类型 ： 0:临时，1：正式员工，2、义工，20：护理人员，21：医护人员，30：外包人员
	stype = 3: 入职时间大于某个时间
	stype = 4：入职时间小于某个时间
	stype = 5: 离职时间大于某个时间
	stype = 6: 离职时间小于某个时间
	stype = 7: 登录账户名
	stype = 8: 注册电话
	stype = 9: 员工工号
	stype = 10: 注册电子邮件
	stype = 11: 部门标识
	stype = 12: 部门名称
	stype = 13: 模糊搜索（姓名、账户、电话，身份证）
	stype = 14: 状态 （1：有效，0：无效）
	stype = 15: 取得资格证书日期（开始日期）
	stype = 16: 取得资格证书日期（结束日期）
	stype = 17: 创建时间（开始日期）
	stype = 18: 创建时间（结束日期）
	stype = 19: 机构名称查询
	stype = 20: 机构标识
*/

func getStaffListCount(ctx context.Context,cate int,siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {

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
			v = " a.name like " + v
		case "2":
			v = " a.staff_type =" + contents[i]
		case "3":
			v = " a.entry_time >= '" + contents[i] + "'"
		case "4":
			v = " a.entry_time < '" + contents[i] + "'"
		case "5":
			v = " a.leave_time >= '" + contents[i] + "'"
		case "6":
			v = " a.leave_time < '" + contents[i] + "'"
		case "7":
			v = " b.name like " + v
		case "8":
			v = " b.phone like " + v
		case "9":
			v = " a.employee_code like " + v
		case "10":
			v = " b.email like " + v
		case "11":
			v = " d.department_id = " + contents[i]
		case "12":
			v = " e.name like " + v
		case "13":
			v = "( a.name like " + v + " or a.papers like " + v + " or a.phone like " + v + " or b.name like " + v + " )"
		case "14":
			v = " a.status =" + contents[i]
		case "15":
			v = " a.qual_cert_date >= '" + contents[i] + "'"
		case "16":
			v = " a.qual_cert_date < '" + contents[i] + "'"
		case "17":
			v = " a.create_time >= '" + contents[i] + "'"
		case "18":
			v = " a.create_time < '" + contents[i] + "'"
		case "19":
			v = " g.org_name like " + v
		case "20":
			v = " a.organization_id = " + contents[i]
		case "21":
			v = " a.papers = " + contents[i]
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
		statement = fmt.Sprintf(SQL_STAFF,imconf.Config.ObjectsdbName, imconf.Config.ObjectsdbName)
	}else{
		statement = fmt.Sprintf(SQL_COUNT_STAFF,imconf.Config.ObjectsdbName, imconf.Config.ObjectsdbName)
	}

	if siteid > 1 {
		statement = fmt.Sprintf(statement+
			"where a.siteid=%d %s order by a.staff_id %s limit ? offset ?", siteid, conditions, sort)
	} else {
		statement = fmt.Sprintf(statement+
			"where 1=1 %s order by a.staff_id %s limit ? offset ?", conditions, sort)

	}
	cnt, err = o.Raw(statement, num, start).Values(&vs)

	if err != nil {
		logger.Error("不能获取员工列表：%v", err.Error())
		return nil, 0, err
	}

	if cate == 2 {
		if cnt > 0 {
			cnt = utils.Convert2Int64(vs[0]["ucount"])
		}
	}

	return vs, cnt, nil
}

func GetStaffGroupLists(ctx context.Context,siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {

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
			v = " a.staff_type =" + contents[i]
		case "3":
			v = " a.entry_time >= '" + contents[i] + "'"
		case "4":
			v = " a.entry_time < '" + contents[i] + "'"
		case "5":
			v = " a.leave_time >= '" + contents[i] + "'"
		case "6":
			v = " a.leave_time < '" + contents[i] + "'"
		case "7":
			v = " b.name like " + v
		case "8":
			v = " b.phone like " + v
		case "9":
			v = " b.job_number like " + v
		case "10":
			v = " b.email like " + v
		case "11":
			v = " d.department_id = " + contents[i]
		case "12":
			v = " e.name like " + v
		case "13":
			v = "( a.name like " + v + " or a.papers like " + v + " or a.phone like " + v + " or b.name like " + v + " )"
		case "14":
			v = " a.status =" + contents[i]
		case "15":
			v = " ifnull(gs.group_id,0) =" + contents[i]

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

	statement = fmt.Sprintf(SQL_STAFF_GROUP,

		imconf.Config.ObjectsdbName, imconf.Config.ObjectsdbName)

	if siteid > 1 {
		statement = fmt.Sprintf(statement+
			"where a.siteid=%d %s order by a.staff_id %s limit ? offset ?", siteid, conditions, sort)
	} else {
		statement = fmt.Sprintf(statement+
			"where 1=1 %s order by a.staff_id %s limit ? offset ?", conditions, sort)

	}
	cnt, err = o.Raw(statement, start, num).Values(&vs)

	if err != nil {
		logger.Error("不能获取员工列表：%v", err.Error())
		return nil, 0, err
	}

	return vs, cnt, nil
}

func GetStaff(ctx context.Context,id int64, siteid int64, appid int64) (interface{}, int64, error) {

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

	statement := fmt.Sprintf(SQL_STAFF,
		imconf.Config.ObjectsdbName, imconf.Config.ObjectsdbName)

	if siteid > 1 {
		cnt, err = o.Raw(statement+
			"where a.siteid =?  and a.staff_id =?", siteid, id).Values(&vs)
	} else {
		cnt, err = o.Raw(statement+
			"where a.staff_id =?", id).Values(&vs)
	}

	if err != nil {
		logger.Error("不能获取员工信息：%v", err.Error())
		return nil, 0, err
	}

	return vs, cnt, nil
}

func PostStaff(ctx context.Context,tracer trace.Tracer,siteid, appid int64, token string, param map[string]interface{}) (id int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v NcStaff

	if param["id"] != nil {
		v.StaffId = utils.Convert2Int64(param["id"])
	}

	if param["siteid"] != nil {
		v.Siteid = utils.Convert2Int64(param["siteid"])
	} else {
		v.Siteid = siteid
	}

	if param["effectivetime"] != nil {
		v.EntryTime,_ = time.Parse("2006-01-02 15:04:05",param["effectivetime"].(string))
	}else{
		v.EntryTime = time.Now()
	}
	if param["expiretime"] != nil {
		v.LeaveTime,_ = time.Parse("2006-01-02 15:04:05",param["expiretime"].(string))
	}

	o := orm.NewOrm()
	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	var role md.Role
	if param["roleid"] != nil {
		role.Id = utils.Convert2Int64(param["roleid"])
		o.Using("objects")
		err = SetSearchPath(o, imconf.Config.ObjectsdbSchema)
		if err != nil {
			logger.Error("不能设置search_path :%v", err)
			return 0, err
		}
		err = o.Read(&role)
		if err != nil && err != orm.ErrNoRows {
			logger.Error("不能获取角色信息: %v", err)
			return 0, err
		}
		param["effectivetime"] = role.EffectiveTime
		param["expiretime"] = role.ExpireTime

		o.Using("default")
		err = SetSearchPath(o, imconf.Config.StaffdbSchema)
		if err != nil {
			logger.Error("不能设置search_path :%v", err)
			return 0, err
		}
	}

	if param["name"] != nil {
		v.Name = param["name"].(string)
	}
	if param["papers"] != nil {
		v.Papers = param["papers"].(string)
	}
	if param["stafftype"] != nil {
		v.StaffType = utils.Convert2Int16(param["stafftype"])
	}
	if param["papertype"] != nil {
		v.PaperType = utils.Convert2Int16(param["papertype"])
	}

	if param["phone"] != nil {
		v.Phone = utils.ConvertToString(param["phone"])
	}
	if param["address"] != nil {
		v.Address = utils.ConvertToString(param["address"])
	}

	if param["signature"] != nil {
		v.Signature = utils.ConvertToString(param["signature"])
	}

	//2018-12-01增加
	if param["expertflag"] != nil {
		v.ExpertFlag = utils.Convert2Int16(param["expertflag"])
	}
	if param["employeecode"] != nil {
		v.EmployeeCode = utils.ConvertToString(param["employeecode"])
	}
	if param["gender"] != nil {
		v.Gender = utils.Convert2Int16(param["gender"])
	}
	if param["nationality"] != nil {
		v.Nationality = utils.ConvertToString(param["nationality"])
	}
	if param["dutiescode"] != nil {
		v.DutiesCode = utils.ConvertToString(param["dutiescode"])
	}
	if param["joblevelcode"] != nil {
		v.JobLevelCode = utils.ConvertToString(param["joblevelcode"])
	}
	if param["birthday"] != nil {
		v.Birthday,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["birthday"]))
	}
	if param["archiveid"] != nil {
		v.ArchiveId = utils.Convert2Int64(param["archiveid"])
	}
	if param["archivebarcode"] != nil {
		v.ArchiveBarcode = utils.ConvertToString(param["archivebarcode"])
	}
	if param["archiveaddress"] != nil {
		v.ArchiveAddress = utils.ConvertToString(param["archiveaddress"])
	}
	if param["formername"] != nil {
		v.FormerName = utils.ConvertToString(param["formername"])
	}
	if param["enname"] != nil {
		v.EnglishName = utils.ConvertToString(param["enname"])
	}
	if param["nativeplace"] != nil {
		v.NativePlace = utils.ConvertToString(param["nativeplace"])
	}
	if param["politicalstatuscode"] != nil {
		v.PoliticalStatusCode = utils.ConvertToString(param["politicalstatuscode"])
	}
	if param["marriagestatuscode"] != nil {
		v.MarriageStatusCode = utils.ConvertToString(param["marriagestatuscode"])
	}
	if param["healthstatuscode"] != nil {
		v.HealthStatusCode = utils.ConvertToString(param["healthstatuscode"])
	}
	if param["expertclasscode"] != nil {
		v.ExpertClassCode = utils.ConvertToString(param["expertclasscode"])
	}
	if param["expertallawance"] != nil {
		v.ExpertAllawance = utils.ConvertToString(param["expertallawance"])
	}
	if param["onjobstatus"] != nil {
		v.OnJobStatus = utils.Convert2Int16(param["onjobstatus"])
	}
	if param["directordeptid"] != nil {
		v.DirectorDeptId = utils.Convert2Int64(param["directordeptid"])
	}
	if param["worktel"] != nil {
		v.WorkTel = utils.ConvertToString(param["worktel"])
	}
	if param["hometel"] != nil {
		v.HomeTel = utils.ConvertToString(param["hometel"])
	}
	if param["postclass"] != nil {
		v.PostClass = utils.ConvertToString(param["postclass"])
	}
	if param["workmajor"] != nil {
		v.WorkMajor = utils.ConvertToString(param["workmajor"])
	}
	if param["majorpost"] != nil {
		v.MajorPost = utils.ConvertToString(param["majorpost"])
	}
	if param["recruitresource"] != nil {
		v.RecruitResource = utils.ConvertToString(param["recruitresource"])
	}
	if param["employmentstarttime"] != nil {
		v.EmploymentStartTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["employmentstarttime"]))
	}
	if param["employmentendtime"] != nil {
		v.EmploymentEndTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["employmentendtime"]))
	}
	if param["joinarmyflag"] != nil {
		v.JoinArmyFlag = utils.Convert2Int16(param["joinarmyflag"])
	}
	if param["joinarmytime"] != nil {
		v.JoinArmyTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["joinarmytime"]))
	}
	if param["workagestarttime"] != nil {
		v.WorkageStartTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["workagestarttime"]))
	}
	if param["workageendremark"] != nil {
		v.WorkageEndRemark = utils.ConvertToString(param["workageendremark"])
	}
	if param["workage"] != nil {
		v.Workage = utils.Convert2Int16(param["workage"])
	}
	if param["postcertid"] != nil {
		v.PostCertId = utils.ConvertToString(param["postcertid"])
	}
	if param["educationcertid"] != nil {
		v.EducationCertId = utils.ConvertToString(param["educationcertid"])
	}
	if param["academicrecord"] != nil {
		v.AcademicRecord = utils.ConvertToString(param["academicrecord"])
	}
	if param["degreen"] != nil {
		v.Degreen = utils.ConvertToString(param["degreen"])
	}
	if param["firstlanguage"] != nil {
		v.FirstLanguage = utils.ConvertToString(param["firstlanguage"])
	}
	if param["firstlanguagelevel"] != nil {
		v.FirstLanguageLevel = utils.ConvertToString(param["firstlanguagelevel"])
	}
	if param["secondlanguage"] != nil {
		v.SecondLanguage = utils.ConvertToString(param["secondlanguage"])
	}
	if param["secondlanguagelevel"] != nil {
		v.SecondLanguageLevel = utils.ConvertToString(param["secondlanguagelevel"])
	}
	if param["joinworktime"] != nil {
		v.JoinWorkTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["joinworktime"]))
	}
	if param["specialty"] != nil {
		v.Specialty = utils.ConvertToString(param["specialty"])
	}

	if param["qualcertname"] != nil {
		v.QualCertName = utils.ConvertToString(param["qualcertname"])
	}
	if param["qualcertcode"] != nil {
		v.QualCertCode = utils.ConvertToString(param["qualcertcode"])
	}
	if param["qualcertdate"] != nil {
		v.QualCertDate,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["qualcertdate"]))
	}

	if param["academy"] != nil {
		v.Academy = utils.ConvertToString(param["academy"])
	}
	if param["nation"] != nil {
		v.Nation = utils.ConvertToString(param["nation"])
	}
	if param["academyspecialty"] != nil {
		v.AcademySpecialty = utils.ConvertToString(param["academyspecialty"])
	}
	if param["graduatedate"] != nil {
		v.GraduateDate,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["graduatedate"]))
	}
	if param["joinpartydate"] != nil {
		v.JoinPartyDate,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["joinpartydate"]))
	}
	if param["specialtylevel"] != nil {
		v.SpecialtyLevel = utils.Convert2Int16(param["specialtylevel"])
	}

	if param["organizationid"] != nil {
		v.OrganizationId = utils.Convert2Int64(param["organizationid"])
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

	v.CreateTime,_ = time.Parse("2006-01-02 15:04:05",utils.GetTimeFormat("2006-01-02 15:04:05"))

	var deptid int64

	if param["departmentid"] != nil {
		deptid = utils.Convert2Int64(param["departmentid"])
	}

	//o := orm.NewOrm()
	o.Begin()

	// 查询身份证是否重复
	var vvs []orm.Params
	cnt, err := o.Raw("select staff_id from nc_staff where papers = ? ", v.Papers).Values(&vvs)
	if err == nil && cnt > 0 {
		o.Rollback()
		return 0, errors.New("身份证ID重复。")
	}

	sid, err := insertAdmin(ctx,tracer,siteid, appid, token, param)
	if err != nil {
		o.Rollback()
		if strings.Contains(err.Error(),"1062") && strings.Contains(err.Error(),"siteid_email") {
			err = errors.New("该邮箱已经被占用,请选择未使用的邮箱。")
		}
		return 0, err
	}

	v.StaffId = sid
	_, err = o.Insert(&v)
	if err != nil {
		o.Rollback()
		logger.Error("不能插入员工信息：%v", err.Error())
		return 0, err
	}

	if deptid > 0 {
		_, err = o.QueryTable("nc_department_staff").
			Filter("staff_id", v.StaffId).
			Update(orm.Params{"status": 0, "injob_flag": 0, "departure_time": v.CreateTime})

		if err != nil {
			o.Rollback()
			logger.Error("不能更新员工部门信息：%v", err.Error())
			return 0, err
		}
		var df NcDepartmentStaff
		df.Siteid = v.Siteid
		df.StaffId = sid
		df.DepartmentId = deptid
		df.InductionTime = v.CreateTime
		df.Status = 1
		df.InjobFlag = 1
		df.CreateTime = v.CreateTime

		_, err = o.Insert(&df)

		if err != nil {
			o.Rollback()
			logger.Error("不能更新员工部门信息：%v", err.Error())
			return 0, err
		}
	}
	o.Commit()

	return sid, err
}

func PutStaff(ctx context.Context,tracer trace.Tracer, siteid, appid, id int64, token string, param map[string]interface{}) (cnt int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v NcStaff

	if param["id"] != nil {
		v.StaffId = utils.Convert2Int64(param["id"])
	}

	if id != v.StaffId {
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

	if param["effectivetime"] != nil {
		v.EntryTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["effectivetime"]))
	}else{
		v.EntryTime,_ = time.Parse("2006-01-02 15:04:05",v.EntryTime.Format("2006-01-02 15:04:05"))
	}
	if param["expiretime"] != nil {
		v.LeaveTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["expiretime"]))
	}else{
		v.LeaveTime,_ = time.Parse("2006-01-02 15:04:05",v.LeaveTime.Format("2006-01-02 15:04:05"))
	}

	var role md.Role
	if param["roleid"] != nil {
		role.Id = utils.Convert2Int64(param["roleid"])
		o.Using("objects")
		err = SetSearchPath(o, imconf.Config.ObjectsdbSchema)
		if err != nil {
			logger.Error("不能设置search_path :%v", err)
			return 0, err
		}
		err = o.Read(&role)
		if err != nil && err != orm.ErrNoRows {
			logger.Error("不能获取角色信息: %v", err)
			return 0, err
		}
		param["effectivetime"] = role.EffectiveTime
		param["expiretime"] = role.ExpireTime

		o.Using("default")
		err = SetSearchPath(o, imconf.Config.StaffdbSchema)
		if err != nil {
			logger.Error("不能设置search_path :%v", err)
			return 0, err
		}

	}

	if param["siteid"] != nil {
		v.Siteid = utils.Convert2Int64(param["siteid"])
	}

	if param["name"] != nil {
		v.Name = param["name"].(string)
	}
	if param["papers"] != nil {
		v.Papers = param["papers"].(string)
	}

	if param["stafftype"] != nil {
		v.StaffType = utils.Convert2Int16(param["stafftype"])
	}
	if param["papertype"] != nil {
		v.PaperType = utils.Convert2Int16(param["papertype"])
	}

	if param["phone"] != nil {
		v.Phone = utils.ConvertToString(param["phone"])
	}
	if param["address"] != nil {
		v.Address = utils.ConvertToString(param["address"])
	}

	if param["signature"] != nil {
		v.Signature = utils.ConvertToString(param["signature"])
	}

	//2018-12-01增加
	if param["expertflag"] != nil {
		v.ExpertFlag = utils.Convert2Int16(param["expertflag"])
	}
	if param["employeecode"] != nil {
		v.EmployeeCode = utils.ConvertToString(param["employeecode"])
	}
	if param["gender"] != nil {
		v.Gender = utils.Convert2Int16(param["gender"])
	}
	if param["nationality"] != nil {
		v.Nationality = utils.ConvertToString(param["nationality"])
	}
	if param["dutiescode"] != nil {
		v.DutiesCode = utils.ConvertToString(param["dutiescode"])
	}
	if param["joblevelcode"] != nil {
		v.JobLevelCode = utils.ConvertToString(param["joblevelcode"])
	}
	if param["birthday"] != nil {
		v.Birthday,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["birthday"]))
	}else{
		v.Birthday,_ = time.Parse("2006-01-02 15:04:05",v.Birthday.Format("2006-01-02 15:04:05"))
	}
	if param["archiveid"] != nil {
		v.ArchiveId = utils.Convert2Int64(param["archiveid"])
	}
	if param["archivebarcode"] != nil {
		v.ArchiveBarcode = utils.ConvertToString(param["archivebarcode"])
	}
	if param["archiveaddress"] != nil {
		v.ArchiveAddress = utils.ConvertToString(param["archiveaddress"])
	}
	if param["formername"] != nil {
		v.FormerName = utils.ConvertToString(param["formername"])
	}
	if param["enname"] != nil {
		v.EnglishName = utils.ConvertToString(param["enname"])
	}
	if param["nativeplace"] != nil {
		v.NativePlace = utils.ConvertToString(param["nativeplace"])
	}
	if param["politicalstatuscode"] != nil {
		v.PoliticalStatusCode = utils.ConvertToString(param["politicalstatuscode"])
	}
	if param["marriagestatuscode"] != nil {
		v.MarriageStatusCode = utils.ConvertToString(param["marriagestatuscode"])
	}
	if param["healthstatuscode"] != nil {
		v.HealthStatusCode = utils.ConvertToString(param["healthstatuscode"])
	}
	if param["expertclasscode"] != nil {
		v.ExpertClassCode = utils.ConvertToString(param["expertclasscode"])
	}
	if param["expertallawance"] != nil {
		v.ExpertAllawance = utils.ConvertToString(param["expertallawance"])
	}
	if param["onjobstatus"] != nil {
		v.OnJobStatus = utils.Convert2Int16(param["onjobstatus"])
	}
	if param["directordeptid"] != nil {
		v.DirectorDeptId = utils.Convert2Int64(param["directordeptid"])
	}
	if param["worktel"] != nil {
		v.WorkTel = utils.ConvertToString(param["worktel"])
	}
	if param["hometel"] != nil {
		v.HomeTel = utils.ConvertToString(param["hometel"])
	}
	if param["postclass"] != nil {
		v.PostClass = utils.ConvertToString(param["postclass"])
	}
	if param["workmajor"] != nil {
		v.WorkMajor = utils.ConvertToString(param["workmajor"])
	}
	if param["majorpost"] != nil {
		v.MajorPost = utils.ConvertToString(param["majorpost"])
	}
	if param["recruitresource"] != nil {
		v.RecruitResource = utils.ConvertToString(param["recruitresource"])
	}
	if param["employmentstarttime"] != nil {
		v.EmploymentStartTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["employmentstarttime"]))
	}else{
		v.EmploymentStartTime,_ = time.Parse("2006-01-02 15:04:05",v.EmploymentStartTime.Format("2006-01-02 15:04:05"))
	}

	if param["employmentendtime"] != nil {
		v.EmploymentEndTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["employmentendtime"]))
	}else{
		v.EmploymentEndTime,_ = time.Parse("2006-01-02 15:04:05",v.EmploymentEndTime.Format("2006-01-02 15:04:05"))
	}
	if param["joinarmyflag"] != nil {
		v.JoinArmyFlag = utils.Convert2Int16(param["joinarmyflag"])
	}
	if param["joinarmytime"] != nil {
		v.JoinArmyTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["joinarmytime"]))
	}else{
		v.JoinArmyTime,_ = time.Parse("2006-01-02 15:04:05",v.JoinArmyTime.Format("2006-01-02 15:04:05"))
	}
	if param["workagestarttime"] != nil {
		v.WorkageStartTime,_= time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["workagestarttime"]))
	}else{
		v.WorkageStartTime,_ = time.Parse("2006-01-02 15:04:05",v.WorkageStartTime.Format("2006-01-02 15:04:05"))
	}
	if param["workageendremark"] != nil {
		v.WorkageEndRemark = utils.ConvertToString(param["workageendremark"])
	}
	if param["workage"] != nil {
		v.Workage = utils.Convert2Int16(param["workage"])
	}
	if param["postcertid"] != nil {
		v.PostCertId = utils.ConvertToString(param["postcertid"])
	}
	if param["educationcertid"] != nil {
		v.EducationCertId = utils.ConvertToString(param["educationcertid"])
	}
	if param["academicrecord"] != nil {
		v.AcademicRecord = utils.ConvertToString(param["academicrecord"])
	}
	if param["degreen"] != nil {
		v.Degreen = utils.ConvertToString(param["degreen"])
	}
	if param["firstlanguage"] != nil {
		v.FirstLanguage = utils.ConvertToString(param["firstlanguage"])
	}
	if param["firstlanguagelevel"] != nil {
		v.FirstLanguageLevel = utils.ConvertToString(param["firstlanguagelevel"])
	}
	if param["secondlanguage"] != nil {
		v.SecondLanguage = utils.ConvertToString(param["secondlanguage"])
	}
	if param["secondlanguagelevel"] != nil {
		v.SecondLanguageLevel = utils.ConvertToString(param["secondlanguagelevel"])
	}
	if param["joinworktime"] != nil {
		v.JoinWorkTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["joinworktime"]))
	}else{
		v.JoinWorkTime,_ = time.Parse("2006-01-02 15:04:05",v.JoinWorkTime.Format("2006-01-02 15:04:05"))
	}
	if param["specialty"] != nil {
		v.Specialty = utils.ConvertToString(param["specialty"])
	}
	if param["qualcertname"] != nil {
		v.QualCertName = utils.ConvertToString(param["qualcertname"])
	}
	if param["qualcertcode"] != nil {
		v.QualCertCode = utils.ConvertToString(param["qualcertcode"])
	}
	if param["qualcertdate"] != nil {
		v.QualCertDate,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["qualcertdate"]))
	}else{
		v.QualCertDate,_ = time.Parse("2006-01-02 15:04:05",v.QualCertDate.Format("2006-01-02 15:04:05"))
	}
	if param["academy"] != nil {
		v.Academy = utils.ConvertToString(param["academy"])
	}
	if param["nation"] != nil {
		v.Nation = utils.ConvertToString(param["nation"])
	}
	if param["academyspecialty"] != nil {
		v.AcademySpecialty = utils.ConvertToString(param["academyspecialty"])
	}
	if param["graduatedate"] != nil {
		v.GraduateDate,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["graduatedate"]))
	}else{
		v.GraduateDate,_ = time.Parse("2006-01-02 15:04:05",v.GraduateDate.Format("2006-01-02 15:04:05"))
	}
	if param["joinpartydate"] != nil {
		v.JoinPartyDate,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["joinpartydate"]))
	}else{
		v.JoinPartyDate,_ = time.Parse("2006-01-02 15:04:05",v.JoinPartyDate.Format("2006-01-02 15:04:05"))
	}
	if param["specialtylevel"] != nil {
		v.SpecialtyLevel = utils.Convert2Int16(param["specialtylevel"])
	}
	if param["organizationid"] != nil {
		v.OrganizationId = utils.Convert2Int64(param["organizationid"])
	}


	v.CreateTime,_ = time.Parse("2006-01-02 15:04:05",v.CreateTime.Format("2006-01-02 15:04:05"))

	//通用字段
	if param["remark"] != nil {
		v.Remark = param["remark"].(string)
	}

	if param["status"] != nil {
		v.Status = utils.Convert2Int16(param["status"])
		if v.Status == 0 {
			param["state"] = 1
		} else if v.Status == 1 {
			param["state"] = 0
		}
	}
	if param["operatorid"] != nil {
		v.Operatorid = utils.Convert2Int64(param["operatorid"])
	}

	var deptid int64

	if param["departmentid"] != nil {
		deptid = utils.Convert2Int64(param["departmentid"])
	}

	o.Begin()

	// 查询身份证是否重复
	var vvs []orm.Params
	cnt, err = o.Raw("select staff_id from nc_staff where papers = ? and staff_id <> ?", v.Papers, v.StaffId).Values(&vvs)
	if err == nil && cnt > 0 {
		o.Rollback()
		return 0, errors.New("身份证ID重复。")
	}


	_, err = updateAdmin(ctx,tracer,siteid, appid, token, param)
	if err != nil {
		logger.Error("不能更新员工信息(admin) : err = %v",err.Error())
		o.Rollback()
		if strings.Contains(err.Error(),"1062") && strings.Contains(err.Error(),"siteid_email") {
			err = errors.New("该邮箱已经被占用,请选择未使用的邮箱。")
		}

		if strings.Contains(err.Error(),"1062") && strings.Contains(err.Error(),"siteid_phone") {
			err = errors.New("电话号码重复。")
		}

		return 0, err
	}

	cnt, err = o.Update(&v)
	if err != nil {
		o.Rollback()
		logger.Error("不能更新员工信息：%v", err.Error())
		return 0, err
	}

	if deptid > 0 {
		_, err = o.QueryTable("nc_department_staff").
			Filter("staff_id", v.StaffId).
			Update(orm.Params{"status": 0, "injob_flag": 0, "departure_time": v.CreateTime})

		if err != nil {
			o.Rollback()
			logger.Error("不能更新员工部门信息：%v", err.Error())
			return 0, err
		}
		var df NcDepartmentStaff
		df.Siteid = v.Siteid
		df.StaffId = v.StaffId
		df.DepartmentId = deptid
		df.InductionTime = v.CreateTime
		df.Status = 1
		df.InjobFlag = 1
		df.CreateTime = v.CreateTime

		_, err = o.Insert(&df)

		if err != nil {
			o.Rollback()
			logger.Error("不能更新员工部门信息：%v", err.Error())
			return 0, err
		}
	}

	o.Commit()

	return cnt, nil
}

func DeleteStaff(ctx context.Context,id int64, param map[string]interface{}) (num int64, err error) {

	o := orm.NewOrm()
	err = SetSearchPath(o, imconf.Config.StaffdbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	var v NcStaff

	var ids []interface{}

	if param["id"] != nil {
		ids = param["id"].([]interface{})
	}

	var cnt int64 = 0

	for _, rid := range ids {
		v.StaffId = utils.Convert2Int64(rid)
		num, err = o.Delete(&v)
		cnt += num
	}

	return cnt, err
}

func insertAdmin(ctx context.Context,tracer trace.Tracer,siteid, appid int64, token string, params map[string]interface{}) (id int64, err error) {
	defer func() {
		if e := recover(); e != nil {
			logger.Error("%v", e)
			err = errors.New("panic error.")
		}
	}()

	params["siteid"] = params["siteid"]

	//id,err = apis.RpcxAddService(appid,siteid,token,"OaJSONRpc","AddAdmin",imconf.Config.RpcxOaBasepath,imconf.Config.ConsulAddress,params)
	id,err = apis.RpcxAdd(ctx,tracer,appid,siteid,token,imconf.Config.Sdt,imconf.Config.Sda,"OaJSONRpc","AddAdmin",imconf.Config.RpcxOaBasepath,client.Failtry,client.RoundRobin,params)
	return id, err
}

func updateAdmin(ctx context.Context,tracer trace.Tracer,siteid, appid int64, token string, params map[string]interface{}) (cnt int64, err error) {
	defer func() {
		if e := recover(); e != nil {
			logger.Error("%v", e)
			err = errors.New("panic error.")
		}
	}()

	params["siteid"] = params["orgid"]

	cnt,err = apis.RpcxUpdate(ctx,tracer,appid,siteid,token,imconf.Config.Sdt,imconf.Config.Sda,"OaJSONRpc","UpdateAdmin",imconf.Config.RpcxOaBasepath,client.Failtry,client.RoundRobin,params)


	return cnt, err
}
