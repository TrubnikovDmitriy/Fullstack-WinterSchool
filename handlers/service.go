package handlers

import (
	"../services"
	"github.com/valyala/fasthttp"
	"github.com/satori/go.uuid"
)


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

