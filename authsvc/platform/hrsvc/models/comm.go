package models

import (
	l4g "github.com/libra9z/log4go"
	"platform/hrsvc/common"

	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"platform/hrsvc/imconf"
)

var logger l4g.Logger

func InitLogger() {
	logger = common.Logger
}

const (
	PAGENUM_MAX = 1000

	PLAN_TYPE_WEEK  = 1
	PLAN_TYPE_MONTH = 2
)

var SQL_STAFF_GROUP = "select a.staff_id as id,a.org_id as orgid,a.name,a.staff_type as stafftype, " +
	"a.paper_type as papertype,a.papers,a.phone,a.address," +
	"a.entry_time as effectivetime,a.leave_time as expiretime,b.passwd,b.image_url as imageurl," +
	"a.create_time as createtime,a.status,a.remark,a.operatorid,a.signature, " +
	"b.name as username,b.email,b.job_number as jobnumber,b.role_id as roleid, " +
	"c.name as rolename,c.expire_time as roleexpire, " +
	"d.injob_flag as injobflag,d.department_id as departmentid,e.name as deptname, " +
	"ifnull(gs.group_id,0) as groupid,g.name as groupname " +
	"from nc_staff a " +
	"left join %s.admin b on a.staff_id=b.id " +
	"left join %s.role c on b.role_id=c.id " +
	"left join nc_department_staff d on a.staff_id=d.staff_id and d.status=1 " +
	"left join nc_department e on d.department_id=e.department_id " +
	"left join nc_group_staff gs on a.staff_id=gs.staff_id and gs.status=1 " +
	"left join nc_group g on gs.group_id=g.group_id "

var SQL_WEEKLY_SCHEDULE = "select a.week_sched_id as id,a.siteid,a.organization_id as organizationid,a.group_id as groupid,a.staff_id as staffid," +
	"a.sched_type as schedtype,a.week_start_date as starttime,a.week_end_date as endtime, " +
	"a.monday,a.tuesday,a.wednesday,a.thursday,a.friday,a.saturday,a.sunday,a.sched_plan_id as schedplanid, " +
	"a.create_time as createtime,a.status,a.remark,a.operatorid, " +
	"b.name as staffname,b.phone,c.name as groupname " +
	"from nc_week_schedules a " +
	"left join nc_staff b on a.staff_id=b.staff_id " +
	"left join nc_group c on a.group_id=c.group_id "

var SQL_COUNT_WEEKLY_SCHEDULE = "select count(*) as ucount " +
	"from nc_week_schedules a " +
	"left join nc_staff b on a.staff_id=b.staff_id " +
	"left join nc_group c on a.group_id=c.group_id "

var SQL_NcOrganization = "select a.org_type as orgtype,a.approval_time as approvaltime,a.branches_num as branchesnum,a.management_area as managementarea,a.org_level as orglevel,a.phone as phone,a.acreage as acreage,a.org_ln as orgln,a.organization_id as id,a.creatorid as creatorid,a.create_time as createtime,a.reg_div_id as regdivid,a.org_code as orgcode,a.membership_level as membershiplevel,a.status as status,a.updaterid as updaterid,a.update_time as updatetime,a.approval_department as approvaldepartment,a.address as address,a.remark as remark,a.reg_code as regcode,a.parent_id as parentid,a.approval_bed_num as approvalbednum,a.legal_person as legalperson, " +
	"a.nature,a.effective_date as effdate,a.expire_date as expdate,a.org_name as orgname,a.manage_cate as managecate,a.reg_name as regname "+
	"from nc_organization a "

var SQL_COUNT_NcOrganization = "select count(*) as ucount " +
	"from nc_organization a "

var SQL_DEPARTMENT_LIST = "select a.department_id as id,a.siteid,a.organization_id as organizationid,a.name,a.parent_id as parentid,a.level,a.create_time as createtime,a.status,a.remark,a.creatorid,a.updaterid,a.update_time as updatetime " +
	"from nc_department a "
var SQL_DEPARTMENT_TREE = "select a.department_id as id,a.siteid,a.organization_id as organizationid,a.name,a.parent_id as parentid,a.level,a.create_time as createtime,a.status,a.remark,a.creatorid,a.updaterid,a.update_time as updatetime " +
	"from nc_department a "

var SQL_COUNT_DEPARTMENT_LIST = "select count(*) as ucount " +
	"from nc_department a "

