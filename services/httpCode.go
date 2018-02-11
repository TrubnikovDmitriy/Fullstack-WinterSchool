package serv

import (
	"github.com/valyala/fasthttp"
	"encoding/json"
	"log"
	"strconv"
)

type ErrorCode struct {
	Code    int    `json:"-"`
	Message string `json:"error_message"`
	Link    string `json:"href"`
}

func NewNotFound() *ErrorCode {
	return &ErrorCode{
		Code:    fasthttp.StatusNotFound,
		Message: fasthttp.StatusMessage(fasthttp.StatusNotFound),
	}
}

func NewBadRequest(message string) *ErrorCode {
	return &ErrorCode{
		Code:    fasthttp.StatusBadRequest,
		Message: message,
		Link:    "Ссылка на документацию к API",
	}
}

func NewServerError(err error) *ErrorCode {
	log.Print(err)
	return &ErrorCode{
		Code:    fasthttp.StatusInternalServerError,
		Message: fasthttp.StatusMessage(fasthttp.StatusInternalServerError),
	}
}

func NewUnauthorized() *ErrorCode {
	return &ErrorCode{
		Code:    fasthttp.StatusUnauthorized,
		Message: "Для совершения указанного действия необходимо зарегистрироваться",
		Link:    GetConfig().Href + "/v1/oauth/authorize",
	}

}

func NewForbidden(message string) *ErrorCode {
	return &ErrorCode{
		Code:    fasthttp.StatusForbidden,
		Message: message,
	}

}

func (err *ErrorCode) String() string {

	printString := "Error code: " + strconv.Itoa(err.Code) +
		"\nError message: " + err.Message + "\n"

	if len(err.Link) != 0 {
		printString += "Link: " + err.Link + "\n"
	}
	return printString + "\n"
}

func (httpCode *ErrorCode) WriteAsJsonResponseTo(ctx *fasthttp.RequestCtx) {
	resp, _ := json.Marshal(httpCode)
	ctx.Write(resp)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.Response.Header.SetBytesV("Access-Control-Allow-Origin", ctx.Request.Header.Peek("origin"))
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	ctx.SetStatusCode(httpCode.Code)
}
