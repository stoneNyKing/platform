package services

import (
	"platform/filesvc/models"
	"strconv"
	"platform/mskit/rest"
	oacerr "platform/pfcomm/errors"
	"context"
	"net/http"
)

type DownloadService struct {
	rest.RestApi
}

func (f *DownloadService) Get(ctx context.Context,r *rest.Request) (interface{}, error) {

	logger.Finest("DownloadService POST function.")

	var maps map[string]interface{}
	maps = make(map[string]interface{})
	maps["result"] = 0

	if !r.IsAuthorized {
		maps["result"] = oacerr.ERR_INVALID_TOKEN
		maps["error"] = oacerr.CommonError(oacerr.ERR_INVALID_TOKEN)
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

	logger.Finest("action = %s,queries = %+v;version=%s,siteid=%d,appid=%d, id=%d", action, r.Queries, r.Version, siteid, appid, id)

	fn, err := models.DownloadFile(r)

	if err != nil {
		maps["result"] = 1
		maps["error"] = err.Error()
	}

	maps["url"] = fn

	return maps, nil

}

func (f *DownloadService) EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {

	if response == nil {
		response = ""
	}

	w.Header().Set("Allow", "HEAD,GET,PUT,DELETE,OPTIONS,POST")

	f.Finish(w)

	p := response.(map[string]interface{})
	http.Redirect(w, f.Request, p["url"].(string), 301)

	//err := json.NewEncoder(w).Encode(response)

	return nil
}
