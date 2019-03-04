package main

import (
	"github.com/martini-contrib/render"
	"net/http"
	"platform/oasvc/config"
	"time"
)

type HttpService struct {
	handler http.Handler
}

func (hh *HttpService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	p := req.URL.Path

	if p == config.Config.Prefix+"/health" {
		resp := []byte("{\"ret\":0,\"msg\":\"ok\"}")
		w.Write(resp)
		return
	}
	hh.handler.ServeHTTP(w, req)
}

func ping(r render.Render, req *http.Request) {
	r.JSON(200, map[string]interface{}{"Ret": 0, "Time": time.Now().Unix()})
}
