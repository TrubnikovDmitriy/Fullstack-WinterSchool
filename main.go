package main

import (
	"./handlers"
	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
)

var router *fasthttprouter.Router

func init() {
	router = fasthttprouter.New()

	// team
	router.GET("/v1/team/:id", handlers.GetTeam)
	router.GET("/v1/team/:id/players", handlers.GetTeamPlayers)

	// player
	router.GET("/v1/player/:id", handlers.GetPlayer)

	// match
	router.GET("/v1/match/:id", handlers.GetMatch)

	// tournament
	router.GET("/v1/tourney/:id", handlers.GetTournamentByID)

	// game
	router.GET("/v1/game/:id", handlers.GetGame)
}

func main() {
	fasthttp.ListenAndServe(":5000", router.Handler)
}
