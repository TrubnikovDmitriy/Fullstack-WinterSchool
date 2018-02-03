package main

import (
	"./handlers"
	"./database"
	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
)

var router *fasthttprouter.Router

func init() {
	router = fasthttprouter.New()

	// persons
	router.POST("/v1/persons", handlers.CreatePerson)
	router.GET("/v1/persons/:person_id", handlers.GetPerson)

	// teams
	router.GET("/v1/teams/:team_id", handlers.GetTeam)
	router.POST("/v1/teams", handlers.CreateTeam)

	// players
	router.GET("/v1/teams/:team_id/players/:player_id", handlers.GetPlayer)
	router.GET("/v1/teams/:team_id/players", handlers.GetTeamPlayers)
	router.POST("/v1/players", handlers.CreatePlayer)

	// matches
	router.GET("/v1/tourney/:tourney_id/matches/:match_id", handlers.GetMatch)

	// tournaments
	router.GET("/v1/tourney/:tourney_id", handlers.GetTournamentByID)
	router.GET("/v1/tourney/:tourney_id/matches", handlers.GetTournamentGrid)
	router.POST("/v1/tourney", handlers.CreateTournament)

	// games
	router.GET("/v1/games/:id", handlers.GetGame)
	router.POST("/v1/games", handlers.CreateGame)

	// test
	router.GET("/test", database.Test)
}

func main() {
	fasthttp.ListenAndServe(":5000", router.Handler)
}
