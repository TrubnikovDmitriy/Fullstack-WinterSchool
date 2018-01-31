package database

import (
	"../models"
	"../services"
	"github.com/valyala/fasthttp"
	"strconv"
)


func GetGameByID(id int) (*models.Game, *services.ErrorCode) {

	const selectGameByID = "SELECT title, about FROM games WHERE id = $1"

	game := models.Game{ID: id}
	master := sharedKeyForReadByTeamID(id)
	err := master.QueryRow(selectGameByID, id).Scan(&game.Title, &game.About)
	if err != nil {
		return nil, checkError(err)
	}

	return &game, nil
}

func CreateGame(game *models.Game) *services.ErrorCode {

	// Валидация
	if !game.Validate() {
		return &services.ErrorCode{
			Code: fasthttp.StatusBadRequest,
			Message: "Request is not valid",
		}
	}

	// Проверка на уникальность имени
	const findTheSameGameTitle = "SELECT id FROM games WHERE title = $1";
	var existingID int

	if master1.QueryRow(findTheSameGameTitle, game.Title).Scan(&existingID) == nil ||
		master2.QueryRow(findTheSameGameTitle, game.Title).Scan(&existingID) == nil {
		return &services.ErrorCode{
			Code: fasthttp.StatusConflict,
			Message: "Game with the same title already exist",
			Link: services.Href + "/games/" + strconv.Itoa(existingID),
		}
	}


	// Генерация ID
	game.ID = getID("SELECT nextval('game_id_seq') FROM generate_series(0, 0);")


	// Ключ шардирования
	master := sharedKeyForWriteByTeamID(game.ID)


	// Добавление
	const createNewGame =
		"INSERT INTO games(id, title, about) VALUES ($1, $2, $3);"

	_, err := master.Exec(createNewGame, game.ID, game.Title, game.About)
	if err != nil {
		return services.NewServerError()
	}

	return nil
}
