package apis

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/boj/redistore"
	"github.com/dchest/captcha"
	"github.com/garyburd/redigo/redis"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"

	"io/ioutil"
	"platform/common/utils"
	"platform/ousvc/config"
	. "platform/ousvc/models"
	//"platform/caller/base"
	"platform/lib/sender"
	. "platform/models"
	"platform/ousvc/dbmodels"
)

var pool *redis.Pool
var UserTraceTab map[string]string

func InitUsers() {
	dbmodels.NewEngine(config.Config.DbDriver, config.Config.DbAddr, config.Config.DbPort, config.Config.DbUser, config.Config.DbPasswd, config.Config.Database, config.Config.ObjectsSchema)
}

func userloger(c martini.Context, req *http.Request, token Token) {
	var userlog dbmodels.UserLog
	userlog.Json = make(map[string]interface{})
	c.Map(&userlog)
	c.Next()
	if userlog.Msg != "" {
		userlog.UserId = token.UserId
		userlog.SiteId = token.SiteId
		userlog.Appid = token.Appid
		userlog.Token = token.Token
		userlog.Level = 1
		userlog.Ip = token.IP
		userlog.Act = req.URL.Path
		dbmodels.AddLog(&userlog)
	}
}

func getUserTraceTab() {
	Try(func() {
		c := pool.Get()
		defer c.Close()
		UserTraceTab = make(map[string]string)
		vals, err := redis.Strings(c.Do("LRANGE", "UserTraceTab", 0, -1))
		if err != nil {
			// fmt.Println("getUserTraceTab error:", err)
		} else {
			// fmt.Println("getUserTraceTab vals:", vals)
			for i, v := range vals {
				UserTraceTab[v] = v
				fmt.Println("[get](%d): %s", i, v)
			}
			// fmt.Println("UserTraceTab -->", UserTraceTab)
		}
	}, func(e interface{}) {
		return
	})
}

func Trace() {
	c := pool.Get()
	defer c.Close()

	psc := redis.PubSubConn{Conn: c}
	psc.Subscribe("TraceNotice")
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			tab := string(v.Data)
			fmt.Printf("%s: message: %s\n", v.Channel, tab)
			if tab == "UserTraceTab" {
				getUserTraceTab()
			}
		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			return
		}
	}
}

func InitUser() {
	dbmodels.NewEngine(config.Config.DbDriver, config.Config.DbAddr, config.Config.DbPort, config.Config.DbUser, config.Config.DbPasswd, config.Config.Database, config.Config.ObjectsSchema)

	UserTraceTab = make(map[string]string)
	pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 600 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Config.SessionStoreIP+":"+config.Config.SessionStorePort)
			if err != nil {
				panic(err)
				return nil, err
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	getUserTraceTab()
	go Trace()
}

