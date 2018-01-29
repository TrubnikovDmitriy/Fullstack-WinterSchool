package handlers

import (
	"../database"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

// GET /v1/player/{id}
func GetPlayer(ctx *fasthttp.RequestCtx) {

	id := ctx.UserValue("id").(string)
	team := database.GetPlayerByID(id)

	ctx.SetStatusCode(200)
	setHeaders(ctx)

	resp, _ := json.Marshal(team)
	ctx.Write(resp)
}
