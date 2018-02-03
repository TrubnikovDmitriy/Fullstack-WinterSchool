package models

import (
	"../services"
	"time"
	"github.com/satori/go.uuid"
	"fmt"
)

type MatchesTreeForm struct {

	LeftChild *MatchesTreeForm  `json:"prev_match_1"`
	RightChild *MatchesTreeForm `json:"prev_match_2"`
	StartTime time.Time         `json:"started"`

	ID uuid.UUID 				`json:"-"`
	LeftID *uuid.UUID 			`json:"-"`
	RightID *uuid.UUID 			`json:"-"`
	arrayMatches []Match 		`json:"-"`
}

func (match *MatchesTreeForm) Validate() *serv.ErrorCode {
	ttl := serv.MaxMatchesInTournament
	return match.recursiveValidate(match.StartTime, &ttl)
}

// Контролируемая рекурсия, чтобы проверить валидность турнирной сетки матчей
func (match *MatchesTreeForm) recursiveValidate(parentTime time.Time, ttl *int) *serv.ErrorCode {
	if *ttl == 0 {
		return serv.NewBadRequest("Too many matches")
	}
	if *ttl != serv.MaxMatchesInTournament {
		if match.StartTime.After(parentTime) || match.StartTime.Equal(parentTime) {
			return serv.NewBadRequest("Next matches begin before previous")
		}
	}
	*ttl--
	if (match.LeftChild == nil) && (match.RightChild == nil) {
		return nil
	}
	if match.LeftChild == match.RightChild {
		return serv.NewBadRequest("Match is identical")
	}
	if (match.LeftChild != nil) && (match.RightChild == nil) {
		return serv.NewBadRequest("Tree is not binary")
	}
	if (match.LeftChild == nil) && (match.RightChild != nil) {
		return serv.NewBadRequest("Tree is not binary")
	}
	err := match.LeftChild.recursiveValidate(match.StartTime, ttl)
	if err != nil {
		return err
	}
	err = match.RightChild.recursiveValidate(match.StartTime, ttl)
	return err
}

// Разворачивает бинарное дерево в массив
func (match *MatchesTreeForm) GetNodesCount() int {
	count := 0
	match.counting(&count)
	return count
}

func (match *MatchesTreeForm) counting(count *int) {
	*count++
	if match.LeftChild != nil {
		match.LeftChild.counting(count)
	}
	if match.RightChild != nil {
		match.RightChild.counting(count)
	}
	return
}

func (match *MatchesTreeForm) CreateArrayMatch() []Match {

	// Дать каждой ноде уникальный ID
	match.setIDs()

	// Разложить в массив и связать между собой указателями на UUID
	var arrayMatches MatchesArrayForm
	match.creatingArrayMatches(&arrayMatches)
	arrayMatches.setNextIDs()

	return arrayMatches.Array
}

func (match *MatchesTreeForm) setIDs() {
	id, err := uuid.NewV4()
	if err != nil {
		fmt.Print(err)
	}
	match.ID = id
	if match.LeftChild != nil {
		match.LeftChild.setIDs()
	}
	if match.RightChild != nil {
		match.RightChild.setIDs()
	}
}

func (match *MatchesTreeForm) creatingArrayMatches(arrayMatches *MatchesArrayForm) {
	var leftID, rightID *uuid.UUID = nil, nil
	if match.LeftChild != nil {
		leftID = &match.LeftChild.ID
	}
	if match.RightChild != nil {
		rightID = &match.RightChild.ID
	}
	arrayMatches.Array = append(arrayMatches.Array, Match {
		ID:         match.ID,
		PrevMatch1: leftID,
		PrevMatch2: rightID,
		StartTime:  &match.StartTime,
	})
	if match.LeftChild != nil {
		match.LeftChild.creatingArrayMatches(arrayMatches)
	}
	if match.RightChild != nil {
		match.RightChild.creatingArrayMatches(arrayMatches)
	}
}

func (arrMatches *MatchesArrayForm) setNextIDs() {
	for _, match := range arrMatches.Array {
		if match.PrevMatch1 != nil && match.PrevMatch2 != nil {
			for i, prevMatch := range arrMatches.Array {
				if prevMatch.ID == *match.PrevMatch1 || prevMatch.ID == *match.PrevMatch2 {
					nextMatch := match.ID
					arrMatches.Array[i].NextMatch = &nextMatch
				}
			}
		}
	}
}