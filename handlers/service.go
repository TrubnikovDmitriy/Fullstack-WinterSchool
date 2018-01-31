package handlers

import (
	"../models"
	"../services"
	"github.com/valyala/fasthttp"
	"strconv"
)

func setHeaders(ctx *fasthttp.RequestCtx) {
	// TODO remove this func
	ctx.SetContentType("application/json; charset=utf-8")
}

func isValid(model models.Validator) bool {
	return model.Validate()
}

func getPathID(strID interface{}) (int, *services.ErrorCode) {
	if strID == nil {
		return 0, &services.ErrorCode{
			Code: fasthttp.StatusBadRequest,
			Message: "Empty path variable",
			Link: "Ссылка на документацию",
		}
	}
	intID, err := strconv.Atoi(strID.(string))
	if err != nil {
		return 0, &services.ErrorCode{
			Code: fasthttp.StatusBadRequest,
			Message: "Incorrect path variable '" + strID.(string) + "'",
		}
	}
	return intID, nil
}

