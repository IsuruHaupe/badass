package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Equipe struct {
	Score       int `json:"Score"`
	FaultNumber int `json:"FaultNumber"`
}
type Badminton struct {
	EquipeA Equipe `json:"EquipeA"`
	EquipeB Equipe `json:"EquipeB"`
}

//value of event
type Point struct {
	Point int `json:"Point"`
}

func ParseEventBadminton(event Event, match Match) []byte {
	// badminton := Badminton{}
	var badminton Badminton
	json.Unmarshal([]byte(match.matchValues), &badminton)
	switch event.EventType {
	case "POINT":
		point := Point{}
		fmt.Println(event.EventValue)
		json.Unmarshal([]byte(event.EventValue), &point)
		fmt.Println("POINT UNMARSHAL : ", point)
		if event.Equipe == "EQUIPEA" {
			badminton.EquipeA.Score += point.Point
		} else {
			badminton.EquipeB.Score += point.Point
		}
	case "FAULT":

	case "FIN_MATCH":

	}
	tmp, err := json.Marshal(badminton)
	if err != nil {
		log.Fatal("error when marshelling in referee.go L.112 : %v", err)
	}
	match.matchValues = string(tmp)
	err = UpdateMatch(db, match)
	if err != nil {
		log.Fatal("Error update match : %v", err)
	}
	return tmp
}
