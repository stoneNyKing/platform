package trace

import (
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation"
	. "github.com/openzipkin/zipkin-go/propagation/b3"
)

// ExtractRpcx will extract a span.Context from the Rpcx Request metadata if
// found in B3 header format.
func ExtractRpcx(md *map[string]string) propagation.Extractor {
	return func() (*model.SpanContext, error) {
		var (
			traceIDHeader      = GetRpcxHeader(md, TraceID)
			spanIDHeader       = GetRpcxHeader(md, SpanID)
			parentSpanIDHeader = GetRpcxHeader(md, ParentSpanID)
			sampledHeader      = GetRpcxHeader(md, Sampled)
			flagsHeader        = GetRpcxHeader(md, Flags)
		)

		return ParseHeaders(
			traceIDHeader, spanIDHeader, parentSpanIDHeader, sampledHeader,
			flagsHeader,
		)
	}
}

// InjectRpcx will inject a span.Context into Rpcx metadata.
func InjectRpcx(md *map[string]string) propagation.Injector {
	return func(sc model.SpanContext) error {
		if (model.SpanContext{}) == sc {
			return ErrEmptyContext
		}

		if sc.Debug {
			setRpcxHeader(md, Flags, "1")
		} else if sc.Sampled != nil {
			// Debug is encoded as X-B3-Flags: 1. Since Debug implies Sampled,
			// we don't send "X-B3-Sampled" if Debug is set.
			if *sc.Sampled {
				setRpcxHeader(md, Sampled, "1")
			} else {
				setRpcxHeader(md, Sampled, "0")
			}
		}

		if !sc.TraceID.Empty() && sc.ID > 0 {
			// set identifiers
			setRpcxHeader(md, TraceID, sc.TraceID.String())
			setRpcxHeader(md, SpanID, sc.ID.String())
			if sc.ParentID != nil {
				setRpcxHeader(md, ParentSpanID, sc.ParentID.String())
			}
		}

		return nil
	}
}

// GetGRPCHeader retrieves the last value found for a particular key. If key is
// not found it returns an empty string.
func GetRpcxHeader(md *map[string]string, key string) string {
	v := (*md)[key]
	return v
}

func setRpcxHeader(md *map[string]string, key, value string) {
	(*md)[key] = value
}