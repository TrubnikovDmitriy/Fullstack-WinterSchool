package handlers

import (
	"../database"
	"github.com/valyala/fasthttp"
)

// GET /v1/tourney/{tourney_id}/matches/{match_id}
func GetMatch(ctx *fasthttp.RequestCtx) {

	tourneyID, err := getPathID(ctx.UserValue("tourney_id").(string))
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
		return
	}
	matchID, err := getPathID(ctx.UserValue("match_id").(string))
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
		return
	}

	match, err := database.GetMatchByID(tourneyID, matchID)

	if err != nil {
		ctx.SetStatusCode(err.Code)
	} else {
		match.WriteAsJsonResponseTo(ctx, fasthttp.StatusOK)
	}
}
