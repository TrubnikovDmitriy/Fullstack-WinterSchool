package handlers

import (
	"../database"
	"../models"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

// GET /v1/teams/{team_id}/players/{player_id}
func GetPlayer(ctx *fasthttp.RequestCtx) {

	teamStrID := ctx.UserValue("team_id").(string)
	playerStrID := ctx.UserValue("player_id").(string)
	teamID, err := checkPathID(teamStrID)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
		return
	}
	playerID, err := checkPathID(playerStrID)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
		return
	}


	player, err := database.GetPlayerByIDs(teamID, playerID)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
	} else {
		player.WriteAsJsonResponseTo(ctx, fasthttp.StatusOK)
	}
}

// GET /v1/team/{team_id}/players
func GetTeamPlayers(ctx *fasthttp.RequestCtx) {

	teamStrID := ctx.UserValue("team_id").(string)
	teamID, err := checkPathID(teamStrID)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
		return
	}

	players, err := database.GetPlayersOfTeam(teamID)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
	} else {
		resp, _ := json.Marshal(players)
		ctx.Write(resp)
		setHeaders(ctx)
		ctx.SetStatusCode(200)
	}
}

// POST /v1/teams/{team_id}/players
func CreatePlayer(ctx *fasthttp.RequestCtx) {

	teamStrID := ctx.UserValue("team_id").(string)
	teamID, err := checkPathID(teamStrID)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
		return
	}

	player := models.Player{TeamID: teamID}
	json.Unmarshal(ctx.PostBody(), &player)

	err = database.CreatePlayer(&player)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
	} else {
		player.WriteAsJsonResponseTo(ctx, fasthttp.StatusCreated)
	}
}