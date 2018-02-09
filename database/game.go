package database

import (
	"../models"
	"../services"
	"github.com/valyala/fasthttp"
	"github.com/satori/go.uuid"
	"log"
	"github.com/jackc/pgx/pgtype"
	"sort"
)


func GetGameByID(id uuid.UUID) (*models.Game, *serv.ErrorCode) {

	const selectGameByID = "SelectGameByID"
	db := sharedKeyForReadByID(id)
	db.Prepare(selectGameByID, "SELECT title, about FROM games WHERE id = $1")

	game := models.Game{ID: id}
	err := db.QueryRow(selectGameByID, id).Scan(&game.Title, &game.About)
	if err != nil {
		return nil, checkError(err)
	}

	return &game, nil
}

func GetGames(limit int, page int) (*models.Games, *serv.ErrorCode)  {

	if limit < 0 || 12 < limit {
		limit = 6
	}
	if page < 0 {
		page = 1
	}

	offset := limit * (page - 1)


	const prepareStatement = "SelectAllGames"
	games := make(models.Games, 0, serv.GetConfig().NumberOfShards * limit / 2)

	for i := 0; i < serv.GetConfig().NumberOfShards; i++ {

		db := choiceMasterSlave(masterConnectionPool[i], slaveConnectionPool[i])
		db.Prepare(prepareStatement, "SELECT id, title, about " +
			"FROM games ORDER BY title LIMIT $1 OFFSET $2")

		rows, err := db.Query(prepareStatement, limit, offset)

		if err != nil {
			log.Print(err)
			rows.Close()
			return nil, serv.NewServerError(err)
		}

		var gameID pgtype.UUID

		for rows.Next() {
			game := models.Game{}
			err = rows.Scan(&gameID, &game.Title, &game.About)
			if err != nil {
				log.Printf("%s : %s", prepareStatement, err)
				rows.Close()
				continue
			}

			game.ID = castUUID(gameID)
			games = append(games, game)
		}
		rows.Close()
	}

	sort.Sort(games)
	if limit > games.Len() {
		limit = games.Len()
	}
	returningGames := games[:limit]

	return &returningGames, nil
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
			Message: "Game with the same title already exists",
			Link: serv.GetConfig().Href + "/games/" + existingID.String(),
		}
	}


	// Генерация ID и ключ шардирования
	game.ID = getID()
	master := sharedKeyForWriteByID(game.ID)

	// Добавление
	const createNewGame = "CreateGame"
	master.Prepare(createNewGame, "INSERT INTO games(id, title, about) VALUES ($1, $2, $3);")

	_, err = master.Exec(createNewGame, game.ID, game.Title, game.About)
	if err != nil {
		return serv.NewServerError(err)
	}

	return nil
}
