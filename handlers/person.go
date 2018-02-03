package handlers

import (
	"../models"
	"../database"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

// POST /v1/persons
func CreatePerson(ctx *fasthttp.RequestCtx) {

	var person models.Person
	json.Unmarshal(ctx.PostBody(), &person)

	err := database.CreatePerson(&person)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
	} else {
		person.WriteAsJsonResponseTo(ctx, fasthttp.StatusCreated)
	}
}

// GET /v1/persons/{id}
func GetPerson(ctx *fasthttp.RequestCtx) {

	id, err := getPathID(ctx.UserValue("person_id").(string))
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
		return
	}

	person, err := database.GetPerson(id)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
	} else {
		person.WriteAsJsonResponseTo(ctx, fasthttp.StatusOK)
	}
}
