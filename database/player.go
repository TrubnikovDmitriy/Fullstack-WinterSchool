package database

import (
	"../models"
	"../services"
	"github.com/valyala/fasthttp"
	"log"
)

func GetPlayerByIDs(teamID int, playerID int) (*models.Player, *services.ErrorCode) {

	const selectPlayerByID =
		"SELECT first_name, last_name, about, team_name FROM players WHERE id = $1;"

	db := sharedKeyForReadByTeamID(teamID)
	player := models.Player{ID: playerID, TeamID: teamID}

	err := db.QueryRow(selectPlayerByID, playerID).
		Scan(&player.FirstName, &player.LastName, &player.About, &player.TeamName)
	if err != nil {
		return nil, checkError(err)
	}

	return &player, nil
}

func GetPlayersOfTeam(teamID int) ([]*models.Player, *services.ErrorCode) {

	const getPlayersByTeamID =
		"SELECT id, first_name, last_name, team_name FROM players WHERE team_id = $1;"

	db := sharedKeyForReadByTeamID(teamID)
	rows, err := db.Query(getPlayersByTeamID, teamID)
	if err != nil {
		return nil, checkError(err)
	}

	var players []*models.Player
	for rows.Next() {
		player := models.Player{TeamID: teamID}
		err = rows.Scan(&player.ID, &player.FirstName, &player.LastName, &player.TeamName)
		if err != nil {
			return nil, checkError(err)
		}
		players = append(players, &player)
	}

	return players, nil
}

func CreatePlayer(player *models.Player) *services.ErrorCode {

	// Валидация
	if !player.Validate() {
		return &services.ErrorCode{
			Code: fasthttp.StatusBadRequest,
			Message: "Request is not valid",
			Link: "TODO ссылку на документацию к API",
		}
	}

	db := sharedKeyForReadByTeamID(player.TeamID)
	const checkTeamExisting = "SELECT team_name FROM teams WHERE id = $1"

	err := db.QueryRow(checkTeamExisting, player.TeamID).Scan(&player.TeamName)
	if err != nil {
		return checkError(err)
	}

	player.ID = getID("SELECT nextval('players_id_seq') FROM generate_series(0, 0);")
	const createPlayer =
		"INSERT INTO players(id, first_name, last_name, about, team_id, team_name) " +
		"VALUES ($1, $2, $3, $4, $5, $6);"

	db = sharedKeyForWriteByID(player.TeamID)
	_, err = db.Exec(createPlayer, player.ID, player.FirstName,
		player.LastName, player.About, player.TeamID, player.TeamName)

	if err != nil {
		log.Print(err)
		return &services.ErrorCode{
			Code: fasthttp.StatusInternalServerError,
			Message: "Request is not valid",
		}
	}

	return nil
}
