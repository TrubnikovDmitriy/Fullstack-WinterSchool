package models

import (
	"github.com/valyala/fasthttp"
	"encoding/json"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (tokens *Tokens) WriteAsJsonResponseTo(ctx *fasthttp.RequestCtx, statusCode int) {
	resp, _ := json.Marshal(tokens)
	ctx.Write(resp)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetStatusCode(statusCode)
}