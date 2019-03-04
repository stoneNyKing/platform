package apis

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	"github.com/boj/redistore"
	"github.com/martini-contrib/sessions"

	"platform/ousvc/config"
	. "platform/ousvc/models"
)

func TokenHander() *martini.ClassicMartini {
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Use(func(c martini.Context, req *http.Request, res http.ResponseWriter, r render.Render) {
		req.Header.Del("Cookie")
		step := 1
		cookies := ""
		rw := res.(martini.ResponseWriter)
		rw.Before(func(martini.ResponseWriter) {
			if step == 2 {
				rw.Header().Set("Set-Cookie", cookies)
			}
			if step == 1 {
				step = 2
				cookies = rw.Header().Get("Set-Cookie")
				as := strings.Split(cookies, ";")
				for _, v := range as {
					as2 := strings.Split(v, "=")
					if as2[0] == "user_session" {
						var result TokenResp
						result.Ret = 0
						result.Token = as2[1]
						r.JSON(200, result)
						return
					}
				}
			}
			// rw.Header().Del("Set-Cookie")
		})
	})

	store, _ := redistore.NewRediStore(10, "tcp", config.Config.SessionStoreIP+":"+config.Config.SessionStorePort, "", []byte(config.Config.SessionKey))
	m.Use(sessions.Sessions("user_session", store))

	m.Post("/generate",
		binding.Bind(TokenGenerateForm{}),
		check_form,
		func(rf TokenGenerateForm, session sessions.Session, req *http.Request, r render.Render) (int, string) {
			appid, _ := strconv.Atoi(req.URL.Query().Get("appid"))
			session.Set("appid", appid)
			return 200, ""
		})

	m.Post("/check_old",
		get_token,
		func(token Token, session sessions.Session, w http.ResponseWriter, req *http.Request, r render.Render) {
			appid, err := strconv.Atoi(req.URL.Query().Get("appid"))

			log.Finest("用户token检查: appid=%d,token=%+v", appid, token)

			if err == nil && appid == token.Appid && appid != 0 {
				r.JSON(200, map[string]interface{}{"Ret": 0, "Userid": token.UserId, "Msg": ""})
				return
			}

			u, _ := url.Parse("http://" + config.Config.HttpAddress)
			req.URL, _ = url.Parse(config.Config.Prefix + "/admin/token/check?appid=" + req.URL.Query().Get("appid") + "&token=" + req.URL.Query().Get("token"))
			proxy := httputil.NewSingleHostReverseProxy(u)
			proxy.ServeHTTP(w, req)
			return
		})
	m.Post("/check",
		check_token,
		func(token Token, session sessions.Session, req *http.Request, r render.Render) {
			r.JSON(200, map[string]interface{}{"Ret": 0, "Userid": token.UserId, "Msg": "user token ok"})
		})
	return m
}
