package admins

import (
	"errors"
	"math/rand"
	"net/http"
	"platform/common/utils"
	"strconv"
	"strings"
	"time"

	"github.com/boj/redistore"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"platform/lib/sender"
	"platform/oasvc/config"
	"platform/oasvc/dbmodels"
	. "platform/oasvc/models"
	oacerr "platform/pfcomm/errors"
)

func InitUser() {
	dbmodels.NewEngine(config.Config.DbDriver, config.Config.DbAddr, config.Config.DbPort, config.Config.DbUser, config.Config.DbPasswd, config.Config.Database, config.Config.ObjectsSchema)
}

func UserHander(debug bool) *martini.ClassicMartini {
	m := martini.Classic()
	m.Use(render.Renderer())
	store, _ := redistore.NewRediStore(10, "tcp", config.Config.SessionStoreIP+":"+config.Config.SessionStorePort, "", []byte(config.Config.SessionKey))
	m.Use(sessions.Sessions("admin_session", store))

	m.Use(check_token)

	if debug == false {
		m.Use(check_power(config.Config.Prefix + "/admin/user"))
	}
	m.Get("/get",
		func(res http.ResponseWriter, req *http.Request, r render.Render) {
			id, err := strconv.ParseInt(req.URL.Query().Get("id"), 10, 64)
			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": err.Error()})
				return
			}

			user, err := dbmodels.GetAdminById(id)
			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": err.Error()})
				return
			}

			if user == nil {
				r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": errors.New("user is not have").Error()})
				return
			}

			// var result ResourceResp
			r.JSON(200, map[string]interface{}{"Ret": 0, "Data": user})
		})

	m.Post("/login",
		binding.Bind(AdminLoginForm{}),
		check_form,
		func(rf AdminLoginForm, req *http.Request, session sessions.Session, r render.Render, token Token) {
			salt := ""
			if token.Token != "" {
				salt = token.Token
			}

			user, err := dbmodels.LoginAdmin(token.SiteId, rf.Name, rf.Passwd, salt)

			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": err.Error()})
				return
			}

			session.Set("userid", user.Id)
			session.Set("siteid", token.SiteId)
			var result IdSiteIdResp
			result.Ret = 0
			result.Id = user.Id
			result.SiteId = user.SiteId
			result.OrganizationId = user.OrganizationId
			//result.Token = token.Token

			r.JSON(200, result)
		})

	m.Post("/resetpasswd",
		binding.Bind(AdminResetPasswdForm{}),
		check_form,
		func(rf AdminResetPasswdForm, req *http.Request, r render.Render) {
			b, err := dbmodels.ResetPasswd(rf.Id, rf.OldPasswd, rf.NewPasswd)

			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": err.Error()})
				return
			}

			if b != true {
				r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": errors.New("passwd set false").Error()})
				return
			}

			var result Resp
			result.Ret = 0
			r.JSON(200, result)
		})

	m.Post("/resetpasswdcode",
		binding.Bind(AdminResetPasswdCodeForm{}),
		check_form,
		func(rf AdminResetPasswdCodeForm, r render.Render, token Token) {
			errors := make([]binding.Error, 0)
			var v bool = false
			var err error
			switch rf.Type {
			case "phone":
				v, err = dbmodels.IsPhoneExist(token.SiteId, rf.Value)
			// case "email":
			// 	v, err = Admins.IsEmailExist(token.SiteId, rf.Value)
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
		func(rf AdminResetPasswdCodeForm, r render.Render, session sessions.Session, token Token) {
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
			// case "email":
			// 	if token.Appid != 5 {
			// 		sender.SendMail(sValue, "找回密码", "您的找回密码确认码为:["+sCode+"]")
			// 	}
			default:
			}

			session.Set("resetpasswdvalue", sValue)
			session.Set("resetpasswdcode", sCode)

			if token.Appid == 5 {
				var result AdminRegCodeResp
				result.Ret = 0
				result.Code = sCode
				r.JSON(200, result)
			} else {
				r.JSON(200, map[string]interface{}{"Ret": 0})
			}
		})

	m.Post("/resetpasswdbycode",
		binding.Bind(AdminResetPasswdFormByCode{}),
		check_form,
		func(rf AdminResetPasswdFormByCode, r render.Render, token Token, session sessions.Session) {

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

			var user *dbmodels.Admin
			var err error

			switch rf.Type {
			case "phone":
				user, err = dbmodels.GetUserByPhone(token.SiteId, rf.Value)
			// case "email":
			// 	user, err = users.GetUserByEmail(token.SiteId, rf.Value)
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

	m.Post("/:action",
		binding.Bind(AdminRegistryForm{}),
		check_form,
		func(rf AdminRegistryForm, params martini.Params, token Token, req *http.Request, r render.Render) {
			var result map[string]interface{}
			var id int64
			result = make(map[string]interface{})
			result["ret"] = 0
			action := params["action"]
			param, err := utils.StructToMap(&rf)
			id, err = PostAdmin(token.SiteId, action, param)
			if err != nil {
				r.JSON(200, map[string]interface{}{"ret": 1, "error": err.Error()})
				return
			}

			result["id"] = id
			r.JSON(200, result)
		})
	m.Put("/:id",
		binding.Bind(AdminModifyForm{}),
		check_form,
		func(rf AdminModifyForm, params martini.Params, token Token, req *http.Request, r render.Render) {
			var result map[string]interface{}
			var cnt int64
			result = make(map[string]interface{})
			result["ret"] = 0
			action := params["id"]
			uid := utils.Convert2Int64(action)
			param, err := utils.StructToMap(&rf)
			cnt, err = PutAdmin(uid, param)
			if err != nil {
				r.JSON(200, map[string]interface{}{"ret": 1, "error": err.Error()})
				return
			}

			result["count"] = cnt
			r.JSON(200, result)
		})
	m.Delete("/:id",
		binding.Bind(AdminDeleteForm{}),
		check_form,
		func(rf AdminDeleteForm, params martini.Params, token Token, req *http.Request, r render.Render) {
			var result map[string]interface{}
			var cnt int64
			result = make(map[string]interface{})
			result["ret"] = 0
			action := params["id"]
			uid := utils.Convert2Int64(action)
			param, err := utils.StructToMap(&rf)
			cnt, err = DeleteAdmin(uid, param)
			if err != nil {
				r.JSON(200, map[string]interface{}{"ret": 1, "error": err.Error()})
				return
			}

			result["count"] = cnt
			r.JSON(200, result)
		})
	m.Get("/:action",
		func(req *http.Request, params martini.Params, r render.Render, token Token) {
			var maps map[string]interface{}
			maps = make(map[string]interface{})
			maps["ret"] = 0

			appid, _ := strconv.ParseInt(req.URL.Query().Get("appid"), 10, 64)
			nums, _ := strconv.ParseInt(req.URL.Query().Get("pagenum"), 10, 64)
			start, _ := strconv.ParseInt(req.URL.Query().Get("pagestart"), 10, 64)
			order := req.URL.Query().Get("orderby")
			sort := req.URL.Query().Get("sort")
			action := params["action"]
			var stypes []string
			var contents []string
			s := req.URL.Query().Get("stype")
			if len(s) > 0 {
				stypes = strings.Split(s, ",")
			}

			c := req.URL.Query().Get("filter")

			cc := c
			if strings.Contains(cc, " and ") ||

				strings.Contains(cc, " or ") {
				err := errors.New("SQL条件有注入风险.")
				maps["ret"] = oacerr.ERR_INVALID_SQLFILTER
				maps["error"] = err.Error()
				r.JSON(200, maps)
				return
			}

			if len(c) > 0 {
				contents = strings.Split(c, ",")
			}
			id := utils.Convert2Int64(action)
			var datas interface{}
			var count int64
			var err error

			switch action {
			case "list":
				datas, count, err = GetAdminLists(token.SiteId, appid, stypes, contents, order, sort, nums, start)
			case "count":
				count, err = GetAdminCount(token.SiteId, appid, stypes, contents, order, sort)
			default:
				if id > 0 {
					datas, count, err = GetAdmin(id, token.SiteId, appid)
				}
			}

			if err != nil {
				maps["ret"] = 1
				maps["error"] = err.Error()
				r.JSON(200, maps)
				return
			}
			if action != "count" {
				maps["data"] = datas
			}
			maps["count"] = count
			r.JSON(200, maps)
		})

	return m
}
