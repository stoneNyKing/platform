package apis

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/drewolson/testflight"
	"github.com/stretchr/testify/assert"

	"platform/lib/helper"
	"platform/ousvc/models"
)

func init() {
}

func GenerateToken(t *testing.T, handler func(appid string, token string)) {
	appid := "16"
	appkey := "238eefffeac907481a6a66d2b28657ce"

	testflight.WithServer(TokenHander(), func(r *testflight.Requester) {
		timeval := strconv.FormatInt(time.Now().Unix(), 10)
		response := r.Post("/generate?appid="+appid, testflight.FORM_ENCODED, "time="+timeval+"&requesttoken="+helper.Md5(appkey+timeval))
		assert.Equal(t, 200, response.StatusCode)

		var tokeresp models.TokenResp
		err := json.Unmarshal(response.RawBody, &tokeresp)
		assert.Nil(t, err, "format error")
		assert.Equal(t, tokeresp.Ret, 0, "Ret error")
		assert.NotEqual(t, tokeresp.Token, "", "Token error")
		handler(appid, tokeresp.Token)
	})
}

func Test_Token_generate(t *testing.T) {
	GenerateToken(t, func(appid string, token string) {
		println(appid, token)

		testflight.WithServer(TokenHander(), func(r *testflight.Requester) {
			response := r.Post("/check?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "")
			assert.Equal(t, 200, response.StatusCode)

			var resp models.Resp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
			assert.Equal(t, resp.Msg, "", "Token error")
		})

		testflight.WithServer(TokenHander(), func(r *testflight.Requester) {
			response := r.Post("/check?appid=1"+appid+"&token="+token, testflight.FORM_ENCODED, "")
			assert.Equal(t, 200, response.StatusCode)

			var resp models.Resp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.NotEqual(t, resp.Ret, 0, "Ret error")
			assert.NotEqual(t, resp.Msg, "", "Token error")
		})

		testflight.WithServer(TokenHander(), func(r *testflight.Requester) {
			response := r.Post("/check?appid="+appid+"&token=1"+token, testflight.FORM_ENCODED, "")
			assert.Equal(t, 200, response.StatusCode)

			var resp models.Resp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.NotEqual(t, resp.Ret, 0, "Ret error")
			assert.NotEqual(t, resp.Msg, "", "Token error")
		})
	})
}
