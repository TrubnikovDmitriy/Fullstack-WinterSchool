package database

import (
	"../models"
	"../services"
	"github.com/valyala/fasthttp"
	"strconv"
	"github.com/satori/go.uuid"
)


func GetTeamByID(teamID uuid.UUID) (*models.Team, *serv.ErrorCode) {

	db := sharedKeyForReadByUUID(teamID)
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

	if !team.Validate() {
		return serv.NewBadRequest("Bad request")
	}

	// Проверка на уникальность имени
	var existingID int
	err := verifyUnique(
		"SELECT id FROM teams WHERE team_name = $1", &existingID, team.Name)
	if err != nil {
		return &serv.ErrorCode{
			Code: fasthttp.StatusConflict,
			Message: "Team with the same name already exist",
			Link: serv.Href + "/teams/" + strconv.Itoa(existingID),
		}
	}

	// Генерация ID и шардирование
	team.ID = getUUID()
	const createTeam = "CreateTeam"
	master := sharedKeyForWriteByUUID(team.ID)
	master.Prepare(createTeam,
		"INSERT INTO teams(id, team_name, about) VALUES ($1, $2, $3);")
	

	// Добавление
	_, err = master.Exec(createTeam, team.ID, team.Name, team.About)
	if err != nil {
		return serv.NewServerError()
	}

	return nil
}
