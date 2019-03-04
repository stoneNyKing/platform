package admins

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/drewolson/testflight"
	"github.com/stretchr/testify/assert"

	. "platform/oasvc/models"
	// "platform/oasvc/dbmodels"
)

func init() {}

func Test_System_Init(t *testing.T) {
	GenerateToken(t, func(appid string, token string) {
		testflight.WithServer(SystemHander(), func(r *testflight.Requester) {
			response := r.Post("/init?appid="+appid+"&token="+token, testflight.FORM_ENCODED, "")
			assert.Equal(t, 200, response.StatusCode)
			fmt.Println(response.Body)
			var resp Resp
			err := json.Unmarshal(response.RawBody, &resp)
			assert.Nil(t, err, "format error")
			assert.Equal(t, resp.Ret, 0, "Ret error")
		})
	})
}
