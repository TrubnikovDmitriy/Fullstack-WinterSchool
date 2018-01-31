package main

import (
	"./handlers"
	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
)

var router *fasthttprouter.Router

func init() {
	router = fasthttprouter.New()

	// teams
	router.GET("/v1/teams/:team_id", handlers.GetTeam)
	router.POST("/v1/teams", handlers.CreateTeam)

	// players
	router.GET("/v1/teams/:team_id/players/:player_id", handlers.GetPlayer)
	router.GET("/v1/teams/:team_id/players", handlers.GetTeamPlayers)
	router.POST("/v1/teams/:team_id/players", handlers.CreatePlayer)

	// matches
	router.GET("/v1/tourney/:tourney_id/matches/:match_id", handlers.GetMatch)

	// tournaments
	router.GET("/v1/tourney/:tourney_id", handlers.GetTournamentByID)

	// games
	router.GET("/v1/games/:id", handlers.GetGame)
	router.POST("/v1/games", handlers.CreateGame)
}

func main() {
	fasthttp.ListenAndServe(":5000", router.Handler)
}
