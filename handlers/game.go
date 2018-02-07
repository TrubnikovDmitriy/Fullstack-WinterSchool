package handlers

import (
	"../database"
	"../models"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

// GET /v1/games/{id}
func GetGame(ctx *fasthttp.RequestCtx) {

	id, err := getPathID(ctx.UserValue("id").(string))
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

// POST /v1/games
func CreateGame(ctx *fasthttp.RequestCtx) {

	game := new(models.Game)
	json.Unmarshal(ctx.PostBody(), game)

	err := database.CreateGame(game)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
	} else {
		game.WriteAsJsonResponseTo(ctx, fasthttp.StatusCreated)
	}
}