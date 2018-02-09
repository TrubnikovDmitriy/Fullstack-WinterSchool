package handlers

import (
	"../database"
	"../models"
	"github.com/valyala/fasthttp"
	"encoding/json"
	"github.com/satori/go.uuid"
)

// GET /v1/teams/{team_id}
func GetTeam(ctx *fasthttp.RequestCtx) {

	teamID, err := getPathID(ctx.UserValue("team_id").(string))
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


	claims, err := GetClaimsFromCookie(ctx)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
		return
	}

	team := new(models.Team)
	json.Unmarshal(ctx.PostBody(), team)

	team.CoachID = uuid.FromStringOrNil(claims["person_id"].(string))
	team.CoachName = claims["first_name"].(string) + " " + claims["last_name"].(string)

	err = database.CreateTeam(team)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
	} else {
		team.WriteAsJsonResponseTo(ctx, fasthttp.StatusCreated)
	}
}