package perfomance

import (
	"testing"
	. "../../tests"
	db "../../database"
	. "../../models"
	"time"
	"github.com/satori/go.uuid"
	"math/rand"
)

func BenchmarkCreateGame(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CreateNewGame()
	}
}

func BenchmarkGetGames(b *testing.B) {
	for i := 0; i < b.N; i++ {
		db.GetGames(i % 16, 3)
	}
}


func BenchmarkGetGameByID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		games, _ := db.GetGames(i % 16, 1)
		for _, game := range *games {
			db.GetGameByID(game.ID)
		}
	}
}


func BenchmarkCreateTournaments(b *testing.B) {
	for i := 0; i < b.N; i++ {

		gameID := CreateNewGame().ID
		timeNow := time.Now()

		for i := 0; i < b.N; i++ {

			id, _ := uuid.NewV4()
			tourney := Tournament {
				Title:        "Tourney title-" + id.String(),
				Started:      timeNow,
				Ended:        timeNow.AddDate(0, 3, 0),
				About:        id.String(),
				OrganizeName: GenerateFirstName() + " " + GenerateLastName(),
				OrganizeID:   id,
				GameID: gameID,
			}
			tourney.MatchTree = &GetNewMatches(rand.Intn(4) + 1)[0]
			db.CreateTournament(&tourney)
		}
	}
}

