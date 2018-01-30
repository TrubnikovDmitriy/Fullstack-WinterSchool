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

func CreateNew(code int) *ErrorCode {
	return &ErrorCode{code, "", ""}
}

func (httpCode *ErrorCode) WriteAsJsonResponse(ctx *fasthttp.RequestCtx) {
	resp, _ := json.Marshal(httpCode)
	ctx.Write(resp)
	ctx.SetStatusCode(httpCode.Code)
}
