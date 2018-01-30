package handlers

import (
	"github.com/valyala/fasthttp"
)

func setHeaders(ctx *fasthttp.RequestCtx) {
	// TODO remove this func
	ctx.SetContentType("application/json; charset=utf-8")
}