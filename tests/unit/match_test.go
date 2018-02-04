package unit

import (
	db "../../database"
	. "../service"
	"testing"
)

func TestMatchGetRoot(t *testing.T) {

	tourney, matches := CreateNewMatches(1)

	match, err := db.GetMatchByID(tourney.ID, matches[0].ID)
	if err != nil {
		t.Errorf("Can't get simple mathch\n%s", err)
	}
	if match.NextMatch != nil {
		t.Errorf("Incorrect next match")
	}
	if match.PrevMatch1 != nil {
		t.Errorf("Incorrect previous match #1")
	}
	if match.PrevMatch2 != nil {
		t.Errorf("Incorrect previous match #2")
	}
}

func TestMatchGetChild(t *testing.T) {

	tourney, matches := CreateNewMatches(3)

	match, err := db.GetMatchByID(tourney.ID, matches[1].ID)
	if err != nil {
		t.Errorf("Can't get simple mathch\n%s", err)
	}
	if *match.NextMatch != matches[0].ID {
		t.Errorf("Incorrect next match")
	}
	if *match.PrevMatch1 != matches[3].ID {
		t.Errorf("Incorrect previous match #1")
	}
	if *match.PrevMatch2 != matches[4].ID {
		t.Errorf("Incorrect previous match #2")
	}
}

