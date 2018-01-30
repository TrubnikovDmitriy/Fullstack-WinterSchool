package handlers

import "github.com/valyala/fasthttp"

func setHeaders(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
}