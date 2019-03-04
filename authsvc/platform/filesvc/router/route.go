package router

import (
	"github.com/go-kit/kit/log"
	"os"
	"platform/filesvc/imconf"
	"platform/filesvc/services"
	"platform/mskit/grace"
	"platform/mskit/rest"
	"platform/mskit/trace"

	"net/http"
	"platform/filesvc/comm"
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

	svc := services.FileService{}
	fsvc := services.FileConfService{}
	dsvc := services.DownloadService{}

	mid := rest.RestMiddleware{Middle: LogMiddleware(logger), Object: logger}

	var downpath,filepath string
	if comm.FilePath != "" {
		filepath = comm.FilePath
	}else{
		filepath = "/var/www/attach"
	}

	if comm.Prefix != "" {
		downpath = 	prefix+"/"+ comm.Prefix +"/*filepath"
	}

	msapp.Router.ServeFiles(downpath,http.Dir(filepath))

	logger.Log("prefix",prefix)

	msapp.RegisterServiceWithTracer(prefix+"/upload", &svc,tracer,logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/file/:action", &svc,tracer,logger, mid)

	msapp.RegisterServiceWithTracer(prefix+"/fconf/:action", &fsvc,tracer,logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/fconf", &fsvc,tracer,logger, mid)

	msapp.RegisterServiceWithTracer(prefix+"/oss/:action", &dsvc,tracer,logger, mid)

	healthsvc := services.HealthCheckService{}

	hmid := rest.RestMiddleware{Middle: NoTokenCheck(logger), Object: logger}
	msapp.RegisterRestService(prefix+"/health", &healthsvc, hmid)

}


