package handlers

import (
	"../database"
	"../models"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

// GET /v1/tourney/{id}
func GetTournamentByID(ctx *fasthttp.RequestCtx) {

	id := ctx.UserValue("tourney_id").(string)
	tourney, err := database.GetTourneyByID(id)

	if err != nil {
		ctx.SetStatusCode(err.Code)
		return
	}

	ctx.SetStatusCode(200)
	setHeaders(ctx)

	resp, _ := json.Marshal(tourney)
	ctx.Write(resp)
}

// POST /v1/tourney
func CreateTournament(ctx *fasthttp.RequestCtx) {

	tournament := models.Tournament{}
	json.Unmarshal(ctx.PostBody(), &tournament)

	err := database.CreateTournament(&tournament)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
	} else {
		tournament.WriteAsJsonResponseTo(ctx, fasthttp.StatusCreated)
	}
}