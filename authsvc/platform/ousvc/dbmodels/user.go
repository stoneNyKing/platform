package dbmodels

import (
	"encoding/json"
	"errors"
	"platform/lib/helper"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var (
	ErrUserAlreadyExist    = errors.New("用户已经存在")
	ErrUserNotExist        = errors.New("用户不存在")
	ErrNameAlreadyUsed     = errors.New("用户名已经存在")
	ErrEmailAlreadyUsed    = errors.New("Email已经存在")
	ErrPhoneAlreadyUsed    = errors.New("手机号码已经存在")
	ErrIdcardAlreadyUsed   = errors.New("身份证号码已经存在")
	ErrRfidAlreadyUsed     = errors.New("Rfid已经存在")
	ErrIMEIAlreadyUsed     = errors.New("IMEI已经存在")
	ErrRfidNotUsed         = errors.New("Rfid不存在")
	ErrWeixinidNotUsed     = errors.New("Weixinid不存在")
	ErrWeixinidAlreadyUsed = errors.New("Weixinid已经存在")
	ErrWarnPasswd          = errors.New("密码错误")
	ErrSiteIdNotExist      = errors.New("租户不存在")
)

type User struct {
	Id           int64  `xorm:"pk autoincr"`
	SiteId       int64  `xorm:"unique not null unique(sietid_name) unique(sietid_idcard) unique(sietid_email) unique(sietid_phone)"`
	Name         string `xorm:"unique not null unique(sietid_name)"`
	Phone        string `xorm:"unique not null unique(sietid_phone)"`
	Email        string `xorm:"unique not null unique(sietid_email)"`
	Idcard       string `xorm:"unique not null unique(sietid_idcard)"`
	Rfid         string `xorm:"'rfid'"`
	Imei         string `xorm:"'imei'"`
	Nickname     string `xorm:"'nickname'"`
	Weixinid     string `xorm:"'weixinid'"`
	Passwd       string `xorm:"not null"`
	Type         int    //保留使用
	IsActive     int
	Created      time.Time `xorm:"created"`
	Updated      time.Time `xorm:"updated"`
	LastLogin    time.Time
	LastLogout   time.Time
	Json         string `xorm:"text"`
	ImageUrl     string
	RegisterType int
	Sscard       string
	Realname     string
}

type NameLog struct {
	SiteId  int64     `xorm:"pk"`
	Type    string    `xorm:"pk"`
	Name    string    `xorm:"pk"`
	Id      int64     `xorm:"not null"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

type UserLog struct {
	Id      int64 `xorm:"pk autoincr"`
	SiteId  int64
	Token   string
	Appid   int
	UserId  int64
	Level   int //0:user 1:admin
	Ip      string
	Act     string
	Msg     string
	Json    map[string]interface{} `xorm:text`
	Created time.Time              `xorm:"created"`
}

type Feedback struct {
	Id      int64 `xorm:"pk autoincr"`
	SiteId  int64
	Appid   int
	Userid  int64     `xorm:"not null"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
	Content string
}

type Appid struct {
	Appid    int64             `xorm:"pk"`
	Appkey   string            `xorm:"not null"`
	Remark   string            `xorm:"text"`
	Status   int               `xorm:"not null"`
	Json     map[string]string `xorm:"json"`
	Created  time.Time         `xorm:"created"`
	Updated  time.Time         `xorm:"updated"`
	BefeFlag int8              `xorm:"not null"`
}

const (
	REGISTER_TYPE_UNKNOW   = 0
	REGISTER_TYPE_PHONE    = 1
	REGISTER_TYPE_NAME     = 2
	REGISTER_TYPE_EMAIL    = 3
	REGISTER_TYPE_IDCARD   = 4
	REGISTER_TYPE_IMEI     = 5
	REGISTER_TYPE_SSCARD   = 6
	REGISTER_TYPE_WEIXINID = 7
)

var (
	// NameRegular   = regexp.MustCompile("^[a-z0-9A-Z\\p{Han}]+(_[a-z0-9A-Z\\p{Han}]+)*$")
	NameRegular     = regexp.MustCompile("^[a-zA-Z\\p{Han}]+[\\.-_a-z0-9A-Z\\p{Han}]{1,15}$")
	WeixinIdRegular = regexp.MustCompile(".*")
	PhoneRegular    = regexp.MustCompile("^(13[0-9]|14[0-9]|15[0-9]|17[0-9]|18[0-9]|19[0-9])\\d{8}$")
	IdcardRegular   = regexp.MustCompile("^(\\d{6})(18|19|20)?(\\d{2})([01]\\d)([0123]\\d)(\\d{3})(\\d|X)?$")
	RfidRegular     = regexp.MustCompile(".*")
	IMEIRegular     = regexp.MustCompile(".*")
	EmailRegular    = regexp.MustCompile("^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+")
)

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
	_, err := orm.Id(self.Appid).MustCols("Appid", "Appkey", "Remark", "Status", "BefeFlag", "Json").Update(self)
	return err
}

