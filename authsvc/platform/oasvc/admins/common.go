package admins

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-martini/martini"
	"github.com/libra9z/log4go"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"

	"fmt"
	"platform/oasvc/common"
	"platform/oasvc/config"
	"platform/oasvc/dbmodels"
)

var log log4go.Logger
var logger log4go.Logger

type (
	Token struct {
		Appid  int
		Token  string
		SiteId int64
		UserId int64
		IP     string
	}
)

func InitLogger() {
	log = common.Logger
	logger = log
	dbmodels.NewEngine(config.Config.DbDriver, config.Config.DbAddr, config.Config.DbPort, config.Config.DbUser, config.Config.DbPasswd, config.Config.Database, config.Config.ObjectsSchema)
}

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

func check_token(c martini.Context, req *http.Request, session sessions.Session, r render.Render) {
	var token Token
	var err error

	strurl := req.URL.Path
	fmt.Printf("req.path=%s\n", strurl)

	token.SiteId = getsiteid(req)
	token.Appid, err = strconv.Atoi(req.URL.Query().Get("appid"))
	if err != nil {
		r.JSON(200, map[string]interface{}{"Ret": 100, "Userid": 0, "Msg": "appid error"})
		return
	}

	token.Token = req.URL.Query().Get("token")
	if token.Token == "" {
		r.JSON(200, map[string]interface{}{"Ret": 101, "Userid": 0, "Msg": "token error"})
		return
	}

	req.Header.Del("Cookie")
	req.Header.Set("Cookie", "admin_session="+token.Token+";")
	v := session.Get("appid")
	if v == nil {
		r.JSON(200, map[string]interface{}{"Ret": 101, "Userid": 0, "Msg": "token Expires"})
		return
	}

	if v.(int) != token.Appid {
		r.JSON(200, map[string]interface{}{"Ret": 101, "Userid": 0, "Msg": "token is not generate for this appid"})
		return
	}

	token.IP = GetIP(req)

	userid := session.Get("userid")
	if userid != nil {
		token.UserId = userid.(int64)
	} else {
		token.UserId = 0
	}
	c.Map(token)
}

func check_power(path string) martini.Handler {
	return func(c martini.Context, req *http.Request, token Token, r render.Render) {

		log.Finest("path=%s", req.URL.Path)

		strurl := path + req.URL.Path

		if strurl == config.Config.Prefix+"/admin/system/init" {
			if token.UserId != 1 {
				r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": "Just admins can init"})
				return
			}
		} else if strings.Contains(strurl, "/admin/user") {
			return
		}

		if dbmodels.CheckPower(token.SiteId, token.UserId, strurl) == false {
			r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": "not have power"})
			return
		}
	}
}

func check_userid(params martini.Params, r render.Render, token Token) {
	userid, err := strconv.ParseInt(params["userid"], 10, 64)
	if err != nil {
		r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": "Not Have UserId"})
		return
	}

	if token.UserId == 0 {
		r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": "Not Login"})
		return
	}

	if userid != token.UserId {
		r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": "param UserId error"})
		return
	}
}

func check_form(errs binding.Errors, r render.Render) {
	if len(errs) > 0 {
		r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": errs})
	}
}
