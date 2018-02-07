package unit

import (
	db "../../database"
	. "../../tests"
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
		t.Fatalf("Created person with empty firstname")
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
		t.Fatalf("Created person with empty lastname")
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
		t.Fatalf("Created person with too short password (ID: %s)\n", person.ID.String())
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpected error\n%s", err)
	}
}

func TestCreatePersonDuplicate(t *testing.T) {

	person := CreateNewPerson()
	err := db.CreatePerson(person)

	if err == nil {
		t.Fatalf("Created the two same persons (ID: %s)\n", person.ID.String())
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
		t.Fatalf("Can't get created person\n%s", err)
	}

	if receivedPerson.FirstName != originalPerson.FirstName {
		t.Errorf("Received person has another first name\n" +
			"Recieved person ID:\t%s,\nOriginal person ID:\t%s\n",
			receivedPerson.ID.String(), originalPerson.ID.String())
	}
	if receivedPerson.LastName != originalPerson.LastName {
		t.Errorf("Received person has another last name\n" +
			"Recieved person name:\t%s,\nOriginal person name:\t%s\n",
			receivedPerson.LastName, originalPerson.LastName)
	}
}
