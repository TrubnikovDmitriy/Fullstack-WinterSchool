package handlers

import (
	"../database"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

// GET /v1/team/{id}
func GetTeam(ctx *fasthttp.RequestCtx) {

	id := ctx.UserValue("id").(string)
	team := database.GetTeamByID(id)

	ctx.SetStatusCode(200)
	setHeaders(ctx)

	resp, _ := json.Marshal(team)
	ctx.Write(resp)
}