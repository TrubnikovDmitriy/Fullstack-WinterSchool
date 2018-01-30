package handlers

import (
	"../database"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

// GET /v1/tourney/{id}
func GetTournamentByID(ctx *fasthttp.RequestCtx) {

	id := ctx.UserValue("id").(string)
	tourney, err := database.GetTournByID(id)

	if err != nil {
		ctx.SetStatusCode(err.Code)
		return
	}

	ctx.SetStatusCode(200)
	setHeaders(ctx)

	resp, _ := json.Marshal(tourney)
	ctx.Write(resp)
}