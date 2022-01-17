package main

import (
	"log"
	"net/http"

	"github.com/segmentio/ksuid"
)

func InitMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case "GET":
		// generate unique ID for this watcher
		matchID := ksuid.New().String()
		// TODO : change this for a general sports
		match := Match{
			id:          matchID,
			equipeA:     "equipeA",
			equipeB:     "equipeB",
			tournament:  0,
			matchValues: InitializeSport("BADMINTON"),
		}
		// create the match in database
		CreateMatch(db, match)
		// send ID of the match to client
		w.Write([]byte(matchID))
	default:
		w.Write([]byte("Unrecognised Query type !"))
		log.Printf("Unrecognised Query type !")
	}
}
