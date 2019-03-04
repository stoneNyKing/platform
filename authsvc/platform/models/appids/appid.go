package appids

import (
	"errors"
	"time"
)

var (
	ErrSiteAlreadyExist   = errors.New("Site already exist")
	ErrSiteNotExist       = errors.New("Site does not exist")
	ErrParentSiteNotExist = errors.New("Parent Site does not exist")
)

type Appid struct {
	Appid   int64     `form:"id" xorm:"int(11) pk not null"` //OutTradeNo 微信对接订单号
	SiteId  int       `xorm:"siteid not null"`
	Appkey  string    	`xorm:"not null"`
	Remark  string    `xorm:"not null"`
	Status  int       `xorm:"tinyint(4) not null default 1"`
	Json    map[string]string					   `xorm:"json"`
	Created time.Time `xorm:"created not null"`
	Updated time.Time `xorm:"updated not null"`
	BefeFlag int8 `xorm:"tinyint(4) not null default 0"`
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
	_, err := orm.Id(self.Appid).MustCols("SiteId", "Appid", "Appkey", "Remark", "Enable","BefeFlag").Update(self)
	return err
}

func Search() ([]Appid, error) {
	list := make([]Appid, 0)
	err := orm.Find(&list)
	return list, err
}
