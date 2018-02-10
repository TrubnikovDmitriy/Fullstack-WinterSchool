package handlers

import (
	"../database"
	"../models"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

// GET /v1/teams/{team_id}/players/{player_id}
func GetPlayer(ctx *fasthttp.RequestCtx) {

	teamID, err := getPathID(ctx.UserValue("team_id").(string))
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
		return
	}
	playerID, err := getPathID(ctx.UserValue("player_id").(string))
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
		return
	}


	player, err := database.GetPlayerByID(teamID, playerID)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
	} else {
		player.WriteAsJsonResponseTo(ctx, fasthttp.StatusOK)
	}
}

// GET /v1/team/{team_id}/players
func GetTeamPlayers(ctx *fasthttp.RequestCtx) {

	teamID, err := getPathID(ctx.UserValue("team_id").(string))
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
		ctx.SetContentType("application/json; charset=utf-8")
		ctx.SetStatusCode(200)
	}
}

// POST /v1/team/{team_id}/players
func CreatePlayer(ctx *fasthttp.RequestCtx) {

	player := new(models.Player)
	json.Unmarshal(ctx.PostBody(), player)

	err := database.AddPlayerInTeam(player)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
	} else {
		player.WriteAsJsonResponseTo(ctx, fasthttp.StatusCreated)
	}
}

// DELETE /v1/team/{team_id}/players
func DeletePlayer(ctx *fasthttp.RequestCtx) {

	player := new(models.Player)
	json.Unmarshal(ctx.PostBody(), player)

	err := database.DeletePlayerFromTeam(player)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
	} else {
		player.WriteAsJsonResponseTo(ctx, fasthttp.StatusCreated)
	}
}