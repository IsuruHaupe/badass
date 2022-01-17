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

func main() {
	IdMatch := "23VQUHJiBkSBOuUAZrrvfXU1mvj"
	//IdMatch := getMatchId()
	flag.Usage = func() {
		io.WriteString(os.Stderr, `Websockets client generator
	Example usage: ./client -ip=172.17.0.1 -conn=10
	`)
		flag.PrintDefaults()
	}
	flag.Parse()

	rand.Seed(time.Now().Unix())
	// creer la connexion websocket
	u := url.URL{Scheme: "ws", Host: *ip + ":8000", Path: "/referee"}
	//u := url.URL{Scheme: "ws", Host: *ip, Path: "/referee"}
	// add referee ID to URL
	params := url.Values{}
	params.Add("IdMatch", "23VQUHJiBkSBOuUAZrrvfXU1mvj")
	//params.Add("IdMatch", IdMatch)
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
		// POINT
		/*Event{IdMatch: IdMatch, Equipe: "EQUIPEA", EventType: "POINT", EventValue: "{\"Point\":1}"},
		Event{IdMatch: IdMatch, Equipe: "EQUIPEA", EventType: "POINT", EventValue: "{\"Point\":-1}"},
		Event{IdMatch: IdMatch, Equipe: "EQUIPEB", EventType: "POINT", EventValue: "{\"Point\":1}"},
		Event{IdMatch: IdMatch, Equipe: "EQUIPEB", EventType: "POINT", EventValue: "{\"Point\":-1}"},*/
		// FAULT
		Event{IdMatch: IdMatch, Equipe: "EQUIPEA", EventType: "FAULT", EventValue: "{\"Player\":\"Isuru\", \"Comment\":\"Imbibe comme une brioche\", \"FaultValue\":1}"},
		Event{IdMatch: IdMatch, Equipe: "EQUIPEA", EventType: "FAULT", EventValue: "{\"Player\":\"Glenn\", \"Comment\":\"C'est le rhum qui prend Glenn\", \"FaultValue\":-1}"},
		Event{IdMatch: IdMatch, Equipe: "EQUIPEB", EventType: "FAULT", EventValue: "{\"Player\":\"Isuru\", \"Comment\":\"Il sent plus rien\", \"FaultValue\":1}"},
		Event{IdMatch: IdMatch, Equipe: "EQUIPEB", EventType: "FAULT", EventValue: "{\"Player\":\"Glenn\", \"Comment\":\"Imbibe comme une brioche\", \"FaultValue\":-1}"},
	}

	for {
		for i := 0; i < len(conns); i++ {
			time.Sleep(tts)
			conn := conns[i]
			// sending message
			for i := 0; i < len(event); i++ {
				time.Sleep(tts)
				body, err := json.Marshal(event[i])
				if err != nil {
					fmt.Println("error when marshelling in referee.go L.112 : %v", err)
				}
				fmt.Println(string(body))
				conn.WriteMessage(websocket.TextMessage, body)
			}
		}
	}
}
