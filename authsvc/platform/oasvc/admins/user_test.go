package admins

import (
	"encoding/json"
	"testing"

	"github.com/drewolson/testflight"
	"github.com/stretchr/testify/assert"

	. "platform/oasvc/models"
	// "platform/oasvc/dbmodels"
)

func init() {}

func Test_Admin_User_registry(t *testing.T) {
	GenerateToken(t, func(appid string, token string) {
		testflight.WithServer(UserHander(true), func(r *testflight.Requester) {
			response := r.Post("/registry?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "roleid=1&name=test&jobnumber=jobnumber&passwd=123456&email=123@123.com&phone=1234567890&description=description&effectiveTime=2014-05-01+1%3a1%3a0&expireTime=1")
			assert.Equal(t, 200, response.StatusCode)
			var resp IdResp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
			assert.Equal(t, resp.Id, 1, "Ret error")
		})
	})
}

func Test_Admin_User_login(t *testing.T) {
	GenerateToken(t, func(appid string, token string) {
		testflight.WithServer(UserHander(true), func(r *testflight.Requester) {
			response := r.Post("/login?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "name=test&passwd=123456")
			assert.Equal(t, 200, response.StatusCode)
			var resp IdResp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
			assert.Equal(t, resp.Id, 1, "Ret error")
		})
	})
}

func Test_Admin_User_get(t *testing.T) {
	GenerateToken(t, func(appid string, token string) {
		testflight.WithServer(UserHander(true), func(r *testflight.Requester) {
			response := r.Get("/get?appid=" + appid + "&token=" + token + "&id=1")
			assert.Equal(t, 200, response.StatusCode)
			var resp AdminResp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
			assert.Equal(t, resp.Data.Id, 1, "Ret error")
			assert.Equal(t, resp.Data.Name, "test", "Ret error")
		})
	})
}

func AutoLogin(t *testing.T, name string, passwd string, handler func(appid string, token string, userid int64)) {
	GenerateToken(t, func(appid string, token string) {
		testflight.WithServer(UserHander(true), func(r *testflight.Requester) {
			response := r.Post("/login?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "name="+name+"&passwd="+passwd)
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

func Test_User_Login(t *testing.T) {
	AutoLogin(t, "test", "123456", func(appid string, token string, userid int64) {
		assert.NotEqual(t, userid, 0, "Ret error")
	})
}

func Test_Admin_User_modify(t *testing.T) {
	GenerateToken(t, func(appid string, token string) {
		testflight.WithServer(UserHander(true), func(r *testflight.Requester) {
			response := r.Post("/modify?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "id=1&roleid=1&name=test23&jobnumber=jobnumber&passwd=654321&email=123@123.com&phone=1234567890&description=description&effectiveTime=2014-05-01+1%3a1%3a0&expireTime=1")
			assert.Equal(t, 200, response.StatusCode)
			var resp Resp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
		})
	})
}

func Test_Admin_User_get2(t *testing.T) {
	GenerateToken(t, func(appid string, token string) {
		testflight.WithServer(UserHander(true), func(r *testflight.Requester) {
			response := r.Get("/get?appid=" + appid + "&token=" + token + "&id=1")
			assert.Equal(t, 200, response.StatusCode)
			var resp AdminResp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
			assert.Equal(t, resp.Data.Id, 1, "Ret error")
			assert.Equal(t, resp.Data.Name, "test", "Ret error")
			assert.Equal(t, resp.Data.Passwd, "654321", "Ret error")
		})
	})
}

func Test_Admin_User_List(t *testing.T) {
	GenerateToken(t, func(appid string, token string) {
		testflight.WithServer(UserHander(true), func(r *testflight.Requester) {
			response := r.Post("/registry?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "roleid=2&name=test12&jobnumber=jobnumbe1r&passwd=123456&email=1231@1213.com&Phone=12314567890&description=description&effectiveTime=2014-05-01+1%3a1%3a0&expireTime=1")
			assert.Equal(t, 200, response.StatusCode)
			var resp IdResp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
			assert.Equal(t, resp.Id, 2, "Ret error")
		})

		testflight.WithServer(UserHander(true), func(r *testflight.Requester) {
			response := r.Get("/list?appid=" + appid + "&token=" + token)
			assert.Equal(t, 200, response.StatusCode)
			var resp AdminListResp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
			assert.Len(t, resp.Datas, 2, "Ret error")
		})

		testflight.WithServer(UserHander(true), func(r *testflight.Requester) {
			response := r.Get("/list?appid=" + appid + "&token=" + token + "&roleid=1")
			assert.Equal(t, 200, response.StatusCode)
			var resp AdminListResp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
			assert.Len(t, resp.Datas, 1, "Ret error")
		})
	})
}

func Test_Admin_User_delete(t *testing.T) {
	GenerateToken(t, func(appid string, token string) {
		testflight.WithServer(UserHander(true), func(r *testflight.Requester) {
			response := r.Post("/delete?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "id=1")
			assert.Equal(t, 200, response.StatusCode)
			var resp Resp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
		})
	})
}
