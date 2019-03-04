package admins

import (
	"errors"
	"strings"
	"time"

	"platform/lib/helper"
)

var (
	ErrAdminAlreadyExist = errors.New("Admin already exist")
	ErrAdminNotExist     = errors.New("Admin does not exist")
	ErrNameAlreadyUsed   = errors.New("Name already used")
	ErrAdminNameIllegal  = errors.New("Admin name contains illegal characters")
	ErrWarnPasswd        = errors.New("Passwd is warn")
	ErrResourceNotExist  = errors.New("Resource does not exist")
	ErrResourceIsParent  = errors.New("Resource have chind Resource")
	ErrRoleNotExist      = errors.New("Role does not exist")
	ErrRoleIsParent      = errors.New("Role have chind Resource")
	ErrSiteIdNotExist    = errors.New("租户不存在")
	ErrAdminIsBlocked    = errors.New("Admin is blocked.")
	ErrAdminExpired      = errors.New("账号已过期。")
)

const (
	STATE_BLOCKED 		= 1
	STATE_NORMAL 		= 0
)

type Resource struct {
	Id            int64		`xorm:"pk"`
	ParentId      int64  `xorm:"index 'parent_id' bigint"`
	Name          string `xorm:"unique not null"`
	Iconid        int
	Level         int    `xorm:"not null"` //0:系统(只有siteid = 0可以看到) 1:应用
	Type          int    `xorm:"not null"` //0:目录 1:页面 2:功能 (目录包含页面, 页面包含功能. 目录不能直接包含功能)
	Url           string `xorm:"not null"`
	Proxy         string
	Description   string
	EffectiveTime time.Time `xorm:"not null"`
	ExpireTime    time.Time `xorm:"not null"`
	Created       time.Time `xorm:"created"`
	Updated       time.Time `xorm:"updated"`
	Icon   		  string
	Appid 		  int64
	PermDesc 	  string
}

type Role struct {
	Id            int64		`xorm:"pk"`
	SiteId        int64		`xorm:"unique(sietid_name)"`
	Name          string    `xorm:"not null unique(sietid_name)"`
	Iconid        int
	Description   string
	EffectiveTime time.Time `xorm:"not null"`
	ExpireTime    time.Time `xorm:"not null"`
	Created       time.Time `xorm:"created"`
	Updated       time.Time `xorm:"updated"`
	Startpage	  string
	WheelFlag 	  int		`xorm:"wheelflag"`
	BlockFlag 	  int		`xorm:"blockflag"`
	OrgId 		  int64
}

type Admin struct {
	Id            int64		`xorm:"pk"`
	SiteId        int64 `xorm:"not null unique(sietid_name) unique(sietid_email) unique(sietid_phone)"`
	RoleId        int64
	Type          int    `xorm:"not null"` //0:用户 1:AppId
	Name          string `xorm:"unique(sietid_name)"`
	JobNumber     string
	Passwd        string `xorm:"not null"`
	Email         string `xorm:"unique(sietid_email)"`
	Phone         string `xorm:"unique(sietid_phone)"`
	Realname      string
	Description   string
	EffectiveTime time.Time `xorm:"not null"`
	ExpireTime    time.Time `xorm:"not null"`
	Created       time.Time `xorm:"created"`
	Updated       time.Time `xorm:"updated"`
	ImageUrl      string 	`xorm:"image_url"`
	State         int    	`xorm:"not null"` //0：正常，1：被阻止登录，默认为0
	Gender 	      int
}

type RoleResource struct {
	Role     			Role     	`xorm:"'role_id' bigint unique(idx_role_resource_role_id)"`
	Resource 			Resource 	`xorm:"'resource_id' bigint"`
	SubSysId 			int64 		`xorm:"not null unique(idx_role_resource_role_id)"`
	StartResourceId 	int64 		`xorm:"not null"`
	ResId 				int64
	Appid  				int64		`xorm:"not null unique(idx_role_resource_role_id)"`
	RoleResourceId		int64 		`xorm:"pk"`
	CreateTime 			string
}


type SiteRes struct {
	ResId			int64		`orm:"pk;auto"`
	ParentId		int64
	ResourceId		int64
	Treeid			string
	Name			string
	StartTime		string
	EndTime			string
	Status			int
	Level			int
	Order			int
	Direction		int
}

type RoleRes struct {
	ResId			int64		`orm:"pk;auto"`
	ParentId		int64
	ResourceId		int64
	Treeid			string
	Name			string
	StartTime		string
	EndTime			string
	Status			int
	Level			int
	Order			int

	//操作权限
	PermSel 		int
	PermAdd 		int
	PermUpd 		int
	PermDel 		int
	PermCancel 		int
	PermAudit 		int

	//数据权限
	PermEval 		string
	PermDoc 		string

}


type SubSys struct {
	SubSysId 		int64  		`orm:"pk;auto"`
	ResId 			int64
	Appid 			int64
	SiteId 			int64
	Subtplid		int64
	Name 			string
	IconUrl 		string
	StartResourceId int64
	Created 		string
	Remark 			string
}

