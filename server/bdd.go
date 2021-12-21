package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

func AddEvent(db *sql.DB, ev Event) (int64, error) {
	result, err := db.Exec("INSERT INTO historique (evenement) VALUES (?)", ev.Event)
	if err != nil {
		return 0, fmt.Errorf("addEvent: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addEvent: %v", err)
	}
	return id, nil
}

func GetAllEvent(db *sql.DB) ([]Event, error) {
	// An albums slice to hold data from returned rows.
	var events []Event

	rows, err := db.Query("SELECT * FROM historique")
	if err != nil {
		return nil, fmt.Errorf("error : %v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var ev Event
		if err := rows.Scan(&ev.Referee.ID, &ev.Event); err != nil {
			return nil, fmt.Errorf("error : %v", err)
		}
		events = append(events, ev)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error : %v", err)
	}

	return events, nil
}

func ConnectToDB() (db *sql.DB) {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "history_of_message",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	return
}
