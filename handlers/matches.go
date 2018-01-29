package handlers

import (
	"../database"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

// GET /v1/match/{id}
func GetMatch(ctx *fasthttp.RequestCtx) {

	id := ctx.UserValue("id").(string)
	match, err := database.GetMatchByID(id)

	if err != nil {
		ctx.SetStatusCode(err.Code)
		return
	}

	ctx.SetStatusCode(200)
	setHeaders(ctx)

	resp, _ := json.Marshal(match)
	ctx.Write(resp)
}
