package models

import "time"

type NcDepartment struct {
	DepartmentId   int64 			`orm:"pk;auto"`
	Siteid         int64
	OrganizationId int64

	Name     string					`orm:"size(255)"`
	ParentId int64
	Level    int16

	Remark     string				`orm:"type(text)"`
	CreateTime time.Time			`orm:"auto_now_add;type(datetime)"`
	UpdateTime time.Time			`orm:"type(datetime)"`
	Status     int16
	Creatorid  int64
	Updaterid  int64
}

type NcGroup struct {
	GroupId        int64 			`orm:"pk;auto"`
	Siteid         int64			`orm:"default(0)"`
	OrganizationId int64			`orm:"default(0)"`

	Name string						`orm:"size(255)"`

	Remark     string				`orm:"type(text)"`
	CreateTime time.Time			`orm:"auto_now_add;type(datetime)"`
	Status     int16
	Operatorid int64
}

type NcStaff struct {
	StaffId        int64 			`orm:"pk;auto"`
	Siteid         int64
	OrganizationId int64

	StaffType int16
	Name      string				`orm:"size(255)"`
	PaperType int16
	Papers    string				`orm:"size(50)"`
	Phone     string				`orm:"size(50)"`
	Address   string				`orm:"size(255)"`

	// 增加离职入职时间
	EntryTime time.Time				`orm:"type(date);null"`
	LeaveTime time.Time				`orm:"type(date);null"`

	Remark     string				`orm:"type(text)"`
	CreateTime time.Time			`orm:"auto_now_add;type(datetime)"`
	Status     int16
	Operatorid int64

	//签名
	Signature string				`orm:"size(255)"`

	//后续增加
	ExpertFlag          int16
	EmployeeCode        string		`orm:"size(100)"`
	Gender              int16
	Nationality         string		`orm:"size(100)"`
	DutiesCode          string		`orm:"size(100)"`
	JobLevelCode        string		`orm:"size(100)"`
	Birthday            time.Time	`orm:"type(date)"`
	ArchiveId           int64
	ArchiveBarcode      string		`orm:"size(100)"`
	ArchiveAddress      string		`orm:"size(255)"`
	FormerName          string		`orm:"size(100)"`
	EnglishName         string		`orm:"size(100)"`
	NativePlace         string		`orm:"size(100)"`
	PoliticalStatusCode string		`orm:"size(100)"`
	MarriageStatusCode  string		`orm:"size(100)"`
	HealthStatusCode    string		`orm:"size(100)"`
	ExpertClassCode     string		`orm:"size(100)"`
	ExpertAllawance     string		`orm:"size(100)"`
	OnJobStatus         int16
	DirectorDeptId      int64
	WorkTel             string		`orm:"size(100)"`
	HomeTel             string		`orm:"size(100)"`
	PostClass           string		`orm:"size(100)"`
	WorkMajor           string		`orm:"size(100)"`
	MajorPost           string		`orm:"size(100)"`
	RecruitResource     string		`orm:"size(100)"`
	EmploymentStartTime time.Time		`orm:"type(date);null"`
	EmploymentEndTime   time.Time		`orm:"type(100);null"`
	JoinArmyFlag        int16
	JoinArmyTime        time.Time	`orm:"type(date)"`
	WorkageStartTime    time.Time	`orm:"type(date)"`
	WorkageEndRemark    string		`orm:"size(500)"`
	Workage             int16
	PostCertId          string		`orm:"size(100)"`
	EducationCertId     string		`orm:"size(100)"`
	AcademicRecord      string		`orm:"size(100)"`
	Degreen             string		`orm:"size(100)"`
	FirstLanguage       string		`orm:"size(100)"`
	FirstLanguageLevel  string		`orm:"size(100)"`
	SecondLanguage      string		`orm:"size(100)"`
	SecondLanguageLevel string		`orm:"size(100)"`
	JoinWorkTime        time.Time	`orm:"type(date)"`
	Specialty           string		`orm:"size(100)"`

	QualCertName 		string		`orm:"size(100)"`
	QualCertCode 		string		`orm:"size(100)"`
	QualCertDate 		time.Time	`orm:"type(date)"`
	Academy				string		`orm:"size(100)"`
	Nation				string		`orm:"size(100)"`
	AcademySpecialty	string		`orm:"size(100)"`
	GraduateDate		time.Time	`orm:"type(date)"`
	JoinPartyDate		time.Time	`orm:"type(date)"`
	SpecialtyLevel		int16
}

type NcWeekSchedules struct {
	WeekSchedId    int64 			`orm:"pk;auto"`
	Siteid         int64
	OrganizationId int64

	GroupId     int64
	StaffId     int64
	SchedPlanId int64
	SchedType   int16

	WeekStartDate time.Time			`orm:"type(date)"`
	WeekEndDate   time.Time			`orm:"type(date)"`

	Monday    int64
	Tuesday   int64
	Wednesday int64
	Thursday  int64
	Friday    int64
	Saturday  int64
	Sunday    int64

	Remark     string				`orm:"type(text)"`
	CreateTime time.Time			`orm:"auto_now_add;type(datetime)"`
	Status     int16
	Operatorid int64
}

