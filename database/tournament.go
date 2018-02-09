package database

import (
	"../models"
	"../services"
	"github.com/valyala/fasthttp"
	"github.com/satori/go.uuid"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx"
)


func GetTourneyByID(id uuid.UUID) (*models.Tournament, *serv.ErrorCode) {

	db := sharedKeyForReadByID(id)
	const selectTourneyByID = "SelectTourneyByID"
	db.Prepare(selectTourneyByID,
		"SELECT id, title, started, ended, about FROM tournaments WHERE id =$1")

	var tourney models.Tournament
	err := db.QueryRow(selectTourneyByID, id).
		Scan(&tourney.ID, &tourney.Title, &tourney.Started, &tourney.Ended, &tourney.About)

	if err != nil {
		return nil, checkError(err)
	}

	return &tourney, nil;
}

func CreateTournament(tourney *models.Tournament) *serv.ErrorCode {

	// Валидация
	errorCode := tourney.Validate()
	if errorCode != nil {
		return errorCode
	}

	// Проверка, что игра с указанным ID существует
	var existingID pgtype.UUID
	db := sharedKeyForReadByID(tourney.GameID)
	pgErr := db.QueryRow(
		"SELECT id FROM games WHERE id = $1", tourney.GameID).Scan(&existingID)
	if pgErr == pgx.ErrNoRows {
		return serv.NewBadRequest("Such game does not exist")
	}
	if pgErr != nil {
		return serv.NewServerError(pgErr)
	}


	// Проверка на уникальность имени
	const findTheSameTournamentTitle = "SELECT id FROM tournaments WHERE title = $1"
	err := verifyUnique(findTheSameTournamentTitle, &existingID, tourney.Title)
	if err != nil {
		return &serv.ErrorCode{
			Code: fasthttp.StatusConflict,
			Message: "Tournament with the same title already exist",
			Link: serv.GetConfig().Href + "/tourney/" + castUUID(existingID).String(),
		}
	}

	// Генерация UUID и ключ шардирования
	tourney.ID = getID()
	master := sharedKeyForWriteByID(tourney.ID)

	// Добавление турнира
	const createNewTournament =
		"INSERT INTO tournaments(id, title, started, ended, about, organize_id, " +
			"organize_name, game_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);"
	_, err = master.Exec(createNewTournament, tourney.ID, tourney.Title,
						tourney.Started, tourney.Ended, tourney.About,
						tourney.OrganizeID, tourney.OrganizeName, tourney.GameID)
	if err != nil {
		return checkError(err)
	}


	// Распарсить дерево матчей в массив
	// (под капотом им еще генерируются uuid и связываются между собой ссылками)
	matches := tourney.MatchTree.CreateArrayMatch()
	return CreateMatches(matches, tourney)
}
