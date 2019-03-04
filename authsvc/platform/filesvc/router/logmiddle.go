package router

import (
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
	"platform/common/utils"
	"platform/filesvc/imconf"
	"platform/mskit/rest"
	"strconv"
	"time"
)


func LogMiddleware(logger log.Logger) rest.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			if request == nil {
				return nil,errors.New("no request avaliable.")
			}

			req := request.(rest.Request)


			//log := logger
			defer func(begin time.Time) {
				logger.Log(
					"method", req.Method,
					"took", time.Since(begin),
				)
			}(time.Now())

			var token string
			tokens := req.GetString("token")
	
			if len(tokens)>0 {
				token = tokens[0]
			}
			var appid int64
			appids := req.GetInt64("appid")
			
			if len(appids)>0 {
				appid = appids[0]
			}

			sa := strconv.FormatInt(appid,10)
			bf := befe[appid]
			var uri string
			if bf == BEFE_USER {
				uri = imconf.Config.OusvcUrl+"/token/check?appid=" + sa +"&token="+token
			}else if bf == BEFE_ADMIN {
				uri = imconf.Config.OasvcUrl+"/admin/token/check?appid=" + sa +"&token="+token
			}else{
				req.SetAuthorized(false)
			}
			
			if imconf.Config.IsAuth {
				req.SetAuthorized( utils.CheckToken(uri,int(appid),token) )
			}else{
				req.SetAuthorized(true)
			}
			
			//log.Finest("appid=%d,auth=%v,token=%s",appid,req.IsAuthorized,token)

			return next(ctx, req)
		}
	}
}
