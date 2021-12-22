package main

import (
	"log"
	"net"
	"net/http"
)

func main() {
	// export credentials in order to connect to the DB
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

	// this go routine waits for updates from the pool of referees
	// and handle them by saving the update in the DB and sends
	// notification to watchers
	go EventController()

	http.HandleFunc("/referee", RefereeWsController)
	http.HandleFunc("/referee/register", RefereeWsController)
	http.HandleFunc("/spectateur", WatcherWsController)
	http.HandleFunc("/live-match", GetLiveMatch)

	if err := http.ListenAndServe("0.0.0.0:8000", nil); err != nil {
		log.Fatal(err)
	}
}
