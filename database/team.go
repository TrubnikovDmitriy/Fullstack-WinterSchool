package database

import (
	"../models"
	"../services"
	"github.com/valyala/fasthttp"
	"strconv"
)


func GetTeamByID(teamID int) (*models.Team, *services.ErrorCode) {

	const selectTeamByID = "SELECT team_name, about FROM teams WHERE id = $1;"

	db := sharedKeyForReadByTeamID(teamID)
	team := models.Team{ID: teamID}

	err := db.QueryRow(selectTeamByID, teamID).Scan(&team.Name, &team.About)
	if err != nil {
		return nil, checkError(err)
	}

	return &team, nil
}

func CreateTeam(team *models.Team) *services.ErrorCode {

	if !team.Validate() {
		return services.NewBadRequest()
	}

	// Проверка на уникальность имени
	const findTheSameTeamName = "SELECT id FROM teams WHERE team_name = $1";
	var existingID int

	if master1.QueryRow(findTheSameTeamName, team.Name).Scan(&existingID) == nil ||
		master2.QueryRow(findTheSameTeamName, team.Name).Scan(&existingID) == nil {
		return &services.ErrorCode{
			Code: fasthttp.StatusConflict,
			Message: "Team with the same name already exist",
			Link: services.Href + "/teams/" + strconv.Itoa(existingID),
		}
	}

	// Генерация ID
	team.ID = getID("SELECT nextval('teams_id_seq') FROM generate_series(0, 0);")
	// Ключ шардирования
	master := sharedKeyForWriteByTeamID(team.ID)

	// Добавление
	const createTeam =
		"INSERT INTO teams(id, team_name, about) VALUES ($1, $2, $3);"

	_, err := master.Exec(createTeam, team.ID, team.Name, team.About)
	if err != nil {
		return services.NewServerError()
	}

	return nil
}
