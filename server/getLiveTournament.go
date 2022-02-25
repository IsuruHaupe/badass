package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetLiveTournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	tournaments, err := getAllTournament(db)
	if err != nil {
		w.Write([]byte("Error when fetching tournaments"))
		fmt.Println("Error when fetching tournaments: %v", err)
	}
	//keys := make([]string, len(tournaments))

	/*i := 0
	for _, tournament := range tournaments {
		keys[i] = tournament.ID
		i++
	}*/

	body, err := json.Marshal(tournaments)
	if err != nil {
		w.Write([]byte("Error when marshelling tournament ID"))
		fmt.Println("error when marshelling : %v", err)
	}
	w.Write(body)
}
