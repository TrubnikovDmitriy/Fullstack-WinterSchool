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
	router.DELETE("/v1/teams/:team_id/players/:player_id", handlers.DeletePlayer)
	router.GET("/v1/teams/:team_id/players", handlers.GetTeamPlayers)
	router.POST("/v1/players", handlers.CreatePlayer)

	// matches
	router.GET("/v1/tourney/:tourney_id/matches/:match_id", handlers.GetMatch)
	router.PUT("/v1/tourney/:tourney_id/matches/:match_id", handlers.UpdateMatch)

	// tournaments
	router.GET("/v1/tourney/:tourney_id", handlers.GetTournamentByID)
	router.GET("/v1/tourney/:tourney_id/matches", handlers.GetTournamentGrid)
	router.POST("/v1/tourney", handlers.CreateTournament)
	router.PUT("/v1/tourney/:tourney_id", handlers.UpdateTournament)

	// games
	router.GET("/v1/games/:game_id/tournaments", handlers.GetTournamentsByGameID)
	router.GET("/v1/games/:game_id", handlers.GetGame)
	router.GET("/v1/games", handlers.GetGames)
	router.POST("/v1/games", handlers.CreateGame)

	// Application server
	router.GET("/v1/app/activate", handlers.ApplicationActivate)
	router.GET("/v1/app/refresh", handlers.ApplicationRefresh)

	// Auth server
	router.GET("/v1/oauth/authorize", handlers.CreateToken)
	router.GET("/v1/oauth/access", handlers.GetToken)
	router.GET("/v1/oauth/refresh", handlers.RefreshToken)
}

func main() {
	fasthttp.ListenAndServe("localhost:5554", router.Handler)
}