func Search() ([]Appid, error) {
	list := make([]Appid, 0)
	err := orm.Find(&list)
	return list, err
}

func IsIdExist(id int64) (bool, error) {

	return orm.Get(&User{Id: id})
}

func IsNameExist(siteid int64, name string) (bool, error) {
	if siteid == 0 {
		return false, ErrSiteIdNotExist
	}

	if len(name) == 0 {
		return false, nil
	}
	return orm.Get(&User{SiteId: siteid, Name: strings.ToLower(name)})
}

func IsPhoneExist(siteid int64, phone string) (bool, error) {
	if siteid == 0 {
		return false, ErrSiteIdNotExist
	}

	if len(phone) == 0 {
		return false, nil
	}
	return orm.Get(&User{SiteId: siteid, Phone: phone})
}

func IsPhoneUniqueExist(siteid int64, phone string) (bool, error) {
	if siteid == 0 {
		return false, ErrSiteIdNotExist
	}

	if len(phone) == 0 {
		return false, nil
	}
	return orm.Get(&User{Phone: phone})
}

// IsEmailUsed returns true if the e-mail has been used.
func IsEmailExist(siteid int64, email string) (bool, error) {
	if siteid == 0 {
		return false, ErrSiteIdNotExist
	}

	if len(email) == 0 {
		return false, nil
	}
	return orm.Get(&User{SiteId: siteid, Email: email})
}

func IsIdcardExist(siteid int64, idcard string) (bool, error) {
	if siteid == 0 {
		return false, ErrSiteIdNotExist
	}

	if len(idcard) == 0 {
		return false, nil
	}
	return orm.Get(&User{SiteId: siteid, Idcard: idcard})
}

func IsIdcardUniqueExist(siteid int64, idcard string) (bool, error) {
	if siteid == 0 {
		return false, ErrSiteIdNotExist
	}

	if len(idcard) == 0 {
		return false, nil
	}
	return orm.Get(&User{Idcard: idcard})
}

func IsRfidExist(siteid int64, rfid string) (bool, error) {
	if siteid == 0 {
		return false, ErrSiteIdNotExist
	}

	if len(rfid) == 0 {
		return false, nil
	}
	return orm.Get(&User{SiteId: siteid, Rfid: rfid})
}

func IsIMEIExist(siteid int64, imei string) (bool, error) {
	if siteid == 0 {
		return false, ErrSiteIdNotExist
	}

	if len(imei) == 0 {
		return false, nil
	}
	return orm.Get(&User{SiteId: siteid, Imei: imei})
}

func IsWeixinidExist(siteid int64, weixinid string) (bool, error) {
	if len(weixinid) == 0 {
		return false, nil
	}
	return orm.Get(&User{Weixinid: weixinid})
}

func IsIdPhoneExist(siteid int64, idcard, phone string) (bool, error) {
	if siteid == 0 {
		return false, ErrSiteIdNotExist
	}
	if len(idcard) == 0 || len(phone) == 0 {
		return false, nil
	}
	return orm.Get(&User{Phone: phone, Idcard: idcard})
}

