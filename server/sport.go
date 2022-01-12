package main

import "fmt"

func ParseEvent(event Event, sport string) Match {
	match, err := getMatch(db, event.idMatch)
	if err != nil {
		fmt.Errorf("Parse event error get match : %v", err)
	}
	switch sport {
	case "BADMINTON":
		return ParseEventBadminton(event, match)
	}
	return Match{}
}
