package main

import (
	"database/sql"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

func InitDB(db *sql.DB) error {
	// Users: drop + create + seed (samme som du havde)
	usersSchema := `
	DROP TABLE IF EXISTS users;

	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);

	INSERT OR IGNORE INTO users (username, email, password) 
	VALUES ('admin', 'keamonk1@stud.kea.dk', '5f4dcc3b5aa765d61d8327deb882cf99');`

	if _, err := db.Exec(usersSchema); err != nil {
		return err
	}

	// Pages: opret tabel hvis ikke findes
	pagesSchema := `
	CREATE TABLE IF NOT EXISTS pages (
		title TEXT PRIMARY KEY UNIQUE,
		url TEXT NOT NULL UNIQUE,
		language TEXT NOT NULL CHECK(language IN ('en', 'da')) DEFAULT 'en',
		last_updated TIMESTAMP,
		content TEXT NOT NULL
	);`

	if _, err := db.Exec(pagesSchema); err != nil {
		return err
	}

	// Seed pages data (brug en transaction + prepared statement)
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT OR IGNORE INTO pages (title, url, language, last_updated, content) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	seedData := []struct {
		Title   string
		URL     string
		Lang    string
		Content string
	}{
		{
			Title:   "Go Basics",
			URL:     "https://go.dev/doc/tutorial/getting-started",
			Lang:    "en",
			Content: "Go is a statically typed, compiled programming language designed at Google. Learn the basics of packages, functions, and goroutines.",
		},
		{
			Title:   "SQL Joins",
			URL:     "https://www.w3schools.com/sql/sql_join.asp",
			Lang:    "en",
			Content: "SQL joins are used to combine rows from two or more tables. Understand INNER JOIN, LEFT JOIN, RIGHT JOIN, and FULL OUTER JOIN.",
		},
		{
			Title:   "Introduktion til Go",
			URL:     "https://go.dev/doc/",
			Lang:    "da",
			Content: "Go (Golang) er et programmeringssprog udviklet af Google. Det er effektivt til backend-systemer og understøtter goroutines til samtidighed.",
		},
		{
			Title:   "SQL Forespørgsler",
			URL:     "https://www.sqlitetutorial.net/",
			Lang:    "da",
			Content: "SQL bruges til at hente og manipulere data i databaser. Eksempler inkluderer SELECT, INSERT, UPDATE og DELETE forespørgsler.",
		},
	}

	for _, p := range seedData {
		if _, err := stmt.Exec(p.Title, p.URL, p.Lang, time.Now(), p.Content); err != nil {
			// Vi logger fejl, men fortsætter med næste række
			log.Printf("Error inserting seed data (%s): %v", p.Title, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
