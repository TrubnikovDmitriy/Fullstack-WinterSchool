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
		Href: serv.Href + "/games/" + game.ID.String(),
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
