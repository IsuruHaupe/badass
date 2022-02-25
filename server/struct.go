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
	IdMatch    string
	Equipe     string
	EventType  string
	EventValue string
}

type Tournament struct {
	ID    string
	Name  string
	Sport string
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
	Id          string
	EquipeA     string
	EquipeB     string
	Tournament  string
	MatchValues string
}

// this map links the file descriptio (fd) in the epoller to the referee ID
var refereeFdToString map[int]string

// this is used to garbage collect unused referee connection every X second
var refereeToRemove map[string]string

// This variable is used to handle referees
var refereeEpoller *Epoll

// transformer en map de map car si on veut supprimer une connection morte de l'array
var referees map[string]map[string]net.Conn

//var referees map[string]*Epoll

var db *sql.DB
