package main

import (
	"encoding/json"
	"fmt"
)

type EquipeFootball struct {
	Score int           `json:"Score"`
	Fault FaultFootball `json:"Fault"`
}

type FaultFootball struct {
	NumberOfRedCard    int `json:"NumberOfRedCard"`
	NumberOfYellowCard int `json:"NumberOfYellowCard"`
}
type Football struct {
	EquipeA         EquipeFootball  `json:"EquipeA"`
	EquipeB         EquipeFootball  `json:"EquipeB"`
	Status          string          `json:"Status"` //("NOT_BEGUN" "FIRST_HALF" "HALF" "SECOND_HALF" "EXTENSION" "PENALTY_SHOOTOUT" "END_MATCH")
	PenaltyShootout PenaltyShootout `json:"PenaltyShootout"`
}

type PenaltyShootout struct {
	ScoreEquipeA int `json:"ScoreEquipeA"`
	ScoreEquipeB int `json:"ScoreEquipeB"`
}

//Event struct :
type Event_Football_Point struct {
	Point int `json:"Point"`
}
type Event_Football_Fault struct {
	FaultValue int `json:"FaultValue"`
}

// function to treat each events for this sport
func ParseEventFootball(event Event, match Match) []byte {
	// badminton := Badminton{}
	var football Football
	json.Unmarshal([]byte(match.matchValues), &football)
	switch event.EventType {
	case "POINT":
		point := Event_Football_Point{}
		json.Unmarshal([]byte(event.EventValue), &point)
		if event.Equipe == "EQUIPEA" {
			football.EquipeA.Score += point.Point
		} else {
			football.EquipeB.Score += point.Point
		}
	case "REDCARD":
		fault := Event_Football_Fault{}
		json.Unmarshal([]byte(event.EventValue), &fault)
		if event.Equipe == "EQUIPEA" {
			// in case of cancel we use fault.FaultValue
			football.EquipeA.Fault.NumberOfRedCard += fault.FaultValue
		} else {
			football.EquipeB.Fault.NumberOfRedCard += fault.FaultValue
		}
	case "YELLOWCARD":
		fault := Event_Football_Fault{}
		json.Unmarshal([]byte(event.EventValue), &fault)
		if event.Equipe == "EQUIPEA" {
			// in case of cancel we use fault.FaultValue
			football.EquipeA.Fault.NumberOfYellowCard += fault.FaultValue
		} else {
			football.EquipeB.Fault.NumberOfYellowCard += fault.FaultValue
		}
	case "POINT_PENALTY_SHOOTOUT":
		point := Event_Football_Point{}
		json.Unmarshal([]byte(event.EventValue), &point)
		if event.Equipe == "EQUIPEA" {
			// in case of cancel we use fault.FaultValue
			football.PenaltyShootout.ScoreEquipeA += point.Point
		} else {
			football.EquipeB.Fault.NumberOfYellowCard += point.Point
		}
	case "BEGIN_MATCH":
		football.Status = "FIRST_HALF"
	case "HALF":
		football.Status = "HALF"
	case "SECOND_HALF":
		football.Status = "SECOND_HALF"
	case "EXTENSION":
		football.Status = "EXTENSION"
	case "PENALTY_SHOOTOUT":
		football.Status = "PENALTY_SHOOTOUT"
	case "END_MATCH":
		football.Status = "END_MATCH"

	}
	tmp, err := json.Marshal(football)
	if err != nil {
		fmt.Println("error when marshelling in football.go: %v", err)
	}
	match.matchValues = string(tmp)
	err = UpdateMatch(db, match)
	if err != nil {
		fmt.Println("Error update match : %v", err)
	}
	return tmp
}

// database related
func InitializeFootball() string {
	football := Football{
		EquipeA: EquipeFootball{
			Score: 0,
			Fault: FaultFootball{
				NumberOfRedCard:    0,
				NumberOfYellowCard: 0,
			},
		},
		EquipeB: EquipeFootball{
			Score: 0,
			Fault: FaultFootball{
				NumberOfRedCard:    0,
				NumberOfYellowCard: 0,
			},
		},
		Status: "NOT_BEGUN",
		PenaltyShootout: PenaltyShootout{
			ScoreEquipeA: 0,
			ScoreEquipeB: 0,
		},
	}
	tmp, err := json.Marshal(football)
	if err != nil {
		fmt.Println("error initialize football struct: %v", err)
	}
	return string(tmp)
}