type Logs struct {
	LogType int //登录日志, 查看日志, 操作日志
	UserId  int64
	Ip      string
	Msg     string
	Created time.Time `xorm:"created"`
	Params  string
}

func InitResource() (*Resource, error) {
	r := &Resource{
		ParentId:    0,
		Name:        "系统后台",
		Iconid:      1,
		Level:       0,
		Type:        0,
		Url:         "/admin/base/index.html",
		Description: "",
		// EffectiveTime: time..Format("2014-09-01"),
		// ExpireTime:    "2014-12-01",
	}
	_, err := orm.Insert(r)
	return r, err
}

//-- Resource
func CreateResource(resource *Resource) (*Resource, error) {
	if resource.Type != 0 && resource.Type != 1 && resource.Type != 2 {
		return nil, errors.New("Type error")
	}

	if resource.ParentId < 1 {
		return nil, errors.New("parentid = 0")
	}

	presource, err := GetResourceById(resource.ParentId)
	if err != nil {
		return nil, err
	}

	if resource.Type == 0 {
		if presource.Type != 0 {
			return nil, errors.New("sub type error, parent type is not 0")
		}
	} else if resource.Type != presource.Type+1 {
		return nil, errors.New("sub type error, < parent type")
	}

	if _, err := orm.Insert(resource); err != nil {
		return nil, err
	}
	return resource, nil
}

func GetResourceById(id int64) (*Resource, error) {
	Resource := new(Resource)
	has, err := orm.Id(id).Get(Resource)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrResourceNotExist
	}
	return Resource, nil
}

func GetResourceByUrl(url string) (*Resource, error) {
	resource := &Resource{Url: url}
	has, err := orm.Get(resource)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrResourceNotExist
	}
	return resource, nil
}

func UpdateResource(resource *Resource) (bool, error) {
	_, err := GetResourceById(resource.Id)
	if err != nil {
		return false, err
	}

	affected, err := orm.Id(resource.Id).Cols("Parent_Id", "Name", "Iconid", "Level", "Type", "Url", "Proxy", "Effectiv_Time", "Expire_Time", "Description").Update(resource)
	return affected == 1, err
}

func DeleteResourceById(id int64) (bool, error) {
	resource := new(Resource)
	total, err := orm.Where("Parent_Id = ?", id).Count(resource)
	if err != nil {
		return false, err
	}
	if total > 0 {
		return false, ErrResourceIsParent
	}
	affected, err := orm.Id(id).Delete(&Resource{Id: id})
	return affected == 1, err
}

func ListResourceById(id int64) ([]Resource, error) {
	var resources = make([]Resource, 0)
	err := orm.Where("Parent_Id = ?", id).Find(&resources)
	return resources, err
}

//--Role

func CreateRole(role *Role) (*Role, error) {
	if _, err := orm.Insert(role); err != nil {
		return nil, err
	}
	return role, nil
}

func GetRoleById(id int64) (*Role, error) {
	role := new(Role)
	has, err := orm.Id(id).Get(role)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrRoleNotExist
	}
	return role, nil
}

func UpdateRole(role *Role) (bool, error) {
	has, err := orm.Get(&Role{Id: role.Id})
	if err != nil {
		return false, err
	}
	if !has {
		return false, ErrRoleNotExist
	}
	affected, err := orm.Id(role.Id).Cols("Name", "Iconid", "EffectiveTime", "ExpireTime", "Description").Update(role)
	return affected == 1, err
}

func DeleteRoleById(id int64) (bool, error) {
	// role := new(Role)
	// total, err := orm.Where("Parent_Id = ?", id).Count(role)
	// if err != nil {
	// 	return false, err
	// }
	// if total > 0 {
	// 	return false, ErrRoleIsParent
	// }
	affected, err := orm.Id(id).Delete(&Role{Id: id})
	return affected == 1, err
}

func ListRole(siteid int64) ([]Role, error) {
	var roles = make([]Role, 0)
	err := orm.Where("site_id = ?", siteid).Asc("Id").Find(&roles)
	return roles, err
}

//--user
func IsAdminNameExist(siteid int64, name string) (bool, error) {
	if len(name) == 0 {
		return false, nil
	}
	return orm.Get(&Admin{SiteId: siteid, Name: strings.ToLower(name)})
}

func IsPhoneExist(siteid int64, phone string) (bool, error) {
	if siteid == 0 {
		return false, ErrSiteIdNotExist
	}

	if len(phone) == 0 {
		return false, nil
	}
	return orm.Get(&Admin{SiteId: siteid, Phone: phone})
}

func GetAdminById(id int64) (*Admin, error) {
	user := new(Admin)
	has, err := orm.Id(id).Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrAdminNotExist
	}
	return user, nil
}

func GetAdminByName(siteid int64, name string) (*Admin, error) {
	if len(name) == 0 {
		return nil, ErrAdminNotExist
	}
	user := &Admin{SiteId: siteid, Name: strings.ToLower(name)}
	has, err := orm.Get(user)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrAdminNotExist
	}
	return user, nil
}

