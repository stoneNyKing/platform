package trace

import (
	"context"
	"fmt"
	"github.com/smallnest/rpcx/share"
	"strings"
)


func Pairs(kv ...string) map[string]string {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: Pairs got the odd number of input pairs for metadata: %d", len(kv)))
	}
	md := map[string]string{}
	var key string
	for i, s := range kv {
		if i%2 == 0 {
			key = strings.ToLower(s)
			continue
		}
		md[key] =  s
	}
	return md
}

// Join joins any number of mds into a single MD.
// The order of values for each key is determined by the order in which
// the mds containing those values are presented to Join.
func Join(mds ...map[string]string) map[string]string {
	out := map[string]string{}
	for _, md := range mds {
		for k, v := range md {
			out[k] = v
		}
	}
	return out
}

type mdIncomingKey struct{}
type mdOutgoingKey struct{}

// NewIncomingContext creates a new context with incoming md attached.
func NewIncomingContext(ctx context.Context, md map[string]string) context.Context {
	return context.WithValue(ctx, mdIncomingKey{}, md)
}

// NewOutgoingContext creates a new context with outgoing md attached. If used
// in conjunction with AppendToOutgoingContext, NewOutgoingContext will
// overwrite any previously-appended metadata.
func NewReqMetaDataContext(ctx context.Context, md map[string]string) context.Context {
	return context.WithValue(ctx,share.ReqMetaDataKey,md)
}

// AppendToOutgoingContext returns a new context with the provided kv merged
// with any existing metadata in the context. Please refer to the
// documentation of Pairs for a description of kv.
func AppendToOutgoingContext(ctx context.Context, kv ...string) context.Context {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: AppendToOutgoingContext got an odd number of input pairs for metadata: %d", len(kv)))
	}
	md, _ := ctx.Value(mdOutgoingKey{}).(rawMD)
	added := make([][]string, len(md.added)+1)
	copy(added, md.added)
	added[len(added)-1] = make([]string, len(kv))
	copy(added[len(added)-1], kv)
	return context.WithValue(ctx, mdOutgoingKey{}, rawMD{md: md.md, added: added})
}

// FromIncomingContext returns the incoming metadata in ctx if it exists.  The
// returned MD should not be modified. Writing to it may cause races.
// Modification should be made to copies of the returned MD.
func FromIncomingContext(ctx context.Context) (md map[string]string, ok bool) {
	md, ok = ctx.Value(mdIncomingKey{}).(map[string]string)
	return
}

// FromOutgoingContextRaw returns the un-merged, intermediary contents
// of rawMD. Remember to perform strings.ToLower on the keys. The returned
// MD should not be modified. Writing to it may cause races. Modification
// should be made to copies of the returned MD.
//
// This is intended for gRPC-internal use ONLY.
func FromOutgoingContextRaw(ctx context.Context) (map[string]string, [][]string, bool) {
	raw, ok := ctx.Value(mdOutgoingKey{}).(rawMD)
	if !ok {
		return nil, nil, false
	}

	return raw.md, raw.added, true
}

// FromOutgoingContext returns the outgoing metadata in ctx if it exists.  The
// returned MD should not be modified. Writing to it may cause races.
// Modification should be made to copies of the returned MD.
func FromOutgoingContext(ctx context.Context) (map[string]string, bool) {
	raw, ok := ctx.Value(mdOutgoingKey{}).(rawMD)
	if !ok {
		return nil, false
	}

	mds := make([]map[string]string, 0, len(raw.added)+1)
	mds = append(mds, raw.md)
	for _, vv := range raw.added {
		mds = append(mds, Pairs(vv...))
	}
	return Join(mds...), ok
}

type rawMD struct {
	md    map[string]string
	added [][]string
}
