package handlers

import (
	"../database"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

// GET /v1/player/{id}
func GetPlayer(ctx *fasthttp.RequestCtx) {

	id := ctx.UserValue("id").(string)
	team, err := database.GetPlayerByID(id)

	if err != nil {
		ctx.SetStatusCode(err.Code)
		return
	}

	ctx.SetStatusCode(200)
	setHeaders(ctx)

	resp, _ := json.Marshal(team)
	ctx.Write(resp)
}