func RegisterUserByName(user *User) (*User, error) {
	if user.SiteId == 0 {
		return nil, ErrSiteIdNotExist
	}

	isExist, err := IsNameExist(user.SiteId, user.Name)
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

func RegisterUserByPhone(user *User) (*User, error) {
	if user.SiteId == 0 {
		return nil, ErrSiteIdNotExist
	}

	isExist, err := IsPhoneExist(user.SiteId, user.Phone)
	if err != nil {
		return nil, err
	} else if isExist {
		return nil, ErrEmailAlreadyUsed
	}

	user.RegisterType = REGISTER_TYPE_PHONE

	// user.Name = strings.ToLower(user.Name)
	if _, err = orm.Insert(user); err != nil {
		return nil, err
	}

	return user, err
}

func RegisterUserByIdcard(user *User) (*User, error) {
	if user.SiteId == 0 {
		return nil, ErrSiteIdNotExist
	}

	isExist, err := IsIdcardExist(user.SiteId, user.Idcard)
	if err != nil {
		return nil, err
	} else if isExist {
		return nil, ErrIdcardAlreadyUsed
	}

	user.RegisterType = REGISTER_TYPE_IDCARD

	// user.Name = strings.ToLower(user.Name)
	if _, err = orm.Insert(user); err != nil {
		return nil, err
	}

	return user, err
}

func RegisterUserByIdPhone(user *User) (*User, error) {
	if user.SiteId == 0 {
		return nil, ErrSiteIdNotExist
	}

	isExist, err := IsIdPhoneExist(user.SiteId, user.Idcard, user.Phone)
	if err != nil {
		return nil, err
	} else if isExist {
		return nil, ErrIdcardAlreadyUsed
	}

	// user.Name = strings.ToLower(user.Name)
	if _, err = orm.Insert(user); err != nil {
		return nil, err
	}

	return user, err
}

func RegisterUserByRfid(user *User) (*User, error) {
	if user.SiteId == 0 {
		return nil, ErrSiteIdNotExist
	}

	isExist, err := IsRfidExist(user.SiteId, user.Rfid)
	if err != nil {
		return nil, err
	} else if isExist {
		return nil, ErrRfidAlreadyUsed
	}

	// user.Name = strings.ToLower(user.Name)
	if _, err = orm.Insert(user); err != nil {
		return nil, err
	}

	return user, err
}

func RegisterUserByIMEI(user *User) (*User, error) {
	if user.SiteId == 0 {
		return nil, ErrSiteIdNotExist
	}

	isExist, err := IsIMEIExist(user.SiteId, user.Imei)
	if err != nil {
		return nil, err
	} else if isExist {
		return nil, ErrIMEIAlreadyUsed
	}

	user.RegisterType = REGISTER_TYPE_IMEI

	// user.Name = strings.ToLower(user.Name)
	if _, err = orm.Insert(user); err != nil {
		return nil, err
	}

	return user, err
}

func RegisterUserByEmail(user *User) (*User, error) {
	if user.SiteId == 0 {
		return nil, ErrSiteIdNotExist
	}

	isExist, err := IsEmailExist(user.SiteId, user.Email)
	if err != nil {
		return nil, err
	} else if isExist {
		return nil, ErrEmailAlreadyUsed
	}

	user.RegisterType = REGISTER_TYPE_EMAIL

	// user.Name = strings.ToLower(user.Name)
	if _, err = orm.Insert(user); err != nil {
		return nil, err
	}

	// user.Passwd = helper.Md5(user.Passwd + strconv.FormatInt(user.Id, 10))
	// orm.Id(user.Id).Cols("passwd").Update(user)
	return user, err
}

func RegisterUserByWeixinid(user *User) (*User, error) {
	if user.SiteId == 0 {
		return nil, ErrSiteIdNotExist
	}

	isExist, err := IsWeixinidExist(user.SiteId, user.Weixinid)
	if err != nil {
		return nil, err
	} else if isExist {
		return nil, ErrEmailAlreadyUsed
	}

	user.RegisterType = REGISTER_TYPE_WEIXINID

	// user.Name = strings.ToLower(user.Name)
	if _, err = orm.Insert(user); err != nil {
		return nil, err
	}

	// user.Passwd = helper.Md5(user.Passwd + strconv.FormatInt(user.Id, 10))
	// orm.Id(user.Id).Cols("passwd").Update(user)
	return user, err
}

func BindUserWeixinidById(user *User) (*User, error) {
	if user.SiteId == 0 {
		return nil, ErrSiteIdNotExist
	}

	isExist, err := IsIdExist(user.Id)
	if err != nil {
		return nil, err
	} else if !isExist {
		return nil, ErrUserNotExist
	}

	// user.Name = strings.ToLower(user.Name)
	if _, err = orm.Where(" site_id=? and id=? ", user.SiteId, user.Id).Cols("weixinid").Update(user); err != nil {
		return nil, err
	}

	return user, err
}

func BindUserSScardByPhone(user *User) (*User, error) {
	if user.SiteId == 0 {
		return nil, ErrSiteIdNotExist
	}

	isExist, err := IsPhoneExist(user.SiteId, user.Phone)
	if err != nil {
		return nil, err
	} else if !isExist {
		return nil, ErrUserNotExist
	}

	// user.Name = strings.ToLower(user.Name)
	if _, err = orm.Where(" site_id=? and phone=? ", user.SiteId, user.Phone).Cols("sscard", "idcard").Update(user); err != nil {
		return nil, err
	}

	return user, err
}
func BindUserSScardByIdcard(user *User) (*User, error) {
	if user.SiteId == 0 {
		return nil, ErrSiteIdNotExist
	}

	isExist, err := IsIdcardExist(user.SiteId, user.Idcard)
	if err != nil {
		return nil, err
	} else if !isExist {
		return nil, ErrUserNotExist
	}

	// user.Name = strings.ToLower(user.Name)
	if _, err = orm.Where(" site_id=? and idcard=? ", user.SiteId, user.Idcard).Cols("sscard", "phone").Update(user); err != nil {
		return nil, err
	}

	return user, err
}

func GetUser(user *User) (*User, error) {
	has, err := orm.Get(user)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrUserNotExist
	}
	return user, nil
}

