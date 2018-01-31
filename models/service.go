package models

import (
	"github.com/valyala/fasthttp"
)

type Link struct {
	Rel string `json:"rel"`
	Href string `json:"href"`
	Action string `json:"action"`
}

// нахуй такие интерфейсы вообще нужны

type Validator interface {
	Validate() bool
}

type Linker interface {
	GenerateLinks()
}

type Writer interface {
	WriteAsJsonResponseTo(ctx *fasthttp.RequestCtx, statusCode int)
}