package database

import (
	"../models"
	"../errors"
)


func GetTeamByID(id string) (*models.Team, *errors.ErrorCode) {

	const SelectTeamByID = "SELECT id, team_name, about FROM teams WHERE id = $1;"

	team := models.Team{}
	err := conn.QueryRow(SelectTeamByID, id).Scan(&team.ID, &team.Name, &team.About)
	if err != nil {
		return nil, checkError(err)
	}

	return &team, nil
}