func GetUserById(siteid int64, id int64) (*User, error) {
	// if siteid == 0 {
	// 	return nil, ErrSiteIdNotExist
	// }

	user := new(User)
	has, err := orm.Id(id).Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrUserNotExist
	}
	return user, nil
}

func GetUserByIdPhone(siteid int64, idcard, phone string) (*User, error) {

	if siteid == 0 {
		return nil, ErrSiteIdNotExist
	}

	user := &User{Phone: phone, Idcard: idcard}
	has, err := orm.Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrUserNotExist
	}
	return user, nil
}

// GetUserByName returns the user object by given name if exists.
func GetUserByName(siteid int64, name string) (*User, error) {
	if siteid == 0 {
		return nil, ErrSiteIdNotExist
	}
	if len(name) == 0 {
		return nil, ErrUserNotExist
	}
	user := &User{SiteId: siteid, Name: strings.ToLower(name)}
	has, err := orm.Get(user)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrUserNotExist
	}
	return user, nil
}

// GetUserEmailsByNames returns a slice of e-mails corresponds to names.
/*
func GetUserEmailsByNames(names []string) []string {
	mails := make([]string, 0, len(names))
	for _, name := range names {
		u, err := GetUserByName(name)
		if err != nil {
			continue
		}
		mails = append(mails, u.Email)
	}
	return mails
}
*/

func GetUserByPhone(siteid int64, phone string) (*User, error) {
	if siteid == 0 {
		return nil, ErrSiteIdNotExist
	}
	if len(phone) == 0 {
		return nil, ErrUserNotExist
	}
	user := &User{SiteId: siteid, Phone: phone}
	has, err := orm.Get(user)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrUserNotExist
	}
	return user, nil
}

func GetUserByPhoneUnique(siteid int64, phone string) (*User, error) {
	if siteid == 0 {
		return nil, ErrSiteIdNotExist
	}
	if len(phone) == 0 {
		return nil, ErrUserNotExist
	}
	user := &User{Phone: phone}
	has, err := orm.Get(user)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrUserNotExist
	}
	return user, nil
}

func GetUserByIdcard(siteid int64, idcard string) (*User, error) {
	if siteid == 0 {
		return nil, ErrSiteIdNotExist
	}
	if len(idcard) == 0 {
		return nil, ErrUserNotExist
	}
	user := &User{SiteId: siteid, Idcard: idcard}
	has, err := orm.Get(user)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrUserNotExist
	}
	return user, nil
}

func GetUserByIdcardUnique(siteid int64, idcard string) (*User, error) {
	if siteid == 0 {
		return nil, ErrSiteIdNotExist
	}
	if len(idcard) == 0 {
		return nil, ErrUserNotExist
	}
	user := &User{Idcard: idcard}
	has, err := orm.Get(user)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrUserNotExist
	}
	return user, nil
}

func GetUserByRfid(siteid int64, rfid string) (*User, error) {
	if siteid == 0 {
		return nil, ErrSiteIdNotExist
	}
	if len(rfid) == 0 {
		return nil, ErrUserNotExist
	}
	user := &User{SiteId: siteid, Rfid: rfid}
	has, err := orm.Get(user)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrUserNotExist
	}
	return user, nil
}

