package rpc

import (
	"context"
	"errors"
	"platform/common/utils"
	"platform/authsvc/models"
	"platform/mskit/trace"
)

func CheckAuth(ctx context.Context,tracer trace.Tracer,appid, siteid int64,token string, data interface{})(interface{},error) {

	logger.Info("调用CheckAuth方法")
	if data == nil {
		return 1,errors.New("parameters is null.")
	}

	param := data.(map[string]interface{})
	logger.Finest("params=%v",param)

	var orgid int64
	if param["id"] != nil {
		orgid = utils.Convert2Int64(param["id"])
	}else{
		orgid = siteid
	}

	var userid int64
	if param["userid"] != nil {
		userid = utils.Convert2Int64(param["userid"])
	}

	_,p, err := models.GetAuth(orgid,userid)

	ret := make(map[string]interface{})
	ret["ret"] = 0

	if err != nil {
		ret["error"] = err.Error()
		ret["ret"] = 1
	}
	if p != nil {
		ret["produces"] = p
	}

	return ret,nil
}


func AddLicense(ctx context.Context,tracer trace.Tracer,appid, siteid int64,token string, data interface{})(interface{},error) {

	logger.Info("调用AddLicense方法")
	if data == nil {
		return 1,errors.New("parameters is null.")
	}

	param := data.(map[string]interface{})
	logger.Finest("params=%v",param)

	ret := make(map[string] interface{})
	ret["ret"] = 0

	id, err := models.PostAppLicense(siteid,param)

	if err != nil {
		ret["error"] = err.Error()
		ret["ret"] = 1
	}
	ret["id"] = id

	return ret,err
}

func DeleteLicense(ctx context.Context,tracer trace.Tracer,appid, siteid int64,token string, data interface{})(interface{},error) {

	logger.Info("调用DeleteLicense方法")
	if data == nil {
		return 1,errors.New("parameters is null.")
	}

	param := data.(map[string]interface{})
	logger.Finest("params=%v",param)

	ret := make(map[string]interface{})
	ret["ret"] = 0

	num, err := models.DeleteAppLicense(param)

	if err != nil {
		ret["error"] = err.Error()
		ret["ret"] = 1
	}
	ret["count"] = num

	return ret,err
}

func UpdateLicense(ctx context.Context,tracer trace.Tracer,appid, siteid int64,token string, data interface{})(interface{},error) {

	logger.Info("调用UpdateLicense方法")
	if data == nil {
		return 1,errors.New("parameters is null.")
	}


	param := data.(map[string]interface{})
	logger.Finest("params=%v",param)

	ret := make(map[string]interface{})
	ret["ret"] = 0

	var orgid int64
	if param["id"] != nil {
		orgid = utils.Convert2Int64(param["id"])
	}

	cnt, err := models.PutAppLicense(orgid, param)

	if err != nil {
		ret["error"] = err.Error()
		ret["ret"] = 1
	}
	ret["count"] = cnt

	return ret,err
}

func GetLicense(ctx context.Context,tracer trace.Tracer,appid, siteid int64,token string, data interface{})(interface{},error) {

	logger.Info("调用GetLicense方法")
	if data == nil {
		return 0,errors.New("parameters is null.")
	}

	param := data.(map[string]interface{})
	logger.Finest("params=%v",param)

	var orgid int64
	if param["orgid"] != nil {
		orgid = utils.Convert2Int64(param["orgid"])
	}else{
		orgid = siteid
	}

	var userid int64
	if param["userid"] != nil {
		userid = utils.Convert2Int64(param["userid"])
	}
	// 查询数据
	license, count, err := models.GetAppLicense(orgid,userid)

	ret := make(map[string]interface{})
	ret["ret"] = 0


	if err != nil {
		ret["error"] = err.Error()
		ret["ret"] = 1
	}
	ret["count"] = count
	ret["data"] = license
	return ret,err
}


func GetLicenseCounts(ctx context.Context,tracer trace.Tracer,appid, siteid int64,token string, data interface{})(interface{},error) {

	logger.Info("调用GetLicenseCounts方法")
	if data == nil {
		return 0,errors.New("parameters is null.")
	}

	param := data.(map[string]interface{})
	logger.Finest("params=%v",param)

	var orgid int64
	if param["orgid"] != nil {
		orgid = utils.Convert2Int64(param["orgid"])
	}else{
		orgid = siteid
	}

	var stypes,contents []string

	if param["stypes"] != nil {
		s := param["stypes"].([]interface{})
		for _,sv := range s {
			stypes = append(stypes,utils.ConvertToString(sv))
		}
	}
	if param["contents"] != nil {
		s := param["contents"].([]interface{})
		for _,sv := range s {
			contents = append(contents,utils.ConvertToString(sv))
		}
	}

	// 查询数据
	license, count, err := models.GetApiLicenseCounts(orgid,appid,stypes,contents,"","",0,0)

	ret := make(map[string]interface{})
	ret["ret"] = 0


	if err != nil {
		ret["error"] = err.Error()
		ret["ret"] = 1
	}
	ret["count"] = count
	ret["data"] = license
	return ret,err
}


func GetPkgServices(ctx context.Context,tracer trace.Tracer,appid, siteid int64,token string, data interface{})(interface{},error) {

	logger.Info("调用GetPkgServices方法")
	if data == nil {
		return 0,errors.New("parameters is null.")
	}

	param := data.(map[string]interface{})
	logger.Finest("params=%v",param)

	var stypes,contents []string
	if param["stypes"] != nil {
		p1 := param["stypes"].([]interface{})

		for _,v := range p1 {
			stypes = append(stypes,utils.ConvertToString(v))
		}
	}
	if param["contents"] != nil {
		p1 := param["contents"].([]interface{})

		for _,v := range p1 {
			contents = append(contents,utils.ConvertToString(v))
		}
	}
	// 查询数据
	ps, count, err := models.GetApiPkgServiceLists(siteid,appid,stypes,contents,"","",0,0)

	ret := make(map[string]interface{})
	ret["ret"] = 0

	if err != nil {
		ret["error"] = err.Error()
		ret["ret"] = 1
	}
	ret["count"] = count
	ret["data"] = ps
	return ret,err
}