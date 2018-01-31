package handlers

import (
	"../database"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

// GET /v1/team/{team_id}
func GetTeam(ctx *fasthttp.RequestCtx) {

	id := ctx.UserValue("team_id").(string)
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