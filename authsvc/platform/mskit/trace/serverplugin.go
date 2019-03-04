package trace

import (
	"context"
	"fmt"
	"github.com/smallnest/rpcx/share"
	"platform/mskit/log"
	zipkin "github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/smallnest/rpcx/protocol"
)

type ZipkinTracePlugin struct {
	tracer 		*zipkin.Tracer
}

func NewZipkinTracePlugin(ziptracer Tracer) *ZipkinTracePlugin {
	z := &ZipkinTracePlugin{
		tracer:ziptracer.GetZipkinTracer(),
	}

	return z
}

func (p *ZipkinTracePlugin) PostReadRequest(ctx context.Context, r *protocol.Message, e error) error {
	config := TracerOptions{
		Tags:      make(map[string]string),
		Name:      "test",
		Logger:    log.Mslog,
		Propagate: true,
	}

	var (
		spanContext model.SpanContext
		name        string
		tags        = make(map[string]string)
	)

	m := ctx.Value(share.ReqMetaDataKey)

	if m==nil {
		//config.logger.Log("err", "unable to retrieve method name: missing rpcx interceptor hook")
	} else {
		fmt.Printf("metadata=%+v\n",m.(map[string]string))
	}

	//fmt.Printf("r.Metadata=%+v\n",r.Metadata)

	if config.Name != "" {
		name = config.Name
	} else if m !=nil {
		//name = rpcMethod
		r.Metadata = m.(map[string]string)
	}

	if config.Propagate {
		spanContext = p.tracer.Extract(ExtractRpcx(&r.Metadata))
		if spanContext.Err != nil {
			config.Logger.Log("err", spanContext.Err)
		}
	}

	span := p.tracer.StartSpan(
		name,
		zipkin.Kind(model.Server),
		zipkin.Tags(config.Tags),
		zipkin.Tags(tags),
		zipkin.Parent(spanContext),
		zipkin.FlushOnFinish(false),
	)
	if span == nil {
		fmt.Printf("不能开始span\n")
	}
	ctx = zipkin.NewContext(ctx, span)
	return nil
}

func (p *ZipkinTracePlugin) PostWriteResponse(ctx context.Context, req *protocol.Message, res *protocol.Message, err error) error {

	if span := zipkin.SpanFromContext(ctx); span != nil {
		if err != nil {
			zipkin.TagError.Set(span, err.Error())
		}
		// calling span.Finish() a second time is a noop, if we didn't get to
		// ClientAfter we can at least time the early bail out by calling it
		// here.
		span.Finish()
		// send span to the Reporter
		span.Flush()
	}

	return nil
}