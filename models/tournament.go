package models

import (
	"../services"
	"time"
	"github.com/valyala/fasthttp"
	"encoding/json"
	"github.com/satori/go.uuid"
)

type Tournament struct {
	ID      uuid.UUID `json:"-"`

	Title   string    `json:"title"`
	Started time.Time `json:"started"`
	Ended   time.Time `json:"ended"`
	About   string    `json:"about"`

	MatchTree *MatchesTreeForm `json:"matches_tree, omitempty"`
	Matches []*Match           `json:"matches, omitempty"`
	Links []Link               `json:"href, omitempty"`
}

func (tourney *Tournament) Validate() *serv.ErrorCode {
	if len(tourney.Title) == 0 {
		return serv.NewBadRequest("Title with zero length")
	}
	if len(tourney.Title) > serv.MaxFieldLength {
		return serv.NewBadRequest("Too long title")
	}
	if tourney.Ended.Before(tourney.Started) {
		return serv.NewBadRequest("Match ended before it begin")
	}
	if tourney.Ended.Equal(tourney.Started) {
		return serv.NewBadRequest("Start and end time of tournament is equal")
	}
	if tourney.MatchTree == nil {
		return serv.NewBadRequest("No matches")
	}
	return tourney.MatchTree.Validate()
}

func (tourney *Tournament) GenerateLinks() {

	tourney.Links = append(tourney.Links, Link {
		Rel: "Ссылка на турнир",
		Href: serv.Href + "/tourney/" + tourney.ID.String(),
		Action: "GET",
	})

	tourney.Links = append(tourney.Links, Link {
		Rel: "Турнирная сетка",
		Href: serv.Href + "/tourney/" + tourney.ID.String() + "/matches",
		Action: "GET",
	})

	tourney.Links = append(tourney.Links, Link {
		Rel: "Команды участницы",
		Href: serv.Href + "/tourney/" + tourney.ID.String() + "/teams",
		Action: "GET",
	})
}

func (tourney *Tournament) WriteAsJsonResponseTo(ctx *fasthttp.RequestCtx, statusCode int) {
	tourney.GenerateLinks()
	resp, _ := json.Marshal(tourney)
	ctx.Write(resp)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetStatusCode(statusCode)
}



