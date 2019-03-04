package models

import (
	"net/http"

	"github.com/martini-contrib/binding"

	"platform/ousvc/dbmodels"
)

type (
	UserRegCodeForm struct {
		Type  string `form:"type" binding:"required"`
		Value string `form:"value" binding:"required"`
	}

	UserResetPasswdCodeForm struct {
		Type  string `form:"type" binding:"required"`
		Value string `form:"value" binding:"required"`
	}

	UserCheckForm struct {
		Type   string `form:"type" binding:"required"`
		Value  string `form:"value" binding:"required"`
		Phone  string `form:"phone"`
		Idcard string `form:"idcard"`
	}

	UserCheckRegcodeForm struct {
		Type  string `form:"type" binding:"required"`
		Value string `form:"value" binding:"required"`
		Code  string `form:"code" binding:"required"`
	}

	UserRegistryForm struct {
		Type   string `form:"type" binding:"required"`
		Value  string `form:"value" binding:"required"`
		Passwd string `form:"passwd" binding:"required"`
		Code   string `form:"code" binding:"required"`
		Name   string `form:"name"`
		Phone  string `form:"phone"`
		Idcard string `form:"idcard"`
	}

	UserBindSscardForm struct {
		Type   string `form:"type" binding:"required"`
		Value  string `form:"value" binding:"required"`
		Sscard string `form:"sscard" binding:"required"`
		Code   string `form:"code" binding:"required"`
		Name   string `form:"name" binding:"required"`
		Idcard string `form:"idcard" binding:"required"`
		Phone  string `form:"phone" binding:"required"`
	}

	UserLoginForm struct {
		Name    string `form:"name" binding:"required"`
		Passwd  string `form:"passwd" binding:"required"`
		Captcha string `form:"captcha"`
	}

	UserBindWxForm struct {
		Phone    string `form:"phone" binding:"required"`
		Code     string `form:"code" binding:"required"`
		Weixinid string `form:"weixinid" binding:"required"`
		Captcha  string `form:"captcha"`
	}

	UserUnBindWxForm struct {
		UserId   int64  `form:"userid" binding:"required"`
		Weixinid string `form:"weixinid" binding:"required"`
	}

	UserLoginByIdForm struct {
		UserId  int64  `form:"userid" binding:"required"`
		Passwd  string `form:"passwd" binding:"required"`
		Captcha string `form:"captcha"`
	}

	UserLoginByIMEIForm struct {
		IMEI string `form:"imei" binding:"required"`
	}
	UserLoginByWeixinForm struct {
		Weixinid string `form:"weixinid" binding:"required"`
	}

	UserPatientLoginForm struct {
		PatientId int64  `form:"patientid" binding:"required"`
		Timestamp string `form:"timestamp" binding:"required"`
		Auth      string `form:"auth" binding:"required"`
	}

	UserResetPasswdForm struct {
		Type   string `form:"type" binding:"required"`
		Value  string `form:"value" binding:"required"`
		Passwd string `form:"passwd" binding:"required"`
		Code   string `form:"code" binding:"required"`
	}

	UserPasswdForm struct {
		Oldpasswd string `form:"oldpasswd" binding:"required"`
		Newpasswd string `form:"newpasswd" binding:"required"`
	}

	UserCloneTokenForm struct {
		UserId    int64  `form:"userid" binding:"required"`
		Timestamp string `form:"timestamp" binding:"required"`
		Auth      string `form:"auth" binding:"required"`
	}

	UserFindByNickName struct {
		NickName string `form:"nickname" binding:"required"`
		Type     string `form:"type" binding:"required"`
		Keys     string `form:"keys" binding:"required"`
	}

	UserCloneToken struct {
		Ret       int
		UserId    int64
		Timestamp string
		Auth      string
	}

	UserIdResp struct {
		Ret    int
		UserId int64
	}

	UserIdSiteIdResp struct {
		Ret      int
		UserId   int64
		SiteId   int64
		Token    string
		NickName string
	}

	UserNickNameResp struct {
		Ret      int
		UserId   int64
		SiteId   int64
		NickName string
	}

	UserCaptchaIdResp struct {
		Ret       int
		CaptchaId string
	}

	UserRegCodeResp struct {
		Ret  int
		Code string
	}

	UserProfileResp struct {
		Ret  int
		Data map[string]interface{}
	}

	UserLogListResp struct {
		Ret   int
		Datas []dbmodels.UserLog
	}

	FeedbackListResp struct {
		Ret   int
		Datas []dbmodels.Feedback
	}
)

