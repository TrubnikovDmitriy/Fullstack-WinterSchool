package unit

import (
	db "../../database"
	. "../service"
	"testing"
	"github.com/valyala/fasthttp"
)


func TestCreatePersonHappyPath(t *testing.T) {
	err := db.CreatePerson(GetNewPerson())
	if err != nil {
		t.Errorf("Error in create simple person\n%s", err)
	}
}

func TestCreatePersonWithoutFirstName(t *testing.T) {
	person := GetNewPerson()
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
	person := GetNewPerson()
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
	person := GetNewPerson()
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

	person := CreateNewPerson()
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

	originalPerson := GetNewPerson()
	db.CreatePerson(originalPerson)

	receivedPerson, err := db.GetPerson(originalPerson.ID)

	if err != nil {
		t.Errorf("Can't get created person\n%s", err)
		return
	}

	if receivedPerson.FirstName != originalPerson.FirstName {
		t.Errorf("Received person has another first name\n" +
			"Recieved person ID:\t%s,\nOriginal person ID:\t%s\n",
			receivedPerson.ID.String(), originalPerson.ID.String())
	}
	if receivedPerson.LastName != originalPerson.LastName {
		t.Errorf("Received person has another last name\n" +
			"Recieved person ID:\t%s,\nOriginal person ID:\t%s\n",
			receivedPerson.ID.String(), originalPerson.ID.String())
	}
	if receivedPerson.Mail != originalPerson.Mail {
		t.Errorf("Received person has another mail name\n" +
			"Recieved person ID:\t%s,\nOriginal person ID:\t%s\n",
			receivedPerson.ID.String(), originalPerson.ID.String())
	}
}
