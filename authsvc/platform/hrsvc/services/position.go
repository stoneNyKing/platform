package services

import (
	"context"
	"encoding/json"
	"errors"
	"platform/hrsvc/models"
	"platform/mskit/rest"
	oacerr "platform/pfcomm/errors"
	"strconv"
	"strings"
)

type PositionService struct {
	rest.RestApi
}

func (f *PositionService) Get(ctx context.Context,r *rest.Request) (interface{}, error) {

	logger.Finest("PositionService get function.")

	var maps map[string]interface{}
	maps = make(map[string]interface{})
	maps["ret"] = 0

	if !r.IsAuthorized {
		maps["ret"] = oacerr.ERR_INVALID_TOKEN
		maps["error"] = "无效的token"
		return maps, nil
	}

	action := r.Params.ByName("action")

	logger.Finest("action = %s,queries = %+v;version=%s", action, r.Queries, r.Version)

	id, err := strconv.ParseInt(action, 10, 64)

	if err != nil {
		id = 0
	}

	if id == 0 && action != "list" && action != "count" {
		err = errors.New("the path is not allowed.")
		maps["ret"] = oacerr.ERR_INVALID_PATH
		maps["error"] = err.Error()
		return maps, nil
	}

	var siteid int64
	siteids := r.GetInt64("site")

	if len(siteids) > 0 {
		siteid = siteids[0]
	}

	var appid int64
	appids := r.GetInt64("appid")

	if len(appids) > 0 {
		appid = appids[0]
	}

	var start, num int64
	var sort, order string
	var stypes []string
	var contents []string

	starts := r.GetInt64("pagestart")

	if len(starts) > 0 {
		start = starts[0]
	}

	nums := r.GetInt64("pagenum")
	if len(nums) > 0 {
		num = nums[0]
	}

	sorts := r.GetString("sort")

	if len(sorts) > 0 {
		sort = sorts[0]
	}

	orders := r.GetString("orderby")
	if len(orders) > 0 {
		order = orders[0]
	}

	s := r.GetString("stype")
	if len(s) > 0 {
		stypes = strings.Split(s[0], ",")
	}

	c := r.GetString("filter")
	cc := strings.ToLower(strings.Join(c, ","))
	if strings.Contains(cc, " and ") ||
		strings.Contains(cc, " or ") {
		err = errors.New("SQL条件有注入风险.")
		maps["ret"] = oacerr.ERR_INVALID_SQLFILTER
		maps["error"] = err.Error()
		return maps, nil
	}

	if len(c) > 0 {
		contents = strings.Split(c[0], ",")
	}

	var ver int

	version, err := strconv.ParseFloat(r.Version, 32)
	if err != nil {
		logger.Error("不能获取版本信息: %v", err.Error())
		ver = 0
	}
	ver = int(version * 10)

	logger.Fine("客户端版本信息: %v", ver)

	var datas interface{}
	var count int64

	switch action {
	case "list":
		datas, count, err = models.GetPositionLists(ctx,siteid, appid, stypes, contents, order, sort, num, start)
	case "count":
		count, err = models.GetPositionCount(ctx,siteid, appid, stypes, contents, order, sort)

	default:
		if id > 0 {
			datas, count, err = models.GetPosition(ctx,id, siteid, appid)
		}
	}

	if err != nil {
		maps["ret"] = 1
		maps["error"] = err.Error()
		return maps, nil
	}

	if action != "count" {
		maps["data"] = datas
	}
	maps["count"] = count

	return maps, nil
}

