package database

import (
	"../models"
	"log"
)

func GetPlayerByID(id string) *models.Player {

	const SelectPlayerByID =
		"SELECT id, first_name, last_name, about FROM players WHERE id = $1;"

	player := models.Player{}

	err := conn.QueryRow(SelectPlayerByID, id).
		Scan(&player.ID, &player.FirstName, &player.LastName, &player.About)

	if err != nil {
		log.Fatal(err)
	}

	return &player
}
