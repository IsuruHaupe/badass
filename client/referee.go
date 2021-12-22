package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/segmentio/ksuid"
)

var (
	ip          = flag.String("ip", "127.0.0.1", "server IP")
	connections = flag.Int("conn", 1, "number of websocket connections")
)

type Event struct {
	Referee RefereeID `json:"referee"`
	Event   string    `json:"event"`
}

type RefereeID struct {
	ID string `json:"ID"`
}

func main() {
	referee := RefereeID{
		ID: ksuid.New().String(),
	}
	body, err := json.Marshal(referee)
	if err != nil {
		fmt.Println("error when marshelling in referee.go L.40 : %v", err)
	}
	flag.Usage = func() {
		io.WriteString(os.Stderr, `Websockets client generator
Example usage: ./client -ip=172.17.0.1 -conn=10
`)
		flag.PrintDefaults()
	}
	flag.Parse()

	rand.Seed(time.Now().Unix())
	//route := [2]string{"/arbitre", "/spectateur"}

	// envoyer l'id du referee
	u := url.URL{Scheme: "http", Host: *ip + ":8000", Path: "/referee/register"}
	// send referee ID
	resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// creer la connexion websocket
	//u := url.URL{Scheme: "ws", Host: *ip + ":8000", Path: route[rand.Intn(len(route))]}
	u = url.URL{Scheme: "ws", Host: *ip + ":8000", Path: "/referee"}
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

	//message := os.Args[1]
	event := []Event{
		Event{Referee: referee, Event: "match created"},
		Event{Referee: referee, Event: "updates on score 1"},
		Event{Referee: referee, Event: "updates on timeout"},
		Event{Referee: referee, Event: "updates on score 2"},
		Event{Referee: referee, Event: "math ended"},
		Event{Referee: referee, Event: "event apres"},
		Event{Referee: referee, Event: "event apres"},
		Event{Referee: referee, Event: "event apres"},
	}

	// envoyer une structure qui contient l'id du referee et son message
	for {
		for i := 0; i < len(conns); i++ {
			time.Sleep(tts)
			conn := conns[i]
			//log.Printf("Spectateur %d sending message", i+1)
			//if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second*5)); err != nil {
			//	fmt.Printf("Failed to receive pong: %v", err)
			//}
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
	/*for i := 0; i < len(conns); i++ {
		time.Sleep(tts)
		conn := conns[i]
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%v", message)))
	}*/
}
