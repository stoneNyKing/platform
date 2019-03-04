package main

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/drewolson/testflight"
	"github.com/stretchr/testify/assert"
	"platform/oasvc/config"

	"platform/lib/helper"
	. "platform/oasvc/models"
)

func init() {}

func Test_Get_Illegal(t *testing.T) {
	testflight.WithServer(PrivateHandler(config.Config.ApiUserPrefix), func(r *testflight.Requester) {
		response := r.Get("/ping?k=<&k2=>")
		assert.Equal(t, 200, response.StatusCode)
		var resp Resp
		err := json.Unmarshal(response.RawBody, &resp)
		assert.Nil(t, err, "format error")
		assert.Equal(t, resp.Ret, 100001, "Ret error")
	})
}

func Test_Post_Illegal(t *testing.T) {
	testflight.WithServer(PrivateHandler(config.Config.Prefix), func(r *testflight.Requester) {
		response := r.Post("/ping?k=1&k2=2", testflight.FORM_ENCODED, "pk=<&pk2=>")
		assert.Equal(t, 200, response.StatusCode)
		var resp Resp
		err := json.Unmarshal(response.RawBody, &resp)
		assert.Nil(t, err, "format error")
		assert.Equal(t, resp.Ret, 100001, "Ret error")
	})
}

func Test__ping(t *testing.T) {
	testflight.WithServer(PrivateHandler(config.Config.Prefix), func(r *testflight.Requester) {
		response := r.Get("/ping")
		assert.Equal(t, 200, response.StatusCode)
		var resp Resp
		err := json.Unmarshal(response.RawBody, &resp)
		assert.Nil(t, err, "format error")
		assert.Equal(t, resp.Ret, 0, "Ret error")
	})
}

func GenerateToken(t *testing.T, handler func(appid string, token string)) {
	appid := "5"
	appkey := "9de0c791f3af4dd1935110b8bff363e2"

	testflight.WithServer(PrivateHandler(config.Config.Prefix), func(r *testflight.Requester) {
		timeval := strconv.FormatInt(time.Now().Unix(), 10)
		response := r.Post("/api/v1/admin/token/generate?appid="+appid, testflight.FORM_ENCODED, "time="+timeval+"&requesttoken="+helper.Md5(appkey+timeval))

		assert.Equal(t, 200, response.StatusCode)

		var tokeresp TokenResp
		err := json.Unmarshal(response.RawBody, &tokeresp)
		assert.Nil(t, err, "format error")
		assert.Equal(t, tokeresp.Ret, 0, "Ret error")
		assert.NotEqual(t, tokeresp.Token, "", "Token error")
		handler(appid, tokeresp.Token)
	})
}

