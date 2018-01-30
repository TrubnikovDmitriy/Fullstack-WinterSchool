package models


type Player struct {
	ID int `json:"id,omitempty"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	About string `json:"about,omitempty"`
	TeamID int `json:"team_id"`
}

func (player *Player) Validation() bool {
	if len(player.FirstName) == 0 {
		return false
	}
	if len(player.LastName) == 0 {
		return false
	}
	if len(player.FirstName) > 50 {
		return false
	}
	if len(player.LastName) > 50 {
		return false
	}
	return true
}
