package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetLiveMatchForTournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	tournamentID := r.URL.Query().Get("tournamentID")
	matchs, err := getMatchForTournament(db, tournamentID)
	if err == nil {
		w.Write([]byte("Error when fetching match for tournament"))
		fmt.Println(err)
	}

	body, err := json.Marshal(matchs)
	if err != nil {
		w.Write([]byte("Error when marshelling match for tournament"))
		fmt.Println("Error when marshelling match for tournament: %v", err)
	}
	w.Write(body)
}
