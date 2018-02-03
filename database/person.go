package database

import (
	"../models"
	"../services"
	"github.com/satori/go.uuid"
	"github.com/jackc/pgx/pgtype"
	"github.com/valyala/fasthttp"
)

func CreatePerson(person *models.Person) *serv.ErrorCode {

	errorCode := person.Validate()
	if errorCode != nil {
		return errorCode
	}

	var existingID pgtype.UUID
	err := verifyUnique("SELECT id FROM persons WHERE mail = $1",
		&existingID, person.Mail)
	if err != nil {
		return &serv.ErrorCode{
			Code: fasthttp.StatusConflict,
			Message: "User with the same mail already exists",
			Link: serv.Href + "/persons/" +
				uuid.FromBytesOrNil(existingID.Bytes[:]).String(),
		}
	}

	person.ID = getUUID()
	master := sharedKeyForWriteByUUID(person.ID)
	const insertPerson = "InsertPerson"
	master.Prepare(insertPerson,
		"INSERT INTO persons(id, first_name, last_name, about, mail, passw) " +
			"VALUES($1, $2, $3, $4, $5, $6);")


	_, err = master.Exec(insertPerson,
			person.ID, person.FirstName, person.LastName,
			person.About, person.Mail, serv.PasswordHashing(person.Password))
	person.Password = ""
	if err != nil {
		return serv.NewServerError(err)
	}

	return nil
}

func GetPerson(id uuid.UUID) (*models.Person, *serv.ErrorCode) {

	const selectPersonByID = "SelectPersonByID"
	db := sharedKeyForReadByUUID(id)
	db.Prepare(selectPersonByID,
		"SELECT first_name, last_name, mail, about FROM persons WHERE id = $1")

	person := models.Person{ID: id}
	err := db.QueryRow(selectPersonByID, id).
		Scan(&person.FirstName, &person.LastName, &person.Mail, &person.About)
	if err != nil {
		return nil, checkError(err)
	}

	return &person, nil
}
