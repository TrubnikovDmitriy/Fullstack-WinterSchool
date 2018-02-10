package database

import (
	"../models"
	"../services"
	"log"
	"github.com/satori/go.uuid"
	"github.com/jackc/pgx/pgtype"
)


func GetPlayerByID(teamID uuid.UUID, playerID uuid.UUID) (*models.Player, *serv.ErrorCode) {

	const selectPlayerByID = "SelectPlayerByID"
	db := sharedKeyForReadByID(teamID)
	db.Prepare(selectPlayerByID,
		"SELECT person_id, nickname, team_name, retire " +
			"FROM players WHERE id = $1;")

	player := models.Player{ID: playerID, TeamID: teamID}
	personID := pgtype.UUID{}
	err := db.QueryRow(selectPlayerByID, player.ID).
		Scan(&personID, &player.Nickname, &player.TeamName, &player.Retire)
	if err != nil {
		return nil, checkError(err)
	}
	player.PersonID = castUUID(personID)

	return &player, nil
}

func GetPlayersOfTeam(teamID uuid.UUID) ([]*models.Player, *serv.ErrorCode) {

	db := sharedKeyForReadByID(teamID)
	const getPlayersByTeamID = "GetPlayersByTeamID"
	db.Prepare(getPlayersByTeamID,
		"SELECT id, person_id, nickname, team_name " +
			"FROM players WHERE team_id = $1 AND retire = FALSE;")

	rows, err := db.Query(getPlayersByTeamID, teamID)
	defer rows.Close()

	if err != nil {
		return nil, checkError(err)
	}

	var players []*models.Player
	var playerID, personID pgtype.UUID
	for rows.Next() {

		player := models.Player{TeamID: teamID}
		err = rows.Scan(&playerID, &personID, &player.Nickname, &player.TeamName)
		if err != nil {
			log.Printf("%s : %s", getPlayersByTeamID, err)
			continue
		}

		player.Retire = false
		player.ID = castUUID(playerID)
		player.PersonID = castUUID(personID)
		player.GenerateLinks()

		players = append(players, &player)
	}

	return players, nil
}