type NcDepartmentStaff struct {
	DepStaffId     int64 				`orm:"pk;auto"`
	Siteid         int64
	OrganizationId int64

	StaffId      int64
	PositionId   int64
	DepartmentId int64

	InductionTime time.Time				`orm:"type(date)"`
	DepartureTime time.Time				`orm:"type(date)"`
	InjobFlag     int16

	Remark     string					`orm:"type(text)"`
	CreateTime time.Time				`orm:"auto_now_add;type(datetime)"`
	Status     int16
}

type NcAttendance struct {
	AttendanceId   int64 				`orm:"pk;auto"`
	Siteid         int64
	OrganizationId int64

	StaffId        int64
	ClockTimeStart time.Time			`orm:"type(datetime)"`
	ClockTimeEnd   time.Time			`orm:"type(datetime)"`
	Clocktime      time.Time			`orm:"type(date)"`

	Remark     string					`orm:"type(text)"`
	Operatorid int64
}
type NcAttendanceApply struct {
	ApplyId        int64 				`orm:"pk;auto"`
	Siteid         int64
	OrganizationId int64

	StaffId              int64
	ApplyTime            time.Time		`orm:"type(datetime)"`
	ApplyType            int16
	AttendanceType       int16
	Reason               string			`orm:"type(text)"`
	AttendanceTimeFirst  time.Time		`orm:"type(datetime)"`
	AttendanceTimeSecond time.Time		`orm:"type(datetime)"`

	Remark      string					`orm:"type(text)"`
	ApplyStatus int16
	Operatorid  int64
}

type NcGroupStaff struct {
	GrpStaffId int64 					`orm:"pk;auto"`

	GroupId  int64
	StaffId  int64
	Position int

	Remark     string					`orm:"type(text)"`
	CreateTime time.Time				`orm:"auto_now_add;type(datetime)"`
	Status     int16
}

type NcPosition struct {
	PositionId int64 					`orm:"pk;auto"`
	Siteid     int64

	Name string							`orm:"size(255)"`

	Remark     string					`orm:"type(text)"`
	CreateTime time.Time				`orm:"auto_now_add;type(datetime)"`
	Status     int16
	Operatorid int64
}

type NcMonthSchedules struct {
	MonthSchedId   int64 				`orm:"pk;auto"`
	Siteid         int64
	OrganizationId int64

	SchedMonth  time.Time				`orm:"type(date)"`
	GroupId     int64
	StaffId     int64
	SchedPlanId int64

	Day1  int64
	Day2  int64
	Day3  int64
	Day4  int64
	Day5  int64
	Day6  int64
	Day7  int64
	Day8  int64
	Day9  int64
	Day10 int64
	Day11 int64
	Day12 int64
	Day13 int64
	Day14 int64
	Day15 int64
	Day16 int64
	Day17 int64
	Day18 int64
	Day19 int64
	Day20 int64
	Day21 int64
	Day22 int64
	Day23 int64
	Day24 int64
	Day25 int64
	Day26 int64
	Day27 int64
	Day28 int64
	Day29 int64
	Day30 int64
	Day31 int64

	Remark     string						`orm:"type(text)"`
	CreateTime time.Time					`orm:"auto_now_add;type(datetime)"`
	Status     int16
	Operatorid int64
}

type NcSchedPlan struct {
	SchedPlanId    int64 					`orm:"pk;auto"`
	Siteid         int64
	OrganizationId int64

	GroupId  int64
	Name     string							`orm:"size(255)"`
	PlanType int16

	Remark     string						`orm:"type(text)"`
	CreateTime time.Time					`orm:"auto_now_add;type(datetime)"`
	Status     int16
	Operatorid int64
}

type NcSchedules struct {
	SchedulesId    int64 					`orm:"pk;auto"`
	Siteid         int64
	OrganizationId int64

	Name           string
	SchedStartTime time.Time				`orm:"type(datetime)"`
	SchedEndTime   time.Time				`orm:"type(datetime)"`

	Remark     string						`orm:"type(text)"`
	CreateTime time.Time					`orm:"auto_now_add;type(datetime)"`
	Status     int16
	Operatorid int64
}

type NcOrganization struct {
	OrganizationId int64 					`orm:"pk;auto"`

	Siteid             int64
	RegCode            string				`orm:"size(100)"`
	ApprovalBedNum     int
	ApprovalDepartment string				`orm:"size(100)"`
	ApprovalTime       time.Time			`orm:"type(datetime)"`
	Acreage            float32
	BranchesNum        int
	ManagementArea     string				`orm:"size(255)"`
	OrgCode            string				`orm:"size(50)"`
	OrgLevel           string				`orm:"size(50)"`
	Address            string				`orm:"size(255)"`
	MembershipLevel    string				`orm:"size(50)"`
	Phone              string				`orm:"size(100)"`
	LegalPerson        string				`orm:"size(50)"`
	OrgType            string				`orm:"size(50)"`
	OrgLn              string				`orm:"size(100)"`
	ParentId           int64

	RegDivId int64
	RegName  		string					`orm:"size(100)"`
	Nature		  	int8
	EffectiveDate 	time.Time				`orm:"type(datetime);null"`
	ExpireDate 		time.Time				`orm:"type(datetime);null"`
	OrgName 		string					`orm:"size(255)"`
	ManageCate 		int8

	CreateTime time.Time					`orm:"auto_now_add;type(datetime)"`
	UpdateTime time.Time					`orm:"type(datetime);null"`
	Status     int16
	Creatorid  int64
	Updaterid  int64
	Remark     string						`orm:"type(text);null"`
}
