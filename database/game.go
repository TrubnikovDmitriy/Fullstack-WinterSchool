package database

import (
	"../models"
	"../services"
	"github.com/valyala/fasthttp"
	"strconv"
	"github.com/jackc/pgx"
)


func GetGameByID(id string) (*models.Game, *services.ErrorCode) {

	const selectGameByID = "SELECT id, title, about FROM games WHERE id = $1"
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, &services.ErrorCode{
			Code: fasthttp.StatusBadRequest,
			Message: "Incorrect path variable '" + id + "'",
		}
	}

	var master *pgx.ConnPool
	if intID % 2 != 0 {
		master = master1
	} else {
		master = master2
	}

	var game models.Game
	err = master.QueryRow(selectGameByID, id).Scan(&game.ID, &game.Title, &game.About)
	if err != nil {
		return nil, checkError(err)
	}

	return &game, nil
}

func CreateNewGame(game *models.Game) *services.ErrorCode {

	// Валидация
	if !game.Validation() {
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
	const incrementGameID =
		"SELECT nextval('game_id_seq') FROM generate_series(0, 0);"
	// Обращаемся к обеим базам, для инкрементирования счетчика ID
	master1.QueryRow(incrementGameID).Scan(&game.ID)
	master2.QueryRow(incrementGameID)


	// Ключ шардирования
	var master *pgx.ConnPool
	if game.ID % 2 != 0 {
		master = master1
	} else {
		master = master2
	}


	// Добавление
	const createNewGame =
		"INSERT INTO games(id, title, about) VALUES ($1, $2, $3);"
	_, err := master.Exec(createNewGame, game.ID, game.Title, game.About)
	if err != nil {
		return &services.ErrorCode{
			Code: fasthttp.StatusBadRequest,
			Message: "Request is not valid",
		}
	}

	game.GenerateLink()
	return nil
}
