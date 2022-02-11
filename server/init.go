package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/segmentio/ksuid"
)

func InitMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	switch r.Method {
	case "GET":
		equipeA := r.URL.Query().Get("equipeA")
		equipeB := r.URL.Query().Get("equipeB")
		tournamentID := r.URL.Query().Get("tournamentID")
		// generate unique ID for this watcher
		matchID := ksuid.New().String()

		// TODO : change this for a general sports
		match := Match{
			id:          matchID,
			equipeA:     equipeA,
			equipeB:     equipeB,
			tournament:  tournamentID,
			matchValues: InitializeSport("BADMINTON"),
		}
		fmt.Println(match)
		// create the match in database
		CreateMatch(db, match)
		// send ID of the match to client
		w.Write([]byte(matchID))
	default:
		w.Write([]byte("Unrecognised Query type !"))
		log.Printf("Unrecognised Query type !")
	}
}

func InitTournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	switch r.Method {
	case "GET":
		tournamentName := r.URL.Query().Get("tournamentName")
		tournamentSport := r.URL.Query().Get("sport")
		// generate unique ID for this watcher
		tournamentID := ksuid.New().String()
		// TODO : change this for a general sports
		tournament := Tournament{
			ID:    tournamentID,
			name:  tournamentName,
			sport: tournamentSport,
		}
		// create the match in database
		CreateTournament(db, tournament)
		// send ID of the match to client
		w.Write([]byte(tournamentID))
	default:
		w.Write([]byte("Unrecognised Query type !"))
		log.Printf("Unrecognised Query type !")
	}
}
