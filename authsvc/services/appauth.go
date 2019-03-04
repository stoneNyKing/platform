package services

import (
	"context"
	"platform/authsvc/models"
	"platform/mskit/rest"
	oacerr "platform/pfcomm/errors"
	"strconv"
	"errors"
)

type  AppAuth struct{
	rest.RestApi
}

func (f *AppAuth) Get (ctx context.Context,r *rest.Request) (interface{}, error){
	logger.Finest("AppAuth get function");

	// 返回值默认为0，鉴权成功
	var maps map[string]interface{}
	maps = make(map[string] interface{})
	maps["ret"] = 0

	if !r.IsAuthorized{
		maps["ret"] = oacerr.ERR_INVALID_TOKEN;
		maps["error"] = "无效的token"
		return maps, nil;
	}

	action := r.Params.ByName("action")
	logger.Finest("action = %s,queries = %+v;version=%s", action, r.Queries, r.Version)

	id, err := strconv.ParseInt(action, 10, 64)
	if err != nil {
		id = 0
	}

	var userid int64
	userids := r.GetInt64("uid")

	if len(userids) > 0 {
		userid = userids[0]
	}

	if id == 0 {
		err = errors.New("the path is not allowed.")
		maps["ret"] = 1
		maps["error"] = err.Error()
		return maps, nil
	}

	ret,p, err := models.GetAuth(id,userid)
	maps["ret"] = ret

	if err != nil {
		maps["error"] = err.Error()
		maps["ret"]  = 1
		return maps,err
	}

	if p != nil {
		maps["produces"] = p
	}

	return maps, nil
}