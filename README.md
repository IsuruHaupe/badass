
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
- [Database](#database)
- [Test](#test)
- [Docker](#docker)
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

Let's say  a referee need to referee a game. First he/she will create a match using the specified route. See [routes](#routes). With the returned unique ID for a match the referee initiate a websocket connection. The referee then use this connection to send live updates of the match. See [events](#event) to check how events should look like. **The referee front-end send the unique ID in the URL query when creating the websocket**.  See [referee.go](client/referee.go) for more details. *Possible example : ws://127.0.0.1:8000/referee?**IdMatch=23PhWzEt2YdyRGM7iJHQ8uiCVwZ***.

* How does watcher requests work ?

Let's say you are a watcher and would like to follow live updates for a specific match. You will need to get the live match being played and get the match ID link to this match using the proper route. See [routes](#routes). When you have your match ID the watcher will initiate a websocket connection to receive updates. Whenever the referee of the match sends an update about the match the watcher is following, the server will forward those updates to him. **The watcher sends the ID of the match (referee ID) he/she wants to get live updates from in the URL query**. *Possible example : ws://127.0.0.1:8000/spectateur?**matchID=23PhWzEt2YdyRGM7iJHQ8uiCVwZ***. See [watcher.go](client/watcher.go) for more details.

We store the referees connections in a map that is controlled by an epoll instance (not available in windows, use docker provided in the repo, see [docker](#docker)) that will save computing ressources while waiting for referee to post messages (see [reference](#references) 1 & 2 for more info about the optimization).

We have another map of map to link the referee to a pool of watchers (the keys are the referee ID and the values are a map of watchers connections). Whenever a referee sends an update, the epoll instance catch it and retrieves the pool of watcher for that referee using this map of map. We then iterate over the pool of connection and send the update to every watcher. Just like a pub/sub broker.

# Improvement - Add a new sport

We work hard to find a way to allow new developpers to code new sports. That is to say the application is agnostic to any sport.
If you want to add a new sport you should treat each case for each event type. Every time a new event is send, we use a switch case in [sports.go](server/sports.go) to disciminate the sport and then we use a switch case to parse the event using specific parser functions. See [badminton.go](server/badminton.go). You should create a new *my_sport.go* and treat each event accordingly.

# Lost of connection 

* **Watcher side** : if the connection is lost from the watcher side, a new unique ID for the watcher is generated and all the previous events are sent to the watcher by the server when reconnecting. When the referee will send new updates, the server will remove the previous connection from the map of connection. The match ID must be supplied in the websocket connection as a query param.

* **Referee side** : if the connection is lost from the referee side, we do not remove the referee ID from the map since every watchers of a match are linked to a match by the referee ID (since he/she is in charge of sending updates, he/she acts like a topic). **It is the duty of the front developer to generate an unique ID and cache it in the frontend in order to resend it via a websocket connection to register the referee again (aka a reconnection) when he/she tries to reconnect to the server.** The pool of watcher is kept intact and the referee can sends update again.

# Referee Garbage Collector 

We use a go routine to periodically remove unused match ID in the list of match ID. Whenever a connection is lost, it is automatically removed from the epoller, this is true for every watcher. However, in the case of referee, the situation where the referee loses his/her connection may happen. In this case the referee's pool of watcher and referee ID in kept alive in the memory of the server. The referee client can then reconnect using a websocket connection and using the same UUID and restart sending events. However after the referee garbage collector is executed, the refereeID and the pool of watcher are removed. Referee and watchers must then reconnect. The cycle used to execute the garbage collection can be changed. **By default it's one hour.**

# Routes 

* */create-match* : TODO
* */create-tournament* : TODO 
* */referee* : This route receives the handshake to instantiate a websocket connection between the server and the referee. pass the refereeID a a string (key = refereeID, key = UUID)
* */spectateur* : This route is used to instantiate the websocket connection for a watcher to receive live event of a specified match. The match ID must be passed in the URL request (key = matchID, key = UUID of the match)
* */live-match* : This route returns the live match that can be followed. Use the result of this GET request to initiate a websocket connection with the server.
* */live-tournament* : TODO

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

* when creating a websocket connection you need to add params in the URL : 
* when creating a websocket connection with a referee : 
    * IdMatch : id of the match 
* when creating a websocket connection with a watcher : 
    * matchID : id of the match 
* when getting every match for a tournament : 
    * tournamentID : id of the tournament


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
* ```go run client.go``` will retrieve live match that can be follow via a GET request to the route */live-match*. Then will take the first entry in the array of ID of match and send a POST request with it to notify the server that the watcher wants to get updates on that particular match. 

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

# References 

* 1 - [Going Infinite, handling 1 millions websockets connections in Go / Eran Yanay](https://www.youtube.com/watch?v=LI1YTFMi8W4&t=1928s)
* 2 - [eranyanay/1m-go-websockets github repo](https://github.com/eranyanay/1m-go-websockets)
* 3 - [mqtt-essentials](https://www.hivemq.com/tags/mqtt-essentials/)
* 4 - [gobwas](https://github.com/gobwas/ws)
* 5 - [gorilla/websocket](https://github.com/gorilla/websocket)
* 6 - [Accessing a relational database](https://go.dev/doc/tutorial/database-access)

# TODO 

* cache the id of the referee in order to reconnect to the existing pool of watchers
* create the database for match (id_Match, event, id_Tournoi, id_equipe1, id_equipe2)
* create the database for tournament (id_Tournoi, nom)
* handle creation of match
    * create team name
    * create sport type
    * register the team in the match database
    * create a "get-summary" route to return summary of a match (takes a matchid) -> returns score, faults, timeout from history DB
    * create a struct to retrieve match event for the get-summary route {id_match, equipe, event_type, value}
*  handle creation of tournament and joining a tournament by sending the tournament ID when creating a match
* Add the referee as a spectator to check if events are being saved correctly