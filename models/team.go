package models


type Team struct {
	ID int `json:"id"`
	Name string `json:"name"`
	About string `json:"About"`
}

func (team *Team) Validate() bool {
	if len(team.Name) == 0 {
		return false
	}
	if len(team.Name) > 100 {
		return false
	}
	return true
}
