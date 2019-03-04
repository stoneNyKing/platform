package router

import (
	"platform/confsvc/imconf"
	"platform/confsvc/services"
	"platform/mskit/grace"
	"platform/mskit/log"
	"platform/mskit/rest"
	"platform/mskit/trace"
	"platform/pfcomm/apis"

)

func InitRoute(prefix string,msapp *grace.MicroService) {

	// Logging domain.
	var logger =log.Mslog

	optracer,_:= apis.CreateTracer(imconf.Config.RecordAddr,imconf.Config.ServiceName,logger,imconf.Config.Debug,
		imconf.Config.ZipkinUrl,imconf.Config.AppdashAddr,imconf.Config.LightstepToken,imconf.Config.KafkaAddress)

	var options []trace.TraceOption
	options = append(options,trace.WithTracerOption(true))
	options = append(options,trace.OpenTracerOption(optracer))
	tracer := trace.NewTracer(options...)

	msapp.SetTracer(tracer)

	sitednssvc := services.SiteDNSService{}
	domainsvc := services.DomainConfService{}
	appidsvc := services.AppidService{}

	mid := rest.RestMiddleware{Middle: LogMiddleware(logger), Object: logger}
	hmid := rest.RestMiddleware{Middle: NoTokenCheck(logger), Object: logger}

	msapp.RegisterServiceWithTracer(prefix+"/dconf/:action", &domainsvc, tracer,logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/dconf", &domainsvc, tracer,logger, mid)

	msapp.RegisterServiceWithTracer(prefix+"/appid/:action", &appidsvc, tracer,logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/appid", &appidsvc, tracer,logger, mid)

	msapp.RegisterServiceWithTracer(prefix+"/site", &sitednssvc,tracer,logger, hmid)
	
	healthsvc := services.HealthCheckService{}

	msapp.RegisterRestService(prefix+"/health", &healthsvc,  hmid)
}


