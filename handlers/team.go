package handlers

import (
	"../database"
	"../models"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

// GET /v1/teams/{team_id}
func GetTeam(ctx *fasthttp.RequestCtx) {

	teamID, err := getPathID(ctx.UserValue("team_id"))
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
	}

	team, err := database.GetTeamByID(teamID)
	if err != nil {
		ctx.SetStatusCode(err.Code)
	} else {
		team.WriteAsJsonResponseTo(ctx, fasthttp.StatusOK)
	}
}

// POST /v1/teams
func CreateTeam(ctx *fasthttp.RequestCtx) {

	var team models.Team
	json.Unmarshal(ctx.PostBody(), &team)

	err := database.CreateTeam(&team)
	if err != nil {
		ctx.SetStatusCode(err.Code)
	} else {
		team.WriteAsJsonResponseTo(ctx, fasthttp.StatusCreated)
	}
}