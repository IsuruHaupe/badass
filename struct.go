package main

import (
	"database/sql"
	"net"
	"sync"
)

type Epoll struct {
	fd          int
	connections map[int]net.Conn
	lock        *sync.RWMutex
}

type Event struct {
	Referee RefereeID
	Event   string
}

type Tournament struct {
	ID    string
	name  string
	sport string
}

type Referee struct {
	ID      string
	name    string
	surname string
}

type RefereeID struct {
	ID string
}

type Match struct {
	ID         string
	equipe1    string
	equipe2    string
	refere     Referee
	tournament Tournament
}

// This variable is used to handle referees
var refereeEpoller *Epoll

// transformer en map de map car si on veut supprimer une connection morte de l'array
var referees map[string]map[string]net.Conn

//var referees map[string]*Epoll

var db *sql.DB
