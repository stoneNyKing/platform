package apis

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/drewolson/testflight"
	"github.com/stretchr/testify/assert"

	"platform/lib/helper"
	"platform/ousvc/dbmodels"
	. "platform/ousvc/models"
)

func init() {}

func Test_User_Init(t *testing.T) {
	users.Truncate()
	Test_Token_generate(t)
}

// resetpasswdcode
// resetpasswd

func Test_User_Check1(t *testing.T) {
	GenerateToken(t, func(appid string, token string) {
		testflight.WithServer(UserHander(), func(r *testflight.Requester) {
			response := r.Post("/check?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "type=email&value=test@test.com")
			assert.Equal(t, 200, response.StatusCode)
			var resp UserIdResp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 110330, "Ret error")
			assert.Equal(t, resp.UserId, 0, "userid error")
		})
	})
}

func Test_User_Registry_byEmail(t *testing.T) {
	GenerateToken(t, func(appid string, token string) {
		testflight.WithServer(UserHander(), func(r *testflight.Requester) {
			response := r.Post("/regcode?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "type=email&value=test@test.com")
			assert.Equal(t, 200, response.StatusCode)

			var coderesp UserRegCodeResp
			err := json.Unmarshal(response.RawBody, &coderesp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, coderesp.Ret, 0, "Ret error")
			assert.NotEqual(t, coderesp.Code, "", "Code error")

			testflight.WithServer(UserHander(), func(r *testflight.Requester) {
				response := r.Post("/registry?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "type=email&value=test@test.com&code="+coderesp.Code+"&passwd=123456")
				assert.Equal(t, 200, response.StatusCode)

				var regresp UserIdResp
				err := json.Unmarshal(response.RawBody, &regresp)
				assert.Nil(t, err, "format error")
				assert.Equal(t, regresp.Ret, 0, "Ret error")
				assert.NotEqual(t, regresp.UserId, 0, "UserId error")
			})
		})
	})
}

func Test_User_Check2(t *testing.T) {
	GenerateToken(t, func(appid string, token string) {
		testflight.WithServer(UserHander(), func(r *testflight.Requester) {
			response := r.Post("/check?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "type=email&value=test@test.com")
			assert.Equal(t, 200, response.StatusCode)
			var resp UserIdResp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
			assert.NotEqual(t, resp.UserId, 0, "userid error")
		})
	})

}

func Test_User_Registry_byPhone(t *testing.T) {
	GenerateToken(t, func(appid string, token string) {
		testflight.WithServer(UserHander(), func(r *testflight.Requester) {
			response := r.Post("/regcode?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "type=phone&value=13812345678")
			assert.Equal(t, 200, response.StatusCode)

			var coderesp UserRegCodeResp
			err := json.Unmarshal(response.RawBody, &coderesp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, coderesp.Ret, 0, "Ret error")
			assert.NotEqual(t, coderesp.Code, "", "Code error")

			testflight.WithServer(UserHander(), func(r *testflight.Requester) {
				response := r.Post("/registry?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "type=phone&value=13812345678&code="+coderesp.Code+"&passwd=123456")
				assert.Equal(t, 200, response.StatusCode)

				var regresp UserIdResp
				err := json.Unmarshal(response.RawBody, &regresp)
				assert.Nil(t, err, "format error")
				assert.Equal(t, regresp.Ret, 0, "Ret error")
				assert.NotEqual(t, regresp.UserId, 0, "UserId error")
			})
		})
	})
}

func GetCaptcha(t *testing.T, appid string, token string, handler func(CaptchaId string)) {
	testflight.WithServer(UserHander(), func(r *testflight.Requester) {
		response := r.Get("/captchaid?appid=" + appid + "&token=" + token)
		assert.Equal(t, 200, response.StatusCode)

		var captcharesp UserCaptchaIdResp
		err := json.Unmarshal(response.RawBody, &captcharesp)
		assert.Nil(t, err, "format error")
		assert.Equal(t, captcharesp.Ret, 0, "Ret error")
		assert.NotEqual(t, captcharesp.CaptchaId, "", "UserId error")
		handler(captcharesp.CaptchaId)
	})
}

func AutoLoginByid(t *testing.T, userid int64, passwd string, handler func(appid string, token string, userid int64)) {
	GenerateToken(t, func(appid string, token string) {
		passwd = helper.Md5(passwd + token)

		GetCaptcha(t, appid, token, func(CaptchaId string) {
			testflight.WithServer(UserHander(), func(r *testflight.Requester) {
				response := r.Post("/loginbyid?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "userid="+strconv.FormatInt(userid, 10)+"&passwd="+passwd+"&captcha="+CaptchaId)
				assert.Equal(t, 200, response.StatusCode)

				var regresp UserIdResp
				err := json.Unmarshal(response.RawBody, &regresp)
				assert.Nil(t, err, "format error")
				assert.Equal(t, regresp.Ret, 0, "Ret error")
				assert.NotEqual(t, regresp.UserId, 0, "UserId error")
				assert.Equal(t, regresp.UserId, userid, "Ret error")
				handler(appid, token, regresp.UserId)
			})
		})
	})
}

func AutoLogin(t *testing.T, name string, passwd string, handler func(appid string, token string, userid int64)) {
	GenerateToken(t, func(appid string, token string) {
		passwd = helper.Md5(passwd + token)

		GetCaptcha(t, appid, token, func(CaptchaId string) {
			testflight.WithServer(UserHander(), func(r *testflight.Requester) {
				response := r.Post("/login?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "name="+name+"&passwd="+passwd+"&captcha="+CaptchaId)
				assert.Equal(t, 200, response.StatusCode)

				var regresp UserIdResp
				err := json.Unmarshal(response.RawBody, &regresp)
				assert.Nil(t, err, "format error")
				assert.Equal(t, regresp.Ret, 0, "Ret error")
				assert.NotEqual(t, regresp.UserId, 0, "UserId error")
				handler(appid, token, regresp.UserId)
			})
		})
	})
}

func Test_User_Login_byId(t *testing.T) {
	GenerateToken(t, func(appid string, token string) {
		testflight.WithServer(UserHander(), func(r *testflight.Requester) {
			response := r.Post("/check?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "type=email&value=test@test.com")
			assert.Equal(t, 200, response.StatusCode)

			var resp UserIdResp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
			assert.NotEqual(t, resp.UserId, 0, "userid error")
			AutoLoginByid(t, resp.UserId, "123456", func(appid string, token string, userid int64) {
				assert.NotEqual(t, userid, 0, "Ret error")
			})
		})
	})
}

func Test_User_Login_byEmail(t *testing.T) {
	AutoLogin(t, "test@test.com", "123456", func(appid string, token string, userid int64) {
		assert.NotEqual(t, userid, 0, "Ret error")
	})
}

func Test_User_Login_byPhone(t *testing.T) {
	AutoLogin(t, "13812345678", "123456", func(appid string, token string, userid int64) {
		assert.NotEqual(t, userid, 0, "Ret error")
	})
}

func Test_User_Passwd(t *testing.T) {
	AutoLogin(t, "13812345678", "123456", func(appid string, token string, userid int64) {
		testflight.WithServer(UserHander(), func(r *testflight.Requester) {
			oldpasswd := "123456"
			response := r.Post("/"+strconv.FormatInt(userid, 10)+"/passwd?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "oldpasswd="+oldpasswd+"&newpasswd=abcdefg")
			assert.Equal(t, 200, response.StatusCode)

			var resp Resp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
		})

		testflight.WithServer(UserHander(), func(r *testflight.Requester) {
			oldpasswd := helper.Md5("abcdefg" + token)
			response := r.Post("/"+strconv.FormatInt(userid, 10)+"/passwd?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "oldpasswd="+oldpasswd+"&newpasswd=123456")
			assert.Equal(t, 200, response.StatusCode)

			var resp Resp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
		})

		testflight.WithServer(UserHander(), func(r *testflight.Requester) {
			oldpasswd := "123456"
			response := r.Post("/"+strconv.FormatInt(userid, 10)+"/passwd?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "oldpasswd="+oldpasswd+"&newpasswd=abcdefg")
			assert.Equal(t, 200, response.StatusCode)

			var resp Resp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
		})

		testflight.WithServer(UserHander(), func(r *testflight.Requester) {
			response := r.Post("/"+strconv.FormatInt(userid, 10)+"/passwd?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "oldpasswd=123456&newpasswd=abcdefg")
			assert.Equal(t, 200, response.StatusCode)

			var resp Resp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 101010, "Ret error")
		})

		testflight.WithServer(UserHander(), func(r *testflight.Requester) {
			response := r.Post("/"+strconv.FormatInt(userid, 10)+"/passwd?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "oldpasswd=abcdefg&newpasswd=123456")
			assert.Equal(t, 200, response.StatusCode)

			var resp Resp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
		})
	})
}

// func Test_User_Logout(t *testing.T) {
// 	data := PostResp(t, "/api/v1/user/3/logout", 0)
// 	println(data)
// }

func Test_User_GetProfile(t *testing.T) {
	AutoLogin(t, "13812345678", "123456", func(appid string, token string, userid int64) {
		testflight.WithServer(UserHander(), func(r *testflight.Requester) {
			response := r.Get("/" + strconv.FormatInt(userid, 10) + "/profile?appid=" + appid + "&token=" + token + "&type=" + "ihealth" + "&keys=" + "[\"Id\",\"Name\",\"Phone\",\"Email\",\"Passwd\",\"Created\",\"Updated\",\"name\"]")
			assert.Equal(t, 200, response.StatusCode)

			var userprofileresp UserProfileResp
			err := json.Unmarshal(response.RawBody, &userprofileresp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, userprofileresp.Ret, 0, "Ret error")

			d := userprofileresp.Data
			assert.NotEqual(t, d["Name"], "name", "Name error")
			assert.Equal(t, d["Phone"], "13812345678", "Phone error")
			assert.NotEqual(t, d["Email"], "email@email.com", "email error")
			assert.Nil(t, d["Passwd"], "Passwd error ")
			assert.Nil(t, d["name"], "name error ")
		})
	})
}

func Test_User_PutProfile(t *testing.T) {
	AutoLogin(t, "13812345678", "123456", func(appid string, token string, userid int64) {
		testflight.WithServer(UserHander(), func(r *testflight.Requester) {
			kv := make(map[string]interface{})
			kv["Name"] = "Value"
			kv["NameX"] = "ValueX"
			kv["name"] = "value"
			datas, err := json.Marshal(kv)
			if err != nil {
				t.Errorf("format error", err)
			}
			response := r.Post("/"+strconv.FormatInt(userid, 10)+"/profile?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "type="+"ihealth"+"&datas="+string(datas))
			assert.Equal(t, 200, response.StatusCode)

			var resp Resp
			err = json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
		})
	})
}

func Test_User_GetProfile2(t *testing.T) {
	AutoLogin(t, "13812345678", "123456", func(appid string, token string, userid int64) {
		testflight.WithServer(UserHander(), func(r *testflight.Requester) {
			response := r.Get("/" + strconv.FormatInt(userid, 10) + "/profile?appid=" + appid + "&token=" + token + "&type=" + "ihealth" + "&keys=" + "[\"Id\",\"Name\",\"Phone\",\"Email\",\"Passwd\",\"Created\",\"Updated\",\"name\"]")
			assert.Equal(t, 200, response.StatusCode)

			var userprofileresp UserProfileResp
			err := json.Unmarshal(response.RawBody, &userprofileresp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, userprofileresp.Ret, 0, "Ret error")

			d := userprofileresp.Data
			assert.NotEqual(t, d["Name"], "Value", "Name error")
			assert.NotEqual(t, d["NameX"], "ValueX", "Name error")
			assert.Equal(t, d["Phone"], "13812345678", "Phone error")
			assert.NotEqual(t, d["Email"], "email@email.com", "email error")
			assert.Nil(t, d["Passwd"], "Passwd error ")
			assert.Equal(t, d["name"], "value", "name error ")
		})
	})
}

func Test_User_ResetPasswd_byEmail(t *testing.T) {
	GenerateToken(t, func(appid string, token string) {
		testflight.WithServer(UserHander(), func(r *testflight.Requester) {
			response := r.Post("/resetpasswdcode?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "type=email&value=test@test.com")
			assert.Equal(t, 200, response.StatusCode)

			var coderesp UserRegCodeResp
			err := json.Unmarshal(response.RawBody, &coderesp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, coderesp.Ret, 0, "Ret error")
			assert.NotEqual(t, coderesp.Code, "", "Code error")

			testflight.WithServer(UserHander(), func(r *testflight.Requester) {
				response := r.Post("/resetpasswd?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "type=email&value=test@test.com&code="+coderesp.Code+"&passwd=654321")
				assert.Equal(t, 200, response.StatusCode)

				var resp Resp
				err := json.Unmarshal(response.RawBody, &resp)
				assert.Nil(t, err, "format error")
				assert.Equal(t, resp.Ret, 0, "Ret error")
			})
		})
	})
}

func Test_User_Logs(t *testing.T) {
	AutoLogin(t, "13812345678", "123456", func(appid string, token string, userid int64) {
		testflight.WithServer(UserHander(), func(r *testflight.Requester) {
			response := r.Get("/" + strconv.FormatInt(userid, 10) + "/logs?appid=" + appid + "&token=" + token + "&rows=20&page=0&level=10")
			assert.Equal(t, 200, response.StatusCode)

			var userlogsresult UserLogListResp
			err := json.Unmarshal(response.RawBody, &userlogsresult)
			assert.Nil(t, err, "format error")
			assert.Equal(t, userlogsresult.Ret, 0, "Ret error")
			assert.NotEqual(t, len(userlogsresult.Datas), 0)
		})
	})
}

func Test_User_Authlongin(t *testing.T) {
	AutoLogin(t, "13812345678", "123456", func(appid string, token string, userid int64) {
		testflight.WithServer(UserHander(), func(r *testflight.Requester) {
			response := r.Get("/" + strconv.FormatInt(userid, 10) + "/getauthlongin?appid=" + appid + "&token=" + token)
			assert.Equal(t, 200, response.StatusCode)

			var result UserCloneToken
			err := json.Unmarshal(response.RawBody, &result)
			assert.Nil(t, err, "format error")
			assert.Equal(t, result.Ret, 0, "Ret error")
			assert.NotEqual(t, result.UserId, 0, "UserId error")
			assert.NotEqual(t, result.Timestamp, "", "Timestamp error")
			assert.NotEqual(t, result.Auth, "", "Auth error")
			GenerateToken(t, func(appid string, token string) {
				response := r.Post("/authlongin?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "userid="+strconv.FormatInt(result.UserId, 10)+"&timestamp="+result.Timestamp+"&auth="+result.Auth)
				assert.Equal(t, 200, response.StatusCode)

				var resp UserIdResp
				err := json.Unmarshal(response.RawBody, &resp)
				assert.Nil(t, err, "format error")
				assert.Equal(t, resp.Ret, 0, "Ret error")
				assert.NotEqual(t, resp.UserId, 0, "Ret error")
			})
		})
	})
}
