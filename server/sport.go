package main

import "fmt"

func ParseEvent(event Event, sport string) []byte {
	match, err := getMatch(db, event.IdMatch)
	if err != nil {
		fmt.Errorf("Parse event error get match : %v", err)
	}
	fmt.Println("MATCH", match)
	switch sport {
	case "BADMINTON":
		return ParseEventBadminton(event, match)
	}
	return nil
}
