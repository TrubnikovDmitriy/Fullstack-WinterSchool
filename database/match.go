package database

import (
	"github.com/satori/go.uuid"
	"context"
	"github.com/jackc/pgx/pgtype"
	"../services"
	"../models"
)

func GetMatchByID(tourneyID uuid.UUID, matchID uuid.UUID) (*models.Match, *serv.ErrorCode) {

	const selectMatchByID = "SelectMatchByID"
	db := sharedKeyForReadByID(tourneyID)
	db.Prepare(selectMatchByID,
		"SELECT team_id_1, team_id_2, " +
			"team_score_1, team_score_2, " +
			"start_time, end_time, link, " +
			"prev_match_id_1, prev_match_id_2, next_match_id " +
			"FROM matches WHERE id = $1")


	match := models.Match{ ID: matchID, TourneyID: tourneyID }
	pgtypeUUID := [5]pgtype.UUID{}
	commonUUID := [5]*uuid.UUID{}

	row := db.QueryRow(selectMatchByID, matchID)
	err := row.Scan(&pgtypeUUID[0], &pgtypeUUID[1],
					&match.FirstTeamScore, &match.SecondTeamScore,
					&match.StartTime, &match.EndTime, &match.Link,
					&pgtypeUUID[2], &pgtypeUUID[3], &pgtypeUUID[4])
	if err != nil {
		return nil, checkError(err)
	} else {
		for i, pgUUID := range pgtypeUUID {
			if pgUUID.Status != pgtype.Null {
				var temp uuid.UUID
				temp, _ = uuid.FromBytes(pgUUID.Bytes[:])
				commonUUID[i] = &temp
			}
		}
		match.FirstTeamID = commonUUID[0]
		match.SecondTeamID = commonUUID[1]
		match.PrevMatch1 = commonUUID[2]
		match.PrevMatch2 = commonUUID[3]
		match.NextMatch = commonUUID[4]
	}

	return &match, nil
}

func CreateMatches(matches []models.Match, tourney *models.Tournament) *serv.ErrorCode {

	master := sharedKeyForWriteByID(tourney.ID)

	const prepareInsert = "insertMatches"
	master.Prepare(prepareInsert,
		"INSERT INTO matches" +
			"(id, tourn_id, prev_match_id_1, " +
			"prev_match_id_2, next_match_id, " +
			"start_time, organize_id) " +
			"VALUES($1, $2, $3, $4, $5, $6, $7);")
	batch := master.BeginBatch()
	defer batch.Close()

	for _, match := range matches {
		batch.Queue(prepareInsert, []interface{}{
			match.ID,
			tourney.ID,
			match.PrevMatch1,
			match.PrevMatch2,
			match.NextMatch,
			match.StartTime,
			tourney.OrganizeID,
		}, nil, nil)
	}

	err := batch.Send(context.Background(), nil)
	if err != nil {
		return checkError(err)
	}

	_, err = batch.ExecResults()
	if err != nil {
		return checkError(err)
	}
	return nil
}

func GetTournamentGrid(id uuid.UUID) (*models.MatchesArrayForm, *serv.ErrorCode) {

	const selectGridByTourneyID = "SelectGridByTourneyID"

	db := sharedKeyForReadByID(id)
	db.Prepare(selectGridByTourneyID,
		"SELECT id, team_id_1, team_id_2, " +
			"team_score_1, team_score_2, " +
			"start_time, end_time, link," +
			"prev_match_id_1, prev_match_id_2, " +
			"next_match_id " +
			"FROM matches WHERE tourn_id = $1")

	var grid models.MatchesArrayForm

	rows, err := db.Query(selectGridByTourneyID, id)
	defer rows.Close()
	if err != nil {
		return nil, checkError(err)
	}

	for rows.Next() {
		pgtypeUUID := [3]pgtype.UUID{}
		commonUUID := [3]*uuid.UUID{}
		m := models.Match{TourneyID: id}
		err = rows.Scan(
			&m.ID, &m.FirstTeamID, &m.SecondTeamID,
			&m.FirstTeamScore, &m.SecondTeamScore,
			&m.StartTime, &m.EndTime, &m.Link,
			&pgtypeUUID[0], &pgtypeUUID[1], &pgtypeUUID[2],
		)
		if err != nil {
			return nil, checkError(err)
		}

		for i, pgUUID := range pgtypeUUID {
			if pgUUID.Status != pgtype.Null {
				var temp uuid.UUID
				temp, _ = uuid.FromBytes(pgUUID.Bytes[:])
				commonUUID[i] = &temp
			}
		}
		m.PrevMatch1 = commonUUID[0]
		m.PrevMatch2 = commonUUID[1]
		m.NextMatch = commonUUID[2]

		grid.Array = append(grid.Array, m)
	}

	return &grid, nil;
}