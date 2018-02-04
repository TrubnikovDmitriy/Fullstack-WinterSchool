package unit

import (
	"../../database"
	"../../models"
	"testing"
	"github.com/satori/go.uuid"
	"strings"
)

func createNewGame() *models.Game {
	id, _ := uuid.NewV1()
	uniqueSuffix := strings.Split(id.String(), "-")[0]
	game := models.Game{
		Title: "Title-" + uniqueSuffix,
		About: "Some text about useful things",
	}

	return &game
}

func TestDatabase(t *testing.T) {
	t.Error("Always error for test Jenkins")
	game := createNewGame()
	err := database.CreateGame(game)
	if err != nil {
		t.Error(err.Message)
	}
}
