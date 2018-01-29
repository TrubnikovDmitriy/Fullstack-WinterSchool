package models


type Player struct {
	ID int `json:"id,omitempty"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	About string `json:"about,omitempty"`
}
