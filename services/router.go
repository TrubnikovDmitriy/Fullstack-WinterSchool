package services

import (
	"github.com/buaazp/fasthttprouter"

	"../handlers"
)


func InitRouter() *fasthttprouter.Router {

	router := fasthttprouter.New()

	// team
	router.GET("/v1/team/:id", handlers.GetTeam)

	// player
	router.GET("/v1/player/:id", handlers.GetPlayer)

	return router
}
