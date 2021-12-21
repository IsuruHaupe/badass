package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gobwas/ws/wsutil"
	"github.com/gorilla/websocket"
)

func RefereeHandler() {
	// l'objectif est que quand on reçoit une MAJ de l'arbitre
	//   on l'enregistre. Puis on l'envoie aux spectateurs concernes
	for {
		connections, err := refereeEpoller.Wait()
		if err != nil {
			log.Printf("Failed to epoll wait %v", err)
			continue
		}
		// on itere sur les arbitres ayant envoye des donnees
		for _, referee := range connections {
			if referee == nil {
				break
			}
			tts := time.Second
			time.Sleep(tts)
			msg, _, err := wsutil.ReadClientData(referee)
			var decodedMsg Event
			err = json.Unmarshal(msg, &decodedMsg)
			//fmt.Println(decodedMsg)
			if err != nil {
				fmt.Printf("Erreur en essayant de lire les donnees du referee : %v", err)
				if err := refereeEpoller.Remove(referee); err != nil {
					log.Printf("Failed to remove %v", err)
				}
				referee.Close()
			} else {
				// sauvegarder l'event dans la bdd
				AddEvent(db, decodedMsg)

				//envoyer les MAJ au watcher
				// on recupere l'id du referee
				refereeID := decodedMsg.Referee.ID
				// on itere sur les connexions et on envoie un message
				poolOfWatchers := referees[refereeID]
				// si pas de spectateur, on sauvegarde les events dans la bdd simplement
				if len(poolOfWatchers) != 0 {
					/*previous_events, err := GetAllEvent(db)
					if err != nil {
						log.Fatal(err)
					}*/
					for watcherID, watcherConn := range poolOfWatchers {
						// recuperer l'event dans la bdd
						//websocket.TextMessage
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
			// envoyer la MAJ au spectateur
			// il faut discriminer les spectateurs en fonction des matchs qu'ils regardent

			// send message example
			//log.Printf("Server sending message")
			//err = wsutil.WriteServerMessage(conn, websocket.TextMessage, []byte(event[rand.Intn(len(event))]))
			//if err != nil {
			//	if err := epoller.Remove(conn); err != nil {
			//		log.Printf("Failed to remove %v", err)
			//	}
			//	conn.Close()
			//}

		}
	}
}

/*func Start() {
	event := [5]string{"match created",
		"updates on score 1",
		"updates on timeout",
		"updates on score 2",
		"math ended"}
	for {
		connections, err := epoller.Wait()
		if err != nil {
			log.Printf("Failed to epoll wait %v", err)
			continue
		}
		for _, conn := range connections {
			if conn == nil {
				break
			}
			tts := time.Second
			time.Sleep(tts)
			// receive message
			#msg, _, err := wsutil.ReadClientData(conn)
			#if err != nil {
			#	if err := epoller.Remove(conn); err != nil {
			#		log.Printf("Failed to remove %v", err)
			#	}
			#	conn.Close()
			#} else {
			#	// This is commented out since in demo usage, stdout is showing messages sent from > 1M connections at very high rate
			#	log.Printf("msg: %s", string(msg))
			#}

			// send message
			log.Printf("Server sending message")
			err = wsutil.WriteServerMessage(conn, websocket.TextMessage, []byte(event[rand.Intn(len(event))]))
			if err != nil {
				if err := epoller.Remove(conn); err != nil {
					log.Printf("Failed to remove %v", err)
				}
				conn.Close()
			}
		}
	}
}*/
