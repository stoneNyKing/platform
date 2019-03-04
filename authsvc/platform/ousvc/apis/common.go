package apis

import (
	"net/http"
	"platform/common/utils"
	"strconv"
	"strings"

	"github.com/go-martini/martini"
	"github.com/libra9z/log4go"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

var log log4go.Logger

type (
	Token struct {
		Appid     int
		Token     string
		SiteId    int64
		UserId    int64 //支持userid和patientid两套token认证系统。
		PatientId int64 //理论上需要拆分token存储。现阶段简单起见，暂不做拆分，分开存储。
		IP        string
	}
)

func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

func GetIP(req *http.Request) string {
	addr := req.Header.Get("X-Real-IP")
	if addr == "" {
		addr = req.Header.Get("X-Forwarded-For")
		if addr == "" {
			addr = req.RemoteAddr
		}
	}
	ip := strings.Split(addr, ":")
	if len(ip) > 0 {
		if ip[0] != "[" {
			return ip[0]
		}
	}

	return strings.Split(req.RemoteAddr, ":")[0]
}

func get_token(c martini.Context, req *http.Request, session sessions.Session, r render.Render) {
	var token Token

	token.SiteId = getsiteid(req)
	token.Appid = 0
	token.UserId = 0
	token.PatientId = 0
	token.IP = GetIP(req)
	token.Token = req.URL.Query().Get("token")
	if token.Token != "" {
		req.Header.Del("Cookie")
		req.Header.Set("Cookie", "user_session="+token.Token+";")
	} else {
		cookies := req.Header.Get("Cookie")
		as := strings.Split(cookies, ";")
		for _, v := range as {
			as2 := strings.Split(v, "=")
			if as2[0] == "user_session" {
				token.Token = as2[1]
				break
			}
		}
	}

	siteid := session.Get("siteid")
	if siteid != nil {
		token.SiteId = siteid.(int64)
	} else {
		siteid = session.Get("site")
		if siteid != nil {
			token.SiteId = utils.Convert2Int64(siteid)
		}
	}
	appid := session.Get("appid")
	if appid != nil {
		token.Appid = appid.(int)
	}
	userid := session.Get("userid")
	if userid != nil {
		token.UserId = userid.(int64)
	}
	patientid := session.Get("patientid")
	if patientid != nil {
		token.PatientId = patientid.(int64)
	}
	c.Map(token)
}

func check_token(c martini.Context, req *http.Request, session sessions.Session, r render.Render) {
	var err error
	var token Token

	token.SiteId = getsiteid(req)
	token.Appid = 0
	token.UserId = 0
	token.PatientId = 0
	token.IP = GetIP(req)
	token.Appid, err = strconv.Atoi(req.URL.Query().Get("appid"))
	if err != nil {
		r.JSON(200, map[string]interface{}{"Ret": 100100, "Userid": 0, "Msg": "appid错误"})
		return
	}

	token.Token = req.URL.Query().Get("token")
	if token.Token != "" {
		req.Header.Del("Cookie")
		req.Header.Set("Cookie", "user_session="+token.Token+";")
	} else {
		cookies := req.Header.Get("Cookie")
		as := strings.Split(cookies, ";")
		for _, v := range as {
			as2 := strings.Split(v, "=")
			if as2[0] == "user_session" {
				token.Token = as2[1]
				break
			}
		}
	}

	if token.Token == "" {
		r.JSON(200, map[string]interface{}{"Ret": 100200, "Userid": 0, "Msg": "没有传递token"})
		return
	}

	v := session.Get("appid")
	if v == nil {
		log.Error("没有获取到session中的appid!")
		r.JSON(200, map[string]interface{}{"Ret": 100210, "Userid": 0, "Msg": "token无效"})
		return
	}

	if v.(int) != token.Appid {
		r.JSON(200, map[string]interface{}{"Ret": 100110, "Userid": 0, "Msg": "token的appid不匹配"})
		return
	}

	userid := session.Get("userid")
	if userid != nil {
		token.UserId = userid.(int64)
	}
	patientid := session.Get("patientid")
	if patientid != nil {
		token.PatientId = patientid.(int64)
	}

	log.Finest("step=2,(check_token)token = %+v ", token)

	c.Map(token)
}

func check_userid(params martini.Params, r render.Render, token Token) {
	userid, err := strconv.ParseInt(params["userid"], 10, 64)

	//return

	if err != nil {
		r.JSON(200, map[string]interface{}{"Ret": 100300, "Msg": "用户ID传递错误"})
		return
	}

	log.Finest("token=%+v", token)
	if token.UserId == 0 && token.PatientId == 0 {
		r.JSON(200, map[string]interface{}{"Ret": 100310, "Msg": "用户没有登录"})
		return
	}

	if userid != token.UserId && userid != token.PatientId {
		r.JSON(200, map[string]interface{}{"Ret": 100320, "Msg": "用户ID和session不一致"})
		return
	}
}

func check_form(errs binding.Errors, r render.Render) {
	if len(errs) > 0 {
		r.JSON(200, map[string]interface{}{"Ret": 100400, "Msg": errs})
	}
}
