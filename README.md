
- [Background](#background)
- [Purpose of this project](#purpose-of-this-project)
- [General purpose of this architecture](#general-purpose-of-this-architecture)
- [Requirements](#requirements)
- [Architecture](#architecture)
- [Improvement - Add a new sport](#improvement---add-a-new-sport)
- [Lost of connection](#lost-of-connection)
- [Referee Garbage Collector](#referee-garbage-collector)
- [Routes](#routes)
- [Event](#event)
- [Sports](#sports)
  * [Badminton](#badminton)
    + [Events badminton](#events-badminton)
- [Database](#database)
- [Test](#test)
- [Docker](#docker)
    + [Note](#note)
- [Deployment](#deployment)
- [Security Issue](#security-issue)
- [References](#references)
- [TODO](#todo)


# Background 

Sport does not yet have the productivity tools of companies and projects!
The organisation of the multi-sport and multi-site tournament does not benefit from a simple and adapted application
for student tournaments, for example, or for sports clubs.

# Purpose of this project 

The purpose of this project is to design and develop a software to facilitate the management of a multi-sport tournament; in particular by communicating in real time to the organisers the results of the various matches entered on mobile devices (telephone, tablet, etc.).

# General purpose of this architecture 

The general purpose of this repo/projet is to transfer any kind of messages from a publisher to a subscriber.

# Requirements 

* Golang >= 1.17
* An Unix OS - (we optimized the handling of pool of connection for the referee using epoll, see references 1 & 2 for more info). If you don't have epoll you can use docker (see [docker section](#docker)).

# Architecture 

The project is build to mimic the process of a pub/sub protocol but over websockets. The server represents a broker and relay published messages by the referee (aka publisher) to watchers (aka subscriber) of a match.

When the server is launched it will wait for incoming websocket connection requests and process incoming requests. See [routes](#routes).
The idea when you want to interact with the system is that you need to follow a couple of rules: 

* You need to create the instance of the match of the tournament in order to get an unique ID
* You use this unique ID as a param in the URL when initiating a websocket connection. 
* When you are a referee you send data and when you are a watcher you receive data.

* How does referee requests work ?

Let's say  a referee need to referee a game. First he/she will create a match using the specified route. See [routes](#routes). With the returned unique ID for a match the referee initiate a websocket connection. The referee then use this connection to send live updates of the match. Every time an event is sent, the new state of the match is sent back to the referee as truth. Use that to update front-end part.

See [events](#event) to check how events should look like. **The referee front-end send the unique ID in the URL query when creating the websocket**.  See [referee.go](client/referee.go) for more details. *Possible example : ws://127.0.0.1:8000/referee?**IdMatch=23PhWzEt2YdyRGM7iJHQ8uiCVwZ***.

* How does watcher requests work ?

Let's say you are a watcher and would like to follow live updates for a specific match. You will need to get the live match being played and get the match ID link to this match using the proper route. See [routes](#routes). When you have your match ID the watcher will initiate a websocket connection to receive updates. Whenever the referee of the match sends an update about the match the watcher is following, the server will forward those updates to him. **The watcher sends the ID of the match (referee ID) he/she wants to get live updates from in the URL query**. *Possible example : ws://127.0.0.1:8000/spectateur?**matchID=23PhWzEt2YdyRGM7iJHQ8uiCVwZ***. See [watcher.go](client/watcher.go) for more details.

We store the referees connections in a map that is controlled by an epoll instance (not available in windows, use docker provided in the repo, see [docker](#docker)) that will save computing ressources while waiting for referee to post messages (see [reference](#references) 1 & 2 for more info about the optimization).

We have another map of map to link the referee to a pool of watchers (the keys are the referee ID and the values are a map of watchers connections). Whenever a referee sends an update, the epoll instance catch it and retrieves the pool of watcher for that referee using this map of map. We then iterate over the pool of connection and send the update to every watcher. Just like a pub/sub broker.

# Improvement - Add a new sport

We worked hard to find a way to allow new developpers to code new sports. That is to say the application is agnostic to any sport.
If you want to add a new sport you should treat each case for each event type in [sports.go](server/sports.go). Every time a new event is received, we use a switch case in [sports.go](server/sports.go) to disciminate the sport and then we use a switch case to parse the event using specific parser functions. See [badminton.go](server/badminton.go). You should create a new *my_sport.go* and treat each event accordingly.

# Lost of connection 

* **Watcher side** : if the connection is lost from the watcher side, a new unique ID for the watcher is generated and all the previous events are sent to the watcher by the server when reconnecting. When the referee will send new updates, the server will remove the previous connection from the map of connection. The match ID must be supplied in the websocket connection as a query param.

* **Referee side** : if the connection is lost from the referee side, we do not remove the referee ID from the map since every watchers of a match are linked to a match by the referee ID (since he/she is in charge of sending updates, he/she acts like a topic). **It is the duty of the front developer to generate an unique ID and cache it in the frontend in order to resend it via a websocket connection to register the referee again (aka a reconnection) when he/she tries to reconnect to the server.** The pool of watcher is kept intact and the referee can sends update again.

# Referee Garbage Collector 

We use a go routine to periodically remove unused match ID in the list of match ID. Whenever a connection is lost, it is automatically removed from the epoller, this is true for every watcher. However, in the case of referee, the situation where the referee loses his/her connection may happen. In this case the referee's pool of watcher and referee ID in kept alive in the memory of the server. The referee client can then reconnect using a websocket connection and using the same UUID and restart sending events. However after the referee garbage collector is executed, the refereeID and the pool of watcher are removed. Referee and watchers must then reconnect. The cycle used to execute the garbage collection can be changed. **By default it's one hour.**

# Routes 

* [*/create-match*](server/init.go) : GET request to receive match ID generated by the server after sending in the query the following params : 
    * equipeA : the name as a string of the first team 
    * equipeB : the name as a string of the second team 
    * tournamentID : the UUID as a string of the tournament (if no tournament don't send anything)
    * Note : the server uses those details to create a match reference in database and send the UUID generated for that match in order for a referee to send updates on this match via events.
* [*/create-tournament*](server/init.go)  : GET request to receive tournament ID generated by the server after sending in the query the following params : 
    * tournamentName: name as a string of the tournament
    * sport : sport name, please be careful as the sport name is case-sensitive and should be one handled by the server, see [sport section](#sports)
    * Note : the server uses those details to create a tournament reference in database and send the UUID generated for that tournament in order for a referee to create match using this tournament ID in the params (param *tournamentID*).
* [*/referee*](server/wsControllers.go) : Websocket request to instantiate a websocket connection between the server and the referee. Pass the following params in the query : 
    * IdMatch : UUID as a string of the match the referee if sending events to.
* [*/spectateur*](server/wsControllers.go) : Websocket request to instantiate the websocket connection for a watcher to receive live event of a specified match. Pass the following params in the query : 
    * matchID : UUID as a string of the match the watcher wants to receive updates from.
* [*/get-live-match*](server/getLiveMatch.go) : Get request that returns the live match that can be followed. Use the result of this GET request to initiate a websocket connection with the server using a specific match ID.
* [*/get-live-match-for-tournament*](server/getLiveMatchForTournament.go) : Get request that returns all the match for a given tournament ID. Pass the following params in the query : 
    * tournamentID : UUID as a string of the tournament
* [*/get-live-tournament*](server/getLiveTournament.go) : Get request that returns the live tournaments. Use the result of this GET request to create a match using this tournament ID.

# Event 

The front-end for the referee will send event represented in JSON format respecting the following pattern : 

* for creating a tournament in the db

```json
{
    tournamentName: "name of the tournament",
    sport: "sport",
}
```


* for creating a match in the db

```json
{
    equipeA: "name of the team A",
    equipeB: "name of the team B",
    tournamentID: "ID of the tournament as an UUID",
}
```

* when creating a websocket connection with a referee : 
    * IdMatch : id of the match 

* when creating a websocket connection with a watcher : 
    * matchID : id of the match 

* when getting every match for a tournament : 
    * tournamentID : id of the tournament

# Sports 
Here is the list of sports that the database can handle : 

## Badminton
 use BADMINTON when creating a tournament with the specified route and sending the sport (ex : sport=BADMINTON). Be careful it is case-sensitive.

### Events badminton 

* POINT : Add a point in one equipe 
    ```
    matchValues : { 
        Point : //(int) number of point (can be negatif if the referee cancel a point)
    }
    ```
* FAULT : Add a fault
    ```
    matchValues : { 
        Player : //(String) player name
        Comment : //(String) comment
        FaultValue : //(String) type of value (carton rouge/ carton jaune)
    }
    ```
* BEGIN_MATCH : change the status of the match
* END_MATCH : change the status of the match

# Database 

We use a mySQL database where every events from a match are stored. When a client connects to a match the aim is to get all the events he/she missed and sends them to him/her.

# Test 

In order to test the system, just launch three terminal and use the following command : 

```
// in the server folder
go run *.go 
// in the client folder
go run referee.go

go run client.go
```

* ```go run *.go``` will create the server and listen for incoming requests. You should see a connected message indicating you succesfully connected to database.
* ```go run referee.go``` will create an unique ID for the referee program and send it via a POST request to the server, the server will register this ID as a referee and will wait for websocket connection. Then the referee program create a websocket connection and start sending fake updates.
* ```go run client.go``` will retrieve live match that can be follow via a GET request to the route */get-live-match*. Then will take the first entry in the array of ID of match and send a POST request with it to notify the server that the watcher wants to get updates on that particular match. 

# Docker 

To run the server in a docker environment use the following command inside the global directory : 

```bash
# run docker compose and free command line
docker-compose up -d 

# query the message from the server instance
docker logs -f server

# you can now run referee.go and watcher.go
# inside client 
go run referee.go 
go run watcher.go
```

### Note 

You will need to uncomment a line in [bdd.go](server/bdd.go) in the ```ConnectToDB``` function. And maybe update the database inside the docker image. And update [referee.go](client/referee.go) and [watcher.go](client/watcher.go) in the client folder.


# Add a new sport (example : Football)

1 - Create a new file with the name of sport (ex : football.go) in the server folder.

2 - Create the structure for a match of this sport as a JSON representation. For example : 
```JSON
{
    "EquipeA" : {
        "Score": int,
        "Fault": {
            "NumberOfRedCard": int,
            "NumberOfYellowCard": int
        }
    },
    "EquipeB" : {
        "Score": int,
        "Fault": {
            "NumberOfRedCard": int,
            "NumberOfYellowCard": int
        }
    },
    "status": string, //("NOT_BEGUN" "FIRST_HALF" "HALF" "SECOND_HALF" "EXTENSION" "PENALTY_SHOOTOUT" "END_MATCH")
    "PenaltyShootout": {
        "ScoreEquipeA":int,
        "ScoreEquipeB":int
    }
}
```

3 - Implement this structure in go and add in the file football.go

```go
type EquipeFootball struct {
	Score int           `json:"Score"`
	Fault FaultFootball `json:"Fault"`
}

type FaultFootball struct {
	NumberOfRedCard    int `json:"NumberOfRedCard"`
	NumberOfYellowCard int `json:"NumberOfYellowCard"`
}
type Football struct {
	EquipeA         EquipeFootball  `json:"EquipeA"`
	EquipeB         EquipeFootball  `json:"EquipeB"`
	Status          string          `json:"Status"` //("NOT_BEGUN" "FIRST_HALF" "HALF" "SECOND_HALF" "EXTENSION" "PENALTY_SHOOTOUT)
	PenaltyShootout PenaltyShootout `json:"PenaltyShootout"`
}

type PenaltyShootout struct {
	ScoreEquipeA int    `json:"ScoreEquipeA"`
	ScoreEquipeB int    `json:"ScoreEquipeB"`
}
```
4 - Create in the file football.go a function for intializing a struct Football and return it in JSON :
``` go
func InitializeFootball() string {
	football := Football{
		EquipeA: EquipeFootball{
			Score: 0,
			Fault: FaultFootball{
				NumberOfRedCard:    0,
				NumberOfYellowCard: 0,
			},
		},
		EquipeB: EquipeFootball{
			Score: 0,
			Fault: FaultFootball{
				NumberOfRedCard:    0,
				NumberOfYellowCard: 0,
			},
		},
		Status: "NOT_BEGUN",
		PenaltyShootout: PenaltyShootout{
			ScoreEquipeA: 0,
			ScoreEquipeB: 0,
		},
	}
	tmp, err := json.Marshal(football)
	if err != nil {
		fmt.Println("error initialize football struct: %v", err)
	}
	return string(tmp)
}
```

5 - Create in the file football.go a function for managing the different events with a switch case :

```go
//Event struct :
type Event_Football_Point struct {
	Point int `json:"Point"`
}
type Event_Football_Fault struct {
	FaultValue int    `json:"FaultValue"`
}
// function to treat each events for this sport
func ParseEventFootball(event Event, match Match) []byte {
	// badminton := Badminton{}
	var football Football
	json.Unmarshal([]byte(match.matchValues), &football)
	switch event.EventType {
	case "POINT":
		point := Event_Football_Point{}
		json.Unmarshal([]byte(event.EventValue), &point)
		if event.Equipe == "EQUIPEA" {
			football.EquipeA.Score += point.Point
		} else {
			football.EquipeB.Score += point.Point
		}
	case "REDCARD":
		fault := Event_Football_Fault{}
		json.Unmarshal([]byte(event.EventValue), &fault)
		if event.Equipe == "EQUIPEA" {
			// in case of cancel we use fault.FaultValue
			football.EquipeA.Fault.NumberOfRedCard += fault.FaultValue
		} else {
			football.EquipeB.Fault.NumberOfRedCard += fault.FaultValue
		}
	case "YELLOWCARD":
		fault := Event_Football_Fault{}
		json.Unmarshal([]byte(event.EventValue), &fault)
		if event.Equipe == "EQUIPEA" {
			// in case of cancel we use fault.FaultValue
			football.EquipeA.Fault.NumberOfYellowCard += fault.FaultValue
		} else {
			football.EquipeB.Fault.NumberOfYellowCard += fault.FaultValue
		}
	case "POINT_PENALTY_SHOOTOUT":
		point := Event_Football_Point{}
		json.Unmarshal([]byte(event.EventValue), &point)
		if event.Equipe == "EQUIPEA" {
			// in case of cancel we use fault.FaultValue
			football.PenaltyShootout.ScoreEquipeA += point.Point
		} else {
			football.EquipeB.Fault.NumberOfYellowCard += point.Point
		}
	case "FIRST_HALF":
		football.Status = "FIRST_HALF"
	case "HALF":
		football.Status = "HALF"
	case "SECOND_HALF":
		football.Status = "SECOND_HALF"
	case "EXTENSION":
		football.Status = "EXTENSION"
	case "PENALTY_SHOOTOUT":
		football.Status = "PENALTY_SHOOTOUT"
	case "END_MATCH":
		football.Status = "END_MATCH"

	}
	tmp, err := json.Marshal(football)
	if err != nil {
		fmt.Println("error when marshelling in football.go: %v", err)
	}
	match.matchValues = string(tmp)
	err = UpdateMatch(db, match)
	if err != nil {
		fmt.Println("Error update match : %v", err)
	}
	return tmp
}
```
6 - In the file sport.go in the function InitializeSport() add on the switch, a case for calling your function for initializing your sport

```go
func InitializeSport(sport string) string {
	switch sport {
	case "BADMINTON":
		return InitializeBadminton()
	case "FOOTBALL":
		return InitializeFootball()
	}
	return ""
}
```
7 - In the file sport.go in the function ParseEvent() add on the switch, a case for calling your function for managing the event of your sport

```go
func ParseEvent(event Event, sport string) []byte {
	// retrieve the match to which this event is tied
	match, err := getMatch(db, event.IdMatch)
	if err != nil {
		fmt.Errorf("Parse event error get match : %v", err)
	}
	// treat each sport accordingly
	switch sport {
	case "BADMINTON":
		return ParseEventBadminton(event, match)
	case "FOOTBALL":
		return ParseEventFootball(event, match)
	}
	return nil
}

```

# Deployment 

We deployed a functionnal version of the back-end in heroku. You can do the same by following the readme in the *production* branch in this repository. 

Basically you can deploy this back-end wherever you want (for example on a rasberry pi or a cloud provider) and allow connection on the device so that referee and watchers interact with each other. The referee and the watchers then only need to know the URL of the website to connect to it from the frontend and interact with the server. For this you will need an internet connection (wi-fi or mobile data).

# Security Issue

A ill intentioned user can easily hack a match that is being played by simply querying the live match route and getting the UUID of the match. Then using this UUID can connect via a websocket connection to the referee route. To that extent new developpers contributing on that project should implement secure authentification of user for example using JWT. We didn't have time to implement it but we still wanted to point out that this is feature that will be needed for real use case of this architecture. Please look a the [todo](#todo) section with useful links to implement it.

# References 

* 1 - [Going Infinite, handling 1 millions websockets connections in Go / Eran Yanay](https://www.youtube.com/watch?v=LI1YTFMi8W4&t=1928s)
* 2 - [eranyanay/1m-go-websockets github repo](https://github.com/eranyanay/1m-go-websockets)
* 3 - [mqtt-essentials](https://www.hivemq.com/tags/mqtt-essentials/)
* 4 - [gobwas](https://github.com/gobwas/ws)
* 5 - [gorilla/websocket](https://github.com/gorilla/websocket)
* 6 - [Accessing a relational database](https://go.dev/doc/tutorial/database-access)

# TODO 

* cache the id of the referee in order to reconnect to the existing pool of watchers (front side)
* Implement secure authentification 
    * https://www.bacancytechnology.com/blog/golang-jwt
    * https://learn.vonage.com/blog/2020/03/13/using-jwt-for-authentication-in-a-golang-application-dr/
    * https://www.sohamkamani.com/golang/session-based-authentication/
    * https://www.sohamkamani.com/golang/password-authentication-and-storage/
* implement TDD 
    * https://github.com/quii/learn-go-with-tests