var SQL_STAFF = "select a.staff_id as id,a.siteid,a.organization_id as organizationid,a.name,a.staff_type as stafftype, " +
	"a.paper_type as papertype,a.papers,a.phone,a.address," +
	"a.entry_time as effectivetime,a.leave_time as expiretime,b.passwd,b.image_url as imageurl," +
	"a.create_time as createtime,a.status,a.remark,a.operatorid,a.signature, " +
	"a.expert_flag as expertflag,a.employee_code as employeecode, a.gender,a.nationality,a.duties_code as dutiescode,a.job_level_code as joblevelcode,a.birthday, " +
	"a.archive_id as archiveid,a.archive_barcode as archivebarcode,a.archive_address as archiveaddress,a.former_name as formername,a.english_name as enname,a.native_place as nativeplace," +
	"a.political_status_code as politicalstatuscode,a.marriage_status_code as marriagestatuscode, a.health_status_code as healthstatuscode,a.expert_class_code as expertclasscode," +
	"a.expert_allawance as expertallawance, a.on_job_status as onjobstatus,a.director_dept_id as directordeptid, a.work_tel as worktel,a.home_tel as hometel,a.post_class as postclass," +
	"a.work_major as workmajor,a.major_post as majorpost,a.recruit_resource as recruitresource,a.employment_start_time as employmentstarttime, a.employment_end_time as employmentendtime," +
	"a.join_army_flag as joinarmyflag,a.join_army_time as joinarmytime,a.workage_start_time as workagestarttime,  a.workage_end_remark as workageendremark,a.workage,a.post_cert_id as postcertid," +
	"a.education_cert_id as educationcertid,a.academic_record as academicrecord,a.degreen,a.first_language as firstlanguage,a.first_language_level as firstlanguagelevel, " +
	"a.second_language as secondlanguage,a.second_language_level as secondlanguagelevel,a.join_work_time as joinworktime,a.specialty, " +
	"a.qual_cert_name as qualcertname,a.qual_cert_code as qualcertcode,a.qual_cert_date as qualcertdate,a.nation,a.academy, " +
	"a.academy_specialty as academyspecialty,a.graduate_date as graduatedate,a.join_party_date as joinpartydate,a.specialty_level as specialtylevel," +
	"b.name as username,b.email,b.job_number as jobnumber,b.role_id as roleid, " +
	"c.name as rolename,c.expire_time as roleexpire, " +
	"d.injob_flag as injobflag,d.department_id as departmentid,e.name as deptname, " +
	"g.org_name as orgname " +
	"from nc_staff a " +
	"left join %s.admin b on a.staff_id=b.id " +
	"left join %s.role c on b.role_id=c.id " +
	"left join nc_department_staff d on a.staff_id=d.staff_id and d.status=1 " +
	"left join nc_department e on d.department_id=e.department_id "+
	"left join nc_organization g on a.organization_id=g.organization_id "

var SQL_COUNT_STAFF = "select count(*) as ucount " +
	"from nc_staff a " +
	"left join %s.admin b on a.staff_id=b.id " +
	"left join %s.role c on b.role_id=c.id " +
	"left join nc_department_staff d on a.staff_id=d.staff_id and d.status=1 " +
	"left join nc_department e on d.department_id=e.department_id "+
	"left join nc_organization g on a.organization_id=g.organization_id "

var SQL_ATTENDANCE = "select a.attendance_id as id,a.siteid,a.organization_id as organizationid,a.staff_id as staffid," +
	"a.clock_time_start as clocktimestart,a.clock_time_end as clocktimeend,a.clocktime, " +
	"a.remark,a.operatorid,b.name as staffname,b.phone " +
	"from nc_attendance a " +
	"left join nc_staff b on a.staff_id=b.staff_id "

var SQL_COUNT_ATTENDANCE = "select count(*) as ucount " +
	"from nc_attendance a " +
	"left join nc_staff b on a.staff_id=b.staff_id "

var SQL_ATTENDANCE_APPLY = "select a.apply_id as id,a.siteid,a.organization_id as organizationid,a.staff_id as staffid,a.apply_time as applytime,a.apply_type as applytype," +
	"a.attendance_type as attendancetype,a.reason,a.attendance_time_first as attendancetimefirst, " +
	"a.attendance_time_second as attendancetimesecond, " +
	"a.apply_status as status,a.remark,a.operatorid, " +
	"b.name as staffname,b.phone " +
	"from nc_attendance_apply a " +
	"left join nc_staff b on a.staff_id=b.staff_id "

var SQL_COUNT_ATTENDANCE_APPLY = "select count(*) as ucount " +
	"from nc_attendance_apply a " +
	"left join nc_staff b on a.staff_id=b.staff_id "

