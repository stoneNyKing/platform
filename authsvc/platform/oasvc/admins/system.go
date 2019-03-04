package admins

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/boj/redistore"
	"github.com/go-martini/martini"
	// "github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"

	"platform/oasvc/config"
	// . "platform/api/models"
	"platform/models/sites"
	"platform/oasvc/dbmodels"
)

type (
	TimeResp struct {
		Ret  int
		Time int64
	}

	SiteProfileResp struct {
		Ret  int
		Data map[string]interface{}
	}
)

func InitSystem() {
	dbmodels.NewEngine(config.Config.DbDriver, config.Config.DbAddr, config.Config.DbPort, config.Config.DbUser, config.Config.DbPasswd, config.Config.Database, config.Config.ObjectsSchema)
	sites.NewEngine(config.Config.DbDriver, config.Config.DbAddr, config.Config.DbPort, config.Config.DbUser, config.Config.DbPasswd, config.Config.Database, config.Config.ObjectsSchema)
}

func getsiteid(req *http.Request) int64 {
	siteid, err := strconv.ParseInt(req.URL.Query().Get("siteid"), 10, 64)
	if err != nil || siteid == 0 {
		siteid, err = strconv.ParseInt(req.URL.Query().Get("site"), 10, 64)
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
			log.Finest("host=%s,port=%d", host, port)
			siteid, err = sites.GetSiteIdByHost(host, port)
			if err != nil {
				log.Error("获取租户标识出错(host=%s,port=%d): %v", host, port, err)
				siteid = 1
			}
		}
	}
	return siteid
}

func SystemHander() *martini.ClassicMartini {
	m := martini.Classic()
	m.Use(render.Renderer())
	store, _ := redistore.NewRediStore(10, "tcp", config.Config.SessionStoreIP+":"+config.Config.SessionStorePort, "", []byte(config.Config.SessionKey))
	m.Use(sessions.Sessions("admin_session", store))

	// m.Use(check_power("/api/v1/admin/system"))
	m.Get("/time",
		func(r render.Render, req *http.Request) {
			var result TimeResp
			result.Ret = 0
			result.Time = time.Now().Unix()
			r.JSON(200, result)
		})

	m.Get("/siteinfo", func(req *http.Request, r render.Render) {
		siteid := getsiteid(req)
		sKeys := req.URL.Query().Get("keys")
		var keys []string
		if err := json.Unmarshal([]byte(sKeys), &keys); err != nil {
			r.JSON(200, map[string]interface{}{"Ret": 1, "Msg": err.Error()})
			return
		}

		d, err := sites.GetProfile(siteid, keys)
		if err != nil {
			r.JSON(200, map[string]interface{}{"Ret": 2, "Msg": err.Error()})
			return
		}

		var result SiteProfileResp
		result.Ret = 0
		result.Data = d
		r.JSON(200, result)
	})

	return m
}
