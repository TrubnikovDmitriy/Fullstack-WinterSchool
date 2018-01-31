package models

import (
	"../services"
	"strconv"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

type Player struct {
	ID int `json:"-"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	About string `json:"about,omitempty"`
	TeamID int `json:"-"`
	TeamName string `json:"team_name"`
	Links []*Link `json:"href,omitempty"`
}

type PlayerTest struct {
	//ID int `json:"-"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	About string `json:"about,omitempty"`
	//TeamID int `json:"-"`
	//TeamName string `json:"team_name"`
	//Links []*Link `json:"href,omitempty"`
}

func (player *Player) Validate() bool {
	if len(player.FirstName) == 0 {
		return false
	}
	if len(player.LastName) == 0 {
		return false
	}
	if len(player.FirstName) > 50 {
		return false
	}
	if len(player.LastName) > 50 {
		return false
	}
	return true
}

func (player *Player) GenerateLinks() {
	if player.TeamID != 0 {
		teamLink := &Link{
			Rel: "Команда",
			Href: services.Href + "/teams/" + strconv.Itoa(player.TeamID),
			Action: "GET",
		}
		player.Links = append(player.Links, teamLink)

		teamPlayersLink := &Link{
			Rel: "Состав команды",
			Href: services.Href + "/teams/" + strconv.Itoa(player.TeamID) + "/players",
			Action: "GET",
		}
		player.Links = append(player.Links, teamPlayersLink)
	}
	if player.TeamID != 0 && player.ID != 0 {
		playerLink := &Link{
			Rel: "Страница игрока",
			Href: services.Href + "/teams/" + strconv.Itoa(player.TeamID) +
				"/players/" + strconv.Itoa(player.ID),
			Action: "GET",
		}
		player.Links = append(player.Links, playerLink)
	}
}

func (player *Player) WriteAsJsonResponseTo(ctx *fasthttp.RequestCtx, statusCode int) {
	player.GenerateLinks()
	resp, _ := json.Marshal(player)
	ctx.Write(resp)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetStatusCode(statusCode)
}
