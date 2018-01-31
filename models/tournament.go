package models

import (
	"../services"
	"time"
	"strconv"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

type Tournament struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Started time.Time `json:"started"`
	Ended time.Time	`json:"ended"`
	About string `json:"about"`
	Matches []*Match `json:"matches,omitempty"`
	TeamsID []int `json:"teams_id,omitempty"`
	Links []*Link `json:"href"`
}

type TournamentCreateForm struct {

}

func (tourney *Tournament) Validate() bool {
	if len(tourney.Title) == 0 {
		return false
	}
	if len(tourney.Title) > services.MaxFieldLength {
		return false
	}
	if tourney.Ended.Before(tourney.Started) {
		return false
	}
	if tourney.Ended.Equal(tourney.Started) {
		return false
	}
	if len(tourney.Matches) > services.MaxMatchesInTournament {
		return false
	}
	if len(tourney.Matches) == 0 {
		return false
	}
	for _, value := range tourney.Matches {
		if !value.Validate() {
			return false
		}
	}
	return true
}

func (tourney *Tournament) GenerateLinks() {
	if tourney.ID != 0 {
		link := &Link{
			Rel: "game",
			Href: services.Href + "/tourney/" + strconv.Itoa(tourney.ID),
			Action: "GET",
		}
		tourney.Links = append(tourney.Links, link)
	}
}

func (tourney *Tournament) WriteAsJsonResponseTo(ctx *fasthttp.RequestCtx, statusCode int) {
	tourney.GenerateLinks()
	resp, _ := json.Marshal(tourney)
	ctx.Write(resp)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetStatusCode(statusCode)
}



