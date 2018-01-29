package handlers

import (
	"../database"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

// GET /v1/team/{id}
func GetTeam(ctx *fasthttp.RequestCtx) {

	id := ctx.UserValue("id").(string)
	team, err := database.GetTeamByID(id)

	if err != nil {
		ctx.SetStatusCode(err.Code)
		return
	}

	ctx.SetStatusCode(200)
	setHeaders(ctx)

	resp, _ := json.Marshal(team)
	ctx.Write(resp)
}

// GET /v1/team/{id}/players
func GetTeamPlayers(ctx *fasthttp.RequestCtx) {

	id := ctx.UserValue("id").(string)
	posts, _ := database.GetPlayersOfTeam(id)

	ctx.SetStatusCode(200)
	setHeaders(ctx)

	resp, _ := json.Marshal(posts)
	ctx.Write(resp)
}