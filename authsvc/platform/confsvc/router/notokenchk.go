package router

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
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
