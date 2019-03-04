package models

import "time"

type SysDomainConf struct {
	DomainId		int64		`orm:"pk;auto" gorm:"primary_key;AUTO_INCREMENT"`
	Siteid			int64
	Domain			string		`orm:"size(50)"`
	Domaingrp		int
	Name			string		`orm:"size(255)"`
	Keyid			int64
	Value			string		`orm:"size(255)"`
	Action 			string		`orm:"size(255)"`
}

type SysSiteConf struct {
	SiteDns    string 			`orm:"pk" gorm:"primary_key;AUTO_INCREMENT"`
	SkinStyle  string			`orm:"type(text)"`
	RootAreaId int64
	Status     int
}

type Appid struct {
	Appid    int64 		`orm:"pk" gorm:"primary_key;AUTO_INCREMENT"`
	SiteId  int64       `gorm:"siteid not null"`
	Appkey     string `gorm:"not null"`
	Remark  string    `gorm:"not null"`
	Status  int       `gorm:"tinyint(4) not null default 1"`
	Json    string
	Created time.Time `gorm:"not null"`
	Updated time.Time `gorm:"not null"`
	BefeFlag int8 `gorm:"type:tinyint(4) not null default 0"`
}

