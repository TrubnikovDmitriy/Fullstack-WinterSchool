package models

import (
	"github.com/valyala/fasthttp"
)

type Link struct {
	Rel 	string `json:"rel"`
	Href 	string `json:"href"`
	Action 	string `json:"action"`
}

type Validator interface {
	Validate() bool
}

type Linker interface {
	GenerateLinks()
}

type Writer interface {
	WriteAsJsonResponseTo(ctx *fasthttp.RequestCtx, statusCode int)
}