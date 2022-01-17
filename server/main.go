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

	http.HandleFunc("/create-match", InitMatch)
	http.HandleFunc("/referee", RefereeWsController)
	http.HandleFunc("/spectateur", WatcherWsController)
	http.HandleFunc("/live-match", GetLiveMatch)
	http.HandleFunc("/", HelloServer)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