var SQL_DEPT_STAFF = "select a.department_id as id,a.siteid,a.organization_id as organizationid,a.name as deptname,a.parent_id as parentid,a.level," +
	"a.create_time as createtime,a.status,a.remark,a.operatorid, " +
	"b.induction_time as inductiontime,b.departure_time as departuretime,b.injob_flag as injobflag, " +
	"b.position_id as positionid, c.name as staffname,c.phone,d.name as posname " +
	"from nc_department a " +
	"left join nc_department_staff b on a.department_id=b.department_id and b.injob_flag =1  " +
	"left join nc_staff c on b.staff_id=c.staff_id " +
	"left join nc_position d on b.position_id= d.position_id "
var SQL_COUNT_DEPT_STAFF = "select count(*) as ucount " +
	"from nc_department a " +
	"left join nc_department_staff b on a.department_id=b.department_id and b.injob_flag =1  " +
	"left join nc_staff c on b.staff_id=c.staff_id " +
	"left join nc_position d on b.position_id= d.position_id "

var SQL_GROUP = "select a.group_id as id,a.siteid,a.organization_id as organizationid,a.name," +
	"a.create_time as createtime,a.status,a.remark,a.operatorid " +
	"from nc_group a "
var SQL_COUNT_GROUP = "select count(*) as ucount " +
	"from nc_group a "

var SQL_GROUP_STAFF = "select a.grp_staff_id as id,a.group_id as groupid,a.siteid,a.organization_id as organizationid,a.staff_id as staffid," +
	"a.create_time as createtime,a.status,a.remark,a.position, " +
	"b.name as groupname,c.name as staffname,c.phone,c.address " +
	"from nc_group_staff a " +
	"left join nc_group b on a.group_id=b.group_id " +
	"left join nc_staff c on a.staff_id=c.staff_id "

var SQL_COUNT_GROUP_STAFF = "select count(*) as ucount " +
	"from nc_group_staff a " +
	"left join nc_group b on a.group_id=b.group_id " +
	"left join nc_staff c on a.staff_id=c.staff_id "

var SQL_POSITION = "select a.position_id as id,a.name, " +
	"a.create_time as createtime,a.status,a.remark,a.operatorid " +
	"from nc_position a "

var SQL_COUNT_POSITION = "select count(*) as ucount " +
	"from nc_position a "

var SQL_SCHEDULE_MONTHLY = "select a.month_sched_id as id,a.siteid,a.organization_id as organizationid,a.sched_month as month,a.staff_id as staffid,a.group_id as groupid," +
	"a.day1,a.day2,a.day3,a.day4,a.day5,a.day6,a.day7,a.day8,a.day9,a.day10,a.day11,a.day12,a.day13,a.day14,a.day15, " +
	"a.day16,a.day17,a.day18,a.day19,a.day20,a.day21,a.day22,a.day23,a.day24,a.day25,a.day26,a.day27,a.day28,a.day29,a.day30,a.day31, " +
	"a.create_time as createtime,a.status,a.remark,a.operatorid,a.sched_plan_id as schedplanid, " +
	"b.name as staffname,b.phone,c.name as groupname " +
	"from nc_month_schedules a " +
	"left join nc_staff b on a.staff_id=b.staff_id " +
	"left join nc_group c on a.group_id=c.group_id "

var SQL_COUNT_SCHEDULE_MONTHLY = "select count(*) as ucount " +
	"from nc_month_schedules a " +
	"left join nc_staff b on a.staff_id=b.staff_id " +
	"left join nc_group c on a.group_id=c.group_id "

var SQL_SCHED_PLAN = "select a.sched_plan_id as id,a.siteid,a.organization_id as organizationid,a.name,a.group_id as groupid,a.plan_type as plantype," +
	"a.create_time as createtime,a.status,a.remark,a.operatorid, " +
	"b.name as groupname " +
	"from nc_sched_plan a " +
	"left join nc_group b on a.group_id=b.group_id "

var SQL_COUNT_SCHED_PLAN = "select count(*) as ucount " +
	"from nc_sched_plan a " +
	"left join nc_group b on a.group_id=b.group_id "

var SQL_SCHEDULES = "select a.schedules_id as id,a.siteid,a.organization_id as organizationid,a.name," +
	"a.sched_start_time as starttime,a.sched_end_time as endtime, " +
	"a.create_time as createtime,a.status,a.remark,a.operatorid " +
	"from nc_schedules a "
var SQL_COUNT_SCHEDULES = "select  count(*) as ucount " +
	"from nc_schedules a "

func SetSearchPath(o orm.Ormer, schema string) (err error) {
	if o == nil {
		return errors.New("orm is nil")
	}

	if imconf.Config.DbDriver == "pgsql" {
		smt := fmt.Sprintf("set search_path to \"%s\";", schema)
		_, err = o.Raw(smt).Exec()
	} else if imconf.Config.DbDriver == "mysql" {
		return nil
	}

	return err
}
