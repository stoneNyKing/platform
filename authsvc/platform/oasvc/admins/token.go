package admins

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	"github.com/boj/redistore"
	"github.com/martini-contrib/sessions"

	"platform/oasvc/config"
	. "platform/oasvc/models"
)

func TokenHander() *martini.ClassicMartini {
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Use(func(c martini.Context, res http.ResponseWriter, r render.Render) {
		step := 1
		rw := res.(martini.ResponseWriter)
		rw.Before(func(martini.ResponseWriter) {
			token := rw.Header().Get("Set-Cookie")
			if token != "" && step == 1 {
				step = 2
				as := strings.Split(token, ";")
				for _, v := range as {
					as2 := strings.Split(v, "=")
					if as2[0] == "admin_session" {
						token = as2[1]
						break
					}
				}

				var result TokenResp
				result.Ret = 0
				result.Token = token
				r.JSON(200, result)
			}
			rw.Header().Del("Set-Cookie")
		})
		c.Next()
	})

	store, _ := redistore.NewRediStore(10, "tcp", config.Config.SessionStoreIP+":"+config.Config.SessionStorePort, "", []byte(config.Config.SessionKey))
	m.Use(sessions.Sessions("admin_session", store))

	m.Post("/generate",
		//binding.Form(TokenGenerateForm{}),
		binding.Bind(TokenGenerateForm{}),
		check_form,
		func(rf TokenGenerateForm, session sessions.Session, req *http.Request, r render.Render) string {
			appid, _ := strconv.Atoi(req.URL.Query().Get("appid"))
			session.Set("appid", appid)
			return ""
		})

	m.Post("/check",
		check_token,
		func(token Token, session sessions.Session, req *http.Request, r render.Render) {
			r.JSON(200, map[string]interface{}{"Ret": 0, "Userid": token.UserId, "Msg": "admin token ok"})
		})
	return m
}
