package models

type Game struct {
	ID int `json:"id"`
	Title string `json:"title"`
	About string `json:"about"`
}
