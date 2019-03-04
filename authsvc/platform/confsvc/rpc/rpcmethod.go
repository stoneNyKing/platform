package rpc

import (
	"context"
	"errors"
	"platform/common/utils"
	"platform/confsvc/models"
	"platform/mskit/trace"
	"reflect"
	"strconv"
)

func GetConf(ctx context.Context,tracer trace.Tracer,appid, siteid int64,token string, data interface{})(interface{},error){
	logger.Info("调用GetConf方法")
	if data == nil {
		return 0,errors.New("parameters is null.")
	}

	param := data.(map[string]interface{})
	logger.Finest("params=%v",param)

	if param["siteid"] != nil {
		siteid = utils.Convert2Int64(param["siteid"])
	}

	var stypes,contents []string

	if param["stypes"] != nil {
		st := param["stypes"].([]interface{})
		for _,v := range st {
			t :=reflect.ValueOf(v)
			switch t.Kind() {
			case reflect.Float64:
				vv := strconv.FormatInt(int64(v.(float64)),10)
				stypes = append(stypes,vv)
			case reflect.Int64:
				vv := strconv.FormatInt(v.(int64),10)
				stypes = append(stypes,vv)
			case reflect.String:
				stypes = append(stypes,v.(string))
			}
		}
	}
	if param["contents"] != nil {
		st := param["contents"].([]interface{})
		for _,v := range st {
			t :=reflect.ValueOf(v)
			switch t.Kind() {
			case reflect.Float64:
				vv := strconv.FormatInt(int64(v.(float64)),10)
				contents = append(contents,vv)
			case reflect.Int64:
				vv := strconv.FormatInt(v.(int64),10)
				contents = append(contents,vv)
			case reflect.String:
				contents = append(contents,v.(string))
			}
		}
	}

	datas,_,err := models.GetDomainConfLists(siteid,appid,stypes,contents,"","",1000,0)

	r := make(map[string]interface{})
	r["data"] = datas
	if err != nil {
		r["error"] = err.Error()
		r["ret"] = 1
	}else{
		r["ret"] = 0
	}
	return r,err
}

func GetAppid(ctx context.Context,tracer trace.Tracer,appid, siteid int64,token string, data interface{})(interface{},error){
	logger.Info("调用GetAppid方法")

	if data == nil {
		return 0,errors.New("parameters is null.")
	}

	param := data.(map[string]interface{})
	logger.Finest("params=%v",param)

	if param["siteid"] != nil {
		siteid = utils.Convert2Int64(param["siteid"])
	}

	var stypes,contents []string

	if param["stypes"] != nil {
		st := param["stypes"].([]interface{})
		for _,v := range st {
			t :=reflect.ValueOf(v)
			switch t.Kind() {
			case reflect.Float64:
				vv := strconv.FormatInt(int64(v.(float64)),10)
				stypes = append(stypes,vv)
			case reflect.Int64:
				vv := strconv.FormatInt(v.(int64),10)
				stypes = append(stypes,vv)
			case reflect.String:
				stypes = append(stypes,v.(string))
			}
		}
	}
	if param["contents"] != nil {
		st := param["contents"].([]interface{})
		for _,v := range st {
			t :=reflect.ValueOf(v)
			switch t.Kind() {
			case reflect.Float64:
				vv := strconv.FormatInt(int64(v.(float64)),10)
				contents = append(contents,vv)
			case reflect.Int64:
				vv := strconv.FormatInt(v.(int64),10)
				contents = append(contents,vv)
			case reflect.String:
				contents = append(contents,v.(string))
			}
		}
	}

	datas,cnt,err := models.GetAppidLists(siteid,appid,stypes,contents,"","",1000,0)

	r := make(map[string]interface{})
	r["count"] = cnt
	r["data"] = datas
	if err != nil {
		r["error"] = err.Error()
		r["ret"] = 1
	}else{
		r["ret"] = 0
	}
	return r,err
}

