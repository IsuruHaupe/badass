package main

func ParseEvent(event Event, match Match, sport string) {
	switch sport {
	case "BADMINTON":
		return ParseEventBadminton(event, match)
	}
}
