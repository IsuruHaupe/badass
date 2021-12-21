package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gobwas/ws"
	"github.com/segmentio/ksuid"
)

/*
This handler is used to accept incoming connection
from watchers. It will handle two type of requests.

GET request  : get request is used to upgrade the http connection
to a websocket connection.

POST request : post request is used to post the refereeID (i.e. match ID)
the watcher wants to follow. During the handshake we cannot send data as
it is a get request and we cannot change that behabior. So we use a second
request to send the match ID to the server.

*/

type matchFollowed struct {
	mu    sync.Mutex
	match Match
}

var matchToFollow matchFollowed

func WatcherWsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Upgrade connection
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			return
		}

		// get history of events
		/*previous_events, err := GetAllEvent(db)
		if err != nil {
			log.Fatal(err)
		}*/
		/*for _, event := range previous_events {
			err = wsutil.WriteServerMessage(conn, websocket.TextMessage, []byte(event.event))
			if err == nil {
				log.Printf("Failed : %v", err)
			}
		}*/

		// generate unique ID for this watcher
		watcherID := ksuid.New().String()
		// add the watcher in the map of the referee
		fmt.Println(matchToFollow.match.ID)
		referees[matchToFollow.match.ID][watcherID] = conn
		if err != nil {
			log.Printf("Failed to add connection to referee : %v", err)
			conn.Close()
		}

		// unlock the variable after we saved the name of the referee
		matchToFollow.mu.Unlock()

		fmt.Printf("Pool de watcher : \n %v \n", referees)

	case "POST":
		// the match to follow is saved globally
		err := json.NewDecoder(r.Body).Decode(&matchToFollow.match)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// we lock the global variable until we saved
		// the id of the referee
		matchToFollow.mu.Lock()
	}
}

/*
This handler is used to accept incoming connection
from referees. It will handle two type of requests.

GET request  : get request is used to upgrade the http connection
to a websocket connection.

POST request : post request is used to post the refereeID.
During the handshake we cannot send data as it is a get
request and we cannot change that behabior. So we use a second
request to send the referee ID to the server.
It is then saved as a new entry in the map 'referees'

*/
func RefereeWsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Upgrade connection
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			return
		}
		// retrieve file descriptor id
		_, err = refereeEpoller.Add(conn)
		if err != nil {
			log.Printf("Failed to add connection to referee : %v", err)
			conn.Close()
		}
		// create a new pool of watchers for this referee
		//watchersEpoller, err := MkEpoll()
		/*if err != nil {
			panic(err)
		}*/
		// link referee id to watchers pool
		//referees[fd] = watchersEpoller
	case "POST":
		var referee RefereeID
		// retrieve referee ID
		err := json.NewDecoder(r.Body).Decode(&referee)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// init empty map of watcher for this referee ID
		referees[referee.ID] = make(map[string]net.Conn)
		fmt.Printf("Referee ID : %+v", referee)
		fmt.Printf("List d'arbitre : \n %v \n", referees)
	}
}

func GetLiveMatch(w http.ResponseWriter, r *http.Request) {
	keys := make([]string, len(referees))

	i := 0
	for k := range referees {
		keys[i] = k
		i++
	}

	body, err := json.Marshal(keys)
	if err != nil {
		fmt.Println("error when marshelling in referee.go L.40 : %v", err)
	}
	w.Write(body)
}
