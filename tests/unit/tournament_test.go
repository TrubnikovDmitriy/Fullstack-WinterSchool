package unit

import (
	db "../../database"
	. "../service"
	"time"
	"testing"
	"github.com/valyala/fasthttp"
)

func TestCreateTourneyHappyPath(t *testing.T) {

	tourney := GetNewTournament()
	matches := GetNewMatches(3)

	tourney.MatchTree = &matches[0]

	err := db.CreateTournament(tourney)
	if err != nil {
		t.Errorf("Can't create simple tournament:\n%s", err)
	}
}

func TestCreateTourneyEmptyMatch(t *testing.T) {

	tourney := GetNewTournament()

	err := db.CreateTournament(tourney)
	if err == nil {
		t.Errorf("Created tournament without matches (ID: %s)", tourney.ID)
		return
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error code for bad request\n%s\n", err)
	}
}

func TestCreateTourneyWithoutTitle(t *testing.T) {

	tourney := GetNewTournament()
	tourney.MatchTree = &GetNewMatches(3)[0]
	tourney.Title = ""

	err := db.CreateTournament(tourney)
	if err == nil {
		t.Errorf("Created tournament without title (ID: %s)", tourney.ID)
		return
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error code for bad request\n%s\n", err)
	}
}

func TestCreateTourneyWithoutAbout(t *testing.T) {

	tourney := GetNewTournament()
	tourney.MatchTree = &GetNewMatches(3)[0]
	tourney.About = ""

	err := db.CreateTournament(tourney)
	if err == nil {
		t.Errorf("Created tournament without about-field (ID: %s)", tourney.ID)
		return
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error code for bad request\n%s\n", err)
	}
}

func TestCreateTourneyWithIncorrectData(t *testing.T) {

	tourney := GetNewTournament()
	tourney.MatchTree = &GetNewMatches(3)[0]
	tourney.Started = tourney.Ended.Add(1)

	err := db.CreateTournament(tourney)
	if err == nil {
		t.Errorf("Created tournament where 'end' before 'start' (ID: %s)", tourney.ID)
		return
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error code for bad request\n%s\n", err)
	}
}

func TestCreateTourneyWithEmptyData(t *testing.T) {

	tourney := GetNewTournament()
	tourney.MatchTree = &GetNewMatches(3)[0]
	tourney.Started = time.Time{}

	err := db.CreateTournament(tourney)
	if err == nil {
		t.Errorf("Created tournament where 'end' before 'start' (ID: %s)", tourney.ID)
		return
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error code for bad request\n%s\n", err)
	}
}

func TestCreateTooManyMatches(t *testing.T) {
	tourney := GetNewTournament()
	matches := GetNewMatches(10)

	tourney.MatchTree = &matches[0]

	err := db.CreateTournament(tourney)
	if err == nil {
		t.Errorf("Created tounaments with matches deep is 10\nID: %s\n", tourney.ID)
		return
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error code for bad request\n%s\n", err)
	}
}

func TestCreateNotBinaryMatchTree(t *testing.T) {
	tourney := GetNewTournament()
	matches := GetNewMatches(3)
	matches[2].RightChild = nil

	tourney.MatchTree = &matches[0]

	err := db.CreateTournament(tourney)
	if err == nil {
		t.Errorf("Created tounaments with matches deep is 10\nID: %s\n", tourney.ID)
		return
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error code for bad request\n%s\n", err)
	}
}

func TestCreateTourneyDuplicate(t *testing.T) {

	tourney := CreateNewTournament()
	db.CreateTournament(tourney)

	err := db.CreateTournament(tourney)
	if err == nil {
		t.Errorf("Created the two same tournaments (ID: %s)\n", tourney.ID.String())
		return
	}
	if err.Code != fasthttp.StatusConflict {
		t.Errorf("Unexpected error for duplicate:\n%s", err)
	}
}

