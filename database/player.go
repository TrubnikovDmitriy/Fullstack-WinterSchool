package database


import (
	"../models"
	"../services"
	"github.com/valyala/fasthttp"
	"log"
	"github.com/satori/go.uuid"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx"
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

func CreatePlayer(player *models.Player) *serv.ErrorCode {

	// Валидация
	errorCode := player.Validate()
	if errorCode != nil {
		return errorCode
	}


	// Проверка приглашений для данного игрока
	const checkInvite = "CheckInvite"
	db := sharedKeyForReadByID(player.PersonID)
	db.Prepare(checkInvite, "SELECT team_name FROM teams " +
		"WHERE person_id = $1 AND team_id = $2")

	err := db.QueryRow(checkInvite, player.PersonID, player.TeamID).Scan(&player.TeamName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &serv.ErrorCode{
				Code: fasthttp.StatusForbidden,
				Message: "You're not invited to the given team, check out the invite list",
				Link: serv.GetConfig().Href + "/persons/" + player.PersonID.String() + "/invite-list",
			}
		}
		log.Print(err)
		return &serv.ErrorCode{ Code:fasthttp.StatusInternalServerError }
	}


	// Генерация ID и шардирование
	player.ID = getID()
	db = sharedKeyForWriteByID(player.TeamID)
	const createPlayer = "CreatePlayer"
	db.Prepare(createPlayer,
		"INSERT INTO players(id, person_id, nickname, team_id, team_name) " +
			"VALUES ($1, $2, $3, $4, $5, $6);")


	// Создание нового игрока
	_, err = db.Exec(createPlayer, player.ID, player.PersonID,
		player.Nickname, player.TeamID, player.TeamName)
	if err != nil {
		return serv.NewServerError(err)
	}


	// Удалить инвайт
	db = sharedKeyForWriteByID(player.PersonID)
	db.Exec("DELETE FROM teams WHERE person_id = $1", player.PersonID)

	return nil
}