func GetUserByIMEI(siteid int64, imei string) (*User, error) {
	if siteid == 0 {
		return nil, ErrSiteIdNotExist
	}
	if len(imei) == 0 {
		return nil, ErrUserNotExist
	}
	user := &User{SiteId: siteid, Imei: imei}
	has, err := orm.Get(user)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrUserNotExist
	}
	return user, nil
}

func GetUserByWeixinid(weixinid string) (*User, error) {
	if len(weixinid) == 0 {
		return nil, ErrUserNotExist
	}
	user := &User{Weixinid: weixinid}
	has, err := orm.Get(user)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrUserNotExist
	}
	return user, nil
}

func GetUserByEmail(siteid int64, email string) (*User, error) {
	if siteid == 0 {
		return nil, ErrSiteIdNotExist
	}
	if len(email) == 0 {
		return nil, ErrUserNotExist
	}
	user := &User{SiteId: siteid, Email: strings.ToLower(email)}
	has, err := orm.Get(user)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrUserNotExist
	}
	return user, nil
}

func LoginUserAll(name, passwd, Salt string) (*User, error) {
	if len(name) == 0 {
		return nil, ErrUserNotExist
	}
	user := new(User)
	rows, err := orm.Where("name = ? or email = ? or phone = ? or idcard = ? or rfid = ?", strings.ToLower(name), name, name, name, name).Rows(user)
	if err != nil {
		return nil, ErrUserNotExist
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(user)
		if err != nil {
			return nil, ErrUserNotExist
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

	return nil, ErrUserNotExist
}

func LoginUserByid(siteid int64, id int64, passwd, Salt string) (*User, error) {
	if siteid <= 0 {
		return nil, ErrSiteIdNotExist
	}

	user, err := GetUserById(siteid, id)
	if err != nil {
		return nil, err
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
	return nil, ErrUserNotExist
}

func LoginUserPlain(siteid int64, name, passwd, Salt string) (*User, error) {
	if siteid <= 0 {
		return nil, ErrSiteIdNotExist
	}

	if len(name) == 0 {
		return nil, ErrUserNotExist
	}
	user := new(User)
	rows, err := orm.Where("name = ? or email = ? or phone = ? or idcard = ? or rfid = ?", strings.ToLower(name), name, name, name, name).And("site_id = ?", siteid).Rows(user)
	if err != nil {
		return nil, ErrUserNotExist
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(user)
		if err != nil {
			return nil, ErrUserNotExist
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

	// user := User{Name: strings.ToLower(name)}
	// has, err := orm.Get(&user)
	// if err != nil {
	// 	return nil, err
	// } else if !has {
	// 	return nil, ErrUserNotExist
	// }

	// if user.Passwd != passwd {
	// 	return nil, ErrUserNotExist
	// }
	// return &user, nil
	return nil, ErrUserNotExist
}

/*
	ltype : 登录类型
			1： 短信验证码登录
			2： 密码登录
*/
func LoginUserByPhoneCode(name, passwd, Salt string, ltype int) (*User, error) {
	// if siteid == 0 {
	// 	return nil, ErrSiteIdNotExist
	// }

	if len(name) == 0 {
		return nil, ErrUserNotExist
	}
	user := new(User)
	rows, err := orm.Where("phone = ? ", name).Rows(user) //.And("site_id = ?", siteid).Rows(user)
	if err != nil {
		return nil, ErrUserNotExist
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(user)
		if err != nil {
			return nil, ErrUserNotExist
		}
		if ltype != 1 {
			if user.Passwd == passwd {
				return user, nil
			}

			if helper.Md5(user.Passwd+Salt) == passwd {
				return user, nil
			}

			if helper.Md5(helper.Md5(user.Passwd)+Salt) == passwd {
				return user, nil
			}
		} else {
			return user, nil
		}
	}

	// user := User{Name: strings.ToLower(name)}
	// has, err := orm.Get(&user)
	// if err != nil {
	// 	return nil, err
	// } else if !has {
	// 	return nil, ErrUserNotExist
	// }

	// if user.Passwd != passwd {
	// 	return nil, ErrUserNotExist
	// }
	// return &user, nil
	return nil, ErrUserNotExist
}

func ResetPasswd(siteid int64, id int64, oldpasswd string, newpasswd string, salt string) (bool, error) {
	user, err := GetUserById(siteid, id)
	if err != nil {
		return false, err
	}

	if user.Passwd != oldpasswd && helper.Md5(user.Passwd+salt) != oldpasswd && helper.Md5(helper.Md5(user.Passwd)+salt) != oldpasswd {
		return false, ErrWarnPasswd
	}

	user.Passwd = newpasswd
	affected, err := orm.Id(id).Where("site_id = ?", siteid).Cols("passwd").Update(user)
	return affected == 1, err
}

func SetPasswd(siteid int64, id int64, newpasswd string) (bool, error) {
	user, err := GetUserById(siteid, id)
	if err != nil {
		return false, err
	}
	user.Passwd = newpasswd
	affected, err := orm.Id(id).Where("site_id = ?", siteid).Cols("passwd").Update(user)
	return affected == 1, err
}

func SetUserName(siteid int64, id int64, name string) (bool, error) {
	v, err := IsNameExist(siteid, name)
	if err != nil {
		return false, err
	}

	if v == true {
		return false, ErrNameAlreadyUsed
	}

	user, err := GetUserById(siteid, id)
	if err != nil {
		return false, err
	}

	if user.Name != "" {
		return false, ErrNameAlreadyUsed
	}

	user.Name = name
	affected, err := orm.Id(id).Where("site_id = ?", siteid).Cols("name").Update(user)
	return affected == 1, err
}

func SetUserName2(siteid int64, id int64, name string) (bool, error) {
	v, err := IsNameExist(siteid, name)
	if err != nil {
		return false, err
	}

	if v == true {
		return false, ErrNameAlreadyUsed
	}

	user, err := GetUserById(siteid, id)
	if err != nil {
		return false, err
	}

	user.Name = name
	affected, err := orm.Id(id).Where("site_id = ?", siteid).Cols("name").Update(user)
	return affected == 1, err
}

func SetUserPhone(siteid int64, id int64, phone string) (bool, error) {
	v, err := IsPhoneExist(siteid, phone)
	if err != nil {
		return false, err
	}

	if v == true {
		return false, ErrPhoneAlreadyUsed
	}

	user, err := GetUserById(siteid, id)
	if err != nil {
		return false, err
	}

	if user.Phone != "" {
		return false, ErrPhoneAlreadyUsed
	}

	user.Phone = phone
	affected, err := orm.Id(id).Where("site_id = ?", siteid).Cols("phone").Update(user)
	return affected == 1, err
}

func SetUserIdcard(siteid int64, id int64, idcard string) (bool, error) {
	v, err := IsIdcardExist(siteid, idcard)
	if err != nil {
		return false, err
	}

	if v == true {
		return false, ErrIdcardAlreadyUsed
	}

	user, err := GetUserById(siteid, id)
	if err != nil {
		return false, err
	}

	if user.Idcard != "" {
		return false, ErrIdcardAlreadyUsed
	}

	user.Idcard = idcard
	affected, err := orm.Id(id).Where("site_id = ?", siteid).Cols("idcard").Update(user)
	return affected == 1, err
}

func SetUserRfid(siteid int64, id int64, rfid string) (bool, error) {
	v, err := IsRfidExist(siteid, rfid)
	if err != nil {
		return false, err
	}

	if v == true {
		return false, ErrRfidAlreadyUsed
	}

	user, err := GetUserById(siteid, id)
	if err != nil {
		return false, err
	}

	if user.Rfid != "" {
		return false, ErrRfidAlreadyUsed
	}

	user.Rfid = rfid
	affected, err := orm.Id(id).Where("site_id = ?", siteid).Cols("rfid").Update(user)
	return affected == 1, err
}

func SetUserRfid2(siteid int64, id int64, rfid string) (bool, error) {
	v, err := IsRfidExist(siteid, rfid)
	if err != nil {
		return false, err
	}

	if v == true {
		return false, ErrRfidAlreadyUsed
	}

	user, err := GetUserById(siteid, id)
	if err != nil {
		return false, err
	}

	user.Rfid = rfid
	affected, err := orm.Id(id).Where("site_id = ?", siteid).Cols("rfid").Update(user)
	return affected == 1, err
}

func UnsetUserRfid(siteid int64, id int64) (bool, error) {
	user, err := GetUserById(siteid, id)
	if err != nil {
		return false, err
	}

	if user.Rfid == "" {
		return false, ErrRfidNotUsed
	}

	user.Rfid = ""
	affected, err := orm.Id(id).Where("site_id = ?", siteid).Cols("rfid").Update(user)
	return affected == 1, err
}

func SetUserWeixinid(siteid int64, id int64, weixinid string) (bool, error) {
	v, err := IsWeixinidExist(siteid, weixinid)
	if err != nil {
		return false, err
	}

	if v == true {
		return false, ErrWeixinidAlreadyUsed
	}

	user, err := GetUserById(siteid, id)
	if err != nil {
		return false, err
	}

	if user.Weixinid != "" {
		return false, ErrWeixinidAlreadyUsed
	}

	user.Weixinid = weixinid
	affected, err := orm.Id(id).Where("site_id = ?", siteid).Cols("weixinid").Update(user)
	return affected == 1, err
}

func UnSetUserWeixinid(siteid int64, id int64, weixinid string) (bool, error) {

	user, err := GetUserById(siteid, id)
	if err != nil {
		return false, err
	}

	user.Weixinid = weixinid
	affected, err := orm.Id(id).Cols("weixinid").Update(user)
	return affected == 1, err
}

func SetUserEmail(siteid int64, id int64, email string) (bool, error) {
	v, err := IsEmailExist(siteid, email)
	if err != nil {
		return false, err
	}

	if v == true {
		return false, ErrEmailAlreadyUsed
	}

	user, err := GetUserById(siteid, id)
	if err != nil {
		return false, err
	}

	if user.Email != "" {
		return false, ErrEmailAlreadyUsed
	}

	user.Email = email
	affected, err := orm.Id(id).Where("site_id = ?", siteid).Cols("email").Update(user)
	return affected == 1, err
}

func GetUserProfileByUserObject(user *User, sType string, keys []string) (map[string]interface{}, error) {
	var d = make(map[string]interface{})
	sType = strings.ToLower(sType)
	var uj map[string]map[string]interface{}
	if user.Json == "" {
		uj = make(map[string]map[string]interface{})
	} else {
		err := json.Unmarshal([]byte(user.Json), &uj)
		if err != nil {
			return nil, err
		}
	}

	if uj["SHARE"] == nil {
		uj["SHARE"] = make(map[string]interface{})
	}
	share := uj["SHARE"]

	if uj[sType] == nil {
		uj[sType] = make(map[string]interface{})
	}
	js := uj[sType]

	r := reflect.ValueOf(user)
	for _, k := range keys {
		if k == "Passwd" {
			continue
		}
		//if k == "Idcard" {
		//	d["idcard"] = user.Idcard
		//}
		//if k == "Phone" {
		//	d["phone"] = user.Phone
		//}
		s := string(k[0])
		if s == strings.ToUpper(s) {
			f := reflect.Indirect(r).FieldByName(k)
			if f.IsValid() == true {
				switch f.Kind() {
				case reflect.String:
					d[k] = f.String()
				case reflect.Int64:
					d[k] = f.Int()
				default:
				}
			} else {
				d[k] = share[k]
			}
		} else {
			d[k] = js[k]
		}
	}

	return d, nil
}

func GetUserProfile(siteid int64, id int64, sType string, keys []string) (map[string]interface{}, error) {
	user, err := GetUserById(siteid, id)
	if err != nil {
		return nil, err
	}
	return GetUserProfileByUserObject(user, sType, keys)
}

func SetUserProfile(siteid int64, id int64, sType string, d map[string]interface{}) (bool, error) {
	sType = strings.ToLower(sType)
	if sType == "share" {
		return false, errors.New("type can't set share")
	}

	user, err := GetUserById(siteid, id)
	if err != nil {
		return false, err
	}

	var uj map[string](map[string]interface{})

	if user.Json == "" {
		uj = make(map[string](map[string]interface{}))
	} else {
		err := json.Unmarshal([]byte(user.Json), &uj)
		if err != nil {
			return false, err
		}
	}

	if uj["SHARE"] == nil {
		uj["SHARE"] = make(map[string]interface{})
	}
	share := uj["SHARE"]

	if uj[sType] == nil {
		uj[sType] = make(map[string]interface{})
	}
	js := uj[sType]

	r := reflect.ValueOf(user)
	for k, v := range d {
		s := string(k[0])
		if s == strings.ToUpper(s) {
			f := reflect.Indirect(r).FieldByName(k)
			if f.IsValid() && f.CanSet() {
				r.Elem().FieldByName(k).Set(reflect.ValueOf(v))
			} else {
				share[k] = v
			}
		} else {
			js[k] = v
		}
	}

	ds, err := json.Marshal(&uj)
	if err != nil {
		return false, err
	}
	user.Json = string(ds)
	affected, err := orm.Id(id).Cols("nickname", "json").Update(user)
	return affected == 1, err

}

///NameId
func GetBind(siteid int64, stype string, name string) (int64, error) {
	if siteid == 0 {
		return 0, ErrSiteIdNotExist
	}

	name2id := &NameLog{SiteId: siteid, Type: strings.ToLower(stype), Name: strings.ToLower(name)}
	has, err := orm.Get(name2id)
	if err != nil {
		return 0, err
	} else if !has {
		return 0, nil
	}
	return name2id.Id, nil
}

func SetBind(name2id *NameLog) (*NameLog, error) {
	id, err := GetBind(name2id.SiteId, name2id.Type, name2id.Name)
	if err != nil {
		return nil, err
	} else if id > 0 {
		return nil, ErrNameAlreadyUsed
	}

	name2id.Type = strings.ToLower(name2id.Type)
	name2id.Name = strings.ToLower(name2id.Name)
	if _, err = orm.Insert(name2id); err != nil {
		return nil, err
	}
	return name2id, err
}

//Feedback
func SubmitFeedback(feedback *Feedback) (*Feedback, error) {
	if feedback.SiteId == 0 {
		return nil, ErrSiteIdNotExist
	}

	if _, err := orm.Insert(feedback); err != nil {
		return nil, err
	}
	return feedback, nil
}

func FeedbackList(siteid int64, row int, page int) ([]Feedback, error) {
	feedbacks := make([]Feedback, 0)
	err := orm.Where("site_id = ? ", siteid).Desc("Id").Limit(row, row*page).Find(&feedbacks)
	return feedbacks, err
}

//list
func UserList(siteid int64, name, phone, card, rfid, idphone string, row int, page int, sType string, keys []string) ([]interface{}, error) {
	var d = make([]interface{}, 0)

	m := orm.Asc("Id").And("site_id = ?", siteid)
	if name != "" {
		m = m.And("name = ?", name)
	}

	if phone != "" {
		m = m.And("phone = ?", phone)
	}

	if card != "" {
		m = m.And("card = ?", card)
	}

	if rfid != "" {
		m = m.And("rfid = ?", rfid)
	}

	if idphone != "" {
		m = m.And("(idcard = ? or phone= ?)", idphone)
	}

	rows, err := m.Limit(row, row*page).Rows(new(User))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := new(User)
		err = rows.Scan(user)
		if err != nil {
			return nil, err
		}
		d1, err := GetUserProfileByUserObject(user, sType, keys)
		if err != nil {
			return nil, err
		}
		d = append(d, d1)
	}
	return d, nil
}

func SearchByNickName(siteid int64, nickname string, sType string, keys []string) ([]interface{}, error) {
	var d = make([]interface{}, 0)

	m := orm.Asc("Id").And("site_id = ?", siteid)
	if nickname != "" {
		m = m.And("nickname like ?", "%"+nickname+"%")
	}
	row := 1000000
	page := 0

	rows, err := m.Limit(row, row*page).Rows(new(User))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := new(User)
		err = rows.Scan(user)
		if err != nil {
			return nil, err
		}
		d1, err := GetUserProfileByUserObject(user, sType, keys)
		if err != nil {
			return nil, err
		}
		d = append(d, d1)
	}
	return d, nil
}

//userlog
func AddLog(userlog *UserLog) (*UserLog, error) {
	if userlog.SiteId == 0 {
		return nil, ErrSiteIdNotExist
	}

	var err error
	if _, err = orm.Insert(userlog); err != nil {
		return nil, err
	}
	return userlog, err
}

func LogList(siteid int64, userid int64, row int, page int, level int) ([]UserLog, error) {
	userlogs := make([]UserLog, 0)
	err := orm.Where("site_id = ? and user_id = ? and level <= ?", siteid, userid, level).Desc("Id").Limit(row, row*page).Find(&userlogs)
	return userlogs, err
}

func (self *User) Update() error {
	_, err := orm.Id(self.Id).Update(self)
	return err
}
