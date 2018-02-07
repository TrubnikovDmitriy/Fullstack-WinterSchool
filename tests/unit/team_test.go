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


func TestCreateTeamHappyPath(t *testing.T) {
	team := GetNewTeam()
	err := db.CreateTeam(team)
	if err != nil {
		t.Error("Happy path for creating team failed: " + err.Message)
	}
}

func TestCreateTeamWithoutName(t *testing.T) {
	team := GetNewTeam()
	team.Name = ""
	err := db.CreateTeam(team)
	if err == nil {
		t.Fatalf("Team without team name has been created")
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error for creating team with incorrect fields:\n%s", err)
	}
}

func TestCreateTeamWithTooLongName(t *testing.T) {
	team := GetNewTeam()
	team.Name = "This team name is " + strings.Repeat("very, ", 10) + "very long"
	err := db.CreateTeam(team)
	if err == nil {
		t.Fatalf("Team with too long team name has been created\nTeam name: '%s'", team.Name)
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error for creating team with incorrect fields:\n%s", err)
	}
}

func TestCreateTeamWithoutAbout(t *testing.T) {
	team := GetNewTeam()
	team.About = ""
	err := db.CreateTeam(team)
	if err == nil {
		t.Fatalf("Team without about-field has been created")
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error for creating team with incorrect fields:\n%s", err)
	}
}

func TestCreateTeamWithTooLongAbout(t *testing.T) {
	team := GetNewTeam()
	team.Name = "This about is " + strings.Repeat("very, ", 10) + "very long"
	err := db.CreateTeam(team)
	if err == nil {
		t.Fatalf("Team with too about text has been created\nteam name: '%s'", team.Name)
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error for creating team with incorrect fields:\n%s", err)
	}
}

func TestCreateTeamDuplicate(t *testing.T) {
	team := CreateNewTeam()
	err := db.CreateTeam(team)
	if err == nil {
		t.Fatalf("Created two identical teams (UUID: %s)", team.ID.String())
	}
	if err.Code != fasthttp.StatusConflict {
		t.Errorf("Unexpectable error for creating the same teams:\n%s", err)
	}
}


func TestGetTeamHappy(t *testing.T) {
	teamOriginal := CreateNewTeam()
	team, err := db.GetTeamByID(teamOriginal.ID)
	if err != nil {
		t.Errorf("Error when getting team:\n%s", err)
	}
	if teamOriginal.Name != team.Name {
		t.Errorf("Names do not match\n" +
			"Original:\t%s\nGetting:\t%s\n", teamOriginal.Name, team.Name)
	}
	if teamOriginal.About != team.About {
		t.Errorf("About-fields do not match\n" +
			"Original:\t%s\nGetting:\t%s\n", teamOriginal.About, team.About)
	}
	if teamOriginal.ID != team.ID {
		t.Errorf("UUIDs do not match\n" +
			"Original:\t%s\nGetting:\t%s\n", teamOriginal.ID.String(), team.ID.String())
	}
}

func TestGetTheAbsentTeam(t *testing.T) {
	id, _ := uuid.NewV4()
	team, err := db.GetTeamByID(id)
	if err == nil {
		if team != nil {
			t.Fatalf("Got a non-existing team\n" +
				"requestID:\t%s\nresponseID:\t%s\n", id.String(), team.ID.String())
		} else {
			t.Fatalf("Got a non-existing team\nrequestID:\t%s\n", id.String())
		}
	}
	if err.Code != fasthttp.StatusNotFound {
		t.Errorf("Unexpectable error for getting non-existing team:\n%s", err)
	}
}

func TestGetFewTeams(t *testing.T) {
	var teams [5]*models.Team
	for i := range teams {
		teams[i] = CreateNewTeam()
	}
	for _, team := range teams {
		getTeam, err := db.GetTeamByID(team.ID)
		if err != nil {
			t.Errorf("Error when getting team:\n%s", err)
		}
		if team.Name != getTeam.Name {
			t.Errorf("Names do not match\n" +
				"Original:\t%s\nGetting:\t%s\n", team.Name, getTeam.Name)
		}
		if team.About != getTeam.About {
			t.Errorf("About-fields do not match\n" +
				"Original:\t%s\nGetting:\t%s\n", team.About, getTeam.About)
		}
		if team.ID != getTeam.ID {
			t.Errorf("UUIDs do not match\n" +
				"Original:\t%s\nGetting:\t%s\n", team.ID.String(), getTeam.ID.String())
		}
	}
}

