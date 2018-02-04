package unit

import (
	db "../../database"
	"../../models"
	"testing"
	"github.com/satori/go.uuid"
	"strings"
	"github.com/valyala/fasthttp"
)

func getNewTeam() *models.Team {
	id, _ := uuid.NewV1()
	uniqueSuffix := strings.Split(id.String(), "-")[0]
	team := models.Team{
		Name: "Name-" + uniqueSuffix,
		About: "A few words about this amazing team",
	}

	return &team
}

func createNewTeam() *models.Team {
	teamOriginal := getNewTeam()

	teamToDataBase := *teamOriginal
	db.CreateTeam(&teamToDataBase)

	teamOriginal.ID = teamToDataBase.ID
	return teamOriginal
}


func TestCreateTeamHappyPath(t *testing.T) {
	team := getNewTeam()
	err := db.CreateTeam(team)
	if err != nil {
		t.Error("Happy path for creating team failed: " + err.Message)
	}
}

func TestCreateTeamWithoutName(t *testing.T) {
	team := getNewTeam()
	team.Name = ""
	err := db.CreateTeam(team)
	if err == nil {
		t.Error("Team without team name has been created")
		return
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error for creating team with incorrect fields:\n%s", err)
	}
}

func TestCreateTeamWithTooLongName(t *testing.T) {
	team := getNewTeam()
	team.Name = "This team name is " + strings.Repeat("very, ", 10) + "very long"
	err := db.CreateTeam(team)
	if err == nil {
		t.Errorf("Team with too long team name has been created\nTeam name: '%s'", team.Name)
		return
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error for creating team with incorrect fields:\n%s", err)
	}
}

func TestCreateTeamWithoutAbout(t *testing.T) {
	team := getNewTeam()
	team.About = ""
	err := db.CreateTeam(team)
	if err == nil {
		t.Error("Team without about-field has been created")
		return
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error for creating team with incorrect fields:\n%s", err)
	}
}

func TestCreateTeamWithTooLongAbout(t *testing.T) {
	team := getNewTeam()
	team.Name = "This about is " + strings.Repeat("very, ", 10) + "very long"
	err := db.CreateTeam(team)
	if err == nil {
		t.Errorf("Team with too about text has been created\nteam name: '%s'", team.Name)
		return
	}
	if err.Code != fasthttp.StatusBadRequest {
		t.Errorf("Unexpectable error for creating team with incorrect fields:\n%s", err)
	}
}

func TestCreateTeamConflict(t *testing.T) {
	team := createNewTeam()
	err := db.CreateTeam(team)
	if err == nil {
		t.Errorf("Created two identical teams (UUID: %s)", team.ID.String())
		return
	}
	if err.Code != fasthttp.StatusConflict {
		t.Errorf("Unexpectable error for creating the same teams:\n%s", err)
	}
}


func TestGetTeamHappy(t *testing.T) {
	teamOriginal := createNewTeam()
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
			t.Errorf("Got a non-existing team\n" +
				"requestID:\t%s\nresponseID:\t%s\n", id.String(), team.ID.String())
		} else {
			t.Errorf("Got a non-existing team\nrequestID:\t%s\n", id.String())
		}
		return
	}
	if err.Code != fasthttp.StatusNotFound {
		t.Errorf("Unexpectable error for getting non-existing team:\n%s", err)
	}
}

func TestGetFewTeams(t *testing.T) {
	var teams [5]*models.Team
	for i := range teams {
		teams[i] = createNewTeam()
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

