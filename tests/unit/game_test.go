package unit

import (
	db "../../database"
	. "../../tests"
	"../../models"
	"testing"
	"github.com/satori/go.uuid"
	"strings"
	"github.com/valyala/fasthttp"
)


func TestCreateGameHappyPath(t *testing.T) {
	game := GetNewGame()
	err := db.CreateGame(game)
	if err != nil {
		t.Error("Happy path for creating game failed: " + err.Message)
	}
}

func TestCreateGameWithoutTitle(t *testing.T) {
	game := GetNewGame()
	game.Title = ""

	err := db.CreateGame(game)
	if err == nil {
		t.Fatalf("Game without title has been created")
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error for creating game with incorrect fields:\n%s", err)
	}
}

func TestCreateGameWithTooLongTitle(t *testing.T) {
	game := GetNewGame()
	game.Title = "This title is " + strings.Repeat("very, ", 10) + "very long"
	err := db.CreateGame(game)
	if err == nil {
		t.Fatalf("Game with too long title has been created\nTitle: '%s'", game.Title)
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error for creating game with incorrect fields:\n%s", err)
	}
}

func TestCreateGameWithoutAbout(t *testing.T) {
	game := GetNewGame()
	game.About = ""
	err := db.CreateGame(game)
	if err == nil {
		t.Fatalf("Game without about-field has been created")
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error for creating game with incorrect fields:\n%s", err)
	}
}

func TestCreateGameWithTooLongAbout(t *testing.T) {
	game := GetNewGame()
	game.Title = "This about is " + strings.Repeat("very, ", 10) + "very long"
	err := db.CreateGame(game)
	if err == nil {
		t.Fatalf("Game with too long about-field has been created\ntitle: '%s'", game.Title)
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error for creating game with incorrect fields:\n%s", err)
	}
}

func TestCreateGameConflict(t *testing.T) {
	game := CreateNewGame()
	err := db.CreateGame(game)
	if err == nil {
		t.Fatalf("Created two identical games (UUID: %s)", game.ID.String())
	}
	if err.Code != fasthttp.StatusConflict {
		t.Errorf("Unexpectable error for creating the same games:\n%s", err)
	}
}


func TestGetGameHappyPath(t *testing.T) {
	gameOriginal := GetNewGame()
	db.CreateGame(gameOriginal)

	game, err := db.GetGameByID(gameOriginal.ID)
	if err != nil {
		t.Fatalf("Error when getting game:\n%s", err)
	}
	if gameOriginal.Title != game.Title {
		t.Errorf("Titles do not match\n" +
			"Original:\t%s\nGetting:\t%s\n", gameOriginal.Title, game.Title)
	}
	if gameOriginal.About != game.About {
		t.Errorf("About-fields do not match\n" +
			"Original:\t%s\nGetting:\t%s\n", gameOriginal.About, game.About)
	}
	if gameOriginal.ID != game.ID {
		t.Errorf("UUIDs do not match\n" +
			"Original:\t%s\nGetting:\t%s\n", gameOriginal.ID.String(), game.ID.String())
	}
}

func TestGetTheAbsentGame(t *testing.T) {

	randomID, _ := uuid.NewV4()
	game, err := db.GetGameByID(randomID)

	if err == nil {
		if game != nil {
			t.Fatalf("Got a non-existing game\n" +
				"requestID:\t%s\nresponseID:\t%s\n", randomID.String(), game.ID.String())
		}
		t.Fatalf("Got a non-existing game\nrequestID:\t%s\n", randomID.String())
	}
	if err.Code != fasthttp.StatusNotFound {
		t.Errorf("Unexpectable error for getting non-existing game:\n%s", err)
	}
}

func TestGetFewGames(t *testing.T) {

	var games [5]*models.Game
	for i := range games {
		games[i] = CreateNewGame()
	}

	for _, game := range games {
		getGame, err := db.GetGameByID(game.ID)
		if err != nil {
			t.Fatalf("Error when getting game:\n%s", err)
		}
		if game.Title != getGame.Title {
			t.Errorf("Titles do not match\n" +
				"Original:\t%s\nGetting:\t%s\n", game.Title, getGame.Title)
		}
		if game.About != getGame.About {
			t.Errorf("About-fields do not match\n" +
				"Original:\t%s\nGetting:\t%s\n", game.About, getGame.About)
		}
		if game.ID != getGame.ID {
			t.Errorf("UUIDs do not match\n" +
				"Original:\t%s\nGetting:\t%s\n", game.ID.String(), getGame.ID.String())
		}
	}
}
