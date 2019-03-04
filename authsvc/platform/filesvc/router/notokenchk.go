package router

import (
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
	"platform/mskit/rest"
)



func NoTokenCheck(logger log.Logger) rest.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			if request == nil {
				return nil,errors.New("no request avaliable.")
			}
		
			req := request.(rest.Request)

			return next(ctx, req)
		}
	}
}
