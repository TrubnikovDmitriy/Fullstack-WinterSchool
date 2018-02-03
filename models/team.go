package models

import (
	"../services"
	"github.com/valyala/fasthttp"
	"encoding/json"
	"github.com/satori/go.uuid"
)

type Team struct {
	ID    uuid.UUID	`json:"-"`
	Name  string  	`json:"team_name"`
	About string  	`json:"about"`
	Links []Link 	`json:"href"`
}

func (team *Team) Validate() bool {
	if len(team.Name) == 0 {
		return false
	}
	if len(team.Name) > serv.MaxFieldLength {
		return false
	}
	return true
}

func (team *Team) GenerateLinks() {

	team.Links = append(team.Links, Link {
		Rel: "Страница команды",
		Href: serv.Href + "/teams/" + team.ID.String(),
		Action: "GET",
	})

	team.Links = append(team.Links, Link {
		Rel: "Состав команды",
		Href: serv.Href + "/teams/" + team.ID.String() + "/players",
		Action: "GET",
	})
}

func (team *Team) WriteAsJsonResponseTo(ctx *fasthttp.RequestCtx, statusCode int) {
	team.GenerateLinks()
	resp, _ := json.Marshal(team)
	ctx.Write(resp)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetStatusCode(statusCode)
}