package unit

import (
	db "../../database"
	. "../../tests"
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
		t.Fatalf("Can't get simple mathch\n%s", err)
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

func TestUpdateMatch(t *testing.T)  {

	tourney, _ := CreateNewMatches(3)
	matchesArray, _ := db.GetTournamentGrid(tourney.ID)

	matches := matchesArray.Array


	for i := range matches {
		UpdateMatch(&matches[i],tourney.OrganizeID)
		upd, err := db.UpdateMatch(&matches[i])
		if err != nil {
			t.Fatalf("Can't update match \n%s", err)
		}

		if upd.EndTime == nil || !upd.EndTime.Equal(*matches[i].EndTime) {
			t.Errorf("Time is not updated")
		}
		if upd.FirstTeamScore != matches[i].FirstTeamScore {
			t.Errorf("First team's score is not updated")
		}
		if upd.SecondTeamScore != matches[i].SecondTeamScore {
			t.Errorf("Second team's score is not updated")
		}
		if upd.FirstTeamID != matches[i].FirstTeamID {
			t.Errorf("First team is not updated")
		}
		if upd.SecondTeamID != matches[i].SecondTeamID {
			t.Errorf("Second team is not updated")
		}
	}

}

