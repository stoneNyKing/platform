package admins

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/drewolson/testflight"
	"github.com/stretchr/testify/assert"

	"platform/lib/helper"
	"platform/oasvc/models"
)

func init() {
}

func GenerateToken(t *testing.T, handler func(appid string, token string)) {
	appid := "5"
	appkey := "9de0c791f3af4dd1935110b8bff363e2"

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
	})
}