func CheckType(stype string, svalue string, sId, sPhone string, errors binding.Errors, req *http.Request) binding.Errors {
	switch stype {
	case "name":
		if !dbmodels.NameRegular.MatchString(svalue) {
			errors = append(errors, binding.Error{
				FieldNames:     []string{"value"},
				Classification: "error",
				Message:        "用户名格式错误",
			})
		}
	case "weixinid":
		if !dbmodels.WeixinIdRegular.MatchString(svalue) {
			errors = append(errors, binding.Error{
				FieldNames:     []string{"value"},
				Classification: "error",
				Message:        "微信id格式错误",
			})
		}
	case "phone":
		if !dbmodels.PhoneRegular.MatchString(svalue) {
			errors = append(errors, binding.Error{
				FieldNames:     []string{"value"},
				Classification: "error",
				Message:        "手机号码格式错误",
			})
		}
	case "idcard":
		if !dbmodels.IdcardRegular.MatchString(svalue) {
			errors = append(errors, binding.Error{
				FieldNames:     []string{"value"},
				Classification: "error",
				Message:        "身份证号码格式错误",
			})
		}
	case "email":
		if !dbmodels.EmailRegular.MatchString(svalue) {
			errors = append(errors, binding.Error{
				FieldNames:     []string{"value"},
				Classification: "error",
				Message:        "email格式错误",
			})
		}
	case "rfid":
		if !dbmodels.RfidRegular.MatchString(svalue) {
			errors = append(errors, binding.Error{
				FieldNames:     []string{"value"},
				Classification: "error",
				Message:        "rfid格式错误",
			})
		}
	case "imei":
		if !dbmodels.IMEIRegular.MatchString(svalue) {
			errors = append(errors, binding.Error{
				FieldNames:     []string{"value"},
				Classification: "error",
				Message:        "imei格式错误",
			})
		}
	case "idphone":
		if !dbmodels.IdcardRegular.MatchString(sId) || !dbmodels.PhoneRegular.MatchString(sPhone) {
			errors = append(errors, binding.Error{
				FieldNames:     []string{"value"},
				Classification: "error",
				Message:        "idcard或者phone格式错误",
			})
		}
	default:
		errors = append(errors, binding.Error{
			FieldNames:     []string{"type"},
			Classification: "error",
			Message:        "类型错误",
		})
	}
	return errors
}

func (rf UserRegCodeForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	errors = CheckType(rf.Type, rf.Value, "", "", errors, req)
	if len(errors) > 0 {
		return errors
	}
	return errors
}

func (rf UserResetPasswdCodeForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	errors = CheckType(rf.Type, rf.Value, "", "", errors, req)
	if len(errors) > 0 {
		return errors
	}
	return errors
}

func (rf UserRegistryForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return CheckType(rf.Type, rf.Value, rf.Idcard, rf.Phone, errors, req)
}

func (rf UserResetPasswdForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return CheckType(rf.Type, rf.Value, "", "", errors, req)
}

func (rf UserCheckForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return CheckType(rf.Type, rf.Value, rf.Idcard, rf.Phone, errors, req)
}

func (rf UserLoginForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	if !dbmodels.NameRegular.MatchString(rf.Name) && !dbmodels.PhoneRegular.MatchString(rf.Name) && !dbmodels.IdcardRegular.MatchString(rf.Name) && !dbmodels.EmailRegular.MatchString(rf.Name) && !dbmodels.RfidRegular.MatchString(rf.Name) {
		return append(errors, binding.Error{
			FieldNames:     []string{"name"},
			Classification: "error",
			Message:        "用户名格式错误",
		})
	}
	return errors
}
