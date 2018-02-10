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


	// Создание промежуточной таблицы game-tourney
	db = sharedKeyForWriteByID(tourney.GameID)
	const createGameTourneyRow =
		"INSERT INTO game_tourney(game_id, tourney_id, started, title) VALUES ($1, $2, $3, $4);"
	_, err = db.Exec(createGameTourneyRow, tourney.GameID, tourney.ID, tourney.Started, tourney.Title)
	if err != nil {
		log.Print(err)
		// Если вдруг что-то пошло не так
		master.Exec("DELETE FROM tournaments WHERE id = $1", tourney.ID)
		return serv.NewServerError(err)
	}


	return CreateMatches(matches, tourney)
}

func GetTournamentsByGameID(gameID uuid.UUID, page int, limit int) (*[]models.Tournament, *serv.ErrorCode) {

	if limit < 0 || 12 < limit {
		limit = 6
	}
	if page < 0 {
		page = 1
	}

	offset := limit * (page - 1)

	const getGamesByGameID = "GetGamesByGameID"
	db := sharedKeyForReadByID(gameID)
	db.Prepare(getGamesByGameID, "SELECT tourney_id, started, title " +
		"FROM game_tourney WHERE game_id = $1 " +
		"ORDER BY started DESC LIMIT $2 OFFSET $3")

	rows, err := db.Query(getGamesByGameID, gameID, limit, offset)
	defer rows.Close()
	if err != nil {
		log.Print(err)
		return nil, serv.NewServerError(err)
	}

	tourneys := make([]models.Tournament, 0, limit)
	var tourneyID pgtype.UUID
	for rows.Next() {
		tourney := models.Tournament{}
		err = rows.Scan(&tourneyID, &tourney.Started, &tourney.Title)
		if err != nil {
			log.Print(err)
			return nil, serv.NewServerError(err)
		}
		tourney.ID = castUUID(tourneyID)
		tourney.GenerateLinks()
		tourneys = append(tourneys, tourney)
	}

	return &tourneys, nil
}
