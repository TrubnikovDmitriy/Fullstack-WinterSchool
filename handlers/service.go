package handlers

import (
	"../services"
	"github.com/valyala/fasthttp"
	"strconv"
	"github.com/satori/go.uuid"
)

func setHeaders(ctx *fasthttp.RequestCtx) {
	// TODO remove this func
	ctx.SetContentType("application/json; charset=utf-8")
}

func getPathID(strID interface{}) (int, *serv.ErrorCode) {
	if strID == nil {
		return 0, &serv.ErrorCode{
			Code: fasthttp.StatusBadRequest,
			Message: "Empty path variable",
			Link: "Ссылка на документацию",
		}
	}
	intID, err := strconv.Atoi(strID.(string))
	if err != nil {
		return 0, &serv.ErrorCode{
			Code: fasthttp.StatusBadRequest,
			Message: "Incorrect path variable '" + strID.(string) + "'",
		}
	}
	return intID, nil
}

func getPathUUID(strID string) (uuid.UUID, *serv.ErrorCode) {
	id, err := uuid.FromString(strID)
	if err != nil {
		return uuid.Nil, &serv.ErrorCode{
			Code: fasthttp.StatusBadRequest,
			Message: "Incorrect path variable: '" + strID + "'",
			Link: "Ссылка на документацию",
		}
	}
	return id, nil
}

