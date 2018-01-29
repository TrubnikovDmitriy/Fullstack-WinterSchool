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

func GetPlayersOfTeam(id string) ([]*models.Player, error) {

	const getPlayersByTeamID =
		"SELECT team_id, id, first_name, last_name FROM players WHERE team_id = $1;"
	rows, err := conn.Query(getPlayersByTeamID, id)
	if err != nil {
		log.Fatal(err)
	}

	var posts []*models.Player
	for rows.Next() {
		post := models.Player{}
		err = rows.Scan(&post.TeamID, &post.ID, &post.FirstName, &post.LastName)
		if err != nil {
			log.Fatal(err)
		}
		posts = append(posts, &post)
	}

	return posts, nil
}