func (f *PositionService) Post(ctx context.Context,r *rest.Request) (interface{}, error) {

	logger.Finest("PositionService POST function.")

	var maps map[string]interface{}
	maps = make(map[string]interface{})
	maps["ret"] = 0

	if !r.IsAuthorized {
		maps["ret"] = oacerr.ERR_INVALID_TOKEN
		maps["error"] = "无效的token"
		return maps, nil
	}

	var siteid int64
	siteids := r.GetInt64("site")

	if len(siteids) > 0 {
		siteid = siteids[0]
	}

	var appid int64
	appids := r.GetInt64("appid")

	if len(appids) > 0 {
		appid = appids[0]
	}

	var token string
	tokens := r.GetString("token")

	if len(tokens) > 0 {
		token = tokens[0]
	}

	action := r.Params.ByName("action")
	id, err := strconv.ParseInt(action, 10, 64)

	if err != nil {
		id = 0
	}

	logger.Finest("action = %s,queries = %+v;version=%s,siteid=%d,appid=%d, id=%d", action, r.Queries, r.Version, siteid, appid, id)

	param := make(map[string]interface{})

	err = json.Unmarshal(r.Body, &param)

	if err != nil {
		logger.Error("将body转为json失败： %v", err.Error())
		maps["ret"] = 1
		maps["error"] = err.Error()
		return maps, nil
	}

	id, err = models.PostPosition(ctx,siteid, appid, token, param)

	if err != nil {
		maps["ret"] = 1
		maps["error"] = err.Error()
	}

	maps["id"] = id

	return maps, nil

}

func (f *PositionService) Put(ctx context.Context,r *rest.Request) (interface{}, error) {

	logger.Finest("PositionService Put function.")

	var maps map[string]interface{}
	maps = make(map[string]interface{})
	maps["ret"] = 0

	if !r.IsAuthorized {
		maps["ret"] = oacerr.ERR_INVALID_TOKEN
		maps["error"] = "无效的token"
		return maps, nil
	}

	var siteid int64
	siteids := r.GetInt64("site")

	if len(siteids) > 0 {
		siteid = siteids[0]
	}

	var appid int64
	appids := r.GetInt64("appid")

	if len(appids) > 0 {
		appid = appids[0]
	}
	var token string
	tokens := r.GetString("token")

	if len(tokens) > 0 {
		token = tokens[0]
	}

	action := r.Params.ByName("action")

	id, err := strconv.ParseInt(action, 10, 64)

	if err != nil {
		id = 0
		maps["ret"] = 1
		maps["error"] = err.Error()
		return maps, nil
	}

	logger.Finest("action = %s,queries = %+v;version=%s,siteid=%d,appid=%d,id=%d", action, r.Queries, r.Version, siteid, appid, id)

	param := make(map[string]interface{})

	err = json.Unmarshal(r.Body, &param)

	if err != nil {
		logger.Error("将body转为json失败： %v", err.Error())
		maps["ret"] = 1
		maps["error"] = err.Error()
		return maps, nil
	}

	cnt, err := models.PutPosition(ctx,siteid, appid, id, token, param)

	if err != nil {
		maps["ret"] = 1
		maps["error"] = err.Error()
	}

	maps["count"] = cnt

	return maps, nil

}

func (f *PositionService) Delete(ctx context.Context,r *rest.Request) (interface{}, error) {

	logger.Finest("PositionService Delete function.")

	var maps map[string]interface{}
	maps = make(map[string]interface{})
	maps["ret"] = 0

	if !r.IsAuthorized {
		maps["ret"] = oacerr.ERR_INVALID_TOKEN
		maps["error"] = "无效的token"
		return maps, nil
	}

	var siteid int64
	siteids := r.GetInt64("site")

	if len(siteids) > 0 {
		siteid = siteids[0]
	}

	var appid int64
	appids := r.GetInt64("appid")

	if len(appids) > 0 {
		appid = appids[0]
	}

	action := r.Params.ByName("action")

	id, err := strconv.ParseInt(action, 10, 64)

	if err != nil {
		id = 0
	}

	logger.Finest("action = %s,queries = %+v;version=%s,siteid=%d,appid=%d,id=%d", action, r.Queries, r.Version, siteid, appid, id)

	param := make(map[string]interface{})

	err = json.Unmarshal(r.Body, &param)

	if err != nil {
		logger.Error("将body转为json失败： %v", err.Error())
		maps["ret"] = 1
		maps["error"] = err.Error()
		return maps, nil
	}

	cnt, err := models.DeletePosition(ctx,id, param)

	if err != nil {
		maps["ret"] = 1
		maps["error"] = err.Error()
	}

	maps["count"] = cnt

	return maps, nil

}