func GetUserByPhone(siteid int64, phone string) (*Admin, error) {
	if siteid == 0 {
		return nil, ErrSiteIdNotExist
	}
	if len(phone) == 0 {
		return nil, ErrAdminNotExist
	}
	user := &Admin{SiteId: siteid, Phone: phone}
	has, err := orm.Get(user)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrAdminNotExist
	}
	return user, nil
}

func CreateAdmin(user *Admin) (*Admin, error) {
	isExist, err := IsAdminNameExist(user.SiteId, user.Name)
	if err != nil {
		return nil, err
	} else if isExist {
		return nil, ErrNameAlreadyUsed
	}

	user.Name = strings.ToLower(user.Name)
	if _, err = orm.Insert(user); err != nil {
		return nil, err
	}

	return user, err
}

func LoginAdmin(siteid int64, name, passwd, Salt string) (*Admin, error) {
	user := new(Admin)
	rows, err := orm.Where("site_id = ? and (name = ? or phone = ?)", siteid, strings.ToLower(name), strings.ToLower(name)).Rows(user)
	if err != nil {
		return nil, ErrAdminNotExist
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(user)
		if err != nil {
			return nil, ErrAdminNotExist
		}
		//判断用户是否被阻止
		if user.State == STATE_BLOCKED {
			return nil,ErrAdminIsBlocked
		}

		//判断账号是否过期
		ts := user.EffectiveTime.Unix()
		te := user.ExpireTime.Unix()
		tn := time.Now().Unix()

		if tn < ts || tn > te {
			return nil,ErrAdminExpired
		}

		if user.Passwd == passwd {
			return user, nil
		}

		if helper.Md5(user.Passwd+Salt) == passwd {
			return user, nil
		}

		if helper.Md5(helper.Md5(user.Passwd)+Salt) == passwd {
			return user, nil
		}
	}
	return nil, ErrAdminNotExist
}

func ResetPasswd(id int64, oldpasswd string, newpasswd string) (bool, error) {
	user, err := GetAdminById(id)
	if err != nil {
		return false, err
	}

	if user.Passwd != oldpasswd {
		return false, ErrWarnPasswd
	}
	user.Passwd = newpasswd
	affected, err := orm.Id(id).Cols("passwd").Update(user)
	return affected == 1, err
}

func SetPasswd(siteid int64, id int64, newpasswd string) (bool, error) {
	user, err := GetAdminById(id)
	if err != nil {
		return false, err
	}
	user.Passwd = newpasswd
	affected, err := orm.Id(id).Where("site_id = ?", siteid).Cols("passwd").Update(user)
	return affected == 1, err
}

func UpdateAdmin(user *Admin) (bool, error) {
	has, err := orm.Get(&Admin{Id: user.Id})
	if err != nil {
		return false, err
	}
	if !has {
		return false, ErrAdminNotExist
	}

	affected, err := orm.Id(user.Id).Cols("Role_Id", "Job_Number", "Passwd", "Email", "Phone", "Description", "EffectiveTime", "ExpireTime").Update(user)
	return affected == 1, err
}

func DeleteAdminById(id int64) (bool, error) {
	affected, err := orm.Id(id).Delete(&Admin{Id: id})
	return affected == 1, err
}

func ListAdmin(siteid, roleid int64, usertype int) ([]Admin, error) {
	admins := make([]Admin, 0)
	m := orm.Where("site_id = ?", siteid).And("type = ?", usertype)
	if roleid > 0 {
		m = m.And("role_Id = ?", roleid)
	}
	err := m.Desc("Id").Find(&admins)
	return admins, err
}

//--RoleResource

func IsRoleResourceExist(roleid int64, resourceid int64) (bool, error) {
	return orm.Get(&RoleResource{Role: Role{Id: roleid}, Resource: Resource{Id: resourceid}})
}

func CreateRoleResource(role_resource *RoleResource) (*RoleResource, error) {
	if role_resource.Resource.Type != 2 {
		return nil, errors.New("roleresource is not bind type 2")
	}

	if _, err := orm.Insert(role_resource); err != nil {
		return nil, err
	}
	return role_resource, nil
}

func DeleteRoleResource(role_resource *RoleResource) (bool, error) {
	affected, err := orm.Delete(role_resource)
	return affected > 0, err
}

func ListRoleResourceByResourceId(resourceid int64) ([]RoleResource, error) {
	var roleresources = make([]RoleResource, 0)
	err := orm.Where("resource_id = ?", resourceid).Find(&roleresources)
	return roleresources, err
}

func ListRoleResourceByRoleId(roleid int64) ([]RoleResource, error) {
	var roleresources = make([]RoleResource, 0)
	err := orm.Where("role_id = ?", roleid).Find(&roleresources)
	return roleresources, err
}

func CheckPower(siteid int64, userid int64, url string) bool {
	sql := "select count(*) as n from resource a, role_resource b, admin c where c.role_id = b.role_id  and  a.id = b.resource_id and  c.site_id= ? and c.id = ? and a.url = ?"
	results, err := orm.Query(sql, siteid, userid, url)
	if err != nil {
		return false
	}

	if string(results[0]["n"]) == "1" {
		return true
	}
	return false
}
