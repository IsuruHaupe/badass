package main

import (
	"encoding/json"
	"fmt"
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
	point int
}

func ParseEventBadminton(event Event, match Match) Match {
	// badminton := Badminton{}
	var badminton Badminton
	json.Unmarshal([]byte(match.matchValues), &badminton)
	switch event.EventType {
	case "POINT":
		point := Point{}
		json.Unmarshal([]byte(event.EventValue), &point)
		if event.Equipe == "EQUIPEA" {
			badminton.EquipeA.Score += point.point
		} else {
			badminton.EquipeB.Score += point.point
		}
	case "FAULT":

	case "FIN_MATCH":

	}
	tmp, err := json.Marshal(badminton)
	fmt.Println(tmp)
	if err != nil {
		fmt.Println("error when marshelling in referee.go L.112 : %v", err)
	}
	match.matchValues = string(tmp)
	return match
}
