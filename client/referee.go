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

type Event struct {
	IdMatch    string `json:"IdMatch"`
	Equipe     string `json:"Equipe"`
	EventType  string `json:"EventType"`
	EventValue string `json:"EventValue"`
}

// Add sport type in GET
func getMatchId() string {
	url := url.URL{Scheme: "http", Host: *ip, Path: "/create-match"}
	//url := url.URL{Scheme: "http", Host: *ip + ":8000", Path: "/create-match"}
	resp, err := http.Get(url.String())

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

func initBadmintonTournament() string {
	// init tournament
	u := url.URL{Scheme: "http", Host: *ip, Path: "/create-tournament"}
	//u := url.URL{Scheme: "http", Host: *ip + ":8000", Path: "/create-tournament"}
	// add params to URL
	params := url.Values{}
	params.Add("tournamentName", "les bourres contre-attaques")
	params.Add("sport", "BADMINTON")
	//params.Add("IdMatch", IdMatch)
	u.RawQuery = params.Encode()

	resp, err := http.Get(u.String())

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

// This function creates a match using the specified route
// The request will return an unique match ID.
func initBadmintonMatch(tournamentID string) string {
	// init match
	u := url.URL{Scheme: "http", Host: *ip, Path: "/create-match"}
	//u := url.URL{Scheme: "http", Host: *ip + ":8000", Path: "/create-match"}
	// add teams to URL
	params := url.Values{}
	params.Add("equipeA", "les bourres")
	params.Add("equipeB", "dikatomik")
	params.Add("tournamentID", tournamentID)
	//params.Add("IdMatch", IdMatch)
	u.RawQuery = params.Encode()

	resp, err := http.Get(u.String())

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

func main() {
	//IdMatch := getMatchId()
	flag.Usage = func() {
		io.WriteString(os.Stderr, `Websockets client generator
	Example usage: ./client -ip=172.17.0.1 -conn=10
	`)
		flag.PrintDefaults()
	}
	flag.Parse()

	rand.Seed(time.Now().Unix())
	// init tournament
	tournamentID := initBadmintonTournament()
	fmt.Println("ID du tournoi : ", tournamentID)

	u := url.URL{Scheme: "ws", Host: *ip, Path: "/referee"}
	//u := url.URL{Scheme: "ws", Host: *ip + ":8000", Path: "/referee"}
	// init multiple match and referee them
	log.Printf("Connecting to %s", u.String())
	var listOfMatch []string
	var conns []*websocket.Conn
	for i := 0; i < *connections; i++ {
		// init match using the tournament ID and save its ID
		listOfMatch = append(listOfMatch, initBadmintonMatch(tournamentID))
		fmt.Println("ID du match : ", listOfMatch[i])
		// add match ID to URL
		params := url.Values{}
		params.Add("IdMatch", listOfMatch[i])
		u.RawQuery = params.Encode()
		// create websocket connection
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			fmt.Println("Failed to connect", i, err)
			break
		}
		conns = append(conns, c)
		defer func() {
			c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))
			time.Sleep(time.Second)
			c.Close()
		}()
	}

	log.Printf("Finished initializing %d connections", len(conns))
	tts := time.Second
	if *connections > 100 {
		tts = time.Millisecond * 5
	}

	//	IdMatch   string `json:"IdMatch"`
	//	Equipe    string `json:"Equipe"`
	//	EventType string `json:"EventType"`
	//	value     string `json:"Value"`
	event := []Event{
		Event{IdMatch: "", Equipe: "", EventType: "BEGIN_MATCH", EventValue: ""},
		// POINT
		Event{IdMatch: "", Equipe: "EQUIPEA", EventType: "POINT", EventValue: "{\"Point\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEA", EventType: "POINT", EventValue: "{\"Point\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEB", EventType: "POINT", EventValue: "{\"Point\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEB", EventType: "POINT", EventValue: "{\"Point\":1}"},
		// FAULT
		Event{IdMatch: "", Equipe: "EQUIPEA", EventType: "FAULT", EventValue: "{\"Player\":\"Isuru\", \"Comment\":\"Imbibe comme une brioche\", \"FaultValue\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEA", EventType: "FAULT", EventValue: "{\"Player\":\"Glenn\", \"Comment\":\"C'est le rhum qui prend Glenn\", \"FaultValue\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEB", EventType: "FAULT", EventValue: "{\"Player\":\"Isuru\", \"Comment\":\"Il sent plus rien\", \"FaultValue\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEB", EventType: "FAULT", EventValue: "{\"Player\":\"Glenn\", \"Comment\":\"Imbibe comme une brioche\", \"FaultValue\":1}"},
		Event{IdMatch: "", Equipe: "", EventType: "END_MATCH", EventValue: ""},
	}

	var matchID string
	for {
		// for each match ID send events
		for i := 0; i < len(conns); i++ {
			// use the corresponding matchID
			matchID = listOfMatch[i]
			time.Sleep(tts)
			conn := conns[i]
			// sending message
			for j := 0; j < len(event); j++ {
				time.Sleep(tts)
				event[j].IdMatch = matchID
				body, err := json.Marshal(event[j])
				if err != nil {
					fmt.Println("error when marshelling in referee.go L.112 : %v", err)
				}
				fmt.Println(string(body))
				conn.WriteMessage(websocket.TextMessage, body)
			}
		}
	}
}
