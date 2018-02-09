package models

import (
	"../services"
	"github.com/valyala/fasthttp"
	"encoding/json"
	"github.com/satori/go.uuid"
)

type Game struct {
	ID uuid.UUID `json:"id"`
	Title string `json:"title"`
	About string `json:"about"`

	OrganizeID uuid.UUID `json:"organize_id"`

	Links []Link `json:"href,omitempty"`
}

func (game *Game) Validate() *serv.ErrorCode {
	err := fieldLengthValidate(game.Title, "title")
	if err != nil {
		return err
	}

	err = fieldLengthValidate(game.About, "about-field")
	if err != nil {
		return err
	}
	return nil
}

func (game *Game) GenerateLinks() {
		game.Links = append(game.Links, Link {
		Rel: "game",
		Href: serv.GetConfig().Href + "/games/" + game.ID.String(),
		Action: "GET",
	})
}

func (game *Game) WriteAsJsonResponseTo(ctx *fasthttp.RequestCtx, statusCode int) {
	game.GenerateLinks()
	resp, _ := json.Marshal(game)
	ctx.Write(resp)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetStatusCode(statusCode)
}



type Games []Game

func (games Games) Len() int {
	return len(games)
}

func (games Games) Swap(i, j int) {
	games[i], games[j] = games[j], games[i]
}

func (games Games) Less(i, j int) bool {
	return games[i].Title < games[j].Title
}

func (games Games) GenerateLinks() {
	for i := range games {
		games[i].GenerateLinks()
	}
}

func (games Games) WriteAsJsonResponseTo(ctx *fasthttp.RequestCtx, statusCode int) {
	games.GenerateLinks()
	resp, _ := json.Marshal(games)
	ctx.Write(resp)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetStatusCode(statusCode)
}
