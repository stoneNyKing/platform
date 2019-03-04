package apis

import (
	"encoding/json"
	"testing"

	"github.com/drewolson/testflight"
	"github.com/stretchr/testify/assert"
)

func Test_System_Time(t *testing.T) {
	testflight.WithServer(SystemHander(), func(r *testflight.Requester) {
		response := r.Get("/time")
		assert.Equal(t, 200, response.StatusCode)
		var regresp TimeResp
		err := json.Unmarshal([]byte(response.Body), &regresp)
		assert.Nil(t, err, "format error")
		assert.Equal(t, regresp.Ret, 0, "Ret error")
		assert.NotEqual(t, regresp.Time, 0, "Time error")
	})
}
