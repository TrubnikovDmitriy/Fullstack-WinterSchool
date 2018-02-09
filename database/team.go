package database

import (
	"../models"
	"../services"
	"github.com/valyala/fasthttp"
	"github.com/satori/go.uuid"
	"github.com/jackc/pgx/pgtype"
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
