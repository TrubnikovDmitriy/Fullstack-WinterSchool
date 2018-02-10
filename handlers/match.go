package handlers

import (
	"../database"
	"../models"
	"github.com/valyala/fasthttp"
	"encoding/json"
	"github.com/satori/go.uuid"
)

// GET /v1/tourney/{tourney_id}/matches/{match_id}
func GetMatch(ctx *fasthttp.RequestCtx) {

	tourneyID, err := getPathID(ctx.UserValue("tourney_id").(string))
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
		return
	}
	matchID, err := getPathID(ctx.UserValue("match_id").(string))
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
		return
	}

	match, err := database.GetMatchByID(tourneyID, matchID)

	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
	} else {
		match.WriteAsJsonResponseTo(ctx, fasthttp.StatusOK)
	}
}


// POST /v1/tourney/{tourney_id}/matches/{match_id}
func UpdateMatch(ctx *fasthttp.RequestCtx) {

	updMatch := new(models.Match)
	json.Unmarshal(ctx.PostBody(), updMatch)

	// Проверка авторизации
	claims, err := GetClaimsFromCookie(ctx)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
		return
	}
	updMatch.OrganizeID = uuid.FromStringOrNil(claims["person_id"].(string))

	// Обновление матча
	match, err := database.UpdateMatch(updMatch)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
	} else {
		match.WriteAsJsonResponseTo(ctx, fasthttp.StatusOK)
	}
}