func TestGetTourney(t *testing.T) {

	original := GetNewTournament()
	original.MatchTree = &GetNewMatches(3)[0]
	err := db.CreateTournament(original)

	received, err := db.GetTourneyByID(original.ID)

	if err != nil {
		t.Errorf("Can't get created tournament\n%s", err)
		return
	}

	if original.Title != received.Title {
		t.Errorf("Received tournament has another title\n" +
			"Recieved:\t%s,\nOriginal:\t%s\n",
			original.Title, received.Title)
	}
	if original.About != received.About {
		t.Errorf("Received tournament has another about-field\n" +
			"Recieved ID:\t%s,\nOriginal ID\t%s\n",
			original.About, received.About)
	}
	if original.Started.YearDay() != received.Started.YearDay() {
		t.Errorf("Received tournament has another time started\n" +
			"Recieved ID:\t%s,\nOriginal ID:\t%s\n",
			original.Started, received.Started)
	}
	if original.Ended.YearDay() != received.Ended.YearDay() {
		t.Errorf("Received tournament has another time ended\n" +
			"Recieved ID:\t%s,\nOriginal ID:\t%s\n",
			original.Started, received.Started)

	}
}

func TestGetTourneyGridSymmetric(t *testing.T) {

	tourney := GetNewTournament()
	matches := GetNewMatches(2)
	tourney.MatchTree = &matches[0]
	db.CreateTournament(tourney)

	arrayMatches, err := db.GetTournamentGrid(tourney.ID)
	if err != nil {
		t.Errorf("Can't get tournament grid:\n%s", err)
		return
	}

	for _, match := range arrayMatches.Array {
		if match.TourneyID != tourney.ID {
			t.Errorf("Matche obtains to another tournament " +
					"(wrong tourney ID = %s)", match.TourneyID)
		}
	}

	if len(matches) != len(arrayMatches.Array) {
		t.Errorf("Wrong number of match (%d != %d)",
			len(matches), len(arrayMatches.Array))
	}

	// Родительская нода
	if tourney.MatchTree.ID != arrayMatches.Array[0].ID {
		t.Errorf("Parent node ID is incorrect")
	}
	if tourney.MatchTree.LeftChild.ID != *arrayMatches.Array[0].PrevMatch1 {
		t.Errorf("Parent's left child ID is incorrect")
	}
	if tourney.MatchTree.RightChild.ID != *arrayMatches.Array[0].PrevMatch2 {
		t.Errorf("Parent's right child ID is incorrect")
	}

	// Левый ребенок
	if tourney.MatchTree.LeftChild.ID != arrayMatches.Array[1].ID {
		t.Errorf("Left node ID is incorrect")
	}
	if arrayMatches.Array[1].PrevMatch1 != nil {
		t.Errorf("Left left node is incorrect")
	}
	if arrayMatches.Array[1].PrevMatch2 != nil {
		t.Errorf("Left right node is incorrect")
	}

	// Правый ребенок
	if tourney.MatchTree.RightChild.ID != arrayMatches.Array[2].ID {
		t.Errorf("Right node ID is incorrect")
	}
	if arrayMatches.Array[2].PrevMatch1 != nil {
		t.Errorf("Left left node is incorrect")
	}
	if arrayMatches.Array[2].PrevMatch2 != nil {
		t.Errorf("Left right node is incorrect")
	}

}

func TestGetTourneyGridAsymmetric(t *testing.T) {

	tourney := GetNewTournament()
	matches := GetNewMatches(4)
	matches[1].RightChild = nil
	matches[1].LeftChild = nil
	tourney.MatchTree = &matches[0]
	db.CreateTournament(tourney)

	arrayMatches, err := db.GetTournamentGrid(tourney.ID)
	if err != nil {
		t.Errorf("Can't get tournament grid:\n%s", err)
		return
	}

	for _, match := range arrayMatches.Array {
		if match.TourneyID != tourney.ID {
			t.Errorf("Matche obtains to another tournament " +
				"(wrong tourney ID = %s)", match.TourneyID)
		}
	}

	if len(arrayMatches.Array) != 9 {
		t.Errorf("Wrong number of match in asymmetrical grid (9 != %d)",
				len(arrayMatches.Array))
	}
}





