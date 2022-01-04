# Background 

Sport does not yet have the productivity tools of companies and projects!
The organisation of the multi-sport and multi-site tournament does not benefit from a simple and adapted application
for student tournaments, for example, or for sports clubs.

# Purpose of this project 

The purpose of this project is to design and develop a software to facilitate the management of a multi-sport tournament; in particular by communicating in real time to the organisers the results of the various matches entered on mobile devices (telephone, tablet, etc.).

# Requirements 

* Golang 
* An Unix OS - (we optimized the handling of pool of connection for the referee using epoll, see references 1 & 2 for more info)

# Architecture 

The project is build to mimic the process of a pub/sub protocol over websockets. The server serves as a broker and relay published messages by the referee (aka publisher) to watchers (aka subscriber) of a match.

When the server is launched it will wait for incoming connection requests, it will treat two kind of connection and two subtype. See [routes](#routes)

* Referee request : 

    * POST request : **the referee frontend generates and sends an unique ID as a string** to identify him/her and the match he/she is judging (the referee ID and the match ID are the same).
    
    * GET request : the referee initiate a websocket by sending a GET request. The referee then use this connection to sends live updates of the match.

* Watcher request : 

    * POST request : the watcher sends the ID of the match (referee ID) he/she wants to get live updates from. Match ID can be retrieved using a GET request to */live-match*.

    * GET resquest : the watcher initiate a websocket connection to receive updates. Whenever the referee of the match sends an update about the match the watcher is following, the server will forward those updates to him.

Since the handshake of the websocket is created through a GET request, we cannot pass the ID of the match the watcher wants to follow at the same time. We need to get a list of the live match we can follow and their respective ID (see [routes](#routes)). Then, we make a POST request to suscribe to the live updates of the match using the dedicated route and ID of the match. After we are subscribe to the match, we need to instantiate the websocket connection to retrieve live updates by using the handshake (aka GET request) with the correct route.

We store the referees connections in a map that is controlled by an epoll instance (not available in windows) that will save ressources while waiting for referee to post messages (see reference 1 & 2 for more info about the optimization).

We have another map of map to link the referee to a pool of watchers (the keys are the referee ID and the values are a map of watchers connections). Whenever a referee sends an update, the epoll instance catch it and retrieves the pool of watcher for that referee using this map of map. We then iterate over the pool of connection and send the update to every watcher. Just like a pub/sub broker.

**Since the websocket handshake is a GET request and we can't do anything about it, we use a POST request before the websocket connection to notify the server of the match ID we want to receive updates on. Then, once the websocket connection is made, we use the match ID previously sent via the GET request to pair the match ID to the watcher ID inside the map of map (the key is the match ID and the value is a map of watcher, the key of the second map is the unique ID of the watcher and the value is the connection to that watcher)**. To do that, we use a mutex to lock the match ID the current watcher wants. Every watcher sending a POST request will have to wait (synchronously) until the previous watcher has instanciated a websocket connection, where the mutex is unlock to make place for a new watcher to send a match he/she wants to follow. **To that end, every client needs to synchronously first sends a POST request with the match ID followed by a websocket connection.**


# Lost of connection 

* **Watcher side** : if the connection is lost from the watcher side, a new unique ID for the watcher is generated and all the previous events are sent to the watcher by the server. When the referee will send new updates, the server will remove the previous connection from the map of connection.

* **Referee side** : if the connection is lost from the referee side, we do not remove the referee ID from the map since every watchers of a match are linked to a match by the referee ID (since he/she is in charge of sending updates, he/she acts like a topic). **It is the duty of the front developer to generate an unique ID and cache it in the frontend in order to resend it via a POST request to register the referee again (aka a reconnection) when he/she tries to reconnect to the server.** The pool of watcher is kept intact and the referee can sends update again.



# Routes 

* */referee* : This route receives the handshake to instantiate a websocket connection between the server and the referee. 
* */referee/register* : This route receives a POST request with the ID of the referee in the form of a JSON format ({ID : id_of_referee_as_a_string}).
* */spectateur/register* : This route receives a POST request with the ID of the match the watcher wants to follow in a JSON format ({ID : id_of_match_as_a_string}).
* */spectateur* : This route is used to instantiate the websocket connection
* */live-match* : This route returns the live match that can be followed. Use the results of this GET request to send a POST request with the ID of the match you want to follow.

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