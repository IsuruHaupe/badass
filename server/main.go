package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		//log.Fatal("$PORT must be set")
		port = "8000"
	}
	// export credentials in order to connect to the DB (non docker)
	SetBDDEnvironmentVariable()
	// connect to the DB using the exported credentials
	db = ConnectToDB()

	// upagrade connection limit for the epoll
	UpgradeConnectionLimit()

	// Start epoll
	var err error
	// epoll for the referee
	// this will handle the load when a lot of referee sends data
	refereeEpoller, err = MkEpoll()
	if err != nil {
		panic(err)
	}
	// we link a referee to a map of unique connection
	// when we receive update from a referee we iterate over the
	// array of watcher (i.e connection) and send them the news
	referees = make(map[string]map[string]net.Conn)
	refereeFdToString = make(map[int]string)
	refereeToRemove = make(map[string]string)
	// this go routine waits for updates from the pool of referees
	// and handle them by saving the update in the DB and sends
	// notification to watchers
	go EventController()
	go RefereeGarbageCollector()

	// this route creates a match in the db using param query
	// and returns an unique match id linking this match in db
	http.HandleFunc("/create-match", InitMatch)
	// this route creates a tournament in the db using the param query
	// and returns an unique tournament id linking this tounament in db
	http.HandleFunc("/create-tournament", InitTournament)
	// this route is use by the referee to initiate a websocket
	// and send data over it
	http.HandleFunc("/referee", RefereeWsController)
	// this route is use by the watcher to initiaite a websocket
	// and receive data over it
	http.HandleFunc("/spectateur", WatcherWsController)
	// GET request to get live match ID
	http.HandleFunc("/live-match", GetLiveMatch)
	// GET request to get all the match for a given tournament ID
	// given in the query param
	http.HandleFunc("/live-tournament", GetLiveMatchForTournament)
	// TODO faire un getTournament pour recuperer l'id de tous les
	// tournois en live
	http.HandleFunc("/", HelloServer)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
