package main

import (
	"database/sql"
	"log"
	"os"
)

var db *sql.DB

func openDatabase(path string) (*sql.DB, error) {
	dbExists := true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		dbExists = false
	}

	database, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if !dbExists {
		if err := InitDB(database); err != nil {
			if cerr := database.Close(); cerr != nil {
				log.Printf("Error closing DB after init failure: %v", cerr)
			}
			return nil, err
		}
	}

	db = database
	return database, nil
}

func closeDatabase() {
	if db == nil {
		return
	}
	if err := db.Close(); err != nil {
		log.Printf("Error closing DB: %v", err)
	}
}
