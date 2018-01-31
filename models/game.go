package models

import (
	"../services"
	"strconv"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

type Game struct {
	ID int `json:"id"`
	Title string `json:"title"`
	About string `json:"about"`
	Links []*Link `json:"href,omitempty"`
}

func (game *Game) Validate() bool {
	if len(game.Title) == 0 {
		return false
	}
	if len(game.Title) > 50 {
		return false
	}
	return true
}

func (game *Game) GenerateLinks() {
	if game.ID != 0 {
		link := &Link{
			Rel: "game",
			Href: services.Href + "/games/" + strconv.Itoa(game.ID),
			Action: "GET",
		}
		game.Links = append(game.Links, link)
	}
}

func (game *Game) WriteAsJsonResponseTo(ctx *fasthttp.RequestCtx, statusCode int) {
	game.GenerateLinks()
	resp, _ := json.Marshal(game)
	ctx.Write(resp)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetStatusCode(statusCode)
}
