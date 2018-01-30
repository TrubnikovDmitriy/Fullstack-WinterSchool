package models

import "time"

type Tournament struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Started time.Time `json:"started"`
	Ended time.Time	`json:"ended"`
	About string `json:"about"`
}