func UserHander() *martini.ClassicMartini {
	m := martini.Classic()
	m.Use(render.Renderer())

	store, _ := redistore.NewRediStore(10, "tcp", config.Config.SessionStoreIP+":"+config.Config.SessionStorePort, "", []byte(config.Config.SessionKey))
	m.Use(sessions.Sessions("user_session", store))

	m.Use(check_token)
	m.Use(userloger)

	m.Get("/session", func(r render.Render, req *http.Request, token Token) {
		r.JSON(200, map[string]interface{}{"Ret": 0, "AppId": token.Appid, "SiteId": token.SiteId, "UserId": token.UserId, "Token": token.Token})
	})

	m.Get("/get", func(r render.Render, req *http.Request, token Token) {
		userid, err := strconv.ParseInt(req.URL.Query().Get("userid"), 10, 64)
		sType := req.URL.Query().Get("type")
		sKeys := req.URL.Query().Get("keys")
		var keys []string
		if err := json.Unmarshal([]byte(sKeys), &keys); err != nil {
			r.JSON(200, map[string]interface{}{"Ret": 102010, "Msg": err.Error()})
			return
		}

		if token.Appid != 7 && token.Appid != 8 {
			keys = []string{"Icon"}
		}

		d, err := dbmodels.GetUserProfile(token.SiteId, userid, sType, keys)
		if err != nil {
			r.JSON(200, map[string]interface{}{"Ret": 102020, "Msg": err.Error()})
			return
		}

		var result UserProfileResp
		result.Ret = 0
		result.Data = d
		r.JSON(200, result)
	})

	m.Get("/captchaid", func(r render.Render, session sessions.Session) {
		var s = captcha.NewLen(4)
		session.Set("captchaId", s)

		var result UserCaptchaIdResp
		result.Ret = 0
		result.CaptchaId = s
		r.JSON(200, result)
	})

	m.Post("/regcode",
		binding.Bind(UserRegCodeForm{}),
		check_form,
		func(rf UserRegCodeForm, r render.Render, token Token) {
			errors := make([]binding.Error, 0)
			var v bool = true
			var err error
			switch rf.Type {
			case "phone":
				v, err = dbmodels.IsPhoneExist(token.SiteId, rf.Value)
			case "email":
				v, err = dbmodels.IsEmailExist(token.SiteId, rf.Value)
			default:
				errors = append(errors, binding.Error{
					FieldNames:     []string{"type"},
					Classification: "error",
					Message:        "Type not accept",
				})
			}

			if err != nil {
				errors = append(errors, binding.Error{
					FieldNames:     []string{"value"},
					Classification: "error",
					Message:        err.Error(),
				})
			}

			if v == true {
				errors = append(errors, binding.Error{
					FieldNames:     []string{"value"},
					Classification: "error",
					Message:        "value is exists",
				})
			}
			check_form(errors, r)
		},
		func(rf UserRegCodeForm, r render.Render, session sessions.Session, token Token) {
			sType := rf.Type
			sValue := rf.Value

			rand.Seed(time.Now().Unix())
			sCode := strconv.Itoa(rand.Intn(900000) + 100000)

			switch sType {
			case "phone":
				if token.Appid != 5 {
					//b, err := sender.SendSms(config.SmsUrl, token.Appid,token.SiteId,token.Token, 50, "registry", token.IP, sValue, "注册码", "您的注册确认码为:["+sCode+"]")
					b, err := sender.RpcxSendSms(log, config.Config.RpcxSmsBasepath, config.Config.ConsulAddress, int64(token.Appid), token.SiteId, token.Token, "registry", token.IP, sValue, "注册码", "您的注册确认码为:["+sCode+"]")
					if err != nil {
						r.JSON(200, map[string]interface{}{"Ret": 110100, "Msg": err.Error()})
						return
					}
					if b != true {
						r.JSON(200, map[string]interface{}{"Ret": 110110, "Msg": "短信发送失败,原因未知"})
						return
					}
				}
			case "email":
				if token.Appid != 5 {
					sender.SendMail(sValue, "注册码", "您的注册确认码为:["+sCode+"]")
				}
			default:
				r.JSON(200, map[string]interface{}{"Ret": 110120, "Msg": "不支持注册码发送"})
				return
			}

			session.Set("regvalue", sValue)
			session.Set("regcode", sCode)

			// fmt.Println(sType, sValue, sCode)

			if token.Appid == 5 {
				var result UserRegCodeResp
				result.Ret = 0
				result.Code = sCode
				r.JSON(200, result)
			} else {
				r.JSON(200, map[string]interface{}{"Ret": 0})
			}
		})

	m.Post("/authcode",
		binding.Bind(UserRegCodeForm{}),
		check_form,
		func(rf UserRegCodeForm, r render.Render, session sessions.Session, token Token) {
			sType := rf.Type
			sValue := rf.Value

			rand.Seed(time.Now().Unix())
			sCode := strconv.Itoa(rand.Intn(900000) + 100000)

			switch sType {
			case "phone":
				if token.Appid != 5 {
					//b, err := sender.SendSms(config.SmsUrl, token.Appid,token.SiteId, token.Token,50, "registry", token.IP, sValue, "注册码", "您的注册确认码为:["+sCode+"]")
					b, err := sender.RpcxSendSms(log, config.Config.RpcxSmsBasepath, config.Config.ConsulAddress, int64(token.Appid), token.SiteId, token.Token, "registry", token.IP, sValue, "注册码", "您的注册确认码为:["+sCode+"]")
					if err != nil {
						r.JSON(200, map[string]interface{}{"Ret": 110100, "Msg": err.Error()})
						return
					}
					if b != true {
						r.JSON(200, map[string]interface{}{"Ret": 110110, "Msg": "短信发送失败,原因未知"})
						return
					}
				}
			case "email":
				if token.Appid != 5 {
					sender.SendMail(sValue, "注册码", "您的注册确认码为:["+sCode+"]")
				}
			default:
				r.JSON(200, map[string]interface{}{"Ret": 110120, "Msg": "不支持注册码发送"})
				return
			}

			session.Set("regvalue", sValue)
			session.Set("regcode", sCode)

			// fmt.Println(sType, sValue, sCode)

			if token.Appid == 5 {
				var result UserRegCodeResp
				result.Ret = 0
				result.Code = sCode
				r.JSON(200, result)
			} else {
				r.JSON(200, map[string]interface{}{"Ret": 0})
			}
		})

	m.Post("/resetpasswdcode",
		binding.Bind(UserResetPasswdCodeForm{}),
		check_form,
		func(rf UserResetPasswdCodeForm, r render.Render, token Token) {
			errors := make([]binding.Error, 0)
			var v bool = false
			var err error
			switch rf.Type {
			case "phone":
				v, err = dbmodels.IsPhoneExist(token.SiteId, rf.Value)
			case "email":
				v, err = dbmodels.IsEmailExist(token.SiteId, rf.Value)
			default:
				errors = append(errors, binding.Error{
					FieldNames:     []string{"type"},
					Classification: "error",
					Message:        "Type not accept",
				})
			}

			if err != nil {
				errors = append(errors, binding.Error{
					FieldNames:     []string{"value"},
					Classification: "error",
					Message:        err.Error(),
				})
			}

			if v != true {
				errors = append(errors, binding.Error{
					FieldNames:     []string{"value"},
					Classification: "error",
					Message:        "value is not exists",
				})
			}
			check_form(errors, r)
		},
		func(rf UserResetPasswdCodeForm, r render.Render, session sessions.Session, token Token) {
			sValue := rf.Value

			rand.Seed(time.Now().Unix())
			sCode := strconv.Itoa(rand.Intn(900000) + 100000)

			switch rf.Type {
			case "phone":
				if token.Appid != 5 {
					//b, err := sender.SendSms(config.SmsUrl, token.Appid,token.SiteId,token.Token, 50, "registry", token.IP, sValue, "找回密码", "您的找回密码确认码为:["+sCode+"]")
					b, err := sender.RpcxSendSms(log, config.Config.RpcxSmsBasepath, config.Config.ConsulAddress, int64(token.Appid), token.SiteId, token.Token, "registry", token.IP, sValue, "找回密码", "您的找回密码确认码为:["+sCode+"]")
					if err != nil {
						r.JSON(200, map[string]interface{}{"Ret": 110210, "Msg": err.Error()})
						return
					}
					if b != true {
						r.JSON(200, map[string]interface{}{"Ret": 110220, "Msg": "短信发送失败,原因未知"})
						return
					}
				}
			case "email":
				if token.Appid != 5 {
					sender.SendMail(sValue, "找回密码", "您的找回密码确认码为:["+sCode+"]")
				}
			default:
			}

			session.Set("resetpasswdvalue", sValue)
			session.Set("resetpasswdcode", sCode)

			if token.Appid == 5 {
				var result UserRegCodeResp
				result.Ret = 0
				result.Code = sCode
				r.JSON(200, result)
			} else {
				r.JSON(200, map[string]interface{}{"Ret": 0})
			}
		})

	m.Post("/check",
		binding.Bind(UserCheckForm{}),
		check_form,
		func(rf UserCheckForm, r render.Render, token Token) {
			sType := rf.Type
			sValue := rf.Value
			sPhone := rf.Phone
			sId := rf.Idcard

			var v bool
			var err error
			var user *dbmodels.User

			switch sType {
			case "name":
				v, err = dbmodels.IsNameExist(token.SiteId, sValue)
				if v == true {
					user, err = dbmodels.GetUserByName(token.SiteId, sValue)
				}
			case "phone":
				v, err = dbmodels.IsPhoneExist(token.SiteId, sValue)
				if v == true {
					user, err = dbmodels.GetUserByPhone(token.SiteId, sValue)
				}
			case "idcard":
				v, err = dbmodels.IsIdcardExist(token.SiteId, sValue)
				if v == true {
					user, err = dbmodels.GetUserByIdcard(token.SiteId, sValue)
				}
			case "email":
				v, err = dbmodels.IsEmailExist(token.SiteId, sValue)
				if v == true {
					user, err = dbmodels.GetUserByEmail(token.SiteId, sValue)
				}
			case "rfid":
				v, err = dbmodels.IsRfidExist(token.SiteId, sValue)
				if v == true {
					user, err = dbmodels.GetUserByRfid(token.SiteId, sValue)
				}
			case "idphone":
				v, err = dbmodels.IsIdPhoneExist(token.SiteId, sId, sPhone)
				if v == true {
					user, err = dbmodels.GetUserByIdPhone(token.SiteId, sId, sPhone)
				}
			case "weixinid":
				v, err = dbmodels.IsWeixinidExist(token.SiteId, sValue)
				if v == true {
					user, err = dbmodels.GetUserByWeixinid(sValue)
				}

			default:
				r.JSON(200, map[string]interface{}{"Ret": 110310, "Msg": "类型:[" + sType + "] 不存在"})
				return
			}

			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 110320, "Msg": err.Error()})
				return
			}
			if v == true {
				r.JSON(200, map[string]interface{}{"Ret": 0, "UserId": user.Id, "SiteId": user.SiteId})
			} else {
				r.JSON(200, map[string]interface{}{"Ret": 110330, "Msg": "数据 [" + sValue + "]" + " 不存在"})
			}

		})
	m.Post("/checkunique",
		binding.Bind(UserCheckForm{}),
		check_form,
		func(rf UserCheckForm, r render.Render, token Token) {
			sType := rf.Type
			sValue := rf.Value

			var v bool
			var err error
			var user *dbmodels.User

			switch sType {
			case "phone":
				v, err = dbmodels.IsPhoneUniqueExist(token.SiteId, sValue)
				if v == true {
					user, err = dbmodels.GetUserByPhoneUnique(token.SiteId, sValue)
				}
			case "idcard":
				v, err = dbmodels.IsIdcardUniqueExist(token.SiteId, sValue)
				if v == true {
					user, err = dbmodels.GetUserByIdcardUnique(token.SiteId, sValue)
				}
			default:
				r.JSON(200, map[string]interface{}{"Ret": 110310, "Msg": "类型:[" + sType + "] 不存在"})
				return
			}

			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 110320, "Msg": err.Error()})
				return
			}
			if v == true {
				r.JSON(200, map[string]interface{}{"Ret": 0, "UserId": user.Id, "SiteId": user.SiteId, "weixinid": user.Weixinid})
			} else {
				r.JSON(200, map[string]interface{}{"Ret": 110330, "Msg": "数据 [" + sValue + "]" + " 不存在"})
			}

		})

	m.Post("/findbynickname", binding.Bind(UserFindByNickName{}),
		check_form,
		func(rf UserFindByNickName, r render.Render, session sessions.Session, token Token) {

			var keys []string
			if err := json.Unmarshal([]byte(rf.Keys), &keys); err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": err.Error()})
				return
			}

			v, err := dbmodels.SearchByNickName(token.SiteId, rf.NickName, rf.Type, keys)
			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": err.Error()})
				return
			} else {
				r.JSON(200, map[string]interface{}{"Ret": 0, "Datas": v})
			}
		})

	m.Post("/checkregcode",
		binding.Bind(UserCheckRegcodeForm{}),
		check_form,
		func(rf UserCheckRegcodeForm, r render.Render, session sessions.Session, token Token, userlog *dbmodels.UserLog) {
			if rf.Type != "imei" {
				_regvalue := session.Get("regvalue")
				if _regvalue == nil {
					r.JSON(200, map[string]interface{}{"Ret": 110410, "Msg": "没有相匹配的注册码,请先获取注册码"})
					return
				}

				_regcode := session.Get("regcode")
				if _regcode == nil {
					r.JSON(200, map[string]interface{}{"Ret": 110420, "Msg": "注册码不在session中,请先获取注册码"})
					return
				}

				if rf.Value != _regvalue.(string) {
					r.JSON(200, map[string]interface{}{"Ret": 110430, "Msg": "注册码不匹配"})
					return
				}

				if rf.Code != _regcode.(string) {
					r.JSON(200, map[string]interface{}{"Ret": 110440, "Msg": "注册码错误"})
					return
				}
			}

			r.JSON(200, map[string]interface{}{"Ret": 0, "Msg": "注册码正确"})
		})

	m.Post("/registry",
		binding.Bind(UserRegistryForm{}),
		check_form,
		func(rf UserRegistryForm, r render.Render, session sessions.Session, token Token, userlog *dbmodels.UserLog) {
			if rf.Type != "imei" {
				_regvalue := session.Get("regvalue")
				if _regvalue == nil {
					r.JSON(200, map[string]interface{}{"Ret": 110410, "Msg": "没有相匹配的注册码,请先获取注册码"})
					return
				}

				_regcode := session.Get("regcode")
				if _regcode == nil {
					r.JSON(200, map[string]interface{}{"Ret": 110420, "Msg": "注册码不在session中,请先获取注册码"})
					return
				}

				if rf.Value != _regvalue.(string) {
					r.JSON(200, map[string]interface{}{"Ret": 110430, "Msg": "注册码不匹配"})
					return
				}

				if rf.Code != _regcode.(string) {
					r.JSON(200, map[string]interface{}{"Ret": 110440, "Msg": "注册码错误"})
					return
				}
			}

			var user *dbmodels.User
			var err error
			switch rf.Type {
			case "name":
				user, err = dbmodels.RegisterUserByName(&dbmodels.User{SiteId: token.SiteId, Name: rf.Value, Passwd: rf.Passwd})
			case "phone":
				user, err = dbmodels.RegisterUserByPhone(&dbmodels.User{SiteId: token.SiteId, Phone: rf.Value, Passwd: rf.Passwd})
			case "idcard":
				user, err = dbmodels.RegisterUserByIdcard(&dbmodels.User{SiteId: token.SiteId, Idcard: rf.Value, Passwd: rf.Passwd})
			case "rfid":
				user, err = dbmodels.RegisterUserByRfid(&dbmodels.User{SiteId: token.SiteId, Rfid: rf.Value, Passwd: rf.Passwd})
			case "imei":
				user, err = dbmodels.RegisterUserByIMEI(&dbmodels.User{SiteId: token.SiteId, Imei: rf.Value, Passwd: "123456"})
			case "email":
				user, err = dbmodels.RegisterUserByEmail(&dbmodels.User{SiteId: token.SiteId, Email: rf.Value, Passwd: rf.Passwd})
			case "idphone":
				user, err = dbmodels.RegisterUserByIdPhone(&dbmodels.User{SiteId: token.SiteId, Idcard: rf.Idcard, Passwd: rf.Passwd, Phone: rf.Phone, Name: rf.Name})

			default:
			}

			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 110450, "Msg": err.Error()})
				return
			}

			session.Set("userid", user.Id)
			session.Set("siteid", token.SiteId)

			userlog.Msg = "注册账号"
			userlog.Json["type"] = rf.Type
			userlog.Json["value"] = rf.Value
			userlog.Json["phone"] = rf.Phone
			userlog.Json["idcard"] = rf.Idcard
			userlog.Json["name"] = rf.Name

			var result UserIdResp
			result.Ret = 0
			result.UserId = user.Id
			r.JSON(200, result)
		})
	m.Post("/bindsscard",
		binding.Bind(UserBindSscardForm{}),
		check_form,
		func(rf UserBindSscardForm, r render.Render, session sessions.Session, token Token, userlog *dbmodels.UserLog) {
			if rf.Type != "imei" {
				_regvalue := session.Get("regvalue")
				if _regvalue == nil {
					r.JSON(200, map[string]interface{}{"Ret": 110410, "Msg": "没有相匹配的注册码,请先获取注册码"})
					return
				}

				_regcode := session.Get("regcode")
				if _regcode == nil {
					r.JSON(200, map[string]interface{}{"Ret": 110420, "Msg": "注册码不在session中,请先获取注册码"})
					return
				}

				if rf.Value != _regvalue.(string) {
					r.JSON(200, map[string]interface{}{"Ret": 110430, "Msg": "注册码不匹配"})
					return
				}

				if rf.Code != _regcode.(string) {
					r.JSON(200, map[string]interface{}{"Ret": 110440, "Msg": "注册码错误"})
					return
				}
			}

			var user *dbmodels.User
			var err error

			switch rf.Type {
			case "phone":
				user, err = dbmodels.BindUserSScardByPhone(&dbmodels.User{SiteId: token.SiteId, Sscard: rf.Sscard, Phone: rf.Phone, Name: rf.Name, Idcard: rf.Idcard})
			case "idcard":
				user, err = dbmodels.BindUserSScardByIdcard(&dbmodels.User{SiteId: token.SiteId, Sscard: rf.Sscard, Phone: rf.Phone, Name: rf.Name, Idcard: rf.Idcard})

			default:
			}

			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 110450, "Msg": err.Error()})
				return
			}

			session.Set("userid", user.Id)
			session.Set("siteid", token.SiteId)

			userlog.Msg = "绑定社保账号"
			userlog.Json["type"] = rf.Type
			userlog.Json["value"] = rf.Value
			userlog.Json["phone"] = rf.Phone
			userlog.Json["idcard"] = rf.Idcard
			userlog.Json["sscard"] = rf.Sscard
			userlog.Json["name"] = rf.Name

			var result UserIdResp
			result.Ret = 0
			result.UserId = user.Id
			r.JSON(200, result)
		})

	m.Post("/resetpasswd",
		binding.Bind(UserResetPasswdForm{}),
		check_form,
		func(rf UserResetPasswdForm, r render.Render, token Token, session sessions.Session) {

			_regvalue := session.Get("resetpasswdvalue")
			if _regvalue == nil {
				r.JSON(200, map[string]interface{}{"Ret": 110510, "Msg": "没有相匹配的重置码,请先获取重置码"})
				return
			}

			_regcode := session.Get("resetpasswdcode")
			if _regcode == nil {
				r.JSON(200, map[string]interface{}{"Ret": 110520, "Msg": "重置码不在session中"})
				return
			}

			if rf.Value != _regvalue.(string) {
				r.JSON(200, map[string]interface{}{"Ret": 110530, "Msg": "提交的重置码不匹配"})
				return
			}

			if rf.Code != _regcode.(string) {
				r.JSON(200, map[string]interface{}{"Ret": 110540, "Msg": "重置码不匹配"})
				return
			}

			var user *dbmodels.User
			var err error

			switch rf.Type {
			case "phone":
				user, err = dbmodels.GetUserByPhone(token.SiteId, rf.Value)
			case "email":
				user, err = dbmodels.GetUserByEmail(token.SiteId, rf.Value)
			default:
			}

			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 110550, "Msg": err.Error()})
				return
			}

			b, err := dbmodels.SetPasswd(token.SiteId, user.Id, rf.Passwd)
			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 110560, "Msg": err.Error()})
				return
			}

			if b != true {
				r.JSON(200, map[string]interface{}{"Ret": 110570, "Msg": "重置密码错误"})
				return
			}

			var result Resp
			result.Ret = 0
			r.JSON(200, result)
		})

	m.Post("/loginbyid",
		binding.Bind(UserLoginByIdForm{}),
		check_form,
		func(rf UserLoginByIdForm, session sessions.Session, r render.Render, token Token, userlog *dbmodels.UserLog) {
			lUserid := rf.UserId
			sPasswd := rf.Passwd

			iscaptcha := AppIds[int64(token.Appid)]["captcha"]

			if token.Appid != 2 && token.Appid != 3 && token.Appid != 4 && iscaptcha != "false" {
				_captchaId := session.Get("captchaId")
				if _captchaId == nil {
					r.JSON(200, map[string]interface{}{"Ret": 110610, "Msg": "验证码不存在,请先获取验证码"})
					return
				}
				captchaId := _captchaId.(string)
				if rf.Captcha == "" || captchaId == "" {
					r.JSON(200, map[string]interface{}{"Ret": 110620, "Msg": "没有提交验证码"})
					return
				}

				// fmt.Println(captchaId, captchaValue)
				if !captcha.VerifyString(captchaId, rf.Captcha) {
					if token.Appid == 5 {
						if captchaId != rf.Captcha {
							r.JSON(200, map[string]interface{}{"Ret": 110630, "Msg": "验证码匹配错误"})
							return
						}
					} else {
						r.JSON(200, map[string]interface{}{"Ret": 110640, "Msg": "验证码匹配错误"})
						return
					}
				}
			}

			if lUserid == 0 || sPasswd == "" {
				r.JSON(200, map[string]interface{}{"Ret": 110650, "Msg": "用户Id或者密码为空"})
				return
			}

			user, err := dbmodels.LoginUserByid(token.SiteId, lUserid, sPasswd, token.Token)
			if err != nil {
				if err == dbmodels.ErrUserNotExist {
					r.JSON(200, map[string]interface{}{"Ret": 110660, "Msg": "用户名和密码匹配失败"})
					return
				}

				r.JSON(200, map[string]interface{}{"Ret": 110670, "Msg": err.Error()})
				return
			}

			session.Set("userid", user.Id)
			session.Set("siteid", user.SiteId)

			userlog.Msg = "登录"
			// userlog.Json["value"] = sValue

			var result UserIdSiteIdResp
			result.Ret = 0
			result.UserId = user.Id
			result.SiteId = user.SiteId
			r.JSON(200, result)
		})

	m.Post("/loginbyimei",
		binding.Bind(UserLoginByIMEIForm{}),
		check_form,
		func(rf UserLoginByIMEIForm, session sessions.Session, r render.Render, token Token, userlog *dbmodels.UserLog) {
			user, err := dbmodels.GetUserByIMEI(token.SiteId, rf.IMEI)
			if err != nil {
				if err == dbmodels.ErrUserNotExist {
					r.JSON(200, map[string]interface{}{"Ret": 110660, "Msg": "没有该imei"})
					return
				}

				r.JSON(200, map[string]interface{}{"Ret": 110670, "Msg": err.Error()})
				return
			}

			session.Set("userid", user.Id)
			session.Set("siteid", user.SiteId)

			userlog.Msg = "IMEI登录"
			userlog.Json["value"] = rf.IMEI

			var result UserNickNameResp
			result.Ret = 0
			result.UserId = user.Id
			result.SiteId = user.SiteId
			result.NickName = user.Nickname
			r.JSON(200, result)
		})

	m.Post("/loginbyweixin",
		binding.Bind(UserLoginByWeixinForm{}),
		check_form,
		func(rf UserLoginByWeixinForm, w http.ResponseWriter, req *http.Request, session sessions.Session, r render.Render, token Token, userlog *dbmodels.UserLog) {
			user, err := dbmodels.GetUserByWeixinid(rf.Weixinid)
			if err != nil {
				if err == dbmodels.ErrUserNotExist {
					r.JSON(200, map[string]interface{}{"Ret": 110660, "Msg": "没有该微信"})
					return
				}

				r.JSON(200, map[string]interface{}{"Ret": 110670, "Msg": err.Error()})
				return
			}

			session.Set("userid", user.Id)
			session.Set("siteid", user.SiteId)
			session.Set("appid", token.Appid)

			userlog.Msg = "微信登录"
			userlog.Json["value"] = rf.Weixinid

			w.Header().Add("userid", "12345")
			var result UserIdSiteIdResp
			result.Ret = 0
			result.UserId = user.Id
			result.SiteId = user.SiteId
			result.NickName = user.Nickname
			r.JSON(200, result)
		})

	m.Post("/login",
		binding.Bind(UserLoginForm{}),
		check_form,
		func(rf UserLoginForm, w http.ResponseWriter, req *http.Request, session sessions.Session, r render.Render, token Token, userlog *dbmodels.UserLog) {
			sValue := rf.Name
			sPasswd := rf.Passwd

			salt := ""
			var err error
			token.Appid, err = strconv.Atoi(req.URL.Query().Get("appid"))

			log.Finest("AppIds = %+v", AppIds)

			salt = token.Token

			iscaptcha := AppIds[int64(token.Appid)]["captcha"]

			if iscaptcha != "false" {
				_captchaId := session.Get("captchaId")
				if _captchaId == nil {
					r.JSON(200, map[string]interface{}{"Ret": 110610, "Msg": "验证码不存在,请先获取验证码"})
					return
				}
				captchaId := _captchaId.(string)
				if rf.Captcha == "" || captchaId == "" {
					r.JSON(200, map[string]interface{}{"Ret": 110620, "Msg": "没有提交验证码"})
					return
				}

				// fmt.Println(captchaId, captchaValue)
				if !captcha.VerifyString(captchaId, rf.Captcha) {
					if token.Appid == 5 {
						if captchaId != rf.Captcha {
							r.JSON(200, map[string]interface{}{"Ret": 110630, "Msg": "验证码匹配错误"})
							return
						}
					} else {
						r.JSON(200, map[string]interface{}{"Ret": 110640, "Msg": "验证码匹配错误"})
						return
					}
				}
			}

			if sValue == "" || sPasswd == "" {
				r.JSON(200, map[string]interface{}{"Ret": 110650, "Msg": "用户名或者密码为空"})
				return
			}

			user, err := dbmodels.LoginUserPlain(token.SiteId, sValue, sPasswd, salt)

			log.Finest("site=%d,name=%s,salt=%s,user=%+v", token.SiteId, sValue, salt, user)

			if err != nil {
				if err == dbmodels.ErrUserNotExist {
					r.JSON(200, map[string]interface{}{"Ret": 110660, "Msg": "用户名和密码匹配失败"})
					return
				}

				r.JSON(200, map[string]interface{}{"Ret": 110670, "Msg": err.Error()})
				return
			}

			session.Set("userid", user.Id)
			session.Set("siteid", user.SiteId)
			session.Set("appid", token.Appid)

			userlog.Msg = "登录"
			userlog.Json["value"] = sValue

			w.Header().Add("userid", "12345")
			var result UserIdSiteIdResp
			result.Ret = 0
			result.UserId = user.Id
			result.SiteId = user.SiteId
			r.JSON(200, result)
		})

	m.Post("/loginbycode",
		binding.Bind(UserLoginForm{}),
		check_form,
		func(rf UserLoginForm, w http.ResponseWriter, req *http.Request, session sessions.Session, r render.Render, token Token, userlog *dbmodels.UserLog) {
			sValue := rf.Name
			sPasswd := rf.Passwd

			_regphone := session.Get("regvalue")
			if _regphone == nil {
				r.JSON(200, map[string]interface{}{"Ret": 110410, "Msg": "没有相匹配的注册码,请先获取注册码"})
				return
			}

			_regcode := session.Get("regcode")
			if _regcode == nil {
				r.JSON(200, map[string]interface{}{"Ret": 110420, "Msg": "注册码不在session中,请先获取注册码"})
				return
			}

			log.Finest("retphone=%v,regcode=%v,svalue=%s,spasswd=%v", _regphone, _regcode, sValue, sPasswd)

			if sValue != _regphone.(string) {
				r.JSON(200, map[string]interface{}{"Ret": 110430, "Msg": "注册码不匹配"})
				return
			}

			if sPasswd != _regcode.(string) {
				r.JSON(200, map[string]interface{}{"Ret": 110440, "Msg": "注册码错误"})
				return
			}

			salt := ""
			var err error
			token.Appid, err = strconv.Atoi(req.URL.Query().Get("appid"))

			log.Finest("AppIds = %+v", AppIds)

			salt = token.Token

			iscaptcha := AppIds[int64(token.Appid)]["captcha"]

			if iscaptcha != "false" {
				_captchaId := session.Get("captchaId")
				if _captchaId == nil {
					r.JSON(200, map[string]interface{}{"Ret": 110610, "Msg": "验证码不存在,请先获取验证码"})
					return
				}
				captchaId := _captchaId.(string)
				if rf.Captcha == "" || captchaId == "" {
					r.JSON(200, map[string]interface{}{"Ret": 110620, "Msg": "没有提交验证码"})
					return
				}

				// fmt.Println(captchaId, captchaValue)
				if !captcha.VerifyString(captchaId, rf.Captcha) {
					if token.Appid == 5 {
						if captchaId != rf.Captcha {
							r.JSON(200, map[string]interface{}{"Ret": 110630, "Msg": "验证码匹配错误"})
							return
						}
					} else {
						r.JSON(200, map[string]interface{}{"Ret": 110640, "Msg": "验证码匹配错误"})
						return
					}
				}
			}

			if sValue == "" || sPasswd == "" {
				r.JSON(200, map[string]interface{}{"Ret": 110650, "Msg": "手机号或者短信验证码为空"})
				return
			}

			user, err := dbmodels.LoginUserByPhoneCode(sValue, sPasswd, salt, 1)

			log.Finest("site=%d,name=%s,salt=%s,user=%+v", token.SiteId, sValue, salt, user)

			if err != nil {
				if err == dbmodels.ErrUserNotExist {
					r.JSON(200, map[string]interface{}{"Ret": 110660, "Msg": "用户名和密码匹配失败"})
					return
				}

				r.JSON(200, map[string]interface{}{"Ret": 110670, "Msg": err.Error()})
				return
			}

			session.Set("userid", user.Id)
			session.Set("siteid", user.SiteId)
			session.Set("appid", token.Appid)

			userlog.Msg = "登录"
			userlog.Json["value"] = sValue

			w.Header().Add("userid", "12345")
			var result UserIdSiteIdResp
			result.Ret = 0
			result.UserId = user.Id
			result.SiteId = user.SiteId
			r.JSON(200, result)
		})

	m.Post("/authlongin",
		binding.Bind(UserCloneTokenForm{}),
		check_form,
		func(rf UserCloneTokenForm, r render.Render, session sessions.Session, token Token) {
			if rf.UserId == 0 {
				r.JSON(200, map[string]interface{}{"Ret": 110100, "Msg": "用户id错误"})
				return
			}
			if rf.Auth == "" {
				r.JSON(200, map[string]interface{}{"Ret": 110110, "Msg": "验证错误"})
				return
			}

			session.Set("userid", rf.UserId)
			session.Set("siteid", token.SiteId)

			var result UserIdResp
			result.Ret = 0
			result.UserId = rf.UserId
			r.JSON(200, result)
		})

	m.Post("/bind_by_weixin",
		binding.Bind(UserLoginForm{}),
		check_form,
		func(rf UserLoginForm, session sessions.Session, r render.Render, token Token, userlog *dbmodels.UserLog) {
			if token.UserId > 0 {
				r.JSON(200, map[string]interface{}{"Ret": 110650, "Msg": "用户已经登录"})
				return
			}

			weixinid := session.Get("weixinid")
			if weixinid == nil {
				r.JSON(200, map[string]interface{}{"Ret": 110650, "Msg": "不存在绑定的微信账号"})
				return
			}
			//

			sValue := rf.Name
			sPasswd := rf.Passwd
			if sValue == "" || sPasswd == "" {
				r.JSON(200, map[string]interface{}{"Ret": 110650, "Msg": "用户名或者密码为空"})
				return
			}

			user, err := dbmodels.LoginUserAll(sValue, sPasswd, token.Token)
			if err != nil {
				if err == dbmodels.ErrUserNotExist {
					r.JSON(200, map[string]interface{}{"Ret": 110660, "Msg": "用户名和密码匹配失败"})
					return
				}

				r.JSON(200, map[string]interface{}{"Ret": 110670, "Msg": err.Error()})
				return
			}

			b, err := dbmodels.SetUserWeixinid(user.SiteId, user.Id, weixinid.(string))
			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 110671, "Msg": err.Error()})
				return
			}
			if b != true {
				r.JSON(200, map[string]interface{}{"Ret": 110672, "Msg": "绑定微信账号失败"})
				return
			}

			session.Set("userid", user.Id)
			session.Set("siteid", user.SiteId)
			session.Delete("weixinid")

			userlog.Msg = "绑定微信账号"
			userlog.Json["value"] = sValue

			var result UserIdResp
			result.Ret = 0
			result.UserId = user.Id
			r.JSON(200, result)
		})

	m.Post("/bindwxbycode",
		binding.Bind(UserBindWxForm{}),
		check_form,
		func(rf UserBindWxForm, w http.ResponseWriter, req *http.Request, session sessions.Session, r render.Render, token Token, userlog *dbmodels.UserLog) {
			sValue := rf.Phone
			sPasswd := rf.Code
			weixinid := rf.Weixinid

			log.Finest("phone=%s,code=%s,weixinid = %s", sValue, sPasswd, weixinid)

			_regphone := session.Get("regvalue")
			if _regphone == nil {
				r.JSON(200, map[string]interface{}{"Ret": 110410, "Msg": "没有相匹配的注册码,请先获取注册码"})
				return
			}

			_regcode := session.Get("regcode")
			if _regcode == nil {
				r.JSON(200, map[string]interface{}{"Ret": 110420, "Msg": "注册码不在session中,请先获取注册码"})
				return
			}

			log.Finest("retphone=%v,regcode=%v", _regphone, _regcode)

			if sValue != _regphone.(string) {
				r.JSON(200, map[string]interface{}{"Ret": 110430, "Msg": "注册码不匹配"})
				return
			}

			if sPasswd != _regcode.(string) {
				r.JSON(200, map[string]interface{}{"Ret": 110440, "Msg": "注册码错误"})
				return
			}

			salt := ""
			var err error
			token.Appid, err = strconv.Atoi(req.URL.Query().Get("appid"))

			log.Finest("AppIds = %+v", AppIds)

			salt = token.Token

			iscaptcha := AppIds[int64(token.Appid)]["captcha"]

			if iscaptcha != "false" {
				_captchaId := session.Get("captchaId")
				if _captchaId == nil {
					r.JSON(200, map[string]interface{}{"Ret": 110610, "Msg": "验证码不存在,请先获取验证码"})
					return
				}
				captchaId := _captchaId.(string)
				if rf.Captcha == "" || captchaId == "" {
					r.JSON(200, map[string]interface{}{"Ret": 110620, "Msg": "没有提交验证码"})
					return
				}

				if !captcha.VerifyString(captchaId, rf.Captcha) {
					if token.Appid == 5 {
						if captchaId != rf.Captcha {
							r.JSON(200, map[string]interface{}{"Ret": 110630, "Msg": "验证码匹配错误"})
							return
						}
					} else {
						r.JSON(200, map[string]interface{}{"Ret": 110640, "Msg": "验证码匹配错误"})
						return
					}
				}
			}

			if sValue == "" || sPasswd == "" {
				r.JSON(200, map[string]interface{}{"Ret": 110650, "Msg": "手机号或者短信验证码为空"})
				return
			}

			user, err := dbmodels.LoginUserByPhoneCode(sValue, sPasswd, salt, 1)

			log.Finest("site=%d,name=%s,salt=%s,user=%+v", token.SiteId, sValue, salt, user)

			if err != nil {
				if err == dbmodels.ErrUserNotExist {
					r.JSON(200, map[string]interface{}{"Ret": 110660, "Msg": "用户名和密码匹配失败"})
					return
				}

				r.JSON(200, map[string]interface{}{"Ret": 110670, "Msg": err.Error()})
				return
			}

			//weixinid := session.Get("weixinid")
			if weixinid == "" {
				r.JSON(200, map[string]interface{}{"Ret": 110650, "Msg": "不存在绑定的微信账号"})
				return
			}
			//

			b, err := dbmodels.SetUserWeixinid(user.SiteId, user.Id, weixinid)
			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 110671, "Msg": err.Error()})
				return
			}
			if b != true {
				r.JSON(200, map[string]interface{}{"Ret": 110672, "Msg": "绑定微信账号失败"})
				return
			}

			session.Set("userid", user.Id)
			session.Set("siteid", user.SiteId)
			session.Set("appid", token.Appid)

			//session.Delete("weixinid")

			userlog.Msg = "绑定微信账号"
			userlog.Json["value"] = sValue

			var result UserIdSiteIdResp
			result.Ret = 0
			result.UserId = user.Id
			result.SiteId = user.SiteId
			r.JSON(200, result)
		})

	m.Post("/unbindweixin",
		binding.Bind(UserUnBindWxForm{}),
		check_form,
		func(rf UserUnBindWxForm, w http.ResponseWriter, req *http.Request, session sessions.Session, r render.Render, token Token, userlog *dbmodels.UserLog) {
			uid := rf.UserId
			weixinid := rf.Weixinid

			log.Finest("userid=%v,weixinid = %s", uid, weixinid)

			if uid <= 0 {
				r.JSON(200, map[string]interface{}{"Ret": 120100, "Msg": "用户id未指定"})
				return
			}

			if ok, _ := dbmodels.IsWeixinidExist(1, weixinid); !ok {
				r.JSON(200, map[string]interface{}{"Ret": 120200, "Msg": "微信id不存在。"})
				return

			}

			b, err := dbmodels.UnSetUserWeixinid(token.SiteId, uid, "")
			if err != nil && err != dbmodels.ErrWeixinidAlreadyUsed {
				r.JSON(200, map[string]interface{}{"Ret": 120010, "Msg": err.Error()})
				return
			}
			if b != true {
				r.JSON(200, map[string]interface{}{"Ret": 120020, "Msg": "解绑微信账号失败"})
				return
			}

			userlog.Msg = "解绑微信账号"
			userlog.Json["value"] = weixinid

			r.JSON(200, map[string]interface{}{"Ret": 0, "Msg": "解绑微信账号成功。"})
		})

	// other id
	m.Post("/otherid_registry_old",
		func(req *http.Request, r render.Render, session sessions.Session, token Token, userlog *dbmodels.UserLog) {
			var err error

			var otherType, otherId, sValue, sPasswd string

			contentType := req.Header.Get("Content-Type")
			if strings.Contains(contentType, "form-urlencoded") {
				otherType = req.FormValue("othertype")
				otherId = req.FormValue("otherid")

				sValue = req.FormValue("name")
				sPasswd = req.FormValue("passwd")

			} else if strings.Contains(contentType, "json") {
				body, err := ioutil.ReadAll(req.Body)
				if err != nil {
					r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": err.Error()})
					return
				}
				param := make(map[string]interface{})
				err = json.Unmarshal(body, &param)
				if err != nil {
					r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": err.Error()})
					return
				}
				sPasswd = utils.ConvertToString(param["passwd"])
				sValue = utils.ConvertToString(param["name"])
				otherType = utils.ConvertToString(param["othertype"])
				otherId = utils.ConvertToString(param["otherid"])

			} else {
				r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": "不能识别的Content-Type类别。"})
				return
			}

			//other
			if otherType == "" || otherId == "" {
				r.JSON(200, map[string]interface{}{"Ret": 110710, "Msg": "otherType or otherId is empty"})
				return
			}

			if sValue == "" || sPasswd == "" {
				r.JSON(200, map[string]interface{}{"Ret": 110720, "Msg": "name or passwd is empty"})
				return
			}

			id, err := dbmodels.GetBind(token.SiteId, otherType, otherId)
			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 110730, "Msg": err.Error()})
				return
			} else if id > 0 {
				r.JSON(200, map[string]interface{}{"Ret": 110740, "Msg": dbmodels.ErrNameAlreadyUsed.Error()})
				return
			}

			//user
			user, err := dbmodels.LoginUserPlain(token.SiteId, sValue, sPasswd, token.Token)
			if err != nil {
				if err == dbmodels.ErrUserNotExist {
					r.JSON(200, map[string]interface{}{"Ret": 110750, "Msg": "Username or password is not correct"})
					return
				}

				r.JSON(200, map[string]interface{}{"Ret": 110760, "Msg": err.Error()})
				return
			}

			name2id, err := dbmodels.SetBind(&dbmodels.NameLog{SiteId: token.SiteId, Type: otherType, Name: otherId, Id: user.Id})
			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 110770, "Msg": err.Error()})
				return
			}

			if name2id.Id != user.Id {
				r.JSON(200, map[string]interface{}{"Ret": 110780, "Msg": "bind id error"})
				return
			}

			session.Set("userid", user.Id)
			session.Set("siteid", token.SiteId)

			userlog.Msg = "第三方账号注册"
			userlog.Json["othertype"] = otherType
			userlog.Json["otherid"] = otherId
			userlog.Json["name"] = sValue

			var result UserIdResp
			result.Ret = 0
			result.UserId = user.Id
			r.JSON(200, result)
		})

	m.Post("/otherid_login",
		func(req *http.Request, r render.Render, session sessions.Session, token Token, userlog *dbmodels.UserLog) {
			var err error

			var otherType, otherId string

			contentType := req.Header.Get("Content-Type")
			if strings.Contains(contentType, "form-urlencoded") {
				otherType = req.FormValue("othertype")
				otherId = req.FormValue("otherid")

			} else if strings.Contains(contentType, "json") {
				body, err := ioutil.ReadAll(req.Body)
				if err != nil {
					r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": err.Error()})
					return
				}
				param := make(map[string]interface{})
				err = json.Unmarshal(body, &param)
				if err != nil {
					r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": err.Error()})
					return
				}
				otherType = utils.ConvertToString(param["othertype"])
				otherId = utils.ConvertToString(param["otherid"])

			} else {
				r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": "不能识别的Content-Type类别。"})
				return
			}

			//other
			if otherType == "" || otherId == "" {
				r.JSON(200, map[string]interface{}{"Ret": 110810, "Msg": "otherType or otherId is empty"})
				return
			}

			id, err := dbmodels.GetBind(token.SiteId, otherType, otherId)
			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 110820, "Msg": err.Error()})
				return
			} else if id == 0 {
				r.JSON(200, map[string]interface{}{"Ret": 110830, "Msg": dbmodels.ErrUserNotExist.Error()})
				return
			}

			session.Set("userid", id)
			session.Set("siteid", token.SiteId)

			userlog.Msg = "第三方账号登录"
			userlog.Json["othertype"] = otherType
			userlog.Json["otherid"] = otherId

			var result UserIdResp
			result.Ret = 0
			result.UserId = id
			r.JSON(200, result)
		})

	//mu
	mu := martini.Classic()
	mu.Use(render.Renderer())
	mu.Use(sessions.Sessions("user_session", store))

	mu.Use(func(c martini.Context, req *http.Request, res http.ResponseWriter, session sessions.Session, r render.Render) {
		userid := session.Get("userid").(int64)
		c.Map(userid)
		c.Next()

		suserid := strconv.FormatInt(userid, 10)
		rw := res.(martini.ResponseWriter)

		// fmt.Println("-->>>>>", suserid, UserTraceTab[suserid], UserTraceTab[suserid] != "")
		if UserTraceTab[suserid] != "" {
			c := pool.Get()
			defer c.Close()
			var traceMsg TraceMsg
			traceMsg.Mtype = "3"
			traceMsg.ID = suserid
			traceMsg.AppID = "1"
			traceMsg.Module = "50"
			traceMsg.Event = req.URL.Path
			traceMsg.Time = time.Now().Format("2006-01-02 15:04:05")
			if rw.Status() == 200 {
				traceMsg.Success = "1"
			} else {
				traceMsg.Success = "0"
			}
			ss, _ := json.Marshal(req)
			traceMsg.Input = string(ss)

			traceMsg.Output = ""
			traceMsg.Notice = ""

			straceMsg, _ := json.Marshal(traceMsg)
			c.Send("LPUSH", "TraceQue", straceMsg)
			// fmt.Println("-->>>>>", suserid)
		}
	})
	mu.Use(check_token)
	mu.Use(userloger)

	mu.Post("/logout", func(r render.Render, session sessions.Session, userlog *dbmodels.UserLog) {
		userlog.Msg = "退出"

		session.Clear()
		r.JSON(200, map[string]interface{}{"Ret": 0})
	})

	mu.Get("/getauthlogin", func(r render.Render, token Token, userlog *dbmodels.UserLog) {
		userlog.Msg = "生成登录二维码"
		var result UserCloneToken
		result.Ret = 0
		result.UserId = token.UserId
		result.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
		result.Auth = "auth_" + strconv.FormatInt(time.Now().Unix(), 10)
		r.JSON(200, result)
	})

	mu.Post("/passwd",
		binding.Bind(UserPasswdForm{}),
		check_form,
		func(rf UserPasswdForm, r render.Render, token Token, userlog *dbmodels.UserLog) {
			v, err := dbmodels.ResetPasswd(token.SiteId, token.UserId, rf.Oldpasswd, rf.Newpasswd, token.Token)
			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 101010, "Msg": err.Error()})
				return
			}
			if v == true {
				r.JSON(200, map[string]interface{}{"Ret": 0})
			} else {
				r.JSON(200, map[string]interface{}{"Ret": 101020, "Msg": "密码重置失败"})
			}

			userlog.Msg = "修改密码"
		})

	mu.Get("/profile", func(req *http.Request, r render.Render, token Token) {
		sType := req.URL.Query().Get("type")
		sKeys := req.URL.Query().Get("keys")
		var keys []string
		if err := json.Unmarshal([]byte(sKeys), &keys); err != nil {
			r.JSON(200, map[string]interface{}{"Ret": 102010, "Msg": err.Error()})
			return
		}
		d, err := dbmodels.GetUserProfile(token.SiteId, token.UserId, sType, keys)
		if err != nil {
			r.JSON(200, map[string]interface{}{"Ret": 102020, "Msg": err.Error()})
			return
		}

		var result UserProfileResp
		result.Ret = 0
		result.Data = d
		r.JSON(200, result)
	})

	mu.Post("/profile", func(req *http.Request, r render.Render, token Token, userlog *dbmodels.UserLog) {

		var sType, sDatas string

		contentType := req.Header.Get("Content-Type")
		if strings.Contains(contentType, "form-urlencoded") {
			sType = req.FormValue("type")
			sDatas = req.FormValue("datas")

		} else if strings.Contains(contentType, "json") {
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": err.Error()})
				return
			}
			param := make(map[string]interface{})
			err = json.Unmarshal(body, &param)
			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": err.Error()})
				return
			}
			sType = utils.ConvertToString(param["type"])
			sDatas = utils.ConvertToString(param["datas"])

		} else {
			r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": "不能识别的Content-Type类别。"})
			return
		}

		datas := make(map[string]interface{})
		if err := json.Unmarshal([]byte(sDatas), &datas); err != nil {
			r.JSON(200, map[string]interface{}{"Ret": 103010, "Msg": "data数据格式错误"})
			return
		}
		v, err := dbmodels.SetUserProfile(token.SiteId, token.UserId, sType, datas)
		if err != nil {
			r.JSON(200, map[string]interface{}{"Ret": 103020, "Msg": err.Error()})
			return
		}
		if v == true {
			r.JSON(200, map[string]interface{}{"Ret": 0})
		} else {
			r.JSON(200, map[string]interface{}{"Ret": 103030, "Msg": "属性设置失败"})
		}

		userlog.Msg = "设置属性"
		userlog.Json["type"] = sType
		userlog.Json["datas"] = sDatas
	})

	mu.Post("/setunique", func(req *http.Request, r render.Render, token Token, userlog *dbmodels.UserLog) {

		var sType, sValue string

		contentType := req.Header.Get("Content-Type")
		if strings.Contains(contentType, "form-urlencoded") {
			sType = req.FormValue("type")
			sValue = req.FormValue("value")

		} else if strings.Contains(contentType, "json") {
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": err.Error()})
				return
			}
			param := make(map[string]interface{})
			err = json.Unmarshal(body, &param)
			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": err.Error()})
				return
			}
			sType = utils.ConvertToString(param["type"])
			sValue = utils.ConvertToString(param["value"])

		} else {
			r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": "不能识别的Content-Type类别。"})
			return
		}

		var v bool
		var err error

		switch sType {
		case "name":
			v, err = dbmodels.SetUserName(token.SiteId, token.UserId, sValue)
		case "phone":
			v, err = dbmodels.SetUserPhone(token.SiteId, token.UserId, sValue)
		case "idcard":
			v, err = dbmodels.SetUserIdcard(token.SiteId, token.UserId, sValue)
		case "rfid":
			v, err = dbmodels.SetUserRfid(token.SiteId, token.UserId, sValue)
		case "email":
			v, err = dbmodels.SetUserEmail(token.SiteId, token.UserId, sValue)
		default:
		}

		if err != nil {
			r.JSON(200, map[string]interface{}{"Ret": 104010, "Msg": err.Error()})
			return
		}
		if v == true {
			r.JSON(200, map[string]interface{}{"Ret": 0})
		} else {
			r.JSON(200, map[string]interface{}{"Ret": 104020, "Msg": "登录名设置失败"})
		}

		userlog.Msg = "设置登陆名"
		userlog.Json["type"] = sType
		userlog.Json["value"] = sValue
	})

	mu.Get("/logs", func(req *http.Request, r render.Render, token Token) {
		rows, _ := strconv.Atoi(req.URL.Query().Get("rows"))
		page, _ := strconv.Atoi(req.URL.Query().Get("page"))
		level, _ := strconv.Atoi(req.URL.Query().Get("level"))

		v, err := dbmodels.LogList(token.SiteId, token.UserId, rows, page, level)
		if err != nil {
			r.JSON(200, map[string]interface{}{"Ret": 105010, "Msg": err.Error()})
			return
		}
		var result UserLogListResp
		result.Ret = 0
		result.Datas = v
		r.JSON(200, result)
	})

	mu.Post("/feedback", func(req *http.Request, r render.Render, token Token) {
		var sContent string

		contentType := req.Header.Get("Content-Type")
		if strings.Contains(contentType, "form-urlencoded") {
			sContent = req.FormValue("content")

		} else if strings.Contains(contentType, "json") {
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": err.Error()})
				return
			}
			param := make(map[string]interface{})
			err = json.Unmarshal(body, &param)
			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": err.Error()})
				return
			}
			sContent = utils.ConvertToString(param["content"])

		} else {
			r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": "不能识别的Content-Type类别。"})
			return
		}

		_, err := dbmodels.SubmitFeedback(&dbmodels.Feedback{Appid: token.Appid, SiteId: token.SiteId, Userid: token.UserId, Content: sContent})
		if err != nil {
			r.JSON(200, map[string]interface{}{"Ret": 108010, "Msg": err.Error()})
			return
		}
		r.JSON(200, map[string]interface{}{"Ret": 0})
	})

	m.Any("/:userid/.*",
		check_userid,
		func(w http.ResponseWriter, r *http.Request, token Token) {
			if p := strings.TrimPrefix(r.URL.Path, "/"+strconv.FormatInt(token.UserId, 10)); len(p) < len(r.URL.Path) {
				r.URL.Path = p
			} else {
				http.NotFound(w, r)
			}
		}, mu.ServeHTTP)
	return m
}
