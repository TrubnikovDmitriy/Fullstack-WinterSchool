package unit

import (
	"testing"
	. "../../tests"
	db "../../database"
	"github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
)


func TestAddPlayerHappyPath(t *testing.T) {
	player := GetNewPlayers(1)[0]
	teamName := player.TeamName
	player.TeamName = ""
	err := db.AddPlayerInTeam(player)
	if err != nil {
		t.Fatalf("Can't create simple player %s", err)
	}

	if teamName != player.TeamName {
		t.Errorf("Team name is incorrect (%s instead of %s)\n",
			player.TeamName, teamName)
	}
}

func TestAddPlayerWithFakePersonID(t *testing.T) {

	player := GetNewPlayers(1)[0]
	player.PersonID, _ = uuid.NewV4()
	err := db.AddPlayerInTeam(player)

	if err == nil {
		t.Fatalf("Added player with fake person ID %s\n", player)
	}

	if err.Code != fasthttp.StatusBadRequest {
		t.Fatalf("Unexpected error for fake\n" +
			"person ID: %d\nmessage: %s\n", err.Code, err.Message)
	}
}

func TestGetEmptyTeam(t *testing.T) {

	team := CreateNewTeam()
	empty, err := db.GetPlayersOfTeam(team.ID)

	if err != nil {
		t.Fatalf("Can't get empty array players of team:\n%s", err)
	}
	if len(empty) != 0 {
		t.Fatalf("Empty team returning non empty list of players\n" +
			"team ID: %s\n", team.ID.String())
	}
}

func TestGetSinglePlayer(t *testing.T) {

	player := CreateNewPlayers(1)[0]
	players, err := db.GetPlayersOfTeam(player.TeamID)
	if err != nil {
		t.Fatalf("Can't get players of team %s", err)
	}

	if len(players) != 1 {
		t.Fatalf("Returning wrong number of players (%d instead of 1)", len(players))
	}

	if players[0].ID != player.ID {
		t.Fatalf("Returning another player (player ID %s instead of %s",
			players[0].ID.String(), player.ID.String())
	}
}

func TestGetFewPlayer(t *testing.T) {

	const numberOfPlayers = 10
	players := CreateNewPlayers(numberOfPlayers)
	players, err := db.GetPlayersOfTeam(players[0].TeamID)
	if err != nil {
		t.Fatalf("Can't get players of team %s", err)
	}

	if len(players) != numberOfPlayers {
		t.Fatalf("Returning wrong number of players (%d instead of %d)",
			len(players), numberOfPlayers)
	}
}

func TestGetFewPlayerAfterDelete(t *testing.T) {

	const numberOfPlayers = 10
	const deletePlayers = 2
	players := CreateNewPlayers(numberOfPlayers)

	db.DeletePlayerFromTeam(players[0])
	db.DeletePlayerFromTeam(players[7])

	players, err := db.GetPlayersOfTeam(players[0].TeamID)
	if err != nil {
		t.Fatalf("Can't get players of team after deleting player\n%s", err)
	}

	if len(players) != numberOfPlayers - deletePlayers {
		t.Fatalf("Returning wrong number of players (%d instead of %d)",
			len(players), numberOfPlayers - deletePlayers)
	}
}
