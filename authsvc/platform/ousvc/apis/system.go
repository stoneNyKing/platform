package apis

import (
	"github.com/boj/redistore"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"net/http"
	"strconv"
	"strings"
	"time"

	"platform/ousvc/config"
	. "platform/ousvc/models"

	"platform/models/sites"
	"platform/ousvc/dbmodels"
)

type (
	TimeResp struct {
		Ret  int
		Time int64
	}
)

func InitSystem() {
	dbmodels.NewEngine(config.Config.DbDriver, config.Config.DbAddr, config.Config.DbPort, config.Config.DbUser, config.Config.DbPasswd, config.Config.Database, config.Config.ObjectsSchema)
}

func getsiteid(req *http.Request) int64 {
	s1 := req.URL.Query().Get("siteid")
	if s1 == "" {
		s1 = req.URL.Query().Get("site")
	}

	siteid, err := strconv.ParseInt(s1, 10, 64)

	if err != nil {
		h := strings.Split(req.Host, ":")
		var host string
		var port int
		if len(h) > 1 {
			host = h[0]
			port, err = strconv.Atoi(h[1])
		} else {
			if len(h) > 0 {
				host = h[0]
				port = 80
			}
		}
		siteid, err = sites.GetSiteIdByHost(host, port)
		if err != nil {
			siteid = 1
		}
	}
	return siteid
}

func SystemHander() *martini.ClassicMartini {
	m := martini.Classic()
	m.Use(render.Renderer())

	store, _ := redistore.NewRediStore(10, "tcp", config.Config.SessionStoreIP+":"+config.Config.SessionStorePort, "", []byte(config.Config.SessionKey))
	m.Use(sessions.Sessions("user_session", store))

	m.Get("/time",
		func(r render.Render, req *http.Request) {
			var result TimeResp
			result.Ret = 0
			result.Time = time.Now().Unix()
			r.JSON(200, result)
		})

	m.Post("/usersite",
		binding.Bind(UserCheckForm{}),
		check_form,
		func(rf UserCheckForm, r render.Render) {
			sType := rf.Type
			sValue := rf.Value
			var err error
			var user *dbmodels.User

			switch sType {
			case "name":
				user = &dbmodels.User{Name: strings.ToLower(sValue)}
			case "phone":
				user = &dbmodels.User{Phone: strings.ToLower(sValue)}
			case "idcard":
				user = &dbmodels.User{Idcard: strings.ToLower(sValue)}
			case "email":
				user = &dbmodels.User{Email: strings.ToLower(sValue)}
			case "rfid":
				user = &dbmodels.User{Rfid: strings.ToLower(sValue)}
			case "imei":
				user = &dbmodels.User{Imei: strings.ToLower(sValue)}
			}

			user, err = dbmodels.GetUser(user)
			if err != nil {
				r.JSON(200, map[string]interface{}{"Ret": 110320, "Msg": err.Error()})
				return
			}
			r.JSON(200, map[string]interface{}{"Ret": 0, "SiteId": user.SiteId})
		})
	return m
}
