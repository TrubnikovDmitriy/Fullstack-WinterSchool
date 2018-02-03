package handlers

import (
	"../database"
	"github.com/valyala/fasthttp"
)

// GET /v1/match/{id}
func GetMatch(ctx *fasthttp.RequestCtx) {

	tourneyID := ctx.UserValue("tourneyID").(string)
	matchID := ctx.UserValue("matchID").(string)
	
	match, err := database.GetMatchByID(matchID, tourneyID)

	if err != nil {
		ctx.SetStatusCode(err.Code)
	} else {
		match.WriteAsJsonResponseTo(ctx, fasthttp.StatusOK)
	}
}
