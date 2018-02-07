package handlers

import (
	"github.com/valyala/fasthttp"
	"../services"
	"../models"
	"../cache"
	"../database"
	"encoding/json"
)

func CreateToken(ctx *fasthttp.RequestCtx) {

	redirectURL := string(ctx.QueryArgs().Peek("redirect"))
	if len(redirectURL) == 0 {
		err := serv.NewBadRequest("Redirect path is absent")
		err.WriteAsJsonResponseTo(ctx)
		return
	}

	oauth := models.OAuth{}
	json.Unmarshal(ctx.PostBody(), &oauth)
	err := database.Auth(&oauth)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
		return
	}

	code := cache.CreateCode(&oauth)
	if code != nil {
		err = &serv.ErrorCode {
			Code: fasthttp.StatusInternalServerError,
			Message: fasthttp.StatusMessage(fasthttp.StatusInternalServerError),
		}
		err.WriteAsJsonResponseTo(ctx)
	}
	ctx.Redirect(redirectURL + "?code=" + code.String(), fasthttp.StatusFound)
}

func GetToken(ctx *fasthttp.RequestCtx) {

}

func RefreshToken(ctx *fasthttp.RequestCtx) {

}
