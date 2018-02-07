package models

import (
	"../services"
	"github.com/satori/go.uuid"
	"encoding/json"
	"github.com/valyala/fasthttp"
)

type Person struct {
	ID uuid.UUID `json:"-"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	About     string `json:"about"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`

	Staff bool   `json:"-"`
	Links []Link `json:"href"`
}

func (person *Person) Validate() *serv.ErrorCode {

	err := fieldLengthValidate(person.FirstName, "first name")
	if err != nil {
		return err
	}
	err = fieldLengthValidate(person.LastName, "last name")
	if err != nil {
		return err
	}
	err = fieldLengthValidate(person.Email, "mail")
	if err != nil {
		return err
	}

	if len(person.Password) < serv.GetConfig().MinPasswordLength {
		return serv.NewBadRequest("Password is too short")
	}

	return nil
}

func (person *Person) GenerateLinks() {

	person.Links = append(person.Links, Link {
		Rel: "Профиль",
		Href: serv.GetConfig().Href + "/persons/" + person.ID.String(),
		Action: "GET",
	})
}

func (person *Person) WriteAsJsonResponseTo(ctx *fasthttp.RequestCtx, statusCode int) {
	person.GenerateLinks()
	resp, _ := json.Marshal(person)
	ctx.Write(resp)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetStatusCode(statusCode)
}