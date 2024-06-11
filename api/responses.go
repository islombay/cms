package api

import "net/http"

type Response struct {
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   []error     `json:"error,omitempty"`
}

var (
	StatusBadRequest = Response{
		Message: "Bad request",
		Code:    http.StatusBadRequest,
	}

	StatusInternal = Response{
		Message: "Internal server error",
		Code:    http.StatusInternalServerError,
	}

	StatusUnauthorized = Response{
		Message: "Unauthorized",
		Code:    http.StatusUnauthorized,
	}

	StatusForbidden = Response{
		Message: "Forbidden",
		Code:    http.StatusForbidden,
	}
)

func makeResponse(data interface{}) Response {
	return Response{
		Code: http.StatusOK,
		Data: data,
	}
}
