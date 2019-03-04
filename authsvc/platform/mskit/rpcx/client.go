package rpcx

import (
	"context"
	"reflect"
	"strings"
	"sync"

	"github.com/smallnest/rpcx/client"

	"github.com/go-kit/kit/endpoint"
)

// Client wraps a RPCx connection and provides a method that implements
// endpoint.Endpoint.
type Client struct {
	client      client.XClient
	failMode 	client.FailMode
	selectMode 	client.SelectMode
	serviceName string
	service		string
	method      string
	sdType      string
	sdAddress   string
	basePath    string
	rpcxReply   reflect.Type
	before      []ClientRequestFunc
	after       []ClientResponseFunc
	finalizer   []ClientFinalizerFunc
}

var ClientPool map[string]*sync.Pool
var lock *sync.Mutex = &sync.Mutex {}

func init(){
	ClientPool = make(map[string]*sync.Pool)
}

// NewClient constructs a usable Client for a single remote endpoint.
// Pass an zero-value protobuf message of the RPC response type as
// the rpcxReply argument.
func NewClient(
	rpcxReply interface{},
	options ...ClientOption,
) *Client {

	c := &Client{
		// We are using reflect.Indirect here to allow both reply structs and
		// pointers to these reply structs. New consumers of the client should
		// use structs directly, while existing consumers will not break if they
		// remain to use pointers to structs.
		rpcxReply: reflect.TypeOf(
			reflect.Indirect(
				reflect.ValueOf(rpcxReply),
			).Interface(),
		),
		before: []ClientRequestFunc{},
		after:  []ClientResponseFunc{},
	}
	for _, option := range options {
		option(c)
	}

	return c
}

// NewClient constructs a usable Client for a single remote endpoint.
// Pass an zero-value protobuf message of the RPC response type as
// the rpcxReply argument.
func NewClientPool(sdtype,sdaddr,basepath, serviceName string,failMode client.FailMode,selectMode client.SelectMode) *sync.Pool {

	clientPool := sync.Pool{New: func() interface{} {
		var cs client.ServiceDiscovery
		switch sdtype {
		case "consul":
			ss := strings.Split(sdaddr, ";")
			cs = client.NewConsulDiscovery(basepath, serviceName, ss, nil)
		case "etcd":
			ss := strings.Split(sdaddr, ";")
			cs = client.NewEtcdDiscovery(basepath, serviceName, ss, nil)
		case "zookeeper":
			ss := strings.Split(sdaddr, ";")
			cs = client.NewZookeeperDiscovery(basepath, serviceName, ss, nil)
		}

		xclient := client.NewXClient(serviceName, failMode, selectMode, cs, client.DefaultOption)
		return xclient
	}}

	ClientPool[serviceName] = &clientPool
	return &clientPool
}

func GetRpcClientPool(serviceName string) *sync.Pool {
	return ClientPool[serviceName]
}

// ClientOption sets an optional parameter for clients.
type ClientOption func(*Client)

// ClientBefore sets the RequestFuncs that are applied to the outgoing RPCx
// request before it's invoked.
func ClientBefore(before ...ClientRequestFunc) ClientOption {
	return func(c *Client) { c.before = append(c.before, before...) }
}

// ClientAfter sets the ClientResponseFuncs that are applied to the incoming
// RPCx response prior to it being decoded. This is useful for obtaining
// response metadata and adding onto the context prior to decoding.
func ClientAfter(after ...ClientResponseFunc) ClientOption {
	return func(c *Client) { c.after = append(c.after, after...) }
}

// ClientFinalizer is executed at the end of every RPCx request.
// By default, no finalizer is registered.
func ClientFinalizer(f ...ClientFinalizerFunc) ClientOption {
	return func(s *Client) { s.finalizer = append(s.finalizer, f...) }
}

// Endpoint returns a usable endpoint that will invoke the RPCx specified by the
// client.
func (c Client) Endpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		if c.finalizer != nil {
			defer func() {
				for _, f := range c.finalizer {
					f(ctx, err)
				}
			}()
		}

		ctx = context.WithValue(ctx, ContextKeyRequestMethod, c.method)

		md := make(map[string]string)
		for _, f := range c.before {
			ctx = f(ctx, &md)
		}

		//ctx = NewReqMetaDataContext(ctx, map[string]string(md))

		traceplugin :=client.OpenTracingPlugin{}
		pc := client.NewPluginContainer()
		pc.Add(traceplugin)
		c.client.SetPlugins(pc)

		req := request.(*RpcRequest)
		rpcxReply := reflect.New(c.rpcxReply).Interface()
		if err = c.client.Call(	ctx, c.service, req, rpcxReply); err != nil {
			return nil, err
		}

		var header, trailer map[string]string
		for _, f := range c.after {
			ctx = f(ctx, header, trailer)
		}

		return rpcxReply, nil
	}
}

func (c *Client)Close() error {
	pc := ClientPool[c.serviceName]
	pc.Put(c.client)
	return nil
}

func (c *Client)GetClientPool() *sync.Pool {

	if pc,ok :=  ClientPool[c.serviceName];ok {
		return pc
	}else{
		lock.Lock()
		defer lock.Unlock()
		if ClientPool[c.serviceName] == nil {
			ClientPool[c.serviceName] = NewClientPool(c.sdType,c.sdAddress,c.basePath,c.serviceName,c.failMode,c.selectMode)
		}
	}

	return ClientPool[c.serviceName]
}


// ClientFinalizerFunc can be used to perform work at the end of a client RPCx
// request, after the response is returned. The principal
// intended use is for error logging. Additional response parameters are
// provided in the context under keys with the ContextKeyResponse prefix.
// Note: err may be nil. There maybe also no additional response parameters depending on
// when an error occurs.
type ClientFinalizerFunc func(ctx context.Context, err error)


func BasePathOption( basepath string) ClientOption {
	return func(c *Client){ c.basePath = basepath}
}
func SdTypeOption( sdtype string) ClientOption {
	return func(c *Client){ c.sdType = sdtype}
}
func SdAddressOption( sdaddress string) ClientOption {
	return func(c *Client){ c.sdAddress = sdaddress}
}
func FailModeOption( failmode client.FailMode) ClientOption {
	return func(c *Client){ c.failMode = failmode}
}
func SelectModeOption( sel client.SelectMode) ClientOption {
	return func(c *Client){ c.selectMode = sel}
}
func MethodOption( method string) ClientOption {
	return func(c *Client){ c.method = method}
}
func ServiceOption( service string) ClientOption {
	return func(c *Client){ c.service = service}
}
func ServiceNameOption( svrname string) ClientOption {
	return func(c *Client){ c.serviceName = svrname}
}
