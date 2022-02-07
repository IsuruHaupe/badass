package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gobwas/ws/wsutil"
	"github.com/gorilla/websocket"
)

/*
This function is used to receive referees updates on match.
We save the update and forward it to watchers
*/
func EventController() {
	for {
		// wait for a referees to send data
		connections, err := refereeEpoller.Wait()
		if err != nil {
			log.Printf("Failed to epoll wait %v", err)
			continue
		}
		// we loop over the referees that sent data
		for _, referee := range connections {
			if referee == nil {
				break
			}
			// TODO: remove the delay
			tts := time.Second
			time.Sleep(tts)
			// read referee data
			msg, _, err := wsutil.ReadClientData(referee)
			// case when referee connection is lost
			if err != nil {
				fmt.Printf("Erreur en essayant de lire les donnees du referee : %v \n", err)
				// remove connection from epoller
				fd, err := refereeEpoller.Remove(referee)
				if err != nil {
					log.Printf("Failed to remove %v", err)
				}
				// close connection
				referee.Close()
				// add the referee ID to be removed
				// we don't remove it now, because the referee might reconnect
				// using the same refereeID and we want to keep the pool of watcher alive
				refereeID := refereeFdToString[fd]
				refereeToRemove[refereeID] = refereeID
			} else {
				var decodedMsg Event
				err = json.Unmarshal(msg, &decodedMsg)
				//fmt.Println(decodedMsg)
				// case when we can't decode the message
				if err != nil {
					fmt.Printf("Erreur en essayant de décoder le message du referee : %v \n", err)
					fd, err := refereeEpoller.Remove(referee)
					if err != nil {
						log.Printf("Failed to remove %v", err)
					}
					// close connection
					referee.Close()
					// add the referee ID to be removed
					// we don't remove it now, because the referee might reconnect
					// using the same refereeID and we want to keep the pool of watcher alive
					refereeID := refereeFdToString[fd]
					refereeToRemove[refereeID] = refereeID
				} else {
					// save the event in the database
					// TODO: save data in a specific table or a specific ID
					//AddEvent(db, decodedMsg)
					match := ParseEvent(decodedMsg, "BADMINTON")

					// retrieve referee ID that sent the update
					IdMatch := decodedMsg.IdMatch
					// retrieve pool of watchers for this match/referee ID
					poolOfWatchers := referees[IdMatch]
					// if no watchers, we just save the data otherwise we loop
					// over the watchers and send them the update
					if len(poolOfWatchers) != 0 {
						for watcherID, watcherConn := range poolOfWatchers {
							// TODO : send previous events only to new watchers
							/*previousEvents, err := GetAllEvent(db)
							if err != nil {
								fmt.Println(err)
							}
							for _, event := range previousEvents {
								body, err := json.Marshal(event)
								if err != nil {
									fmt.Println("error when marshelling event to be send to watcher : %v", err)
								}
								err = wsutil.WriteServerMessage(watcherConn, websocket.TextMessage, body)
								if err != nil {
									delete(poolOfWatchers, watcherID)
									if _, ok := poolOfWatchers[watcherID]; ok {
										log.Printf("Failed to remove %v", err)
									}
									watcherConn.Close()
								}
							}*/

							// send new data
							err = wsutil.WriteServerMessage(watcherConn, websocket.TextMessage, match)
							// handle when connection is dead
							// delete the watcher from the map
							// and close connection
							if err != nil {
								delete(poolOfWatchers, watcherID)
								if _, ok := poolOfWatchers[watcherID]; ok {
									log.Printf("Failed to remove %v", err)
								}
								watcherConn.Close()
							}
						}
						//fmt.Printf("Previous events : \n %v\n", previous_events)
						/*for _, watcher := range poolOfWatchers.connections {
							// recuperer l'event dans la bdd
							for _, event := range previous_events {
								err = wsutil.WriteServerMessage(watcher, websocket.TextMessage, []byte(event.event))
								if err != nil {
									if err := poolOfWatchers.Remove(watcher); err != nil {
										log.Printf("Failed to remove %v", err)
									}
									watcher.Close()
								}
							}
						}*/
					}
				}
			}
		}
	}
}
