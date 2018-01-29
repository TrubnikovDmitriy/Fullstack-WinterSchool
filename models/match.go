package models

import "time"

type Match struct {
	ID int `json:"id"`
	Passed bool `json:"passed"`
	FirstTeamID *int `json:"first_team_id,omitempty"`
	SecondTeamID *int `json:"second_team_id,omitempty"`
	FirstTeamScore int `json:"first_team_score"`
	SecondTeamScore int `json:"second_team_score"`
	Link string `json:"link"`
	StartTime *time.Time `json:"start_time"`
	EndTime *time.Time `json:"end_time,omitempty"`
}
