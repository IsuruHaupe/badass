package main

import "fmt"

func ParseEvent(event Event) []byte {
	// retrieve the match to which this event is tied
	match, err := getMatch(db, event.IdMatch)
	tournament, err := GetTournament(db, match.Tournament)
	if err != nil {
		fmt.Errorf("Parse event error get match : %v", err)
	}
	// treat each sport accordingly
	switch tournament.sport {
	case "BADMINTON":
		return ParseEventBadminton(event, match)
	case "FOOTBALL":
		return ParseEventFootball(event, match)
	}
	return nil
}

// function to init a sport in the database
func InitializeSport(sport string) string {
	switch sport {
	case "BADMINTON":
		return InitializeBadminton()
	case "FOOTBALL":
		return InitializeFootball()
	}
	return ""
}
