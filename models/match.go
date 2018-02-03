package models

import (
	"time"
	"github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

type Match struct {

	ID uuid.UUID `json:"id"`

	FirstTeamID     *int `json:"first_team_id"`
	SecondTeamID    *int `json:"second_team_id"`
	FirstTeamScore  int  `json:"first_team_score"`
	SecondTeamScore int  `json:"second_team_score"`

	Link       string     `json:"link"`
	StartTime  *time.Time `json:"start_time"`
	EndTime    *time.Time `json:"end_time"`

	TourneyID  uuid.UUID  `json:"-"`
	PrevMatch1 *uuid.UUID `json:"prev_match_id_1"`
	PrevMatch2 *uuid.UUID `json:"prev_match_id_2"`
	NextMatch  *uuid.UUID `json:"next_match_id"`

	Links      []Link    `json:"href"`
}

func (match *Match) Validate() bool {
	if match.StartTime == nil {
		return false
	}
	if (match.PrevMatch1 != nil) && (match.PrevMatch2 == nil) {
		return false
	}
	if (match.PrevMatch1 == nil) && (match.PrevMatch2 != nil) {
		return false
	}
	return true
}

func (match *Match) GenerateLinks() {
	//if game.ID != 0 {
	//	link := &Link{
	//		Rel: "game",
	//		Href: services.Href + "/games/" + strconv.Itoa(game.ID),
	//		Action: "GET",
	//	}
	//	game.Links = append(game.Links, link)
	//}
}

func (match *Match) WriteAsJsonResponseTo(ctx *fasthttp.RequestCtx, statusCode int) {
	match.GenerateLinks()
	resp, _ := json.Marshal(match)
	ctx.Write(resp)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetStatusCode(statusCode)
}