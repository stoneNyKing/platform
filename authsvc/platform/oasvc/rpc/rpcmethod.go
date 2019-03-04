package rpc

import (
	"context"
	"errors"
	"platform/common/utils"
	"platform/mskit/trace"
	"platform/oasvc/admins"
	"reflect"
	"strconv"
)

func AddAdmin(ctx context.Context,tracer trace.Tracer, appid, siteid int64, token string, data interface{}) (interface{}, error) {
	logger.Info("调用AddAdmin方法")
	if data == nil {
		return 0, errors.New("parameters is null.")
	}

	param := data.(map[string]interface{})
	logger.Finest("params=%v", param)

	if param["siteid"] != nil {
		siteid = utils.Convert2Int64(param["siteid"])
	}

	id, err := admins.PostAdmin(siteid, "add", param)

	r := make(map[string]interface{})
	r["id"] = id
	if err != nil {
		r["error"] = err.Error()
		r["ret"] = 1
	} else {
		r["ret"] = 0
	}
	return r, err

}

func UpdateAdmin(ctx context.Context,tracer trace.Tracer, appid, siteid int64, token string, data interface{}) (interface{}, error) {
	logger.Info("调用UpdateAdmin方法")
	if data == nil {
		return 0, errors.New("parameters is null.")
	}

	param := data.(map[string]interface{})
	logger.Finest("params=%v", param)

	var id int64
	if param["id"] != nil {
		id = utils.Convert2Int64(param["id"])
	} else {
		return 0, errors.New("没有携带id")
	}

	cnt, err := admins.PutAdmin(id, param)
	r := make(map[string]interface{})
	r["count"] = cnt
	if err != nil {
		r["error"] = err.Error()
		r["ret"] = 1
	} else {
		r["ret"] = 0
	}
	return r, err

}

func ReadOrCreateAdmin(ctx context.Context,tracer trace.Tracer, appid, siteid int64, token string, data interface{}) (interface{}, error) {
	logger.Info("调用ReadOrCreateAdmin方法")
	if data == nil {
		return 0, errors.New("parameters is null.")
	}

	param := data.(map[string]interface{})
	logger.Finest("params=%v", param)

	if param["siteid"] != nil {
		siteid = utils.Convert2Int64(param["siteid"])
	}

	id, err := admins.PostAdmin(siteid, "readorcreate", param)
	r := make(map[string]interface{})
	r["id"] = id
	if err != nil {
		r["error"] = err.Error()
		r["ret"] = 1
	} else {
		r["ret"] = 0
	}
	return r, err

}

func DeleteAdmin(ctx context.Context,tracer trace.Tracer, appid, siteid int64, token string, data interface{}) (interface{}, error) {
	logger.Info("调用DeleteAdmin方法")
	if data == nil {
		return 0, errors.New("parameters is null.")
	}

	param := data.(map[string]interface{})
	logger.Finest("params=%v", param)

	var id int64
	if param["id"] != nil {
		id = utils.Convert2Int64(param["id"])
	} else {
		return 0, errors.New("没有携带id")
	}

	cnt, err := admins.DeleteAdmin(id, param)
	r := make(map[string]interface{})
	r["count"] = cnt
	if err != nil {
		r["error"] = err.Error()
		r["ret"] = 1
	} else {
		r["ret"] = 0
	}
	return r, err

}

func GetAdminList(ctx context.Context,tracer trace.Tracer, appid, siteid int64, token string, data interface{}) (interface{}, error) {
	logger.Info("调用GetAppid方法")

	if data == nil {
		return 0, errors.New("parameters is null.")
	}

	param := data.(map[string]interface{})
	logger.Finest("params=%v", param)

	if param["siteid"] != nil {
		siteid = utils.Convert2Int64(param["siteid"])
	}

	var stypes, contents []string

	if param["stypes"] != nil {
		st := param["stypes"].([]interface{})
		for _, v := range st {
			t := reflect.ValueOf(v)
			switch t.Kind() {
			case reflect.Float64:
				vv := strconv.FormatInt(int64(v.(float64)), 10)
				stypes = append(stypes, vv)
			case reflect.Int64:
				vv := strconv.FormatInt(v.(int64), 10)
				stypes = append(stypes, vv)
			case reflect.String:
				stypes = append(stypes, v.(string))
			}
		}
	}
	if param["contents"] != nil {
		st := param["contents"].([]interface{})
		for _, v := range st {
			t := reflect.ValueOf(v)
			switch t.Kind() {
			case reflect.Float64:
				vv := strconv.FormatInt(int64(v.(float64)), 10)
				contents = append(contents, vv)
			case reflect.Int64:
				vv := strconv.FormatInt(v.(int64), 10)
				contents = append(contents, vv)
			case reflect.String:
				contents = append(contents, v.(string))
			}
		}
	}

	datas, cnt, err := admins.GetAdminLists(siteid, appid, stypes, contents, "", "", 1000, 0)

	r := make(map[string]interface{})
	r["count"] = cnt
	r["data"] = datas
	if err != nil {
		r["error"] = err.Error()
		r["ret"] = 1
	} else {
		r["ret"] = 0
	}
	return r, err
}
