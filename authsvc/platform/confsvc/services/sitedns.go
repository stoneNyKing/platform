package services

import (
	"context"
	"errors"
	"platform/confsvc/models"
	"platform/mskit/rest"
	oacerr "platform/pfcomm/errors"
	"strconv"
	"strings"
)

type SiteDNSService struct {
	rest.RestApi
}

func (f *SiteDNSService) Get(ctx context.Context,r *rest.Request) (interface{}, error) {

	logger.Finest("SiteDNSService get function.")

	var maps map[string]interface{}
	maps = make(map[string]interface{})
	maps["ret"] = 0


	action := r.Params.ByName("action")

	logger.Finest("action = %s,queries = %+v;version=%s", action, r.Queries, r.Version)

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

	var stypes []string
	var contents []string
	var sitedns string

	s := r.GetString("stype")
	if len(s) > 0 {
		stypes = strings.Split(s[0], ",")
	}

	c := r.GetString("filter")
	cc := strings.ToLower(strings.Join(c, ","))
	if strings.Contains(cc, " and ") ||
		strings.Contains(cc, " or ") {
		err := errors.New("SQL条件有注入风险.")
		maps["ret"] = oacerr.ERR_INVALID_SQLFILTER
		maps["error"] = err.Error()
		return maps, nil
	}

	if len(c) > 0 {
		contents = strings.Split(c[0], ",")
	}
	d := r.GetString("dns")
	if len(d) > 0 {
		sitedns = d[0]
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

	datas, count, err = models.GetSiteDNSConf(siteid, appid, stypes,contents,sitedns)

	if err != nil {
		maps["ret"] = 1
		maps["error"] = err.Error()
		return maps, nil
	}
	maps["data"] = datas
	maps["count"] = count

	return maps, nil
}
