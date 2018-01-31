package models

import (
	"../services"
	"time"
)

type Match struct {
	ID int `json:"id"`
	FirstTeamID *int `json:"first_team_id,omitempty"`
	SecondTeamID *int `json:"second_team_id,omitempty"`
	FirstTeamScore int `json:"first_team_score"`
	SecondTeamScore int `json:"second_team_score"`
	Link string `json:"link"`
	StartTime *time.Time `json:"start_time"`
	EndTime *time.Time `json:"end_time,omitempty"`
	TourneyID int `json:"-"`
	PrevMatch1 *int `json:"prev_match_id_1,omitempty"`
	PrevMatch2 *int `json:"prev_match_id_2,omitempty"`
	NextMatch *int `json:"next_match_id,omitempty"`
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




type MatchTreeForm struct {
	LeftChild *MatchTreeForm `json:"prev_match_1"`
	RightChild *MatchTreeForm `json:"prev_match_2"`
	StartTime time.Time `json:"start_time"`
}

func (match *MatchTreeForm) Validate() bool {
	ttl := services.MaxMatchesInTournament
	return match.recursiveValidate(match.StartTime, &ttl)
}

// Контролируемая рекурсия, чтобы проверить
// не начались ли вышестоящие турниры раньше по времени
func (match *MatchTreeForm) recursiveValidate(parentTime time.Time, ttl *int) bool {
	if *ttl == 0 {
		return false
	}
	if *ttl != services.MaxMatchesInTournament {
		if match.StartTime.After(parentTime) || match.StartTime.Equal(parentTime) {
			return false
		}
	}
	*ttl--
	if (match.LeftChild == nil) && (match.RightChild == nil) {
		return true
	}
	if match.LeftChild == match.RightChild {
		return false
	}
	if (match.LeftChild != nil) && (match.RightChild == nil) {
		return false
	}
	if (match.LeftChild == nil) && (match.RightChild != nil) {
		return false
	}
	flag := match.LeftChild.recursiveValidate(match.StartTime, ttl)
	if !flag {
		return false
	}
	flag = match.RightChild.recursiveValidate(match.StartTime, ttl)
	return flag
}