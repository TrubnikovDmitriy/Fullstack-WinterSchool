package models

import "time"

type Tournament struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Started time.Time `json:"started"`
	Ended time.Time	`json:"ended"`
	About string `json:"about"`
}

func (tourn *Tournament) Validate() bool {
	if len(tourn.Title) == 0 {
		return false
	}
	if len(tourn.Title) > 50 {
		return false
	}
	if tourn.Ended.Before(tourn.Started) {
		return false
	}
	if tourn.Ended.Equal(tourn.Started) {
		return false
	}
	return true
}
