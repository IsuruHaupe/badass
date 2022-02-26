package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetLiveMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	matchs, err := getAllMatch(db)
	if err != nil {
		w.Write([]byte("Error when fetching match"))
		fmt.Println(err)
	}

	body, err := json.Marshal(matchs)
	if err != nil {
		w.Write([]byte("Error when marshelling live match"))
		fmt.Println("Error when marshelling live match : %v", err)
	}
	w.Write(body)
}
