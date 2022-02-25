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
	//ip = flag.String("ip", "warm-dusk-64603.herokuapp.com", "server IP")
	ip          = flag.String("ip", "127.0.0.1", "server IP")
	connections = flag.Int("conn", 5, "number of websocket connections")
)

type Event struct {
	IdMatch    string `json:"IdMatch"`
	Equipe     string `json:"Equipe"`
	EventType  string `json:"EventType"`
	EventValue string `json:"EventValue"`
}

// Add sport type in GET
func getMatchId() string {
	//url := url.URL{Scheme: "http", Host: *ip, Path: "/create-match"}
	url := url.URL{Scheme: "http", Host: *ip + ":8000", Path: "/create-match"}
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
	//u := url.URL{Scheme: "http", Host: *ip, Path: "/create-tournament"}
	u := url.URL{Scheme: "http", Host: *ip + ":8000", Path: "/create-tournament"}
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

func initFootballTournament() string {
	// init tournament
	//u := url.URL{Scheme: "http", Host: *ip, Path: "/create-tournament"}
	u := url.URL{Scheme: "http", Host: *ip + ":8000", Path: "/create-tournament"}
	// add params to URL
	params := url.Values{}
	params.Add("tournamentName", "les bourres contre-attaques")
	params.Add("sport", "FOOTBALL")
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
	//u := url.URL{Scheme: "http", Host: *ip, Path: "/create-match"}
	u := url.URL{Scheme: "http", Host: *ip + ":8000", Path: "/create-match"}
	// add teams to URL
	params := url.Values{}
	params.Add("equipeA", "les bourres")
	params.Add("equipeB", "dikatomik")
	params.Add("tournamentID", tournamentID)
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

func initFootballMatch(tournamentID string) string {
	// init match
	//u := url.URL{Scheme: "http", Host: *ip, Path: "/create-match"}
	u := url.URL{Scheme: "http", Host: *ip + ":8000", Path: "/create-match"}
	// add teams to URL
	params := url.Values{}
	params.Add("equipeA", "les bourres")
	params.Add("equipeB", "dikatomik")
	params.Add("tournamentID", tournamentID)
	params.Add("sport", "FOOTBALL")
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
	//tournamentID := initBadmintonTournament()
	tournamentID := initFootballTournament()
	fmt.Println("ID du tournoi : ", tournamentID)

	//u := url.URL{Scheme: "ws", Host: *ip, Path: "/referee"}
	u := url.URL{Scheme: "ws", Host: *ip + ":8000", Path: "/referee"}
	// init multiple match and referee them
	log.Printf("Connecting to %s", u.String())
	var listOfMatch []string
	var conns []*websocket.Conn
	for i := 0; i < *connections; i++ {
		// init match using the tournament ID and save its ID
		//listOfMatch = append(listOfMatch, initBadmintonMatch(tournamentID))
		listOfMatch = append(listOfMatch, initFootballMatch(tournamentID))
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
	/*event := []Event{
		// START MATCH
		Event{IdMatch: "", Equipe: "", EventType: "BEGIN_MATCH", EventValue: ""},
		// POINT
		Event{IdMatch: "", Equipe: "EQUIPEA", EventType: "POINT", EventValue: "{\"Point\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEA", EventType: "POINT", EventValue: "{\"Point\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEB", EventType: "POINT", EventValue: "{\"Point\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEB", EventType: "POINT", EventValue: "{\"Point\":1}"},
		// FAULT
		Event{IdMatch: "", Equipe: "EQUIPEA", EventType: "FAULT", EventValue: "{\"FaultValue\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEA", EventType: "FAULT", EventValue: "{\"FaultValue\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEB", EventType: "FAULT", EventValue: "{\"FaultValue\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEB", EventType: "FAULT", EventValue: "{\"FaultValue\":1}"},
		// END SET
		Event{IdMatch: "", Equipe: "", EventType: "END_SET", EventValue: ""},
		// START NEW SET
		Event{IdMatch: "", Equipe: "", EventType: "NEW_SET", EventValue: ""},
		// POINT
		Event{IdMatch: "", Equipe: "EQUIPEA", EventType: "POINT", EventValue: "{\"Point\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEA", EventType: "POINT", EventValue: "{\"Point\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEB", EventType: "POINT", EventValue: "{\"Point\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEB", EventType: "POINT", EventValue: "{\"Point\":1}"},
		// END SET
		Event{IdMatch: "", Equipe: "", EventType: "END_SET", EventValue: ""},
		// END MATCH
		Event{IdMatch: "", Equipe: "", EventType: "END_MATCH", EventValue: ""},
	}*/

	event := []Event{
		// START MATCH
		Event{IdMatch: "", Equipe: "", EventType: "BEGIN_MATCH", EventValue: ""},
		// POINT
		Event{IdMatch: "", Equipe: "EQUIPEA", EventType: "POINT", EventValue: "{\"Point\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEA", EventType: "POINT", EventValue: "{\"Point\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEB", EventType: "POINT", EventValue: "{\"Point\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEB", EventType: "POINT", EventValue: "{\"Point\":1}"},
		// FAULT
		Event{IdMatch: "", Equipe: "EQUIPEA", EventType: "YELLOW_CARD", EventValue: "{\"FaultValue\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEA", EventType: "RED_CARD", EventValue: "{\"FaultValue\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEB", EventType: "YELLOW_CARD", EventValue: "{\"FaultValue\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEB", EventType: "RED_CARD", EventValue: "{\"FaultValue\":1}"},
		// PAUSE
		Event{IdMatch: "", Equipe: "", EventType: "HALF", EventValue: ""},
		// START OF SECOND HALF
		Event{IdMatch: "", Equipe: "", EventType: "SECOND_HALF", EventValue: ""},
		// START OF EXTENSION
		Event{IdMatch: "", Equipe: "", EventType: "EXTENSION", EventValue: ""},
		// POINT
		Event{IdMatch: "", Equipe: "EQUIPEA", EventType: "POINT", EventValue: "{\"Point\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEA", EventType: "POINT", EventValue: "{\"Point\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEB", EventType: "POINT", EventValue: "{\"Point\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEB", EventType: "POINT", EventValue: "{\"Point\":1}"},
		// END SET
		Event{IdMatch: "", Equipe: "", EventType: "PENALTY_SHOOTOUT", EventValue: ""},
		// POINT
		Event{IdMatch: "", Equipe: "EQUIPEA", EventType: "POINT_PENALTY_SHOOTOUT", EventValue: "{\"Point\":1}"},
		Event{IdMatch: "", Equipe: "EQUIPEA", EventType: "POINT_PENALTY_SHOOTOUT", EventValue: "{\"Point\":1}"},
		// END MATCH
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
					fmt.Println("error when marshelling in referee.go : %v", err)
				}
				//fmt.Println(string(body))
				conn.WriteMessage(websocket.TextMessage, body)

				// read data from the server i.e new state of the match after update
				_, reader, err := conn.NextReader()
				if err == nil {
					bts, err := ioutil.ReadAll(reader)
					if err != nil {
						log.Printf("erreur lors de la lecture des données")
					}
					log.Printf("Message from server : %s", string(bts))
				}
			}
		}
	}
}
