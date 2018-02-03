package handlers

import (
	"../database"
	"../models"
	"github.com/valyala/fasthttp"
	"encoding/json"
)


// GET /v1/tourney/{id}
func GetTournamentByID(ctx *fasthttp.RequestCtx) {

	id, err := getPathID(ctx.UserValue("tourney_id").(string))
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
		return
	}
	tourney, err := database.GetTourneyByID(id)

	if err != nil {
		ctx.SetStatusCode(err.Code)
	} else {
		tourney.WriteAsJsonResponseTo(ctx, fasthttp.StatusOK)
	}
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

// GET /v1/tourney/{id}/matches
func GetTournamentGrid(ctx *fasthttp.RequestCtx) {

	id, err := getPathID(ctx.UserValue("tourney_id").(string))
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
		return
	}

	grid, err := database.GetTournamentGrid(id)
	if err != nil {
		ctx.SetStatusCode(err.Code)
	} else {
		grid.WriteAsJsonResponseTo(ctx, fasthttp.StatusOK)
	}
}
