package services

import (
	"context"
	"encoding/json"
	"errors"
	"platform/authsvc/models"
	"platform/mskit/rest"
	oacerr "platform/pfcomm/errors"
	"strconv"
	"strings"
)

type  ApiPackage struct{
	rest.RestApi
}

func (f *ApiPackage) Get (ctx context.Context,r *rest.Request) (interface{}, error) {
	logger.Finest("ApiPackage get function");

	// 返回值默认为1，鉴权失败
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

	if id == 0 && action != "list" && action != "count" {
		err = errors.New("the path is not allowed.")
		maps["ret"] = oacerr.ERR_INVALID_PATH
		maps["error"] = err.Error()
		return maps, nil
	}

	// 获取site(orgid) 以及加密需要的apikey和orgcode
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
	cc := strings.ToLower(strings.Join(c,","))
	if strings.Contains(cc," and ") ||
		strings.Contains(cc," or ") {
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

	logger.Finest("客户端版本信息: %v", ver)

	var data interface{}
	var count int64

	switch action {
	case "list":
		data, count, err = models.GetApiPackageLists(siteid, appid, stypes,contents,order, sort, num, start)
	case "count":
		count, err = models.GetApiPackageCount(siteid, appid, stypes, contents, order, sort)
	default:
		if id > 0{
			data, count, err = models.GetApiPackage(id)
		}
	}

	if err != nil{
		maps["ret"] = 1
		maps["error"] = err.Error()
		return maps, nil
	}
	maps["count"] = count
	if action != "count" {
		maps["data"] = data
	}

	return maps, nil
}

func (f *ApiPackage) Post (ctx context.Context,r *rest.Request) (interface{}, error) {
	logger.Finest("ApiPackage Post function");

	var maps map[string] interface{}
	maps = make(map[string] interface{})
	maps["ret"] = 0

	if !r.IsAuthorized {
		maps["ret"] = oacerr.ERR_INVALID_TOKEN
		maps["error"] = "无效的token"
		return  maps, nil
	}

	param := make(map[string]interface{})

	err := json.Unmarshal(r.Body, &param)

	if err != nil {
		logger.Error("将body转为json失败： %v", err.Error())
		maps["ret"] = 1
		maps["error"] = err.Error()
		return maps, nil
	}

	id, err := models.PostApiPackage(param)

	if err != nil{
		maps["ret"] = 1
		maps["error"] = err.Error()
		return maps, nil
	}
	maps["id"] = id

	return maps, nil
}

func (f *ApiPackage) Put (ctx context.Context,r *rest.Request) (interface{}, error) {
	logger.Finest("ApiPackage Put function");

	var maps map[string] interface{}
	maps = make(map[string] interface{})
	maps["ret"] = 0

	if !r.IsAuthorized {
		maps["ret"] = oacerr.ERR_INVALID_TOKEN
		maps["error"] = "无效的token"
		return  maps, nil
	}

	action := r.Params.ByName("action")
	id, err := strconv.ParseInt(action, 10, 64)
	if err != nil {
		maps["ret"] = 1
		maps["error"] = err.Error()
		return maps, nil
	}

	param := make(map[string]interface{})

	err = json.Unmarshal(r.Body, &param)

	if err != nil {
		logger.Error("将body转为json失败： %v", err.Error())
		maps["ret"] = 1
		maps["error"] = err.Error()
		return maps, nil
	}

	cnt, err := models.PutApiPackage(id, param)

	if err != nil{
		maps["ret"] = 1
		maps["error"] = err.Error()
		return maps, nil
	}
	maps["count"] = cnt

	return maps, nil
}

func (f *ApiPackage) Delete (ctx context.Context,r *rest.Request) (interface{}, error) {
	logger.Finest("ApiPackage Delete function");

	var maps map[string] interface{}
	maps = make(map[string] interface{})
	maps["ret"] = 0

	if !r.IsAuthorized {
		maps["ret"] = oacerr.ERR_INVALID_TOKEN
		maps["error"] = "无效的token"
		return  maps, nil
	}

	param := make(map[string]interface{})

	err := json.Unmarshal(r.Body, &param)

	if err != nil {
		logger.Error("将body转为json失败： %v", err.Error())
		maps["ret"] = 1
		maps["error"] = err.Error()
		return maps, nil
	}

	num, err := models.DeleteApiPackage(param)

	if err != nil{
		maps["ret"] = 1
		maps["error"] = err.Error()
		return maps, nil
	}
	maps["count"] = num

	return maps, nil
}
