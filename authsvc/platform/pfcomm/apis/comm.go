package apis

import (
	"context"
	"encoding/json"
	"errors"
	hoisie "github.com/libra9z/hoisie-redis"
	"github.com/smallnest/rpcx/client"
	credis "platform/common/redis"
	"platform/common/utils"
	"platform/mskit/rpcx"
	"platform/mskit/trace"
	"reflect"

	"fmt"
)

func TopicSubscribe(topic string, msgto chan<- hoisie.Message) {

	sub := make(chan string, 1)
	sub <- topic

	messages := make(chan hoisie.Message, 0)
	go credis.Subscribe(sub, nil, nil, nil, messages)

	defer close(sub)
	defer close(messages)

	for {
		select {
		case msg := <-messages:
			fmt.Printf("接收到订阅(topic=%s)消息.\n", topic)
			if msg.Channel == topic {
				msgto <- msg
			}
		}
	}
	return
}

func RpcxAddService(appid, siteid int64, token, servicename, methodname, basepath, consuladdr string, param map[string]interface{}) (id int64, err error) {

	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("%v\n", e)
			err = errors.New("panic error.")
		}
	}()

	req := new(rpcx.RpcRequest)
	ret := new(rpcx.RpcResponse)
	r := make(map[string]interface{})
	r["jsonrpc"] = "2.0"
	r["method"] = methodname
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

	rpcx.RpcCallWithConsul(basepath, consuladdr, servicename, "Services", 0, req, ret)

	fmt.Printf("response: %v\n", ret.Ret)
	var rpcret map[string]interface{}
	if ret.Ret != "" {
		err = json.Unmarshal([]byte(ret.Ret), &rpcret)

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

func RpcxUpdateService(appid, siteid int64, token, servicename, methodname, basepath, consuladdr string, param map[string]interface{}) (count int64, err error) {

	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("%v\n", e)
			err = errors.New("panic error.")
		}
	}()

	req := new(rpcx.RpcRequest)
	ret := new(rpcx.RpcResponse)
	r := make(map[string]interface{})
	r["jsonrpc"] = "2.0"
	r["method"] = methodname
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

	rpcx.RpcCallWithConsul(basepath, consuladdr, servicename, "Services", 0, req, ret)

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

func RpcxGetService(appid, siteid int64, token, servicename, methodname, basepath, consuladdr string, param map[string]interface{}) (count int64, data interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("%v\n", e)
			err = errors.New("panic error.")
		}
	}()

	req := new(rpcx.RpcRequest)
	ret := new(rpcx.RpcResponse)
	r := make(map[string]interface{})
	r["jsonrpc"] = "2.0"
	r["method"] = methodname
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

	rpcx.RpcCallWithConsul(basepath, consuladdr, servicename, "Services", 0, req, ret)

	var rpcret map[string]interface{}
	if ret.Ret != "" {
		err = json.Unmarshal([]byte(ret.Ret), &rpcret)

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
					if vr["data"] != nil {
						data = vr["data"]
					}

				}

			}
		}
	} else {
		return 0, nil, errors.New("server not response")
	}

	return
}

func RpcxGet(ctx context.Context,tracer trace.Tracer, appid,siteid int64, token string,
		sdtype,sdaddr string,servicename, methodname, basepath string,
		failmode client.FailMode,selectmode client.SelectMode,
		param map[string]interface{}) (count int64, data interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("%v\n", e)
			err = errors.New("panic error.")
		}
	}()

	req := new(rpcx.RpcRequest)
	//ret := new(rpcx.RpcResponse)
	r := make(map[string]interface{})
	r["jsonrpc"] = "2.0"
	r["method"] = methodname
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

	ret,err := rpcx.RpcxCall(ctx,tracer,sdtype,sdaddr,basepath, servicename, "Services",methodname, failmode,selectmode, req)

	var rpcret map[string]interface{}
	if ret.Ret != "" {
		err = json.Unmarshal([]byte(ret.Ret), &rpcret)

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
					if vr["data"] != nil {
						data = vr["data"]
					}

				}

			}
		}
	} else {
		return 0, nil, errors.New("server not response")
	}

	return
}

func RpcxAdd(ctx context.Context,tracer trace.Tracer, appid,siteid int64, token string,
	sdtype,sdaddr string,servicename, methodname, basepath string,
	failmode client.FailMode,selectmode client.SelectMode,
	param map[string]interface{}) (id int64, err error) {

	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("%v\n", e)
			err = errors.New("panic error.")
		}
	}()

	req := new(rpcx.RpcRequest)
	ret := new(rpcx.RpcResponse)
	r := make(map[string]interface{})
	r["jsonrpc"] = "2.0"
	r["method"] = methodname
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

	//rpcx.RpcCallWithConsul(basepath, consuladdr, servicename, "Services", 0, req, ret)
	ret,err = rpcx.RpcxCall(ctx,tracer,sdtype,sdaddr,basepath, servicename, "Services",methodname, failmode,selectmode, req)

	fmt.Printf("response: %v\n", ret.Ret)
	var rpcret map[string]interface{}
	if ret.Ret != "" {
		err = json.Unmarshal([]byte(ret.Ret), &rpcret)

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

func RpcxUpdate(ctx context.Context,tracer trace.Tracer,appid, siteid int64, token string,
		sdtype,sdaddr string,
		servicename,methodname, basepath string,
		failmode client.FailMode,selectmode client.SelectMode,
		param map[string]interface{}) (count int64, err error) {

	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("%v\n", e)
			err = errors.New("panic error.")
		}
	}()

	req := new(rpcx.RpcRequest)
	//ret := new(rpcx.RpcResponse)
	r := make(map[string]interface{})
	r["jsonrpc"] = "2.0"
	r["method"] = methodname
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
	req.WithTracer = true

	ret,err := rpcx.RpcxCall(ctx,tracer,sdtype,sdaddr,basepath, servicename, "Services", methodname,failmode,selectmode, req)

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
