package handlers

import (
	"../database"
	"../models"
	"../services"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

// GET /v1/games/{game_id}
func GetGame(ctx *fasthttp.RequestCtx) {

	id, err := getPathID(ctx.UserValue("game_id").(string))
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
		return
	}
	game, err := database.GetGameByID(id)

	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
 	} else {
 		game.WriteAsJsonResponseTo(ctx, fasthttp.StatusOK)
	}
}

// GET /v1/games
func GetGames(ctx *fasthttp.RequestCtx) {

	limitStr := ctx.QueryArgs().Peek("limit")
	pageStr := ctx.QueryArgs().Peek("page")

	limit := getIntFromBytes(limitStr, 6)
	page := getIntFromBytes(pageStr, 1)

	games, err := database.GetGames(limit, page)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
		return
	}
	games.WriteAsJsonResponseTo(ctx, fasthttp.StatusOK)
}

// POST /v1/games
func CreateGame(ctx *fasthttp.RequestCtx) {

	claims, err := GetClaimsFromCookie(ctx)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
		return
	}
	if claims["staff"].(string) != "true" {
		err = serv.NewForbidden("You must be administrator to perform this action")
		err.WriteAsJsonResponseTo(ctx)
		return
	}

	game := new(models.Game)
	json.Unmarshal(ctx.PostBody(), game)

	err = database.CreateGame(game)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
	} else {
		game.WriteAsJsonResponseTo(ctx, fasthttp.StatusCreated)
	}
}