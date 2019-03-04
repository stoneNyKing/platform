package services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"platform/common/utils"
	"platform/filesvc/models"
	"platform/mskit/rest"
	oacerr "platform/pfcomm/errors"
	"strconv"
	"strings"
)

type FileService struct {
	rest.RestApi
	Id 		int64
}

func (f *FileService) Get(ctx context.Context,r *rest.Request) (interface{}, error) {

	logger.Finest("FileService get function.")

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

	if id == 0 && action != "list" && action != "count"  {
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

	logger.Finest("客户端版本信息: %v", ver)

	var datas interface{}
	var count int64

	switch action {
	case "list":
		datas, count, err = models.GetFileLists(siteid, appid, stypes, contents, order, sort, num, start)
	case "count":
		count, err = models.GetFileCount(siteid, appid, stypes, contents, order, sort)
	default:
		if id > 0 {
			f.Id = id
			datas, count, err = models.GetFile(id, siteid, appid)
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

	var vs interface{}

	if id >0 {
		vs = datas
	}else {
		vs = maps
	}

	return vs, nil
}

func (f *FileService) Post(ctx context.Context,r *rest.Request) (interface{}, error) {

	logger.Finest("FileService POST function.")

	var maps map[string]interface{}
	maps = make(map[string]interface{})


	if !r.IsAuthorized {
		maps["result"] = oacerr.ERR_INVALID_TOKEN
		maps["error"] = oacerr.CommonError(oacerr.ERR_INVALID_TOKEN)
		return maps,nil
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

	var category int64
	categorys := r.GetInt64("category")

	if len(categorys) > 0 {
		category = categorys[0]
	}

	logger.Finest("action = %s,queries = %+v;version=%s,siteid=%d,appid=%d, id=%d", action, r.Queries, r.Version, siteid, appid, id)
	var ver int

	version := utils.Convert2Float32(r.Version)
	ver = int(version * 10)

	logger.Finest("客户端版本信息: %v", ver)


	fn, err := models.UploadFile(r,int(category))

	if err != nil {
		maps["result"] = 1
		maps["error"] = err.Error()
	}

	if ver < 20 {
		maps["result"] = 0
		maps["url"] = fn
	}else{
		maps["ret"] = 0
		maps["data"] = fn
	}

	return maps, nil

}


func (f *FileService) Put(ctx context.Context,r *rest.Request) (interface{}, error) {

	logger.Finest("DictChargeInsService Put function.")

	var maps map[string]interface{}
	maps = make(map[string]interface{})
	maps["ret"] = 0

	if !r.IsAuthorized {
		maps["ret"] = oacerr.ERR_INVALID_TOKEN
		maps["error"] = "无效的token"
		return maps,nil
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

	var err error
	id := utils.Convert2Int64(action)

	var token string
	tokens := r.GetString("token")
	if len(tokens) > 0 {
		token = tokens[0]
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
	var cnt int64

	cnt, err = models.PutFileInfo(id,appid,siteid,token, param)

	if err != nil {
		maps["ret"] = 1
		maps["error"] = err.Error()
	}

	maps["count"] = cnt

	return maps, nil

}

func (f *FileService) Finish(w http.ResponseWriter) error {

	if w == nil {
		return errors.New("writer is nil ")
	}
	if f.Id > 0 {
		w.Header().Set("Content-Type", "image/jpg")
	}
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Origin,Accept,Content-Range,Content-Description,Content-Disposition,X-Requested-With")
	w.Header().Add("Access-Control-Allow-Methods", "PUT,GET,POST,DELETE,OPTIONS")

	return nil
}
