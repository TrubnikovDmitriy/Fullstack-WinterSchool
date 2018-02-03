package handlers

import (
	"../database"
	"../models"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

// GET /v1/teams/{team_id}
func GetTeam(ctx *fasthttp.RequestCtx) {

	teamID, err := getPathUUID(ctx.UserValue("team_id").(string))
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
	}

	team, err := database.GetTeamByID(teamID)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
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
		err.WriteAsJsonResponseTo(ctx)
	} else {
		team.WriteAsJsonResponseTo(ctx, fasthttp.StatusCreated)
	}
}