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
	result, err := db.Exec("INSERT INTO events (eventMatch, eventType, equipe, idMatch) VALUES (?, ?, ?, ?)", ev.EventValue, ev.EventType, ev.Equipe, ev.IdMatch)
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
func CreateTournament(db *sql.DB, tr Tournament) error {
	_, err := db.Exec("INSERT INTO tournament (id,nameTournament,sport) VALUES (?,?,?)", tr.ID, tr.name, tr.sport)
	if err != nil {
		return fmt.Errorf("Create Tournament: %v", err)
	}
	return nil
}

func GetTournament(db *sql.DB, idTournament string) (Tournament, error) {
	rows, err := db.Query("SELECT * from  tournament  where id = ? ", idTournament)
	if err != nil {
		return Tournament{}, fmt.Errorf("error : %v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var tournament Tournament
		if err := rows.Scan(&tournament.ID, &tournament.name, &tournament.sport); err != nil {
			return Tournament{}, fmt.Errorf("error : %v", err)
		}
		return tournament, nil
	}
	return Tournament{}, fmt.Errorf("error : %v", err)
}

/*func GetAllEvent(db *sql.DB) ([]Event, error) {
	// An albums slice to hold data from returned rows.
	var events []Event

	rows, err := db.Query("SELECT * FROM events")
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
}*/

//Add a new match in db
func CreateMatch(db *sql.DB, m Match) error {
	_, err := db.Exec("INSERT INTO matchs (id, equipeA, equipeB, matchValues, idTournament) VALUES (?,?,?,?,?)", m.Id, m.EquipeA, m.EquipeB, m.MatchValues, m.Tournament)
	if err != nil {
		return fmt.Errorf("Create matchs error: %v", err)
	}

	return nil
}

//Add a new match in db
func UpdateMatch(db *sql.DB, m Match) error {
	_, err := db.Exec("UPDATE matchs SET matchValues = (?) where id= (?)", m.MatchValues, m.Id)
	if err != nil {
		return fmt.Errorf("update matchs: %v", err)
	}
	return nil
}

func getMatch(db *sql.DB, idMatch string) (Match, error) {
	rows, err := db.Query("SELECT * from  matchs  where id = ? ", idMatch)
	if err != nil {
		return Match{}, fmt.Errorf("error : %v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var match Match
		if err := rows.Scan(&match.Id, &match.EquipeA, &match.EquipeB, &match.Tournament, &match.MatchValues); err != nil {
			return Match{}, fmt.Errorf("error : %v", err)
		}
		return match, nil
	}
	return Match{}, fmt.Errorf("error : %v", err)

}

func getAllTournament(db *sql.DB) ([]Tournament, error) {
	var tournaments []Tournament = make([]Tournament, 0)
	rows, err := db.Query("SELECT * from  tournament")
	if err != nil {
		return nil, fmt.Errorf("error : %v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var tournament Tournament
		if err := rows.Scan(&tournament.ID, &tournament.name, &tournament.sport); err != nil {
			return nil, fmt.Errorf("error : %v", err)
		}
		tournaments = append(tournaments, tournament)
	}
	return tournaments, nil

}

func getMatchForTournament(db *sql.DB, tournamentID string) ([]Match, error) {
	var matchs []Match = make([]Match, 0)
	rows, err := db.Query("SELECT * from  matchs  where idTournament = ? ", tournamentID)
	if err != nil {
		return nil, fmt.Errorf("error : %v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var match Match
		if err := rows.Scan(&match.Id, &match.EquipeA, &match.EquipeB, &match.Tournament, &match.MatchValues); err != nil {
			return nil, fmt.Errorf("error : %v", err)
		}
		matchs = append(matchs, match)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error : %v", err)
	}

	return matchs, fmt.Errorf("error : %v", err)

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
		/*db, err = sql.Open("mysql", "root:mypassword@tcp(db:3306)/history_of_message")
		if err != nil {
			log.Panic(err)
		}*/
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
