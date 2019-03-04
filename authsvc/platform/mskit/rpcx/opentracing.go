package rpcx

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/smallnest/rpcx/share"
	"platform/mskit/log"
	"platform/mskit/trace"
)


func RpcxClientOpenTracing(tracer trace.Tracer, options ...trace.TracerOption) ClientOption {
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
			ctx = context.WithValue(ctx,share.ReqMetaDataKey,*md)
			if span := opentracing.SpanFromContext(ctx); span != nil {
				// There's nothing we can do with an error here.
				if err := tracer.GetOpenTracer().Inject(span.Context(), opentracing.TextMap, *md); err != nil {
					config.Logger.Log("err", err)
				}
			}
			return ctx
		},
	)

	clientAfter := ClientAfter(
		func(ctx context.Context, _ map[string]string, _ map[string]string) context.Context {
			if span := opentracing.SpanFromContext(ctx); span != nil {
				span.Finish()
			}

			return ctx
		},
	)

	clientFinalizer := ClientFinalizer(
		func(ctx context.Context, err error) {
			if span := opentracing.SpanFromContext(ctx); span != nil {
				span.Finish()
			}
		},
	)

	return func(c *Client) {
		clientBefore(c)
		clientAfter(c)
		clientFinalizer(c)
	}

}
