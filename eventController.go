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
				if err := refereeEpoller.Remove(referee); err != nil {
					log.Printf("Failed to remove %v", err)
				}
				referee.Close()
			} else {
				var decodedMsg Event
				err = json.Unmarshal(msg, &decodedMsg)
				//fmt.Println(decodedMsg)
				// case when we can't decode the message
				if err != nil {
					fmt.Printf("Erreur en essayant de décoder le message du referee : %v \n", err)
					if err := refereeEpoller.Remove(referee); err != nil {
						log.Printf("Failed to remove %v", err)
					}
					referee.Close()
				} else {
					// save the event in the database
					// TODO: save data in a specific table or a specific ID
					AddEvent(db, decodedMsg)

					// retrieve referee ID that sent the update
					refereeID := decodedMsg.Referee.ID
					// retrieve pool of watchers for this match/referee ID
					poolOfWatchers := referees[refereeID]
					// if no watchers, we just save the data otherwise we loop
					// over the watchers and send them the update
					if len(poolOfWatchers) != 0 {
						for watcherID, watcherConn := range poolOfWatchers {
							// TODO : send previous events only to new watchers
							/*previousEvents, err := GetAllEvent(db)
							if err != nil {
								log.Fatal(err)
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
							err = wsutil.WriteServerMessage(watcherConn, websocket.TextMessage, msg)
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
