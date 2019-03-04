package rpc

import (
	"context"
	"errors"
	"platform/common/utils"
	"platform/hrsvc/models"
	"platform/mskit/trace"
)

func AddStaff(ctx context.Context,tracer trace.Tracer,appid, siteid int64, token string, data interface{}) (interface{}, error) {
	logger.Info("调用AddStaff方法")
	if data == nil {
		return 0, errors.New("parameters is null.")
	}

	param := data.(map[string]interface{})
	logger.Finest("params=%v", param)

	if param["siteid"] != nil {
		siteid = utils.Convert2Int64(param["siteid"])
	}

	id, err := models.PostStaff(ctx,tracer,siteid, appid, token, param)
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

func UpdateStaff(ctx context.Context,tracer trace.Tracer,appid, siteid int64, token string, data interface{}) (interface{}, error) {
	logger.Info("调用UpdateStaff方法")
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

	cnt, err := models.PutStaff(ctx,tracer,siteid, appid, id, token, param)
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

func AddDepartment(ctx context.Context,tracer trace.Tracer,appid, siteid int64, token string, data interface{}) (interface{}, error) {
	logger.Info("调用AddDepartment方法")
	if data == nil {
		return 0, errors.New("parameters is null.")
	}

	param := data.(map[string]interface{})
	logger.Finest("params=%v", param)

	if param["siteid"] != nil {
		siteid = utils.Convert2Int64(param["siteid"])
	}

	id, err := models.PostDept(ctx,siteid, appid, token, param)
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

func UpdateDepartment(ctx context.Context,tracer trace.Tracer,appid, siteid int64, token string, data interface{}) (interface{}, error) {
	logger.Info("调用UpdateDepartment方法")
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

	cnt, err := models.PutDept(ctx,siteid, appid, id, token, param)
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

func AddOrg(ctx context.Context,tracer trace.Tracer,appid, siteid int64, token string, data interface{}) (interface{}, error) {
	logger.Info("调用AddOrg方法")
	if data == nil {
		return 0, errors.New("parameters is null.")
	}

	param := data.(map[string]interface{})
	logger.Finest("params=%v", param)

	if param["siteid"] != nil {
		siteid = utils.Convert2Int64(param["siteid"])
	}

	id, err := models.PostOrganization(ctx,siteid, appid, token, param)
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

func UpdateOrg(ctx context.Context,tracer trace.Tracer,appid, siteid int64, token string, data interface{}) (interface{}, error) {
	logger.Info("调用UpdateOrg方法")
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

	cnt, err := models.PutOrganization(ctx,siteid, appid, id, token, param)
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
