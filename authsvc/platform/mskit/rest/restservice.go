package rest

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/libra9z/httprouter"
	"net/http"
)

var (
	// ErrTwoZeroes is an arbitrary business rule for the Add method.
	ErrTwoZeroes = errors.New("can't sum two zeroes")

	// ErrIntOverflow protects the Add method. We've decided that this error
	// indicates a misbehaving service and should count against e.g. circuit
	// breakers. So, we return it directly in endpoints, to illustrate the
	// difference. In a real service, this probably wouldn't be the case.
	ErrIntOverflow = errors.New("integer overflow")

	// ErrMaxSizeExceeded protects the Concat method.
	ErrMaxSizeExceeded = errors.New("result exceeds maximum size")
)

func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError
	msg := err.Error()

	switch err {
	case ErrTwoZeroes, ErrMaxSizeExceeded, ErrIntOverflow:
		code = http.StatusBadRequest
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorWrapper{Error: msg})
}

type errorWrapper struct {
	Error string `json:"error"`
}

type RestService interface {
	Get(context.Context, *Request) (interface{}, error)
	Post(context.Context, *Request) (interface{}, error)
	Delete(context.Context, *Request) (interface{}, error)
	Put(context.Context, *Request) (interface{}, error)
	Head(context.Context, *Request) (interface{}, error)
	Patch(context.Context, *Request) (interface{}, error)
	Options(context.Context, *Request) (interface{}, error)
	Trace(context.Context, *Request) (interface{}, error)

	//response relate interface
	SetRouter(router *httprouter.Router)
	GetErrorResponse() interface{}
	DecodeRequest(context.Context, *http.Request) (request interface{}, err error)
	EncodeResponse(context.Context, http.ResponseWriter, interface{}) error
}
