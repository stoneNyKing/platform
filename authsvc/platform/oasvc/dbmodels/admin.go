package dbmodels

import (
	"errors"
	"strings"
	"time"

	"platform/lib/helper"
)

var (
	ErrAdminAlreadyExist = errors.New("Admin already exist")
	ErrAdminNotExist     = errors.New("账号不存在或者密码错误。")
	ErrNameAlreadyUsed   = errors.New("Name already used")
	ErrAdminNameIllegal  = errors.New("Admin name contains illegal characters")
	ErrWarnPasswd        = errors.New("Passwd is warn")
	ErrResourceNotExist  = errors.New("Resource does not exist")
	ErrResourceIsParent  = errors.New("Resource have chind Resource")
	ErrRoleNotExist      = errors.New("Role does not exist")
	ErrRoleIsParent      = errors.New("Role have chind Resource")
	ErrSiteIdNotExist    = errors.New("租户不存在")
	ErrAdminIsBlocked    = errors.New("账号已失效.")
	ErrAdminExpired      = errors.New("账号已过期。")
)

const (
	STATE_BLOCKED = 1
	STATE_NORMAL  = 0
)

type Admin struct {
	Id             int64 `xorm:"pk autoincr"`
	SiteId         int64 `xorm:"not null unique(sietid_name) unique(sietid_email) unique(sietid_phone)"`
	OrganizationId int64 `xorm:"not null"`
	RoleId         int64
	Type           int    `xorm:"not null"` //0:用户 1:AppId
	Name           string `xorm:"unique(sietid_name)"`
	JobNumber      string
	Passwd         string `xorm:"not null"`
	Email          string `xorm:"unique(sietid_email)"`
	Phone          string `xorm:"unique(sietid_phone)"`
	Realname       string
	Description    string
	EffectiveTime  time.Time `xorm:"not null"`
	ExpireTime     time.Time `xorm:"not null"`
	Created        time.Time `xorm:"created"`
	Updated        time.Time `xorm:"updated"`
	ImageUrl       string    `xorm:"image_url"`
	State          int       `xorm:"not null"` //0：正常，1：被阻止登录，默认为0
	Gender         int
}

type AdminLog struct {
	AdminLogId int64 `xorm:"pk autoincr"`
	SiteId     int64 `xorm:"not null unique(sietid_adminid)"`
	Appid      int64
	AdminId    int64 `xorm:"not null unique(sietid_adminid)"`
	Token      string
	Ip         string
	Msg        string `xorm:"varchar(1024)"`
	Level      int16
	Action     string    `xorm:"varchar(255)"`
	Json       string    `xorm:"text"`
	Created    time.Time `xorm:"created"`
}

type Appid struct {
	Appid    int64 		`xorm:"pk"`
	Appkey     string `xorm:"not null"`
	Remark  string    `xorm:"text"`
	Status  int       `xorm:"not null"`
	Json    map[string]string					   `xorm:"json"`
	Created        time.Time `xorm:"created"`
	Updated        time.Time `xorm:"updated"`
	BefeFlag int8 			`xorm:"not null"`
}

func (self *Appid) Get() (bool, error) {
	has, err := orm.Get(self)
	return has, err
}

func (self *Appid) Insert() error {
	_, err := orm.InsertOne(self)
	return err
}

func (self *Appid) Delete() error {
	_, err := orm.Id(self.Appid).Delete(self)
	return err
}

func (self *Appid) Update() error {
	_, err := orm.Id(self.Appid).MustCols("Appid", "Appkey", "Remark", "Status","BefeFlag","Json").Update(self)
	return err
}

func Search() ([]Appid, error) {
	list := make([]Appid, 0)
	err := orm.Find(&list)
	return list, err
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
			return nil, ErrAdminIsBlocked
		}

		//判断账号是否过期
		ts := user.EffectiveTime.Unix()
		te := user.ExpireTime.Unix()
		tn := time.Now().Unix()

		if tn < ts || tn > te {
			return nil, ErrAdminExpired
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
