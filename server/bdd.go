package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func AddEvent(db *sql.DB, ev Event) (int64, error) {
	result, err := db.Exec("INSERT INTO history (eventMatch, eventType, equipe, idMatch) VALUES (?)", ev.eventMatch, ev.eventType, ev.equipe, ev.idMatch)
	if err != nil {
		return 0, fmt.Errorf("addEvent: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addEvent: %v", err)
	}
	return id, nil
}

//Add a new tournment in db
func CreateTournament(db *sql.DB, tr Tournament) (int64, error) {
	result, err := db.Exec("INSERT INTO tournament (nameTournament,sport) VALUES (?,?)", tr.name, tr.sport)
	if err != nil {
		return 0, fmt.Errorf("Create Tournament: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Create Tournament: %v", err)
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

//Add a new match in db
func CreateMatch(db *sql.DB, m Match) (int64, error) {
	result, err := db.Exec("INSERT INTO matchs (equipeA,equipeB) VALUES (?,?)", m.equipeA, m.equipeB)
	if err != nil {
		return 0, fmt.Errorf("Create matchs: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Create matchs: %v", err)
	}
	return id, nil
}

func getMatch(db *sql.DB, idMatch string) (Match, error) {
	rows, err := db.Query("SELECT * from  matchs  where id = (?) ", idMatch)
	if err != nil {
		return Match{}, fmt.Errorf("error : %v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var match Match
		if err := rows.Scan(&match.id, &match.equipeA, &match.equipeB, &match.matchValues); err != nil {
			return Match{}, fmt.Errorf("error : %v", err)
		}
		return match, nil
	}
	return Match{}, fmt.Errorf("error : %v", err)

}

// https://stackoverflow.com/questions/39281594/error-1698-28000-access-denied-for-user-rootlocalhost
// for problems with mysql
func ConnectToDB() (db *sql.DB) {
	var err error
	// Get a database handle.
	if os.Getenv("ENV") == "PROD" {
		// TODO : find a way to parse URL from heroku
		// schema : DATABASE_URL='user:pass@tcp(hostname:3306)/your_heroku_database'
		DATABASE_URL := "b8afd730e14ddf:660528a6@tcp(us-cdbr-east-05.cleardb.net:3306)/heroku_142de0a726b37cc"
		db, err = sql.Open("mysql", DATABASE_URL)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		cfg := mysql.Config{
			User:   os.Getenv("DBUSER"),
			Passwd: os.Getenv("DBPASS"),
			Net:    "tcp",
			Addr:   "127.0.0.1:3306",
			DBName: "history_of_message",
		}
		db, err = sql.Open("mysql", cfg.FormatDSN())
		if err != nil {
			log.Fatal(err)
		}
		// for docker env
		//db, err = sql.Open("mysql", "root:mypassword@tcp(db:3306)/testdb")
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
