package services

import (
	"context"
	"platform/mskit/rest"
)

type HealthCheckService struct {
	rest.RestApi
}

func (f *HealthCheckService) Get(ctx context.Context,r *rest.Request) (interface{}, error) {

	//logger.Finest("HealthCheckService get function.")

	var maps map[string]interface{}
	maps = make(map[string]interface{})
	maps["ret"] = 0

	maps["msg"] = "ok"
	
	return maps, nil
}	