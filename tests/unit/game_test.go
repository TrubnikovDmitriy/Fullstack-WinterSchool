package unit

import (
	db "../../database"
	"../../models"
	"testing"
	"github.com/satori/go.uuid"
	"strings"
	"github.com/valyala/fasthttp"
)

func getNewGame() *models.Game {
	id, _ := uuid.NewV1()
	uniqueSuffix := strings.Split(id.String(), "-")[0]
	game := models.Game{
		Title: "Title-" + uniqueSuffix,
		About: "Some text about useful things",
	}

	return &game
}

func createNewGame() *models.Game {
	gameOriginal := getNewGame()

	gameToDataBase := *gameOriginal
	db.CreateGame(&gameToDataBase)

	gameOriginal.ID = gameToDataBase.ID
	return gameOriginal
}


func TestCreateGameHappyPath(t *testing.T) {
	game := getNewGame()
	err := db.CreateGame(game)
	if err != nil {
		t.Error("Happy path for creating game failed: " + err.Message)
	}
}

func TestCreateGameWithoutTitle(t *testing.T) {
	game := getNewGame()
	game.Title = ""
	err := db.CreateGame(game)
	if err == nil {
		t.Error("Game without title has been created")
		return
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error for creating game with incorrect fields:\n%s", err)
	}
}

func TestCreateGameWithTooLongTitle(t *testing.T) {
	game := getNewGame()
	game.Title = "This title is " + strings.Repeat("very, ", 10) + "very long"
	err := db.CreateGame(game)
	if err == nil {
		t.Errorf("Game with too long title has been created\nTitle: '%s'", game.Title)
		return
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error for creating game with incorrect fields:\n%s", err)
	}
}

func TestCreateGameWithoutAbout(t *testing.T) {
	game := getNewGame()
	game.About = ""
	err := db.CreateGame(game)
	if err == nil {
		t.Error("Game without about-field has been created")
		return
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error for creating game with incorrect fields:\n%s", err)
	}
}

func TestCreateGameWithTooLongAbout(t *testing.T) {
	game := getNewGame()
	game.Title = "This about is " + strings.Repeat("very, ", 10) + "very long"
	err := db.CreateGame(game)
	if err == nil {
		t.Errorf("Game with too about text has been created\ntitle: '%s'", game.Title)
		return
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error for creating game with incorrect fields:\n%s", err)
	}
}

func TestCreateGameConflict(t *testing.T) {
	game := createNewGame()
	err := db.CreateGame(game)
	if err == nil {
		t.Errorf("Created two identical games (UUID: %s)", game.ID.String())
		return
		return
	}
	if err.Code != fasthttp.StatusConflict {
		t.Errorf("Unexpectable error for creating the same games:\n%s", err)
	}
}


func TestGetGameHappy(t *testing.T) {
	gameOriginal := createNewGame()
	game, err := db.GetGameByID(gameOriginal.ID)
	if err != nil {
		t.Errorf("Error when getting game:\n%s", err)
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
	id, _ := uuid.NewV4()
	game, err := db.GetGameByID(id)
	if err == nil {
		if game != nil {
			t.Errorf("Got a non-existing game\n" +
				"requestID:\t%s\nresponseID:\t%s\n", id.String(), game.ID.String())
		} else {
			t.Errorf("Got a non-existing game\nrequestID:\t%s\n", id.String())
		}
		return
	}
	if err.Code != fasthttp.StatusNotFound {
		t.Errorf("Unexpectable error for getting non-existing game:\n%s", err)
	}
}

func TestGetFewGames(t *testing.T) {
	var games [5]*models.Game
	for i := range games {
		games[i] = createNewGame()
	}
	for _, game := range games {
		getGame, err := db.GetGameByID(game.ID)
		if err != nil {
			t.Errorf("Error when getting game:\n%s", err)
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
