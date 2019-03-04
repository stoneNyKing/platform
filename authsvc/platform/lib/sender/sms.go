package sender

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/libra9z/log4go"
	"github.com/franela/goreq"
	"fmt"
	"platform/mskit/rpcx"
	"platform/common/utils"
	"reflect"
)

type M struct {
	Result int    `json:"result"`
	Reason string `json:"reason"`
}

type Sms struct {
	Receivers []string `json:"receivers"`
	Content   string   `json:"content"`
}

// curl -d "[{\"receivers\":[\"18952017328123\"],\"content\": \"测试消息下发1\"}]" "http://114.215.237.162:18088/sms/push?appId=1&modId=10&type=register"
func SendSms(url string, appid int,siteid int64,token string, modid int, stype string, srcip string, tophone string, title string, body string) (bool, error) {
	log := log4go.Global
	surl := url + "?appid=" + strconv.Itoa(appid) + "&modId=" + strconv.Itoa(modid) + "&type=" + stype + "&srcIp=" + srcip+"&site="+
		strconv.Itoa(int(siteid))+"&token="+token
	log.Info("SendSms:", appid, modid, stype, srcip, tophone, title, body)

	var ds []Sms
	ds = make([]Sms, 0)
	var d Sms
	d.Receivers = append(d.Receivers, tophone)
	d.Content = body
	ds = append(ds, d)

	ss, err := json.Marshal(ds)
	if err != nil {
		log.Info("SendSms:", err)
	}

	log.Info(surl, string(ss))
	res, err := goreq.Request{
		Method:      "POST",
		Accept:      "application/json; charset=utf-8",
		ContentType: "application/json; charset=utf-8",
		Uri:         surl,
		Body:        ss,
	}.Do()

	if err != nil {
		log.Error("SendSms:", err)
	}

	var m M
	res.Body.FromJsonTo(&m)
	log.Info("SendSms:", res.StatusCode, m.Result, m.Reason)
	if m.Result == 0 {
		log.Error("SendSms return ok:")
		return true, nil
	} else {
		log.Error("SendSms return error:", m.Result, m.Reason)
		return false, errors.New(m.Reason)
	}
}


func RpcxSendSms(log log4go.Logger,basepath,consuladdr string, appid int64,siteid int64,token string,  stype string, srcip string, tophone string, title string, body string) (b bool, err error) {

	defer func(){
		if e := recover();e !=nil {
			fmt.Printf("Rpcx请求错误: %v\n",e)
			err = errors.New("panic error.")
		}
	}()

	log.Info("SendSms: appid=%v,stype=%v,srcip=%v,phone=%v,title=%v,body=%v ", appid,  stype, srcip, tophone, title, body)

	var ds []Sms
	ds = make([]Sms, 0)
	var d Sms
	d.Receivers = append(d.Receivers, tophone)
	d.Content = body
	ds = append(ds, d)

	var count int64

	param := make(map[string]interface{})

	param["modid"] = 50
	param["remoteaddr"] = srcip
	param["stype"] = stype
	param["action"] = title
	param["parameters"] = ds

	req := new(rpcx.RpcRequest)
	ret := new(rpcx.RpcResponse)
	r := make(map[string]interface{})
	r["jsonrpc"] = "2.0"
	r["method"] = "SendSms"
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

	rpcx.RpcCallWithConsul(basepath,consuladdr, "SmsJSONRpc", "Services", 0, req, ret)

	log.Finest("response: %v", ret.Ret)

	var rpcret map[string]interface{}
	if ret.Ret != "" {
		err = json.Unmarshal([]byte(ret.Ret),&rpcret)
		if err != nil {
			return
		}
		if rpcret != nil {
			v := rpcret
			if v["result"]!= nil {
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
	}else{
		return false,errors.New("server not response")
	}
	if count >0 {
		b = true
	}
	return b,nil

}
