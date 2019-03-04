package rpcx

import (
	"context"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"platform/mskit/log"
	"platform/mskit/trace"
)


func RpcxClientZipkinTrace(tracer trace.Tracer, options ...trace.TracerOption) ClientOption {
	config := trace.TracerOptions{
		Tags:      make(map[string]string),
		Name:      "",
		Logger:    log.Mslog,
		Propagate: true,
	}

	for _, option := range options {
		option(&config)
	}

	clientBefore := ClientBefore(
		func(ctx context.Context, md *map[string]string) context.Context {
			var (
				spanContext model.SpanContext
				name        string
			)

			if config.Name != "" {
				name = config.Name
			} else {
				name = ctx.Value(ContextKeyRequestMethod).(string)
			}

			if parent := zipkin.SpanFromContext(ctx); parent != nil {
				spanContext = parent.Context()
			}

			span := tracer.GetZipkinTracer().StartSpan(
				name,
				zipkin.Kind(model.Client),
				zipkin.Tags(config.Tags),
				zipkin.Parent(spanContext),
				zipkin.FlushOnFinish(false),
			)

			if config.Propagate {
				if err := trace.InjectRpcx(md)(span.Context()); err != nil {
					config.Logger.Log("err", err)
				}
			}

			return zipkin.NewContext(ctx, span)
		},
	)

	clientAfter := ClientAfter(
		func(ctx context.Context, _ map[string]string, _ map[string]string) context.Context {
			if span := zipkin.SpanFromContext(ctx); span != nil {
				span.Finish()
			}

			return ctx
		},
	)

	clientFinalizer := ClientFinalizer(
		func(ctx context.Context, err error) {
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
		},
	)

	return func(c *Client) {
		clientBefore(c)
		clientAfter(c)
		clientFinalizer(c)
	}

}
