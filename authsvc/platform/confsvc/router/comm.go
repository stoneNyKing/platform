package router

import (
	"github.com/astaxie/beego/orm"
	"platform/common/utils"
	"platform/confsvc/models"
)

const(
	BEFE_USER	 = 1
	BEFE_ADMIN	 = 2
)

var befe map[int64]int

func InitBefe() {
	befe = GetAppidBefeFlag()
}

func GetAppidBefeFlag()(map[int64]int) {
	p := make(map[string]interface{})
	p["siteid"] = 1

	data,_,err := models.GetAppidLists(1,1,nil,nil,"","",1000,0)

	if err != nil {
		panic(err)
	}
	var list []orm.Params
	if data != nil {
		list = data.([]orm.Params)
	}
	AppIds := make(map[int64]int)
	for _, v := range list {
		if utils.Convert2Int(v["status"]) == 1 {
			if v["befeflag"] == nil {
				AppIds[utils.Convert2Int64(v["appid"])] = 0
			} else {
				AppIds[utils.Convert2Int64(v["appid"])] = utils.Convert2Int(v["befeflag"])
			}
		}
	}

	return AppIds
}