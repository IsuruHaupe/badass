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

When the server is launched it will wait for incoming connection requests, it will treat two kind of connection and two subtype.

Referee request : 

    * POST request : the referee sends an unique ID to identify him and the match he/she if judging (the referee ID and the match ID are the same).
    
    * GET request : the referee initiate a websocket by sending a get request. The referee then use this connection to sends live updates of the match.

Watcher request : 

    * POST request : the watcher sends the ID of the match (referee ID) he/she wants to get live update

    * GET resquest : the watcher initiate a websocket connection to receive updates. Whenever the referee of the match sends an update about the match the watcher is following, the server will forward those updates to him.

Since the handshake of the websocket is created through a GET request, we cannot pass the ID of the match the watcher wants to follow at the same time. We then need to get a list of the live match we can follow and their respective ID. Then, we make a POST request to suscribe to the live updates of the match using the dedicated route and ID of the match. After we are subscribe to the match, we need to instantiate the websocket connection to retrieve live updates by using the handshake (aka GET request) with the correct route.

We store the referees connection are stored in a map that is controlled by an epoll instance that will save ressources while waiting for referee to post messages (see reference 1 & 2 for more info about the optimization).

We have another map of map to link the referee to a pool of watchers (the keys are the referee ID and the values are a map of watchers connection). Whenever a referee sends an update, the epoll instance catch it and retrieves the pool of watcher for that referee using this map of map. We then iterate over the pool of connection and send the update to every watcher. Just like a pub/sub broker.

**Since the websocket handshake is a GET request and we can't do anything about it, we use a POST request before the websocket connection to notify the server of the match ID we want to receive updates on. Then, once the websocket connection is made we use this match ID previously sent to pair the match ID to the watcher ID inside the map of map**. To do that, we use a mutex to lock the match ID the current watcher wants. Every watcher sending a POST request will have to wait (synchronously) until the previous watcher has instanciated a websocket connection, where the mutex is unlock to make place for a new watcher to send a match he/she wants to follow. **To that end, every client needs to synchronously first sends a POST request with the match ID followed by a websocket connection.** 



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