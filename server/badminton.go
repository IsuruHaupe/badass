package main

import (
	"encoding/json"
	"fmt"
)

type Set struct {
	EquipeA Equipe `json:"EquipeA"`
	EquipeB Equipe `json:"EquipeB"`
	Status  string `json:"Status"`
}

type Equipe struct {
	Score       int `json:"Score"`
	FaultNumber int `json:"FaultNumber"`
}
type Badminton struct {
	Sets   []Set  `json:"Sets"`
	Status string `json:"Status"`
}

//value of event
type Point struct {
	Point int `json:"Point"`
}

// TODO utiliser joueur et commentaire
type Fault struct {
	FaultValue int `json:"FaultValue"`
}

// function to treat each events for this sport
func ParseEventBadminton(event Event, match Match) []byte {
	var badminton Badminton
	json.Unmarshal([]byte(match.MatchValues), &badminton)
	switch event.EventType {
	case "POINT":
		point := Point{}
		json.Unmarshal([]byte(event.EventValue), &point)
		if event.Equipe == "EQUIPEA" {
			// update last sets score
			numberSet := len(badminton.Sets)
			badminton.Sets[numberSet-1].EquipeA.Score += point.Point
		} else {
			// update last sets score
			numberSet := len(badminton.Sets)
			badminton.Sets[numberSet-1].EquipeB.Score += point.Point
		}
	case "FAULT":
		fault := Fault{}
		json.Unmarshal([]byte(event.EventValue), &fault)
		if event.Equipe == "EQUIPEA" {
			// update last sets score
			numberSet := len(badminton.Sets)
			// in case of cancel we use fault.FaultValue
			badminton.Sets[numberSet-1].EquipeA.FaultNumber += fault.FaultValue
		} else {
			// update last sets score
			numberSet := len(badminton.Sets)
			badminton.Sets[numberSet-1].EquipeB.FaultNumber += fault.FaultValue
		}
	case "BEGIN_MATCH":
		badminton.Status = "SET_1"
	case "END_MATCH":
		badminton.Status = "END_MATCH"
	case "NEW_SET":
		badminton.Sets = append(badminton.Sets, Set{
			EquipeA: Equipe{
				Score:       0,
				FaultNumber: 0,
			},
			EquipeB: Equipe{
				Score:       0,
				FaultNumber: 0,
			},
			Status: "SET_IN_PROGRESS",
		})
		badminton.Status = "SET_" + fmt.Sprintf("%d", len(badminton.Sets))
	case "END_SET":
		// update last sets status
		numberSet := len(badminton.Sets)
		badminton.Sets[numberSet-1].Status = "END_SET"

	}
	tmp, err := json.Marshal(badminton)
	if err != nil {
		fmt.Println("error when marshelling in badminton.go: %v", err)
	}
	match.MatchValues = string(tmp)
	err = UpdateMatch(db, match)
	if err != nil {
		fmt.Println("Error update match : %v", err)
	}
	return tmp
}

// database related
func InitializeBadminton() string {
	// create a badminton match and create a new set
	badminton := Badminton{
		Sets: []Set{
			Set{
				EquipeA: Equipe{
					Score:       0,
					FaultNumber: 0,
				},
				EquipeB: Equipe{
					Score:       0,
					FaultNumber: 0,
				},
				Status: "SET_IN_PROGRESS",
			},
		},
		Status: "NOT_BEGN",
	}
	tmp, err := json.Marshal(badminton)
	if err != nil {
		fmt.Println("error initialize badminton struct: %v", err)
	}
	return string(tmp)

}
