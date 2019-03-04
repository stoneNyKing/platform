package apis

import (
	"encoding/json"
	"fmt"
	"errors"
	"platform/common/utils"
	"platform/mskit/rpcx"
)

func CreateOpEventRpcx(basepath,consuladdr string,appid,siteid int64,token string,param map[string]interface{}) (id int64,err error){

	defer func(){
		if e:=recover();e!=nil {
			err = errors.New(" panic ")
			id=0
			return
		}
	}()

	req := new(rpcx.RpcRequest)
	ret := new(rpcx.RpcResponse)
	r := make(map[string]interface{})
	r["jsonrpc"] = "2.0"
	r["method"] = "AddStaffOpLog"
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

	rpcx.RpcCallWithConsul(basepath,consuladdr, "OamJSONRpc", "Services", 0, req, ret)

	fmt.Printf("response: %v\n", ret.Ret)
	var rpcret map[string]interface{}
	if ret.Ret != "" {
		err = json.Unmarshal([]byte(ret.Ret),&rpcret)
		if err != nil {
			return
		}
		if rpcret["id"] != nil {
			id = utils.Convert2Int64(rpcret["id"])
			err = nil
		}
	}else{
		return 0,errors.New("server not response")
	}

	return
}

func SendOpEventRpcx(basepath,consuladdr string,appid,siteid int64,token string,objectid int64,objecttype int,
	optype int,ovalue,nvalue string,userid,operatorid int64,staffdbname string)(id int64,err error) {

	param := make(map[string]interface{})
	param["objectid"] = objectid
	param["objecttype"] = objecttype
	param["operatetype"] = optype
	param["userid"] = userid
	param["operatorid"] = operatorid
	param["originvalue"] = ovalue
	param["newvalue"] = nvalue
	param["siteid"] = siteid


	return CreateOpEventRpcx(basepath,consuladdr,appid,siteid,token,param)
}

func SendOpEvent(appid,siteid int64,token string,urlprefix string,objectid int64,objecttype int,
	optype int,ovalue,nvalue string,userid,operatorid int64,staffdbname string)(id int64,err error) {

	param := make(map[string]interface{})
	param["objectid"] = objectid
	param["objecttype"] = objecttype
	param["operatetype"] = optype
	param["userid"] = userid
	param["operatorid"] = operatorid
	param["originvalue"] = ovalue
	param["newvalue"] = nvalue
	param["siteid"] = siteid


	return CreateOpEvent(appid,siteid,token,urlprefix,param)
}

func CreateOpEvent(appid,siteid int64,token string,urlprefix string,param map[string]interface{}) (id int64,err error){
	defer func(){
		if e := recover();e !=nil {
			fmt.Printf("%v\n",e)
			err = errors.New("panic error.")
		}
	}()

	str,err := json.Marshal(param)
	if err != nil {
		return 0,err
	}
	uri := fmt.Sprintf("%s/349/staffop?appid=%d&site=%d&token=%s",urlprefix,appid,siteid,token)

	fmt.Printf("uri=%s\nbody json: %s\n",uri,string(str))

	var vr interface{}
	vr,err = utils.ServicePost(uri,string(str))

	if err != nil {
		return 0,err
	}
	vm := vr.(map[string]interface{})
	if vm != nil && vm["ret"] !=nil {
		ret := utils.Convert2Int(vm["ret"])
		if ret == 1 {
			return 0,errors.New("cannot create event.")
		}

		id = utils.Convert2Int64(vm["id"])
	}

	return id,err
}



func CreateApiInvokeEventRpcx(basepath,consuladdr string,appid,siteid int64,token string,param map[string]interface{}) (id int64,err error){

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
	r["method"] = "AddApiInvokeLog"
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

	rpcx.RpcCallWithConsul(basepath,consuladdr, "OamJSONRpc", "Services", 0, req, ret)

	fmt.Printf("response: %v\n", ret.Ret)
	var rpcret map[string]interface{}
	if ret.Ret != "" {
		err = json.Unmarshal([]byte(ret.Ret),&rpcret)
		if err != nil {
			return
		}
		if rpcret["id"] != nil {
			id = utils.Convert2Int64(rpcret["id"])
			err = nil
		}
	}else{
		return 0,errors.New("server not response")
	}

	return
}



func RpcxApiInvokeLog(basepath,consuladdr string,appid,siteid int64,token string,sid int,successflag int,
	servicename string,ip,path string,userid int64 )(id int64,err error) {

	param := make(map[string]interface{})
	param["userid"] = userid
	param["successflag"] = successflag
	param["servicename"] = servicename
	param["sid"] = sid
	param["path"] = path
	param["ip"] = ip
	param["siteid"] = siteid


	return CreateApiInvokeEventRpcx(basepath,consuladdr,appid,siteid,token,param)
}
