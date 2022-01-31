package main

import (
	"encoding/json"
	"log"
)

type Equipe struct {
	Score       int `json:"Score"`
	FaultNumber int `json:"FaultNumber"`
}
type Badminton struct {
	EquipeA Equipe `json:"EquipeA"`
	EquipeB Equipe `json:"EquipeB"`
	Status  string `json:"Status"`
}

//value of event
type Point struct {
	Point int `json:"Point"`
}

// TODO utiliser joueur et commentaire
type Fault struct {
	Player     string `json:"Player"`
	Comment    string `json:"Comment"`
	FaultValue int    `json:"FaultValue"`
}

// function to treat each events for this sport
func ParseEventBadminton(event Event, match Match) []byte {
	// badminton := Badminton{}
	var badminton Badminton
	json.Unmarshal([]byte(match.matchValues), &badminton)
	switch event.EventType {
	case "POINT":
		point := Point{}
		json.Unmarshal([]byte(event.EventValue), &point)
		if event.Equipe == "EQUIPEA" {
			badminton.EquipeA.Score += point.Point
		} else {
			badminton.EquipeB.Score += point.Point
		}
	case "FAULT":
		fault := Fault{}
		json.Unmarshal([]byte(event.EventValue), &fault)
		if event.Equipe == "EQUIPEA" {
			// in case of cancel we use fault.FaultValue
			badminton.EquipeA.FaultNumber += fault.FaultValue
		} else {
			badminton.EquipeB.FaultNumber += fault.FaultValue
		}
	case "BEGIN_MATCH":
		badminton.Status = "MATCH_IN_PROGRESS"
	case "END_MATCH":
		badminton.Status = "END_MATCH"

	}
	tmp, err := json.Marshal(badminton)
	if err != nil {
		log.Fatal("error when marshelling in referee.go: %v", err)
	}
	match.matchValues = string(tmp)
	err = UpdateMatch(db, match)
	if err != nil {
		log.Fatal("Error update match : %v", err)
	}
	return tmp
}

// database related
func InitializeBadminton() string {
	badminton := Badminton{
		EquipeA: Equipe{
			Score:       0,
			FaultNumber: 0,
		},
		EquipeB: Equipe{
			Score:       0,
			FaultNumber: 0,
		},
		Status: "CREATE",
	}
	tmp, err := json.Marshal(badminton)
	if err != nil {
		log.Fatal("error initialize badminton struct: %v", err)
	}
	return string(tmp)

}
