package main

import (
	"bytes"
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
	ip          = flag.String("ip", "127.0.0.1", "server IP")
	connections = flag.Int("conn", 1, "number of websocket connections")
)
var matchs []string

type Match struct {
	ID string `json:"ID"`
}

func main() {
	//go initReferee()
	u := url.URL{Scheme: "http", Host: *ip + ":8000", Path: "/live-match"}
	fmt.Println(u)
	getLiveMatch(u.String())
	fmt.Println(matchs)
	initWatcher(matchs[0])
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
		fmt.Println("error when marshelling in client.go L.4 8 : %v", err)
	}
}

func initWatcher(matchID string) {

	flag.Usage = func() {
		io.WriteString(os.Stderr, `Websockets client generator Example usage: ./client -ip=172.17.0.1 -conn=10`)
		flag.PrintDefaults()
	}
	flag.Parse()

	rand.Seed(time.Now().Unix())
	//route := [2]string{"/arbitre", "/spectateur"}
	//POST
	match := Match{
		ID: matchID,
	}

	body, err := json.Marshal(match)
	if err != nil {
		fmt.Println("error when marshelling in client.go L.68 : %v", err)
	}
	// send referee ID
	u := url.URL{Scheme: "http", Host: *ip + ":8000", Path: "/spectateur/register"}
	resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//WEBSOCKET
	u = url.URL{Scheme: "ws", Host: *ip + ":8000", Path: "/spectateur"}
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
			time.Sleep(time.Second)		c.Close()
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
			//log.Printf("Spectateur %d sending message", i+1)
			//if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second*5)); err != nil {
			//	fmt.Printf("Failed to receive pong: %v", err)
			//}
			// receiving message
			// decoder les messages avec un unmarshal
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

/*func initReferee() {
	flag.Usage = func() {
		io.WriteString(os.Stderr, `Websockets client generator
Example usage: ./client -ip=172.17.0.1 -conn=10
`)
		flag.PrintDefaults()
	}
	flag.Parse()

	rand.Seed(time.Now().Unix())
	//route := [2]string{"/arbitre", "/spectateur"}

	//u := url.URL{Scheme: "ws", Host: *ip + ":8000", Path: route[rand.Intn(len(route))]}
	u := url.URL{Scheme: "ws", Host: *ip + ":8000", Path: "/arbitre"}
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
	event := []Event{
		{event: "match created"},
		{event: "updates on score 1"},
		{event: "updates on timeout"},
		{event: "updates on score 2"},
		{event: "math ended"},
		{event: "event apres"},
		{event: "event apres"},
		{event: "event apres"},
	}
	for {
		for i := 0; i < len(conns); i++ {
			time.Sleep(tts)
			conn := conns[i]
			//log.Printf("Spectateur %d sending message", i+1)
			//if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second*5)); err != nil {
			//	fmt.Printf("Failed to receive pong: %v", err)
			//}
			// sending message
			for i := 0; i < 8; i++ {
				conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%v", event[i])))
			}
		}
	}
}*/
