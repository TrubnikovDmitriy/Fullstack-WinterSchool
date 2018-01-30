package handlers

import (
	"../database"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

// GET /v1/game/{id}
func GetGame(ctx *fasthttp.RequestCtx) {

	id := ctx.UserValue("id").(string)
	game, err := database.GetGameByID(id)

	if err != nil {
		ctx.SetStatusCode(err.Code)
		return
	}

	ctx.SetStatusCode(200)
	setHeaders(ctx)

	resp, _ := json.Marshal(game)
	ctx.Write(resp)
}
