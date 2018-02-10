package tests

import (
	db "../database"
	. "../models"
	"github.com/satori/go.uuid"
	"math"
	"time"
	"strings"
	"github.com/liderman/text-generator"
	"math/rand"
)

// Функция для создания дерева матчей заданной глубины (односторонняя связь parent->child),
// возвращает массив из связанных между собой элементов (нулевой элемент - корень дерева)
func GetNewMatches(deep int) []MatchesTreeForm {

	nodesCount := int(math.Pow(2, float64(deep))) - 1
	nodes := make([]MatchesTreeForm, nodesCount)
	times := time.Now()

	for i := 1; i <= nodesCount; i++ {
		if i <= nodesCount / 2 {
			nodes[i - 1].LeftChild = &nodes[i * 2 - 1]
			nodes[i - 1].RightChild = &nodes[i * 2]
		}
		nodes[i - 1].StartTime = times.Add(time.Duration(nodesCount - i))
	}
	return nodes
}

// Функция для создания дерева заданной глубины, возвращает корневую ноду
func CreateBinaryTree(deep int) MatchesTreeForm {

	nodesCount := int(math.Pow(2, float64(deep))) - 1
	nodes := make([]MatchesTreeForm, nodesCount)
	times := time.Now()

	for i := 1; i <= nodesCount; i++ {
		if i <= nodesCount / 2 {
			nodes[i - 1].LeftChild = &nodes[i * 2 - 1]
			nodes[i - 1].RightChild = &nodes[i * 2]
		}
		nodes[i - 1].StartTime = times.Add(time.Duration(nodesCount - i))
	}
	return nodes[0]
}


func GetNewTournament() *Tournament {

	id, _ := uuid.NewV4()
	postfix := strings.Split(id.String(), "-")
	timeNow := time.Now()

	tourney := Tournament {
		Title: "Tournament title - " + postfix[0],
		Started: timeNow,
		Ended: timeNow.AddDate(0, 3,0),
		About: postfix[4],
		OrganizeName: GenerateFirstName() + " " + GenerateLastName(),
		OrganizeID: id,
		GameID: CreateNewGame().ID,
	}

	return &tourney
}

func CreateNewTournament() *Tournament {

	originalTourney := GetNewTournament()
	originalTourney.MatchTree = &GetNewMatches(3)[0]

	tourneyForDatabase := *originalTourney
	db.CreateTournament(&tourneyForDatabase)

	return originalTourney
}


func GetNewTeam() *Team {

	coach := CreateNewPerson()

	id, _ := uuid.NewV1()
	uniqueSuffix := strings.Split(id.String(), "-")[0]
	team := Team {
		Name:      "Name-" + uniqueSuffix,
		About:     "A few words about this amazing team",
		CoachID:   coach.ID,
		CoachName: coach.FirstName + " " + coach.LastName,
	}

	return &team
}

func CreateNewTeam() *Team {
	teamOriginal := GetNewTeam()

	teamToDataBase := *teamOriginal
	db.CreateTeam(&teamToDataBase)

	teamOriginal.ID = teamToDataBase.ID
	return teamOriginal
}


func GetNewPerson() *Person {

	id, _ := uuid.NewV4()
	postfixes := strings.Split(id.String(), "-")

	person := Person {
		FirstName: GenerateFirstName(),
		LastName:  GenerateLastName(),
		Email:     GenerateEmail(),
		Password:  postfixes[0],
	}

	return &person
}

func CreateNewPerson() *Person {

	original := GetNewPerson()
	forDatabase := *original
	db.CreatePerson(&forDatabase)

	original.ID = forDatabase.ID
	return original
}


func GetNewGame() *Game {
	id, _ := uuid.NewV1()
	uniqueSuffix := strings.Split(id.String(), "-")[0]
	game := Game {
		Title: "Game-title-" + uniqueSuffix,
		About: "Some text about useful things " + id.String(),
	}

	return &game
}

func CreateNewGame() *Game {
	gameOriginal := GetNewGame()

	gameToDataBase := *gameOriginal
	db.CreateGame(&gameToDataBase)

	gameOriginal.ID = gameToDataBase.ID
	return gameOriginal
}


func CreateNewMatches(deep int) (*Tournament, []MatchesTreeForm) {

	tourney := GetNewTournament()
	matches := GetNewMatches(deep)

	tourney.MatchTree = &matches[0]
	db.CreateTournament(tourney)

	return tourney, matches
}

func UpdateMatch(match *Match, organizeID uuid.UUID) {

	const magicNumber = 42

	endTime := match.StartTime.Add(magicNumber * time.Minute)
	match.EndTime = &endTime

	match.FirstTeamScore = rand.Intn(magicNumber)
	match.SecondTeamScore = rand.Intn(magicNumber)

	match.FirstTeamID = &CreateNewTeam().ID
	match.SecondTeamID = &CreateNewTeam().ID

	match.OrganizeID = organizeID
}

func ConvertToMatch(match *MatchesTreeForm, tourneyID uuid.UUID) *Match {
	return &Match {
		ID: match.ID,
		StartTime: &match.StartTime,
		TourneyID: tourneyID,
	}
}


func GetNewOAuth(person *Person, scope int) *OAuth {
	id, _ := uuid.NewV4()
	return &OAuth{
		Email: person.Email,
		Password: person.Password,
		AppID: id,
		Scope: scope,
	}
}

func CreateNewOAuth(scope int) *OAuth {
	person := CreateNewPerson()
	oauth := GetNewOAuth(person, scope)

	db.Auth(oauth)

	return oauth
}


func GetNewPlayers(numberOfPlayers int) []*Player {

	team := CreateNewTeam()

	players := make([]*Player, numberOfPlayers)
	for i := range players {
		players[i] = &Player{
			Nickname: GenerateLastName(),
			TeamID: team.ID,
			TeamName: team.Name,
			PersonID: CreateNewPerson().ID,
		}
	}

	return players
}

func CreateNewPlayers(numberOfPlayers int) []*Player {

	players := GetNewPlayers(numberOfPlayers)
	for _, player := range players {
		player.TeamName = ""
		db.AddPlayerInTeam(player)
	}

	return players
}





func GenerateEmail() string {
	id, _ := uuid.NewV4()
	postfixes := strings.Split(id.String(), "-")
	return tg.Generate(LastNameTemplate) + "_" + postfixes[1] + tg.Generate(MailsTemplate)
}
func GenerateFirstName() string {
	return tg.Generate(FirstNameTemplate)
}
func GenerateLastName() string {
	return tg.Generate(LastNameTemplate)
}

var tg = text_generator.New()
const FirstNameTemplate string = "{Vasya|Peter|Nikita|Sasha|Dmitriy|Enakentiy|John|Masha|Natasha|Tony}"
const LastNameTemplate string = "{Silaev|Kuzmin|Krasnov|Pupkin|Smirnov|Gorbenko|Trubnikov|Smirnova}"
const MailsTemplate string = "{@mail.ru|@yandex.ru|@gmail.com|@rambler.com}"

