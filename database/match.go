package database

import (
	"../models"
	"../errors"
)

func GetMatchByID(id string) (*models.Match, *errors.ErrorCode) {

	const getMatchByID =
		"SELECT passed, id, first_team_id, second_team_id, first_team_score, second_team_score, " +
			"link, start_time, end_time FROM matches WHERE id = $1"
	var match models.Match

	row := conn.QueryRow(getMatchByID, id)
	err := row.Scan(&match.Passed, &match.ID, &match.FirstTeamID, &match.SecondTeamID,
					&match.FirstTeamScore, &match.SecondTeamScore,
					&match.Link, &match.StartTime, &match.EndTime)
	if err != nil {
		return nil, checkError(err)
	}

	return &match, nil
}