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
	router.GET("/v1/teams/:id", handlers.GetTeam)
	router.GET("/v1/teams/:id/players", handlers.GetTeamPlayers)

	// player
	router.GET("/v1/players/:id", handlers.GetPlayer)

	// match
	router.GET("/v1/matches/:id", handlers.GetMatch)

	// tournament
	router.GET("/v1/tourney/:id", handlers.GetTournamentByID)

	// game
	router.GET("/v1/games/:id", handlers.GetGame)
	router.POST("/v1/games", handlers.CreateGame)
}

func main() {
	fasthttp.ListenAndServe(":5000", router.Handler)
}
