package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetLiveMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	keys := make([]string, len(referees))

	i := 0
	for k := range referees {
		keys[i] = k
		i++
	}

	body, err := json.Marshal(keys)
	if err != nil {
		w.Write([]byte("Error when marshelling live match"))
		fmt.Println("Error when marshelling live match : %v", err)
	}
	w.Write(body)
}
