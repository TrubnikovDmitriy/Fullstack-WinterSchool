package unit

import (
	db "../../database"
	"../../models"
	"strings"
	"github.com/satori/go.uuid"
	"github.com/liderman/text-generator"
	"testing"
	"github.com/valyala/fasthttp"
)

func getNewPerson() *models.Person {

	var FirstNameTemplate string = "{Vasya|Peter|Nikita|Sasha|Dmitriy|Enakentiy|John|Masha|Natasha|Tony}"
	var LastNameTemplate string = "{Silaev|Kuzmin|Krasnov|Sitnikov|Smirnov|Gorbenko|Trubnikov|Smirnova}"
	var MailsTemplate string = "{@mail.ru|@yandex.ru|@gmail.com|@rambler.com}"

	tg := text_generator.New()
	id, _ := uuid.NewV4()
	mail := strings.Split(id.String(), "-")[1]
	password := strings.Split(id.String(), "-")[0]

	person := models.Person {
		FirstName: tg.Generate(FirstNameTemplate),
		LastName: tg.Generate(LastNameTemplate),
		Mail: tg.Generate(LastNameTemplate) + "_" + mail + tg.Generate(MailsTemplate),
		Password: password,
	}

	return &person
}

func createNewPerson() *models.Person {

	original := getNewPerson()
	forDatabase := *original
	db.CreatePerson(&forDatabase)

	return original
}


func TestCreatePersonHappyPath(t *testing.T) {
	err := db.CreatePerson(getNewPerson())
	if err != nil {
		t.Errorf("Error in create simple person\n%s", err)
	}
}

func TestCreatePersonWithoutFirstName(t *testing.T) {
	person := getNewPerson()
	person.FirstName = ""
	err := db.CreatePerson(person)
	if err == nil {
		t.Error("Created person with empty firstname")
		return
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpected error:\n%s", err)
	}
}

func TestCreatePersonWithoutLastName(t *testing.T) {
	person := getNewPerson()
	person.LastName = ""
	err := db.CreatePerson(person)
	if err == nil {
		t.Error("Created person with empty lastname")
		return
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpected error:\n%s", err)
	}
}

func TestCreatePersonWithShortPassword(t *testing.T) {
	person := getNewPerson()
	person.Password = "foo"
	err := db.CreatePerson(person)
	if err == nil {
		t.Errorf("Created person with too short password (ID: %s)\n", person.ID.String())
		return
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpected error\n%s", err)
	}
}

func TestCreatePersonDuplicate(t *testing.T) {

	person := createNewPerson()
	err := db.CreatePerson(person)

	if err == nil {
		t.Errorf("Created the two same persons (ID: %s)\n", person.ID.String())
		return
	}
	if err.Code != fasthttp.StatusConflict {
		t.Errorf("Unexpected error for duplicate:\n%s", err)
	}
}

func TestGetPerson(t *testing.T) {

	originalPerson := getNewPerson()
	db.CreatePerson(originalPerson)

	receivedPerson, err := db.GetPerson(originalPerson.ID)

	if err != nil {
		t.Errorf("Can't get created person\n%s", err)
		return
	}

	if receivedPerson.FirstName != originalPerson.FirstName {
		t.Errorf("Received person has another first name\n" +
			"Recieved person ID:\t%s\n,Original person ID\t%s\n",
			receivedPerson.ID.String(), originalPerson.ID.String())
	}
	if receivedPerson.LastName != originalPerson.LastName {
		t.Errorf("Received person has another last name\n" +
			"Recieved person ID:\t%s\n,Original person ID\t%s\n",
			receivedPerson.ID.String(), originalPerson.ID.String())
	}
	if receivedPerson.Mail != originalPerson.Mail {
		t.Errorf("Received person has another mail name\n" +
			"Recieved person ID:\t%s\n,Original person ID\t%s\n",
			receivedPerson.ID.String(), originalPerson.ID.String())
	}
}
