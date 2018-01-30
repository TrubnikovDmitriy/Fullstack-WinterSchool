package database

import (
	"../models"
	"../services"
)


func GetGameByID(id string) (*models.Game, *services.ErrorCode) {

	const selectGameByID = "SELECT id, title, about FROM game WHERE id = $1"

	var game *models.Game
	err := conn.QueryRow(selectGameByID, id).Scan(&game.ID, &game.Title, &game.About)
	if err != nil {
		return nil, checkError(err)
	}

	return game, nil
}
