package database

import (
	"../models"
	"../services"
	"github.com/valyala/fasthttp"
	"container/heap"
	"strconv"
)

func GetTourneyByID(id string) (*models.Tournament, *services.ErrorCode) {

	const selectTourneyByID = "SELECT id, title, started, ended, about " +
							"FROM tournament WHERE id =$1"

	var tourney models.Tournament
	err := master1.QueryRow(selectTourneyByID, id).
		Scan(&tourney.ID, &tourney.Title, &tourney.Started, &tourney.Ended, &tourney.About)

	if err != nil {
		return nil, checkError(err)
	}

	return &tourney, nil;
}

func CreateTournament(tourney *models.Tournament) *services.ErrorCode {

	// Валидация
	if !tourney.Validate() {
		return &services.ErrorCode {
			Code: fasthttp.StatusBadRequest,
			Message: "Bad request body",
		}
	}


	// Проверка на уникальность имени
	const findTheSameTournamentTitle = "SELECT id FROM tournaments WHERE title = $1";
	var existingID int
	if master1.QueryRow(findTheSameTournamentTitle, tourney.Title).Scan(&existingID) == nil ||
		master2.QueryRow(findTheSameTournamentTitle, tourney.Title).Scan(&existingID) == nil {
		return &services.ErrorCode{
			Code: fasthttp.StatusConflict,
			Message: "Tournament with the same title already exist",
			Link: services.Href + "/tourney/" + strconv.Itoa(existingID),
		}
	}


	// Генерация ID
	tourney.ID = getID("SELECT nextval('game_id_seq') FROM tournament_series(0, 0);")
	// Ключ шардирования
	master := sharedKeyForWriteByID(tourney.ID)


	// Добавление
	const createNewGame =
		"INSERT INTO tournaments(id, title, started, about) VALUES ($1, $2, $3, $4);"

	_, err := master.Exec(createNewGame, tourney.ID,
			tourney.Title, tourney.Started, tourney.About)
	if err != nil {
		return services.NewServerError()
	}


	return nil
}