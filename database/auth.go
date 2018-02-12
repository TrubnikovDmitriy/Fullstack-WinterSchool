package database

import (
	"../models"
	"../services"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx"
)

func Auth(oauth *models.OAuth) *serv.ErrorCode {

	errorCode := oauth.Validate()
	if errorCode != nil {
		return errorCode
	}

	const auth = "Auth"
	db := sharedKeyForReadByMail(oauth.Email)
	db.Prepare(auth, "SELECT person_id, pass, staff FROM auth WHERE email = $1")

	var pgUUID pgtype.UUID
	var passwordHash []byte
	err := db.QueryRow(auth, oauth.Email).Scan(&pgUUID, &passwordHash, &oauth.Staff)
	if err == pgx.ErrNoRows || !serv.PasswordEqual(oauth.Password, passwordHash) {
		return serv.NewNotFound()
	}
	if err != nil {
		return checkError(err, db)
	}

	oauth.PersonID = castUUID(pgUUID)
	const getInfoForAuth = "GetInfoForAuth"
	db = sharedKeyForReadByID(oauth.PersonID)
	db.Prepare(getInfoForAuth, "SELECT first_name, last_name FROM persons WHERE id = $1")

	err = db.QueryRow(getInfoForAuth, oauth.PersonID).Scan(&oauth.FirstName, &oauth.LastName)
	if err != nil {
		return checkError(err, db)
	}

	return nil
}