package database

import (
	"../models"
	"log"
)


func GetTeamByID(id string) *models.Team {

	const SelectTeamByID = "SELECT id, team_name, about FROM teams WHERE id = $1;"

	team := models.Team{}
	err := conn.QueryRow(SelectTeamByID, id).Scan(&team.ID, &team.Name, &team.About)
	if err != nil {
		log.Fatal(err)
	}

	return &team
}
