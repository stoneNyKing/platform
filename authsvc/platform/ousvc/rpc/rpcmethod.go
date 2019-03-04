package rpc

import (
	"context"
	"errors"
	"platform/common/utils"
	"platform/mskit/trace"
	"platform/ousvc/apis"
	"platform/ousvc/dbmodels"
	"reflect"
	"strconv"
)

func CheckUser(ctx context.Context, tracer trace.Tracer, appid, siteid int64, token string, data interface{}) (interface{}, error) {
	logger.Info("调用CheckUser方法")
	if data == nil {
		return 0, errors.New("parameters is null.")
	}

	param := data.(map[string]interface{})
	logger.Finest("params=%v", param)

	var stype, value, sId, sPhone string
	if param["type"] != nil {
		stype = param["type"].(string)
	}
	if param["value"] != nil {
		value = param["value"].(string)
	}
	if param["idcard"] != nil {
		sId = param["idcard"].(string)
	}
	if param["phone"] != nil {
		sPhone = param["phone"].(string)
	}

	var v bool
	var err error
	var user *dbmodels.User
	var r interface{}

	switch stype {
	case "name":
		v, err = dbmodels.IsNameExist(siteid, value)
		if v == true {
			user, err = dbmodels.GetUserByName(siteid, value)
		}
	case "phone":
		v, err = dbmodels.IsPhoneExist(siteid, value)
		if v == true {
			user, err = dbmodels.GetUserByPhone(siteid, value)
		}
	case "idcard":
		v, err = dbmodels.IsIdcardExist(siteid, value)
		if v == true {
			user, err = dbmodels.GetUserByIdcard(siteid, value)
		}
	case "email":
		v, err = dbmodels.IsEmailExist(siteid, value)
		if v == true {
			user, err = dbmodels.GetUserByEmail(siteid, value)
		}
	case "rfid":
		v, err = dbmodels.IsRfidExist(siteid, value)
		if v == true {
			user, err = dbmodels.GetUserByRfid(siteid, value)
		}
	case "idphone":
		v, err = dbmodels.IsIdPhoneExist(siteid, sId, sPhone)
		if v == true {
			user, err = dbmodels.GetUserByIdPhone(siteid, sId, sPhone)
		}
	case "weixinid":
		v, err = dbmodels.IsWeixinidExist(siteid, value)
		if v == true {
			user, err = dbmodels.GetUserByWeixinid(value)
		}

	default:
		r = map[string]interface{}{"Ret": 110310, "Msg": "类型:[" + stype + "] 不存在"}
	}

	if err != nil {
		r = map[string]interface{}{"Ret": 110320, "Msg": err.Error()}
	}
	if v == true {
		r = map[string]interface{}{"Ret": 0, "UserId": user.Id}
	} else {
		r = map[string]interface{}{"Ret": 110330, "Msg": "数据 [" + value + "]" + " 不存在"}
	}

	return r, nil
}

func AddUser(ctx context.Context, tracer trace.Tracer, appid, siteid int64, token string, data interface{}) (interface{}, error) {
	logger.Info("调用AddUser方法")
	if data == nil {
		return 0, errors.New("parameters is null.")
	}

	param := data.(map[string]interface{})
	logger.Finest("params=%v", param)

	if param["siteid"] != nil {
		siteid = utils.Convert2Int64(param["siteid"])
	}

	id, err := apis.PostUser(siteid, "add", param)

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

func ReadOrCreateUser(ctx context.Context, tracer trace.Tracer, appid, siteid int64, token string, data interface{}) (interface{}, error) {
	logger.Info("调用ReadOrCreateUser方法")
	if data == nil {
		return 0, errors.New("parameters is null.")
	}

	param := data.(map[string]interface{})
	logger.Finest("params=%v", param)

	if param["siteid"] != nil {
		siteid = utils.Convert2Int64(param["siteid"])
	}

	id, err := apis.PostUser(siteid, "readorcreate", param)

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

func UpdateUser(ctx context.Context, tracer trace.Tracer, appid, siteid int64, token string, data interface{}) (interface{}, error) {
	logger.Info("调用UpdateUser方法")
	if data == nil {
		return 0, errors.New("parameters is null.")
	}

	param := data.(map[string]interface{})
	logger.Finest("params=%v", param)

	var id int64
	if param["id"] != nil {
		siteid = utils.Convert2Int64(param["siteid"])
	} else {
		return 0, errors.New("没有携带id")
	}

	cnt, err := apis.PutUser(id, param)
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

func DeleteUser(ctx context.Context, tracer trace.Tracer, appid, siteid int64, token string, data interface{}) (interface{}, error) {
	logger.Info("调用DeleteUser方法")
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

	cnt, err := apis.DeleteUser(id, param)
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

func GetUserList(ctx context.Context, tracer trace.Tracer, appid, siteid int64, token string, data interface{}) (interface{}, error) {
	logger.Info("调用GetUserList方法")

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

	datas, cnt, err := apis.GetUserLists(siteid, appid, stypes, contents, "", "", 1000, 0)

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
