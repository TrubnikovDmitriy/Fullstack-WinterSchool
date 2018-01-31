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

func checkPathID(strID string) (int, *services.ErrorCode) {
	intID, err := strconv.Atoi(strID)
	if err != nil {
		return 0, &services.ErrorCode{
			Code: fasthttp.StatusBadRequest,
			Message: "Incorrect path variable '" + strID + "'",
		}
	}
	return intID, nil
}