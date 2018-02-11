package database

import (
	"../models"
	"../services"
	"github.com/valyala/fasthttp"
	"github.com/satori/go.uuid"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx"
	"log"
)


func GetTeamByID(teamID uuid.UUID) (*models.Team, *serv.ErrorCode) {

	db := sharedKeyForReadByID(teamID)
	const selectTeamByID = "SelectTeamByID"
	db.Prepare(selectTeamByID,"SELECT team_name, about FROM teams WHERE id = $1;")

	team := models.Team{ID: teamID}
	err := db.QueryRow(selectTeamByID, teamID).Scan(&team.Name, &team.About)
	if err != nil {
		return nil, checkError(err)
	}

	return &team, nil
}

func CreateTeam(team *models.Team) *serv.ErrorCode {

	errorCode := team.Validate()
	if errorCode != nil {
		return errorCode
	}

	// Проверка на уникальность имени
	var existingID pgtype.UUID
	err := verifyUnique(
		"SELECT id FROM teams WHERE team_name = $1", &existingID, team.Name)
	if err != nil {
		return &serv.ErrorCode{
			Code: fasthttp.StatusConflict,
			Message: "Team with the same name already exist",
			Link: serv.GetConfig().Href + "/teams/" + castUUID(existingID).String(),
		}
	}

	// Генерация ID и шардирование
	team.ID = getID()
	const createTeam = "CreateTeam"
	master := sharedKeyForWriteByID(team.ID)
	master.Prepare(createTeam, "INSERT INTO " +
		"teams(id, team_name, about, coach_id, coach_name) VALUES ($1, $2, $3, $4, $5);")
	

	// Добавление
	_, err = master.Exec(createTeam, team.ID, team.Name, team.About, team.CoachID, team.CoachName)
	if err != nil {
		return serv.NewServerError(err)
	}

	return nil
}

func AddPlayerInTeam(player *models.Player) *serv.ErrorCode {

	// Валидация
	errCode := player.Validate()
	if errCode != nil {
		return errCode
	}


	// Проверка, что человек играющий за этого персонажа существует
	const checkPlayerExist = "CheckPlayerExist"
	db := sharedKeyForReadByID(player.PersonID)
	db.Prepare(checkPlayerExist, "SELECT id FROM persons WHERE id = $1")

	err := db.QueryRow(checkPlayerExist, player.PersonID).Scan(new(pgtype.UUID))
	if err == pgx.ErrNoRows {
		return serv.NewBadRequest("Such person does not exist")
	}
	if err != nil {
		log.Print(err)
		return serv.NewServerError(err)
	}


	// Шардирование по ключу UUID команды
	master := sharedKeyForWriteByID(player.TeamID)
	const checkTeamExist = "CheckTeamExist"
	const addPlayerInTeam = "AddPlayerInTeam"
	master.Prepare(checkTeamExist,
		"SELECT team_name FROM teams WHERE id = $1")
	master.Prepare(addPlayerInTeam,
		"INSERT INTO players(id, person_id, nickname, team_id, team_name) " +
			"VALUES($1, $2, $3, $4, $5);")


	// Проверка, что команда существует
	err = master.QueryRow(checkTeamExist, player.TeamID).Scan(&player.TeamName)
	if err == pgx.ErrNoRows {
		return serv.NewBadRequest("Such person does not exist")
	}
	if err != nil {
		log.Print(err)
		return serv.NewServerError(err)
	}


	// Создание нового игрока
	player.ID = getID()
	_, err = master.Exec(addPlayerInTeam, player.ID, player.PersonID,
		player.Nickname, player.TeamID, player.TeamName)
	if err != nil {
		return serv.NewServerError(err)
	}


	return nil;
}

func DeletePlayerFromTeam(player *models.Player) *serv.ErrorCode {

	if player.ID == uuid.Nil || player.TeamID == uuid.Nil {
		return serv.NewBadRequest("One of sent IDs is incorrect")
	}

	const deletePlayerFromTeam = "DeletePlayerFromTeam"
	db := sharedKeyForWriteByID(player.TeamID)
	db.Prepare(deletePlayerFromTeam, "UPDATE players SET retire = TRUE WHERE id = $1")

	_, err := db.Exec(deletePlayerFromTeam, player.ID)
	if err != nil {
		log.Print(err)
		return serv.NewServerError(err)
	}

	return nil
}
