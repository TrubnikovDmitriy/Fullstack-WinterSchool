package models

import (
	"../services"
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


func fieldLengthValidate(field string, fieldName string) *serv.ErrorCode {

	if len(field) == 0 {
		return serv.NewBadRequest("The " + fieldName + " is missing")
	}
	if len(field) > serv.MaxFieldLength {
		return serv.NewBadRequest("The " + fieldName + " is too long")
	}
	return nil
}