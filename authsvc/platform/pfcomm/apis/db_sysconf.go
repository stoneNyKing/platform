package apis

import (
	"context"
	"encoding/json"
	"github.com/smallnest/rpcx/client"
	"platform/common/utils"
)

func GetAppids(sdt,sda,confBasepath string)(map[int64]map[string]string) {
	p := make(map[string]interface{})
	p["siteid"] = 1

	_,data,err := RpcxGet(context.Background(),nil,1,0,"",sdt,sda,"ConfJSONRpc","GetAppid",
		confBasepath,client.Failtry,client.RoundRobin,p)

	if err != nil {
		panic(err)
	}

	var list []interface{}
	if data != nil {
		list = data.([]interface{})
	}else{
		panic("没有配置appid")
	}

	AppIds := make(map[int64]map[string]string)
	for _, v := range list {
		param := v.(map[string]interface{})
		if utils.Convert2Int(param["status"]) == 1 {
			if param["json"] == nil {
				AppIds[utils.Convert2Int64(param["appid"])] = make(map[string]string)
			} else {
				m := make(map[string]string)
				err = json.Unmarshal([]byte(utils.ConvertToString(param["json"])),&m)
				if err != nil {
					panic(err)
				}
				AppIds[utils.Convert2Int64(param["appid"])] = m
			}
			AppIds[utils.Convert2Int64(param["appid"])]["key"] = utils.ConvertToString(param["appkey"])
		}
	}

	return AppIds
}


func GetAppidBefeFlag(confBasepath,consulAddress string)(map[int64]int) {
	p := make(map[string]interface{})
	p["siteid"] = 1

	_,data,err := RpcxGetService(0,1,"","ConfJSONRpc","GetAppid",
		confBasepath,consulAddress,p)

	if err != nil {
		panic(err)
	}

	var list []interface{}
	if data != nil {
		list = data.([]interface{})
	}else{
		panic("没有配置appid")
	}

	AppIds := make(map[int64]int)
	for _, v := range list {
		param := v.(map[string]interface{})
		if utils.Convert2Int(param["status"]) == 1 {
			if param["befeflag"] == nil {
				AppIds[utils.Convert2Int64(param["appid"])] = 0
			} else {
				AppIds[utils.Convert2Int64(param["appid"])] = utils.Convert2Int(param["befeflag"])
			}
		}
	}

	return AppIds
}


func RetrieveAppidBefeFlag(	sdtype,sdaddr string,confBasepath string )(map[int64]int) {
	p := make(map[string]interface{})
	p["siteid"] = 1

	_,data,err := RpcxGet(context.Background(),nil,0,1,"",sdtype,sdaddr,"ConfJSONRpc","GetAppid",
			confBasepath,client.Failtry,client.RoundRobin,p)

	if err != nil {
		panic(err)
	}

	var list []interface{}
	if data != nil {
		list = data.([]interface{})
	}else{
		panic("没有配置appid")
	}

	AppIds := make(map[int64]int)
	for _, v := range list {
		param := v.(map[string]interface{})
		if utils.Convert2Int(param["status"]) == 1 {
			if param["befeflag"] == nil {
				AppIds[utils.Convert2Int64(param["appid"])] = 0
			} else {
				AppIds[utils.Convert2Int64(param["appid"])] = utils.Convert2Int(param["befeflag"])
			}
		}
	}

	return AppIds
}