package database

import (
	"github.com/satori/go.uuid"
	"github.com/jackc/pgx/pgtype"
	"github.com/valyala/fasthttp"
	"../services"
	"../models"
	"github.com/jackc/pgx"
)

func CreatePerson(person *models.Person) *serv.ErrorCode {

	errorCode := person.Validate()
	if errorCode != nil {
		return errorCode
	}

	const checkUniqueEmail = "CheckUniqueEmail"
	authDB := sharedKeyForWriteByMail(person.Email)
	authDB.Prepare(checkUniqueEmail, "SELECT person_id FROM auth WHERE email = $1")

	var existingID pgtype.UUID
	err := authDB.QueryRow(checkUniqueEmail, person.Email).Scan(&existingID)
	if err == nil {
		return &serv.ErrorCode {
			Code: fasthttp.StatusConflict,
			Message: "User with the same mail already exists",
			Link: serv.GetConfig().Href + "/persons/" +
				uuid.FromBytesOrNil(existingID.Bytes[:]).String(),
		}
	}
	if err != pgx.ErrNoRows {
		return checkError(err, authDB)
	}

	person.ID = getID()
	personDB := sharedKeyForWriteByID(person.ID)

	const insertPerson = "InsertPerson"
	personDB.Prepare(insertPerson,
		"INSERT INTO persons(id, first_name, last_name, about) VALUES($1, $2, $3, $4);")

	const insertAuth = "InsertAuth"
	authDB.Prepare(insertAuth,
		"INSERT INTO auth(email, pass, person_id) VALUES($1, $2, $3);")


	_, err = personDB.Exec(insertPerson, person.ID, person.FirstName, person.LastName, person.About)
	if err != nil {
		return checkError(err, personDB)
	}

	_, err = authDB.Exec(insertAuth, person.Email, serv.PasswordHashing(person.Password), person.ID)
	person.Password = ""
	if err != nil {
		personDB.Exec("DELETE FROM persons WHERE id=$1", person.ID)
		return checkError(err, personDB)
	}

	return nil
}

func GetPerson(id uuid.UUID) (*models.Person, *serv.ErrorCode) {

	const selectPersonByID = "SelectPersonByID"
	db := sharedKeyForReadByID(id)
	db.Prepare(selectPersonByID,
		"SELECT first_name, last_name, about FROM persons WHERE id = $1")

	person := models.Person{ID: id}
	err := db.QueryRow(selectPersonByID, id).
		Scan(&person.FirstName, &person.LastName, &person.About)
	if err != nil {
		return nil, checkError(err, db)
	}

	return &person, nil
}
