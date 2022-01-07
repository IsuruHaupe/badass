package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// )

// func tournamentController(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	switch r.Method {
// 	case "GET":
// 		//Check if this tournament exist 
// 		keys := make([]string, len(referees))

// 		i := 0
// 		for k := range referees {
// 			keys[i] = k
// 			i++
// 		}

// 		body, err := json.Marshal(keys)
// 		if err != nil {
// 			fmt.Println("error when marshelling in referee.go L.40 : %v", err)
// 		}
// 		w.Write(body)
// 	case "POST":

// }