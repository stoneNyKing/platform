package rest

import ()

type Response struct {
	Data   map[string]interface{} `json:"data"`
	Method string
}

func NewResponse() *Response {
	return &Response{Data: make(map[string]interface{})}
}

func (r *Response) GetErrorResponse(resp interface{}) interface{} {

	return nil
}

func (r *Response) GetSuccessResponse(resp interface{}) interface{} {

	return nil
}
