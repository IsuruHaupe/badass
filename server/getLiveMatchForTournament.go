package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GetLiveMatchForTournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	tournamentID := r.URL.Query().Get("tournamentID")
	matchs, err := getMatchForTournament(db, tournamentID)
	if err == nil {
		log.Fatal(err)
	}

	body, err := json.Marshal(matchs)
	if err != nil {
		fmt.Println("error when marshelling in referee.go L.40 : %v", err)
	}
	w.Write(body)
}
