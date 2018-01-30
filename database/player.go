package database

import (
	"../models"
	"../services"
)

func GetPlayerByID(id string) (*models.Player, *services.ErrorCode) {

	const selectPlayerByID =
		"SELECT id, first_name, last_name, about FROM players WHERE id = $1;"

	player := models.Player{}

	err := conn.QueryRow(selectPlayerByID, id).
		Scan(&player.ID, &player.FirstName, &player.LastName, &player.About)

	if err != nil {
		return nil, checkError(err)
	}

	return &player, nil
}

func GetPlayersOfTeam(id string) ([]*models.Player, *services.ErrorCode) {

	const getPlayersByTeamID =
		"SELECT team_id, id, first_name, last_name FROM players WHERE team_id = $1;"
	rows, err := conn.Query(getPlayersByTeamID, id)
	if err != nil {
		return nil, checkError(err)
	}

	var posts []*models.Player
	for rows.Next() {
		post := models.Player{}
		err = rows.Scan(&post.TeamID, &post.ID, &post.FirstName, &post.LastName)
		if err != nil {
			return nil, checkError(err)
		}
		posts = append(posts, &post)
	}

	return posts, nil
}