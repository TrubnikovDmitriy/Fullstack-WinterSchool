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
	Link string `json:"href,omitempty"`
}

func (game *Game) Validation() bool {
	if len(game.Title) == 0 {
		return false
	}
	if len(game.Title) > 50 {
		return false
	}
	return true
}

func (game *Game) GenerateLink() {
	if game.ID != 0 {
		game.Link = services.Href + "/games/" + strconv.Itoa(game.ID)
	}
}

func (game *Game) WriteAsJsonTo(ctx *fasthttp.RequestCtx) {
	resp, _ := json.Marshal(game)
	ctx.Write(resp)
	ctx.SetContentType("application/json; charset=utf-8")
}
