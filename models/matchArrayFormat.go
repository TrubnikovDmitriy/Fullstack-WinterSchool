package models

import (
	"../services"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

type MatchesArrayForm struct {
	Array []Match `json:"grid"`
}


func (array *MatchesArrayForm) GenerateLinks() {
	for i, match := range array.Array {
		matchLink := Link{
			Rel: "Матч",
			Href: serv.GetConfig().Href + "/tourney/" + match.TourneyID.String() +
							"/matches/" + match.ID.String(),
			Action: "GET",
		}
		array.Array[i].Links = append(array.Array[i].Links, matchLink)
	}
}

func (array *MatchesArrayForm) WriteAsJsonResponseTo(ctx *fasthttp.RequestCtx, statusCode int) {
	array.GenerateLinks()
	resp, _ := json.Marshal(array)
	ctx.Write(resp)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetStatusCode(statusCode)
}