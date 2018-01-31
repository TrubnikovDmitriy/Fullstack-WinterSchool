package services

import (
	"github.com/valyala/fasthttp"
	"encoding/json"
)

type ErrorCode struct {
	Code int `json:"-"`
	Message string `json:"error_message"`
	Link string `json:"href"`
}

func NewNotFound() *ErrorCode {
	return &ErrorCode{
		Code: fasthttp.StatusNotFound,
		Message: fasthttp.StatusMessage(fasthttp.StatusNotFound),
	}
}

func NewBadRequest() *ErrorCode {
	return &ErrorCode{
		Code: fasthttp.StatusBadRequest,
		Message: "Request is not valid",
		Link: "Ссылка на документацию к API",
	}
}

func NewServerError() *ErrorCode {
	return &ErrorCode{
		Code: fasthttp.StatusInternalServerError,
		Message: fasthttp.StatusMessage(fasthttp.StatusInternalServerError),
	}
}

func (httpCode *ErrorCode) WriteAsJsonResponseTo(ctx *fasthttp.RequestCtx) {
	resp, _ := json.Marshal(httpCode)
	ctx.Write(resp)
	ctx.SetStatusCode(httpCode.Code)
}
