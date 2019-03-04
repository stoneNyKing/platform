package trace

import (
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
)

type Tracer interface {
	GetOpenTracer()(opentracing.Tracer)
	GetZipkinTracer()(*zipkin.Tracer)
}

type trace struct {
	withTracer 			bool
	withZipkinTracer 	bool
	tracer 				opentracing.Tracer
	zipkinTracer		*zipkin.Tracer
}

type TraceOption func(*trace)

func WithTracerOption(istracer bool) TraceOption {
	return func(t *trace){t.withTracer = istracer}
}
func WithZipkinTracerOption(istracer bool) TraceOption {
	return func(t *trace){t.withZipkinTracer = istracer}
}

func OpenTracerOption(tracer opentracing.Tracer) TraceOption {
	return func(t *trace){t.tracer = tracer}
}
func ZipkinTracerOption(tracer *zipkin.Tracer) TraceOption {
	return func(t *trace){t.zipkinTracer = tracer}
}
func(t *trace)GetOpenTracer() opentracing.Tracer{
	return t.tracer
}
func(t *trace)GetZipkinTracer() *zipkin.Tracer{
	return t.zipkinTracer
}

func NewTracer(options... TraceOption) Tracer {
	t := &trace{}
	for _,option := range options{
		option(t)
	}

	return t
}