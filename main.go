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
	router.GET("/v1/games/:game_id", handlers.GetGame)
	router.POST("/v1/games", handlers.CreateGame)
	router.GET("/v1/games", handlers.GetGames)

	// Application server
	router.GET("/v1/app/activate", handlers.ApplicationActivate)
	router.GET("/v1/app/refresh", handlers.ApplicationRefresh)

	// Auth server
	router.GET("/v1/oauth/authorize", handlers.CreateToken)
	router.GET("/v1/oauth/access", handlers.GetToken)
	router.GET("/v1/oauth/refresh", handlers.RefreshToken)

	// TODO delete test
	router.GET("/test", database.Test)
	router.GET("/test2", database.Test2)
}

func main() {
	fasthttp.ListenAndServe(":5555", router.Handler)
}
