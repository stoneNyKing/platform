package apis

import (
	"encoding/json"
	"errors"
	"fmt"
	"platform/common/utils"
	"platform/mskit/rpcx"
	"reflect"
)

func RpcxCreateOrUpdateGupUser(appid, siteid int64, token, basepath, consuladdr string, param map[string]interface{}) (id int64, err error) {
	defer func(){
		if e:=recover();e!=nil {
			fmt.Errorf("panic = %v",e)
			err = errors.New(" panic ")
			id=0
			return
		}
	}()


	req := new(rpcx.RpcRequest)
	ret := new(rpcx.RpcResponse)
	r := make(map[string]interface{})
	r["jsonrpc"] = "2.0"
	r["method"] = "ReadOrCreateUser"
	r["params"] = param
	r["id"] = "1"

	re, err := json.Marshal(&r)
	if err != nil {
		return
	}
	req.Req = string(re)
	req.Appid = appid
	req.SiteId = siteid
	req.Token = token

	rpcx.RpcCallWithConsul(basepath, consuladdr, "GupJSONRpc", "Services", 0, req, ret)

	fmt.Printf("response: %v\n", ret.Ret)
	var rpcret map[string]interface{}
	if ret.Ret != "" {
		err = json.Unmarshal([]byte(ret.Ret), &rpcret)
		if err != nil {
			return
		}
		if rpcret != nil {
			v := rpcret
			if v["result"] != nil {
				t := reflect.ValueOf(v["result"])
				switch t.Kind() {
				case reflect.Int64:
					id = v["result"].(int64)
				case reflect.Int:
					id = int64(v["result"].(int))
				case reflect.Map:
					vr := v["result"].(map[string]interface{})
					if vr["id"] != nil {
						id = utils.Convert2Int64(vr["id"])
					}
					if vr["error"] != nil {
						err = errors.New(vr["error"].(string))
					}
				}

			}
		}
	} else {
		return 0, errors.New("server not response")
	}

	return
}

func RpcxAddGupUser(appid, siteid int64, token, basepath, consuladdr string, param map[string]interface{}) (id int64, err error) {
	defer func(){
		if e:=recover();e!=nil {
			fmt.Errorf("panic = %v",e)
			err = errors.New(" panic ")
			id=0
			return
		}
	}()


	req := new(rpcx.RpcRequest)
	ret := new(rpcx.RpcResponse)
	r := make(map[string]interface{})
	r["jsonrpc"] = "2.0"
	r["method"] = "AddUser"
	r["params"] = param
	r["id"] = "1"

	re, err := json.Marshal(&r)
	if err != nil {
		return
	}
	req.Req = string(re)
	req.Appid = appid
	req.SiteId = siteid
	req.Token = token

	rpcx.RpcCallWithConsul(basepath, consuladdr, "GupJSONRpc", "Services", 0, req, ret)

	fmt.Printf("response: %v\n", ret.Ret)
	var rpcret map[string]interface{}
	if ret.Ret != "" {
		err = json.Unmarshal([]byte(ret.Ret), &rpcret)
		if err != nil {
			return
		}
		if rpcret != nil {
			v := rpcret
			if v["result"] != nil {
				t := reflect.ValueOf(v["result"])
				switch t.Kind() {
				case reflect.Int64:
					id = v["result"].(int64)
				case reflect.Int:
					id = int64(v["result"].(int))
				case reflect.Map:
					vr := v["result"].(map[string]interface{})
					if vr["id"] != nil {
						id = utils.Convert2Int64(vr["id"])
					}
					if vr["error"] != nil {
						err = errors.New(vr["error"].(string))
					}
				}

			}
		}
	} else {
		return 0, errors.New("server not response")
	}

	return
}

func RpcxUpdateGupUser(appid, siteid int64, token, basepath, consuladdr string, param map[string]interface{}) (count int64, err error) {

	defer func(){
		if e:=recover();e!=nil {
			fmt.Errorf("panic = %v",e)
			err = errors.New(" panic ")
			return
		}
	}()


	req := new(rpcx.RpcRequest)
	ret := new(rpcx.RpcResponse)
	r := make(map[string]interface{})
	r["jsonrpc"] = "2.0"
	r["method"] = "UpdateUser"
	r["params"] = param
	r["id"] = "1"

	re, err := json.Marshal(&r)
	if err != nil {
		return
	}
	req.Req = string(re)
	req.Appid = appid
	req.SiteId = siteid
	req.Token = token

	rpcx.RpcCallWithConsul(basepath, consuladdr, "GupJSONRpc", "Services", 0, req, ret)

	fmt.Printf("response: %v\n", ret.Ret)
	var rpcret map[string]interface{}
	if ret.Ret != "" {
		err = json.Unmarshal([]byte(ret.Ret), &rpcret)
		if err != nil {
			return
		}
		if rpcret != nil {
			v := rpcret
			if v["result"] != nil {
				t := reflect.ValueOf(v["result"])
				switch t.Kind() {
				case reflect.Int64:
					count = v["result"].(int64)
				case reflect.Int:
					count = int64(v["result"].(int))
				case reflect.Map:
					vr := v["result"].(map[string]interface{})
					if vr["count"] != nil {
						count = utils.Convert2Int64(vr["count"])
					}
					if vr["error"] != nil {
						err = errors.New(vr["error"].(string))
					}
				}

			}
		}
	} else {
		return 0, errors.New("server not response")
	}

	return
}

