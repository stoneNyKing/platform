package rest

import "github.com/go-kit/kit/endpoint"

type Middleware func(endpoint.Endpoint) endpoint.Endpoint

type RestMiddleware struct {
	Middle Middleware
	Object interface{}
}

func (rm *RestMiddleware) GetMiddleware() func(interface{}) Middleware {
	return func(inter interface{}) Middleware {
		return rm.Middle
	}
}
