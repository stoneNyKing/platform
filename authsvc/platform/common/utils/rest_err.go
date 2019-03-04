package utils

import (
	"net/http"
)

func RestResponse(w http.ResponseWriter, code int) {
	var body string
	switch code {
	case 200:
		body = "OK"
	case 201:
		body = "Created"
	case 202:
		body = "Accepted"
	case 204:
		body = "no content"
	case 400:
		body = "invalid request"
	case 401:
		body = "Unauthorized"
	case 403:
		body = "Forbidden"
	case 404:
		body = "Not found"
	case 406:
		body = "Not acceptable"
	case 410:
		body = "Gone"
	case 422:
		body = "Unprocesable entity"
	case 500:
		body = "internal server error"
	}

	http.Error(w, body, code)
}
func RestResponseString(code int) string {
	var body string
	switch code {
	case 200:
		body = "OK"
	case 201:
		body = "Created"
	case 202:
		body = "Accepted"
	case 204:
		body = "no content"
	case 400:
		body = "invalid request"
	case 401:
		body = "Unauthorized"
	case 403:
		body = "Forbidden"
	case 404:
		body = "Not found"
	case 406:
		body = "Not acceptable"
	case 410:
		body = "Gone"
	case 422:
		body = "Unprocesable entity"
	case 500:
		body = "internal server error"
	}

	return body
}
