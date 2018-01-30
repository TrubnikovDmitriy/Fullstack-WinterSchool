package handlers

import (
	"../database"
	"../models"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

// GET /v1/games/{id}
func GetGame(ctx *fasthttp.RequestCtx) {

	id := ctx.UserValue("id").(string)
	game, err := database.GetGameByID(id)

	if err != nil {
		err.WriteAsJsonResponse(ctx)
		return
	}

	game.WriteAsJsonTo(ctx)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

// POST /v1/games
func CreateGame(ctx *fasthttp.RequestCtx) {

	game := models.Game{}
	json.Unmarshal(ctx.PostBody(), &game)

	err := database.CreateNewGame(&game)
	if err != nil {
		err.WriteAsJsonResponse(ctx)
		return
	}

	game.WriteAsJsonTo(ctx)
	ctx.SetStatusCode(fasthttp.StatusCreated)
}