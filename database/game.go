package database

import (
	"../models"
	"../services"
	"github.com/valyala/fasthttp"
	"github.com/satori/go.uuid"
)


func GetGameByID(id uuid.UUID) (*models.Game, *serv.ErrorCode) {

	const selectGameByID = "SelectGameByID"
	db := sharedKeyForReadByUUID(id)
	db.Prepare(selectGameByID, "SELECT title, about FROM games WHERE id = $1")

	game := models.Game{ID: id}
	err := db.QueryRow(selectGameByID, id).Scan(&game.Title, &game.About)
	if err != nil {
		return nil, checkError(err)
	}

	return &game, nil
}

func CreateGame(game *models.Game) *serv.ErrorCode {

	// Валидация
	errorCode := game.Validate()
	if errorCode != nil {
		return errorCode
	}

	// Проверка на уникальность имени
	var existingID uuid.UUID
	err := verifyUnique("SELECT id FROM games WHERE title = $1", &existingID, game.Title)
	if err != nil {
		return &serv.ErrorCode{
			Code: fasthttp.StatusConflict,
			Message: "Game with the same title already exist",
			Link: serv.Href + "/games/" + existingID.String(),
		}
	}


	// Генерация ID и ключ шардирования
	game.ID = getUUID()
	master := sharedKeyForWriteByUUID(game.ID)

	// Добавление
	const createNewGame =
		"INSERT INTO games(id, title, about) VALUES ($1, $2, $3);"

	_, err = master.Exec(createNewGame, game.ID, game.Title, game.About)
	if err != nil {
		return serv.NewServerError()
	}

	return nil
}
