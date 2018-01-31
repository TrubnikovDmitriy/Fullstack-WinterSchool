package handlers

import (
	"../database"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

// GET /v1/match/{id}
func GetMatch(ctx *fasthttp.RequestCtx) {

	tourney_id := ctx.UserValue("tourney_id").(string)
	match_id := ctx.UserValue("match_id").(string)



	match, err := database.GetMatchByID(match_id, tourney_id)

	if err != nil {
		ctx.SetStatusCode(err.Code)
		return
	}

	ctx.SetStatusCode(200)
	setHeaders(ctx)

	resp, _ := json.Marshal(match)
	ctx.Write(resp)
}
