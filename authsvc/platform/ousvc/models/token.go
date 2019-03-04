package models

import (
	"fmt"
	"math"
	"net/http"
	"platform/ousvc/config"
	"platform/ousvc/dbmodels"
	"strconv"
	"time"

	"github.com/martini-contrib/binding"

	"platform/common/utils"
	"platform/lib/helper"
)

var AppIds map[int64]map[string]string

func InitAppids() {
	//AppIds = apis.GetAppids(config.Config.RpcxConfBasepath,config.Config.ConsulAddress)
	dbmodels.NewEngine(config.Config.DbDriver, config.Config.DbAddr, config.Config.DbPort, config.Config.DbUser, config.Config.DbPasswd, config.Config.Database, config.Config.ObjectsSchema)
	list, err := dbmodels.Search()
	if err != nil {
		panic(err)
	}

	AppIds = make(map[int64]map[string]string)
	for _, v := range list {
		if v.Status == 1 {
			if v.Json == nil {
				AppIds[v.Appid] = make(map[string]string)
			} else {

				AppIds[v.Appid] = v.Json
			}
			AppIds[v.Appid]["key"] = v.Appkey
		}
	}

	fmt.Printf("appids = %+v\n", AppIds)

}

type (
	TokenGenerateForm struct {
		Time         int64  `form:"time" binding:"required"`
		Requesttoken string `form:"requesttoken" binding:"required"`
	}

	TokenResp struct {
		Ret   int
		Token string
	}
)

func (rf TokenGenerateForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	appids := req.URL.Query().Get("appid")

	if math.Abs(float64(time.Now().Unix()-rf.Time)) >= 5*60 {
		return append(errors, binding.Error{
			FieldNames:     []string{"time"},
			Classification: "error",
			Message:        "时间戳错误超时",
		})
	}
	appid := utils.Convert2Int64(appids)
	appkey := AppIds[appid]["key"]
	if appkey == "" {
		return append(errors, binding.Error{
			FieldNames:     []string{"appid"},
			Classification: "error",
			Message:        "appid不存在",
		})
	}

	rightoken := helper.Md5(appkey + strconv.FormatInt(rf.Time, 10))
	// println("right key:"+rightoken, "RequestToken:", rf.Requesttoken)
	if rf.Requesttoken != rightoken {
		return append(errors, binding.Error{
			FieldNames:     []string{"requesttoken"},
			Classification: "error",
			Message:        "requesttoken验证错误",
		})
	}
	return errors
}
