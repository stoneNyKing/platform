package router

import (
	"platform/authsvc/imconf"
	"platform/authsvc/services"
	"platform/mskit/grace"
	"platform/mskit/rest"
	"platform/mskit/trace"

	"github.com/go-kit/kit/log"
	"os"
	"platform/pfcomm/apis"
)

func InitRoute(prefix string,msapp *grace.MicroService) {

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger,"ts", log.DefaultTimestampUTC)
		logger = log.With(logger,"caller", log.DefaultCaller)
	}

	optracer,_:= apis.CreateTracer(imconf.Config.RecordAddr,imconf.Config.ServiceName,logger,imconf.Config.Debug,
		imconf.Config.ZipkinUrl,imconf.Config.AppdashAddr,imconf.Config.LightstepToken,imconf.Config.KafkaAddress)

	var options []trace.TraceOption
	options = append(options,trace.WithTracerOption(true))
	options = append(options,trace.OpenTracerOption(optracer))
	tracer := trace.NewTracer(options...)

	msapp.SetTracer(tracer)

	license := services.AppLicense{}
	auth := services.AppAuth{}
	pkg := services.ApiPackage{}
	pkgservice := services.ApiPackageService{}
	service := services.ApiService{}


	mid := rest.RestMiddleware{Middle: LogMiddleware(logger), Object: logger}

	msapp.RegisterServiceWithTracer(prefix+"/license/:action", &license,tracer,logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/license", &license,tracer,logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/auth/:action", &auth,tracer,logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/auth", &auth,tracer,logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/package/:action", &pkg,tracer,logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/package", &pkg,tracer,logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/pkgservice/:action", &pkgservice,tracer,logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/pkgservice", &pkgservice,tracer,logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/service/:action", &service,tracer,logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/service", &service,tracer,logger, mid)


	healthsvc := services.HealthCheckService{}
	
	hmid := rest.RestMiddleware{Middle: NoTokenCheck(logger), Object: logger}
	msapp.RegisterRestService(prefix+"/health", &healthsvc, hmid)
}

