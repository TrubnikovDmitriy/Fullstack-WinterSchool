package models

import (
	"github.com/satori/go.uuid"
	"../services"
)

type OAuth struct {

	Email     string    `json:"email"`
	Password  string    `json:"password"`

	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	PersonID  uuid.UUID `json:"person_id"`
	Staff     bool      `json:"staff"`

	AppID      uuid.UUID  `json:"app_id"`
	Scope      int        `json:"scope"` // 1 - Read-Only, 2 - Read-Write
	AccessCode *uuid.UUID `json:"code"`
}

func (oauth *OAuth) Validate() *serv.ErrorCode {

	if len(oauth.Email) == 0 {
		return serv.NewBadRequest("Email is empty")
	}
	if len(oauth.Password) == 0 {
		return serv.NewBadRequest("Password is empty")
	}
	if oauth.AppID == uuid.Nil {
		return serv.NewBadRequest("Application ID are not specified")
	}
	if oauth.Scope == 0 {
		return serv.NewBadRequest("Access privileges (scope) are not specified")
	}
	return nil
}