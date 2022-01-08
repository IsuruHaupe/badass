package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func AddEvent(db *sql.DB, ev Event) (int64, error) {
	result, err := db.Exec("INSERT INTO history (evenement) VALUES (?)", ev.Event)
	if err != nil {
		return 0, fmt.Errorf("addEvent: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addEvent: %v", err)
	}
	return id, nil
}

func AddTournament(db *sql.DB, tr Tournament) (int64, error) {
	result, err := db.Exec("INSERT INTO tournament (tournament) VALUES (?)", tr.name)
	if err != nil {
		return 0, fmt.Errorf("AddTournament: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddTournament: %v", err)
	}
	return id, nil
}

func AddArbitre(db *sql.DB, tr Tournament) (int64, error) {
	result, err := db.Exec("INSERT INTO tournament (tournament) VALUES (?)", tr.name)
	if err != nil {
		return 0, fmt.Errorf("AddTournament: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddTournament: %v", err)
	}
	return id, nil
}

func GetAllEvent(db *sql.DB) ([]Event, error) {
	// An albums slice to hold data from returned rows.
	var events []Event

	rows, err := db.Query("SELECT * FROM history")
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

// https://stackoverflow.com/questions/39281594/error-1698-28000-access-denied-for-user-rootlocalhost
// for problems with mysql
func ConnectToDB() (db *sql.DB) {
	var opts string
	var err error
	// Get a database handle.
	if os.Getenv("ENV") == "PROD" {
		//db, err = sql.Open("mysql", "root:mypassword@tcp(db:3306)/testdb")
		// TODO : find a way to parse URL from heroku
		opts = "zrf4tp5q8lnqwbg5:dlir6epzfdl15g2c@lmc8ixkebgaq22lo.chr7pe7iynqr.eu-west-1.rds.amazonaws.com:3306/ub59fgelo956gbfv"
	} else {
		/*cfg := mysql.Config{
			User:   os.Getenv("DBUSER"),
			Passwd: os.Getenv("DBPASS"),
			Net:    "tcp",
			Addr:   "127.0.0.1:3306",
			DBName: "history_of_message",
		}
		opts = cfg.FormatDSN()
		*/

		opts = "root:mypassword@tcp(db:3306)/testdb"
	}

	db, err = sql.Open("mysql", opts)
	if err != nil {
		log.Fatal(err)
	}
	// MySQL server isn't fully active yet.
	// Block until connection is accepted. This is a docker problem with v3 & container doesn't start
	// up in time.
	for db.Ping() != nil {
		fmt.Println("Attempting connection to db")
		time.Sleep(5 * time.Second)
	}
	fmt.Println("Connected !")
	return
}