func AutoLogin(t *testing.T, name string, passwd string, handler func(appid string, token string, userid int64)) {
	GenerateToken(t, func(appid string, token string) {
		testflight.WithServer(PrivateHandler(config.Config.Prefix), func(r *testflight.Requester) {
			response := r.Post("/api/v1/admin/user/login?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "name="+name+"&passwd="+passwd)
			assert.Equal(t, 200, response.StatusCode)
			var regresp IdResp
			err := json.Unmarshal(response.RawBody, &regresp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, regresp.Ret, 0, "Ret error")
			assert.NotEqual(t, regresp.Id, 0, "UserId error")
			handler(appid, token, regresp.Id)
		})
	})
}

func Test_Admin_init(t *testing.T) {
	GenerateToken(t, func(appid string, token string) {
		testflight.WithServer(PrivateHandler(config.Config.Prefix), func(r *testflight.Requester) {
			response := r.Post("/api/v1/admin/system/init?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "")
			assert.Equal(t, 200, response.StatusCode)

			var regresp Resp
			err := json.Unmarshal(response.RawBody, &regresp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, regresp.Ret, 0, "Ret error")
		})
	})
}

func Test_User_Login(t *testing.T) {
	AutoLogin(t, "admin", "123456", func(appid string, token string, userid int64) {
		assert.NotEqual(t, userid, 0, "Ret error")
	})
}

func Test_User_Tree(t *testing.T) {
	AutoLogin(t, "admin", "123456", func(appid string, token string, userid int64) {
		testflight.WithServer(PrivateHandler(config.Config.Prefix), func(r *testflight.Requester) {
			response := r.Get("/api/v1/admin/resource/tree?appid=" + appid + "&token=" + token)
			assert.Equal(t, 200, response.StatusCode)
			var resp Resp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
		})
	})
}

func Test_User_Tree2(t *testing.T) {
	AutoLogin(t, "guest", "123456", func(appid string, token string, userid int64) {
		testflight.WithServer(PrivateHandler(config.Config.Prefix), func(r *testflight.Requester) {
			response := r.Get("/api/v1/admin/resource/tree?appid=" + appid + "&token=" + token)
			assert.Equal(t, 200, response.StatusCode)
			var resp Resp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 1, "Ret error")
		})
	})
}

func Test_Admin_User_tree(t *testing.T) {
	AutoLogin(t, "admin", "123456", func(appid string, token string, userid int64) {
		testflight.WithServer(PrivateHandler(config.Config.Prefix), func(r *testflight.Requester) {
			response := r.Get("/api/v1/admin/user/tree?appid=" + appid + "&token=" + token)
			assert.Equal(t, 200, response.StatusCode)

			var resp ResourceTreeResp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
			// assert.Equal(t, resp.Data.Id, 1, "Ret error")
			// assert.Equal(t, resp.Data.Name, "test", "Ret error")
		})
	})
}

func Test_Guest_User_tree(t *testing.T) {
	AutoLogin(t, "guest", "123456", func(appid string, token string, userid int64) {
		testflight.WithServer(PrivateHandler(config.Config.Prefix), func(r *testflight.Requester) {
			response := r.Get("/api/v1/admin/user/tree?appid=" + appid + "&token=" + token)
			assert.Equal(t, 200, response.StatusCode)

			var resp ResourceTreeResp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
			// assert.Equal(t, resp.Data.Id, 1, "Ret error")
			// assert.Equal(t, resp.Data.Name, "test", "Ret error")
		})
	})
}

func Test_Admin_porxy_ping(t *testing.T) {
	AutoLogin(t, "admin", "123456", func(appid string, token string, userid int64) {
		testflight.WithServer(PrivateHandler(config.Config.Prefix), func(r *testflight.Requester) {
			response := r.Get("/api/v1/admin/proxy/yanglao/ping?appid=" + appid + "&token=" + token)
			assert.Equal(t, 200, response.StatusCode)

			var resp Resp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
			// assert.Equal(t, resp.Data.Id, 1, "Ret error")
			// assert.Equal(t, resp.Data.Name, "test", "Ret error")
		})
	})
}

func Test_Guest_porxy_ping(t *testing.T) {
	AutoLogin(t, "guest", "123456", func(appid string, token string, userid int64) {
		testflight.WithServer(PrivateHandler(config.Config.Prefix), func(r *testflight.Requester) {
			response := r.Get("/api/v1/admin/proxy/yanglao/ping?appid=" + appid + "&token=" + token)
			assert.Equal(t, 200, response.StatusCode)

			var resp Resp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
			// assert.Equal(t, resp.Data.Id, 1, "Ret error")
			// assert.Equal(t, resp.Data.Name, "test", "Ret error")
		})
	})
}

func Test_Admin_porxy_version(t *testing.T) {
	AutoLogin(t, "admin", "123456", func(appid string, token string, userid int64) {
		testflight.WithServer(PrivateHandler(config.Config.Prefix), func(r *testflight.Requester) {
			response := r.Get("/api/v1/admin/proxy/yanglao/version?appid=" + appid + "&token=" + token)
			assert.Equal(t, 200, response.StatusCode)

			assert.Equal(t, response.Body, "", "Ret error")
			// assert.Equal(t, resp.Data.Id, 1, "Ret error")
			// assert.Equal(t, resp.Data.Name, "test", "Ret error")
		})
	})
}

func Test_Guest_porxy_version(t *testing.T) {
	AutoLogin(t, "guest", "123456", func(appid string, token string, userid int64) {
		testflight.WithServer(PrivateHandler(config.Config.Prefix), func(r *testflight.Requester) {
			response := r.Get("/api/v1/admin/proxy/yanglao/version?appid=" + appid + "&token=" + token)
			assert.Equal(t, 200, response.StatusCode)

			var resp Resp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 1, "Ret error")
			// assert.Equal(t, resp.Data.Id, 1, "Ret error")
			// assert.Equal(t, resp.Data.Name, "test", "Ret error")
		})
	})
}

func Test_Guest_porxy_callcenter_version(t *testing.T) {
	AutoLogin(t, "guest", "123456", func(appid string, token string, userid int64) {
		testflight.WithServer(PrivateHandler(config.Config.Prefix), func(r *testflight.Requester) {
			response := r.Get("/api/v1/admin/proxy/callcenter/version?appid=" + appid + "&token=" + token)
			assert.Equal(t, 200, response.StatusCode)

			var resp Resp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 1, "Ret error")
			// assert.Equal(t, resp.Data.Id, 1, "Ret error")
			// assert.Equal(t, resp.Data.Name, "test", "Ret error")
		})
	})
}
