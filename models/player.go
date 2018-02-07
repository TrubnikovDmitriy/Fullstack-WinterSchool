package models

import (
	"../services"
	"github.com/valyala/fasthttp"
	"encoding/json"
	"github.com/satori/go.uuid"
)


type Player struct {

	ID       uuid.UUID `json:"-"`
	PersonID uuid.UUID `json:"-"`
	Nickname string    `json:"nickname"`

	TeamID   uuid.UUID `json:"-"`
	TeamName string    `json:"team_name"`
	Retire   bool      `json:"retire"`

	Links     []Link  `json:"href,omitempty"`
}


func (player *Player)  Validate() *serv.ErrorCode {

	err := fieldLengthValidate(player.Nickname, "nickname")
	if err != nil {
		return err
	}
	return nil
}

func (player *Player) GenerateLinks() {

	href := serv.GetConfig().Href

	player.Links = append(player.Links, Link {
		Rel: "Команда",
		Href: href + "/teams/" + player.TeamID.String(),
		Action: "GET",
	})

	player.Links = append(player.Links, Link {
		Rel: "Состав команды",
		Href: href + "/teams/" + player.TeamID.String() + "/players",
		Action: "GET",
	})
	player.Links = append(player.Links, Link {
		Rel: "Страница игрока",
		Href: href + "/teams/" + player.TeamID.String() +
			"/players/" + player.ID.String(),
		Action: "GET",
	})
}

func (player *Player) WriteAsJsonResponseTo(ctx *fasthttp.RequestCtx, statusCode int) {
	player.GenerateLinks()
	resp, _ := json.Marshal(player)
	ctx.Write(resp)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetStatusCode(statusCode)
}
