package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var (
	ip = flag.String("ip", "warm-dusk-64603.herokuapp.com", "server IP")
	//ip          = flag.String("ip", "127.0.0.1", "server IP")
	connections = flag.Int("conn", 1, "number of websocket connections")
)

type Match struct {
	Id          string
	EquipeA     string
	EquipeB     string
	Tournament  string
	MatchValues string
}

type Tournament struct {
	ID    string
	Name  string
	Sport string
}

var matchs []string
var tournament []Match
var tournaments []Tournament

func main() {
	//u := url.URL{Scheme: "http", Host: *ip + ":8000", Path: "/get-live-match"}
	u := url.URL{Scheme: "http", Host: *ip, Path: "/get-live-match"}
	fmt.Println(u)
	getLiveMatch(u.String())

	/*u = url.URL{Scheme: "http", Host: *ip + ":8000", Path: "/get-live-match-for-tournament"}
	fmt.Println(u)

	params := url.Values{}
	params.Add("tournamentID", "25b879seX7reJaK7Ts1G7LWHzGU")
	u.RawQuery = params.Encode()

	getLiveTournament(u.String())
	fmt.Println("TOUNRNAMENT", tournament)

	u = url.URL{Scheme: "http", Host: *ip + ":8000", Path: "/get-live-tournament"}
	getAllLiveTournament(u.String())
	fmt.Println("ALL TOURNAMENT", tournaments)
	*/
	//fmt.Println(matchs)
	initWatcher(matchs[0])
}

func getAllLiveTournament(url string) {
	resp, err := http.Get(url)
	if err != nil {

		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &tournaments)
	if err != nil {
		fmt.Println("error when marshelling in watcher.go : %v", err)
	}
}

func getLiveTournament(url string) {
	resp, err := http.Get(url)
	if err != nil {

		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &tournament)
	if err != nil {
		fmt.Println("error when marshelling in watcher.go : %v", err)
	}
}

func getLiveMatch(url string) {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, &matchs)
	if err != nil {
		fmt.Println("error when marshelling in client.go : %v", err)
	}
}

func initWatcher(matchID string) {

	flag.Usage = func() {
		io.WriteString(os.Stderr, `Websockets client generator Example usage: ./client -ip=172.17.0.1 -conn=10`)
		flag.PrintDefaults()
	}
	flag.Parse()

	rand.Seed(time.Now().Unix())
	//WEBSOCKET
	//u := url.URL{Scheme: "ws", Host: *ip + ":8000", Path: "/spectateur"}
	u := url.URL{Scheme: "ws", Host: *ip, Path: "/spectateur"}
	// add match ID to URL
	params := url.Values{}
	params.Add("matchID", matchID)
	u.RawQuery = params.Encode()
	log.Printf("Connecting to %s", u.String())

	var conns []*websocket.Conn
	for i := 0; i < *connections; i++ {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			fmt.Println("Failed to connect", i, err)
			break
		}
		conns = append(conns, c)
		defer func() {
			c.WriteControl(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
				time.Now().Add(time.Second))
			time.Sleep(time.Second)
			c.Close()
		}()
	}

	log.Printf("Finished initializing %d connections", len(conns))
	tts := time.Second
	if *connections > 100 {
		tts = time.Millisecond * 5
	}
	for {
		for i := 0; i < len(conns); i++ {
			time.Sleep(tts)
			conn := conns[i]
			_, reader, err := conn.NextReader()
			if err == nil {
				bts, err := ioutil.ReadAll(reader)
				if err != nil {
					log.Printf("erreur lors de la lecture des donn??es")
				}
				log.Printf("Message from server : %s", string(bts))
			}
		}
	}
}
