package models

import (
	"../services"
	"github.com/valyala/fasthttp"
	"encoding/json"
	"github.com/satori/go.uuid"
)

type Team struct {
	ID      uuid.UUID `json:"-"`
	Name    string    `json:"team_name"`
	About   string    `json:"about"`

	CoachID   uuid.UUID `json:"coach_id"`
	CoachName string    `json:"coach_name"`

	Links []Link 	`json:"href"`
}

func (team *Team) Validate() *serv.ErrorCode {

	err := fieldLengthValidate(team.Name, "team name")
	if err != nil {
		return err
	}
	err = fieldLengthValidate(team.About, "about-field")
	if err != nil {
		return err
	}
	return nil
}

func (team *Team) GenerateLinks() {

	team.Links = append(team.Links, Link {
		Rel: "Страница команды",
		Href: serv.GetConfig().Href + "/teams/" + team.ID.String(),
		Action: "GET",
	})

	team.Links = append(team.Links, Link {
		Rel: "Состав команды",
		Href: serv.GetConfig().Href + "/teams/" + team.ID.String() + "/players",
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