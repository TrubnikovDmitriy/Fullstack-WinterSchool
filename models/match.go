package models

import (
	"../services"
	"time"
	"github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

type Match struct {

	ID uuid.UUID `json:"-"`

	FirstTeamID     *uuid.UUID `json:"first_team_id"`
	SecondTeamID    *uuid.UUID `json:"second_team_id"`
	FirstTeamScore  int        `json:"first_team_score"`
	SecondTeamScore int        `json:"second_team_score"`

	Link       string     `json:"link"`
	StartTime  *time.Time `json:"start_time"`
	EndTime    *time.Time `json:"end_time"`

	TourneyID  uuid.UUID  `json:"-"`
	PrevMatch1 *uuid.UUID `json:"prev_match_id_1"`
	PrevMatch2 *uuid.UUID `json:"prev_match_id_2"`
	NextMatch  *uuid.UUID `json:"next_match_id"`

	Links      []Link    `json:"href"`
}

func (match *Match) Validate() *serv.ErrorCode {

	return nil
}

func (match *Match) GenerateLinks() {

	href := serv.GetConfig().Href

	match.Links = append(match.Links, Link {
		Rel: "Турнирная сетка",
		Href: href + "/tourney/" + match.TourneyID.String() + "/matches",
		Action: "GET",
	})

	match.Links = append(match.Links, Link {
		Rel: "Ключевые события матча",
		Href: href + "/tourney/" + match.TourneyID.String() +
				"/matches/" + match.ID.String() + "/timeline",
		Action: "GET",
	})

	match.Links = append(match.Links, Link {
		Rel: "Оставить комментарий",
		Href: href + "/tourney/" + match.TourneyID.String() +
				"/matches/" + match.ID.String() + "/timeline",
		Action: "POST",
	})

	if match.FirstTeamID != nil {
		match.Links = append(match.Links, Link{
			Rel:    "Ссылка на команду №1",
			Href:   href + "/teams/" + match.FirstTeamID.String(),
			Action: "GET",
		})
	}
	if match.SecondTeamID != nil {
		match.Links = append(match.Links, Link{
			Rel:    "Ссылка на команду №2",
			Href:  	href + "/teams/" + match.SecondTeamID.String(),
			Action: "GET",
		})
	}

	if match.PrevMatch1 != nil {
		match.Links = append(match.Links, Link{
			Rel:    "Ссылка на предыдущий матч №1",
			Href:   href + "/tourney/" + match.TourneyID.String() +
					"/matches/" + match.PrevMatch1.String(),
			Action: "GET",
		})
	}
	if match.PrevMatch2 != nil {
		match.Links = append(match.Links, Link{
			Rel:    "Ссылка на предыдущий матч  №2",
			Href:   href + "/tourney/" + match.TourneyID.String() +
					"/matches/" + match.PrevMatch1.String(),
			Action: "GET",
		})
	}
	if match.NextMatch != nil {
		match.Links = append(match.Links, Link{
			Rel:    "Ссылка на следующий матч",
			Href:   href + "/tourney/" + match.TourneyID.String() +
					"/matches/" + match.NextMatch.String(),
			Action: "GET",
		})
	}
}

func (match *Match) WriteAsJsonResponseTo(ctx *fasthttp.RequestCtx, statusCode int) {
	match.GenerateLinks()
	resp, _ := json.Marshal(match)
	ctx.Write(resp)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetStatusCode(statusCode)
}

func (match *Match) Equal(other *Match) bool {
	return match.ID != other.ID &&
		match.TourneyID != other.TourneyID &&
		match.FirstTeamScore != other.FirstTeamScore &&
		match.FirstTeamID != other.FirstTeamID &&
		match.SecondTeamScore != other.SecondTeamScore &&
		match.SecondTeamID != other.SecondTeamID
}