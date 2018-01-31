package models

import (
	"../services"
	"strconv"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

type Team struct {
	ID int `json:"-"`
	Name string `json:"team_name"`
	About string `json:"about"`
	Links []*Link `json:"href"`
}

func (team *Team) Validate() bool {
	if len(team.Name) == 0 {
		return false
	}
	if len(team.Name) > services.MaxFieldLength {
		return false
	}
	return true
}

func (team *Team) GenerateLinks() {
	if team.ID != 0 {
		teamLink := &Link {
			Rel: "Страница команды",
			Href: services.Href + "/teams/" + strconv.Itoa(team.ID),
			Action: "GET",
		}
		team.Links = append(team.Links, teamLink)

		teamPlayersLink := &Link {
			Rel: "Состав команды",
			Href: services.Href + "/teams/" + strconv.Itoa(team.ID) + "/players",
			Action: "GET",
		}
		team.Links = append(team.Links, teamPlayersLink)
	}
}

func (team *Team) WriteAsJsonResponseTo(ctx *fasthttp.RequestCtx, statusCode int) {
	team.GenerateLinks()
	resp, _ := json.Marshal(team)
	ctx.Write(resp)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetStatusCode(statusCode)
}