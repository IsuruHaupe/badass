package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/gorilla/websocket"
	"github.com/segmentio/ksuid"
)

/*
This handler is used to accept incoming connection
from watchers.

GET request  : get request is used to upgrade the http connection
to a websocket connection. It receives the match ID in the query param.

*/

func WatcherWsController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	switch r.Method {
	case "GET":
		// retrieve ID of the referee
		matchID := r.URL.Query().Get("matchID")
		fmt.Println("matchID =>", matchID)

		// Upgrade connection
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			log.Printf("Error when creating WebSocket connection : %v", err)
		}

		// get history of events
		/*previous_events, err := GetAllEvent(db)
		if err != nil {
			fmt.Println(err)
		}*/
		/*for _, event := range previous_events {
			err = wsutil.WriteServerMessage(conn, websocket.TextMessage, []byte(event.event))
			if err == nil {
				log.Printf("Failed : %v", err)
			}
		}*/

		// generate unique ID for this watcher
		watcherID := ksuid.New().String()
		referees[matchID][watcherID] = conn
		if err != nil {
			log.Printf("Failed to add connection to referee : %v", err)
			conn.Close()
		}

		fmt.Printf("Pool de watcher : \n %v \n", referees)
		err = wsutil.WriteServerMessage(conn, websocket.TextMessage, []byte("Connection succeed !\n"))
		// handle when connection is dead
		// delete the watcher from the map
		// and close connection
		if err != nil {
			poolOfWatchers := referees[matchID]
			delete(poolOfWatchers, watcherID)
			if _, ok := poolOfWatchers[watcherID]; ok {
				log.Printf("Failed to remove %v", err)
			}
			conn.Close()
		}
	default:
		w.Write([]byte("Unrecognised Query type !"))
		log.Printf("Unrecognised Query type !")
	}
}

/*
This handler is used to accept incoming connection
from referees it receives the referee ID a query param
to register him.
*/

func RefereeWsController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case "GET":
		// retrieve ID of the referee
		IdMatch := r.URL.Query().Get("IdMatch")
		fmt.Println("IdMatch =>", IdMatch)

		if _, ok := referees[IdMatch]; !ok {
			// init empty map of watcher for this referee ID
			referees[IdMatch] = make(map[string]net.Conn)
		} else {
			// reconnection of the referee
			// remove the refereeID from the pool of refereeID to remove in refereeToRemove
			delete(refereeToRemove, IdMatch)
		}
		fmt.Printf("Liste d'arbitres : \n %v \n", referees)

		// Upgrade connection
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			return
		}
		// retrieve file descriptor id
		fd, err := refereeEpoller.Add(conn)
		if err != nil {
			log.Printf("Failed to add connection to referee : %v", err)
			conn.Close()
		}

		// link the file descriptor to the refereeID
		refereeFdToString[fd] = IdMatch

		// check if match has been created previously
		_, err = getMatch(db, IdMatch)
		if err != nil {
			fmt.Println("ERROR : ", err)
			err = wsutil.WriteServerMessage(conn, websocket.TextMessage, []byte("ERROR : Match not found in database, please make sure to create the match before connecting to it using 'create-match' route.\n"))
			// handle when connection is dead
			// delete the watcher from the map
			// and close connection
			if err != nil {
				delete(referees, IdMatch)
				if _, ok := referees[IdMatch]; ok {
					log.Printf("Failed to remove %v", err)
				}
				_, err := refereeEpoller.Remove(conn)
				if err != nil {
					log.Printf("Failed to remove %v", err)
				}
				conn.Close()
			}
		} else {
			err = wsutil.WriteServerMessage(conn, websocket.TextMessage, []byte("Connection succeed !\n"))
			// handle when connection is dead
			// delete the watcher from the map
			// and close connection
			if err != nil {
				delete(referees, IdMatch)
				if _, ok := referees[IdMatch]; ok {
					log.Printf("Failed to remove %v", err)
				}
				_, err := refereeEpoller.Remove(conn)
				if err != nil {
					log.Printf("Failed to remove %v", err)
				}
				conn.Close()
			}
		}
	default:
		w.Write([]byte("Unrecognised Query type !"))
		log.Printf("Unrecognised Query type !")
	}
}
