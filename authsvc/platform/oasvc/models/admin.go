package models

import (
	"time"
	// "github.com/martini-contrib/binding"

	"platform/oasvc/dbmodels"
)

type (
	AdminLoginForm struct {
		Name   string `form:"name" binding:"required"`
		Passwd string `form:"passwd" binding:"required"`
	}

	AdminResetPasswdForm struct {
		Id        int64  `form:"id" binding:"required"`
		OldPasswd string `form:"oldpasswd" binding:"required"`
		NewPasswd string `form:"newpasswd" binding:"required"`
	}

	AdminRegistryForm struct {
		RoleId        interface{} `form:"roleid" json:"roleid"`
		SiteId        interface{} `form:"siteid" binding:"required" json:"siteid"`
		Type          interface{} `form:"type" json:"type"`
		Name          string      `form:"name" binding:"required" json:"name"`
		JobNumber     string      `form:"jobnumber" json:"jobnumber"`
		Passwd        string      `form:"passwd" binding:"required" json:"passwd"`
		Email         string      `form:"email" json:"email"`
		Phone         string      `form:"phone" json:"phone"`
		Realname      string      `form:"realname" json:"realname"`
		Description   string      `form:"description" json:"description"`
		EffectiveTime time.Time   `form:"effectiveTime" json:"effectivetime"`
		ExpireTime    time.Time   `form:"expireTime" json:"expiretime"`

		Created  time.Time   `form:"created" json:"created"`
		Updated  time.Time   `form:"updated" json:"updated"`
		ImageUrl string      `form:"imageurl" json:"imageurl"`
		State    interface{} `form:"gender" json:"gender"` //0：正常，1：被阻止登录，默认为0
		Gender   interface{}
	}

	AdminModifyForm struct {
		Id            interface{} `form:"id" binding:"required" json:"id"`
		SiteId        interface{} `form:"siteid" json:"siteid,omitempty"`
		RoleId        interface{} `form:"roleid" json:"roleid,omitempty"`
		Type          interface{} `form:"type" json:"type,omitempty"`
		Name          string      `form:"name" json:"name,omitempty"`
		JobNumber     string      `form:"jobnumber" json:"jobnumber,omitempty"`
		Passwd        string      `form:"passwd" json:"passwd,omitempty"`
		Email         string      `form:"email" json:"email,omitempty"`
		Phone         string      `form:"phone" json:"phone,omitempty"`
		Realname      string      `form:"realname" json:"realname,omitempty"`
		Description   string      `form:"description" json:"description,omitempty"`
		EffectiveTime time.Time   `form:"effectiveTime" json:"effectivetime,omitempty"`
		ExpireTime    time.Time   `form:"expireTime" json:"expiretime,omitempty"`

		Created  time.Time   `form:"created" json:"created,omitempty"`
		Updated  time.Time   `form:"updated" json:"updated,omitempty"`
		ImageUrl string      `form:"imageurl" json:"imageurl,omitempty"`
		State    interface{} `form:"state" json:"state,omitempty"` //0：正常，1：被阻止登录，默认为0
		Gender   interface{} `form:"gender" json:"gender,omitempty"`
	}

	AdminResetPasswdCodeForm struct {
		Type  string `form:"type" binding:"required"`
		Value string `form:"value" binding:"required"`
	}

	AdminResetPasswdFormByCode struct {
		Type   string `form:"type" binding:"required"`
		Value  string `form:"value" binding:"required"`
		Passwd string `form:"passwd" binding:"required"`
		Code   string `form:"code" binding:"required"`
	}

	AdminDeleteForm struct {
		Id []interface{} `form:"id" binding:"required" json:"id"`
	}

	AdminRegCodeResp struct {
		Ret  int
		Code string
	}

	AdminResp struct {
		Ret  int
		Data dbmodels.Admin
	}

	AdminListResp struct {
		Ret   int
		Datas []dbmodels.Admin
	}
)
