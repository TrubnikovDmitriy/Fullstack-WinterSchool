package database

import (
	"../models"
	"../services"
)

func GetTournByID(id string) (*models.Tournament, *services.ErrorCode) {

	const selectTournByID = "SELECT id, title, started, ended, about " +
							"FROM tournament WHERE id =$1"

	var tourney models.Tournament
	err := master1.QueryRow(selectTournByID, id).
		Scan(&tourney.ID, &tourney.Title, &tourney.Started, &tourney.Ended, &tourney.About)

	if err != nil {
		return nil, checkError(err)
	}

	return &tourney, nil;
}