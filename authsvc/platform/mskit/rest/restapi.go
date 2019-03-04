package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"github.com/libra9z/httprouter"
	"io/ioutil"
	"net/http"
	"strings"
)

type RestApi struct {
	Request *http.Request
	Router  *httprouter.Router
	Counter		metrics.Counter
	Gauge 		metrics.Gauge
	Histogram 	metrics.Histogram
}

// Get adds a request function to handle GET request.
func (c *RestApi) SetRouter(r *httprouter.Router) {
	c.Router = r
}

// Get adds a request function to handle GET request.
func (c *RestApi) Get(ctx context.Context, r *Request) (interface{}, error) {
	return nil, nil
}

// Post adds a request function to handle POST request.
func (c *RestApi) Post(ctx context.Context, r *Request) (interface{}, error) {
	return nil, nil
}

// Delete adds a request function to handle DELETE request.
func (c *RestApi) Delete(ctx context.Context, r *Request) (interface{}, error) {
	return nil, nil
}

// Put adds a request function to handle PUT request.
func (c *RestApi) Put(ctx context.Context, r *Request) (interface{}, error) {
	return nil, nil
}

// Head adds a request function to handle HEAD request.
func (c *RestApi) Head(ctx context.Context, r *Request) (interface{}, error) {
	return nil, nil
}

// Patch adds a request function to handle PATCH request.
func (c *RestApi) Patch(ctx context.Context, r *Request) (interface{}, error) {
	return nil, nil
}

// Options adds a request function to handle OPTIONS request.
func (c *RestApi) Options(ctx context.Context, r *Request) (interface{}, error) {
	return nil, nil
}

// Options adds a request function to handle OPTIONS request.
func (c *RestApi) Trace(ctx context.Context, r *Request) (interface{}, error) {
	return nil, nil
}

// GetErrorResponse adds a restservice used for endpoint.
func (c *RestApi) GetErrorResponse() interface{} {
	resp := NewResponse()
	resp.Data["ret"] = 1
	resp.Data["error"] = errors.New("Not allowed.")
	return resp
}

// DecodeRequest adds a restservice used for endpoint.
/*
需要在nginx上配置
proxy_set_header Remote_addr $remote_addr;
*/
func (c *RestApi) DecodeRequest(_ context.Context, r *http.Request) (request interface{}, err error) {

	c.Request = r

	req := Request{Queries: make(map[string]interface{})}

	req.Method = r.Method

	if c.Router == nil {
		fmt.Printf("no router set.\n")
		return nil, errors.New("no router set.")
	}

	_, req.Params, _ = c.Router.Lookup(r.Method, r.URL.EscapedPath())

	values := r.URL.Query()

	accept := r.Header.Get("Accept")
	ss := strings.Split(accept, ";")

	for _, s := range ss {
		sv := strings.Split(s, "=")

		if len(sv) > 1 && strings.TrimSpace(sv[0]) == "version" {
			req.Version = sv[1]
		}
	}

	for k, v := range values {
		req.Queries[k] = v
	}

	ip := r.Header.Get("X-Real-IP")

	if ip == "" {
		req.RemoteAddr = r.RemoteAddr
	} else {
		req.RemoteAddr = ip
	}

	req.OriginRequest = r

	if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		req.Body, err = ioutil.ReadAll(r.Body)

		if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			req.ContentType = CONTENT_TYPE_JSON
		} else if strings.Contains(r.Header.Get("Content-Type"), "application/xml") ||
			strings.Contains(r.Header.Get("Content-Type"), "text/xml") {
			req.ContentType = CONTENT_TYPE_XML
		} else if strings.Contains(r.Header.Get("Content-Type"), "x-www-form-urlencoded") {
			req.ContentType = CONTENT_TYPE_FORM
		}
	} else {
		req.ContentType = CONTENT_TYPE_MULTIFORM
	}

	return req, nil
}

func (c *RestApi) Prepare(r *Request) (interface{}, error) {
	return nil, nil
}

/*
*该方法是在response返回之前调用，用于增加一下个性化的头信息
 */
func (c *RestApi) Finish(w http.ResponseWriter) error {

	if w == nil {
		return errors.New("writer is nil ")
	}

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Origin,Accept,Content-Range,Content-Description,Content-Disposition")
	w.Header().Add("Access-Control-Allow-Methods", "PUT,GET,POST,DELETE,OPTIONS")

	return nil
}

// EncodeResponse adds a restservice used for endpoint.
func (c *RestApi) EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {

	if response == nil {
		response = ""
	}

	w.Header().Set("Allow", "HEAD,GET,PUT,DELETE,OPTIONS,POST")

	c.Finish(w)

	err := json.NewEncoder(w).Encode(response)

	return err
}
