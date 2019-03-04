package apis

import (
	"github.com/go-kit/kit/log"
	"github.com/lightstep/lightstep-tracer-go"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	stdzipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go/reporter"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/openzipkin/zipkin-go/reporter/kafka"
	appdashot "github.com/sourcegraph/appdash/opentracing"
	"os"
	"sourcegraph.com/sourcegraph/appdash"
	"strings"
)

func CreateTracer(recordAddr string, serviceName string, logger log.Logger, Debug bool, Zipkin_address, AppdashAddr, LightstepToken, kafkaaddr string) (stdopentracing.Tracer, *zipkin.Tracer) {
	// Tracing domain.
	var zipkinTracer *zipkin.Tracer
	var tracer stdopentracing.Tracer
	tracer = stdopentracing.GlobalTracer()

	{
		if Zipkin_address != "" || kafkaaddr != "" {
			logger := log.With(logger, "tracer", "Zipkin")

			zipkinTracer = getZipkinTracer(recordAddr,Zipkin_address,kafkaaddr,serviceName,logger)

			tracer = getKafkaOpenTracer(recordAddr,kafkaaddr,serviceName,logger)

		} else if AppdashAddr != "" {
			logger := log.With(logger, "tracer", "Appdash")
			logger.Log("addr", AppdashAddr)
			tracer = appdashot.NewTracer(appdash.NewRemoteCollector(AppdashAddr))
		} else if LightstepToken != "" {
			logger := log.With(logger, "tracer", "LightStep")
			logger.Log() // probably don't want to print out the token :)
			tracer = lightstep.NewTracer(lightstep.Options{
				AccessToken: LightstepToken,
			})
		} else {
			logger := log.With(logger, "tracer", "none")
			logger.Log()
			tracer = stdopentracing.GlobalTracer() // no-op
		}
	}

	return tracer, zipkinTracer
}

func getKafkaOpenTracer(recordAddr, kafkaaddr,serviceName string,logger log.Logger) stdopentracing.Tracer {

	if kafkaaddr == "" {
		return nil
	}

	collector, err := stdzipkin.NewKafkaCollector(
		strings.Split(kafkaaddr, ","),
		stdzipkin.KafkaLogger(logger),
	)
	if err != nil {
		logger.Log("err", err)
		return nil
	}
	tracer, err := stdzipkin.NewTracer(
		stdzipkin.NewRecorder(collector, false, recordAddr, serviceName),
	)
	if err != nil {
		logger.Log("err", err)
		return nil
	}

	stdopentracing.SetGlobalTracer(tracer)

	return tracer
}

func getZipkinTracer(recordAddr, Zipkin_address,kafkaaddr,serviceName string,logger log.Logger) *zipkin.Tracer{

	var zipkinTracer *zipkin.Tracer
	var reporter reporter.Reporter
	var err error

	useNoopTracer := (Zipkin_address == "" && kafkaaddr == "")
	reporter = zipkinhttp.NewReporter(Zipkin_address)
	collector := "Native"
	colladdr := Zipkin_address

	if kafkaaddr != "" {
		reporter, err = kafka.NewReporter(
			strings.Split(kafkaaddr, ","),
		)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		collector = "Kafka"
		colladdr = kafkaaddr
	}
	zEP, _ := zipkin.NewEndpoint(serviceName, recordAddr)
	zipkinTracer, err = zipkin.NewTracer(
		reporter, zipkin.WithLocalEndpoint(zEP), zipkin.WithSharedSpans(true), zipkin.WithNoopTracer(useNoopTracer),
	)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}
	if !useNoopTracer {
		logger.Log("tracer", "Zipkin", "type", collector, "URL", colladdr)
	}

	return zipkinTracer
}