package main

import (
	"database/sql"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type Page struct {
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Language    string    `json:"language"`
	LastUpdated time.Time `json:"last_updated"`
	Content     string    `json:"content"`
}

// ---- Function variables (can be replaced in tests) ----

var (
	InsertUserQuery        func(db *sql.DB, username, email, password string) (int64, error)
	GetUserIDQuery         func(db *sql.DB, username string) (int, error)
	GetUserByIDQuery       func(db *sql.DB, userID string) (int, string, string, string, error)
	GetUserByUsernameQuery func(db *sql.DB, username string) (int, string, string, string, error)
	SearchPagesQuery       func(db *sql.DB, searchTerm, language string) ([]Page, error)
)

// ---- Real implementations ----

func realInsertUserQuery(db *sql.DB, username, email, password string) (int64, error) {
	query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	res, err := db.Exec(query, username, email, password)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func realGetUserIDQuery(db *sql.DB, username string) (int, error) {
	query := "SELECT id FROM users WHERE username = ?"
	var id int
	err := db.QueryRow(query, username).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func realGetUserByIDQuery(db *sql.DB, userID string) (int, string, string, string, error) {
	query := "SELECT * FROM users WHERE id = ?"
	row := db.QueryRow(query, userID)
	var id int
	var username, email, password string
	err := row.Scan(&id, &username, &email, &password)
	if err != nil {
		return 0, "", "", "", err
	}
	return id, username, email, password, nil
}

func realGetUserByUsernameQuery(db *sql.DB, username string) (int, string, string, string, error) {
	query := "SELECT * FROM users WHERE username = ?"
	row := db.QueryRow(query, username)
	var id int
	var dbUsername, email, password string
	err := row.Scan(&id, &dbUsername, &email, &password)
	if err != nil {
		return 0, "", "", "", err
	}
	return id, dbUsername, email, password, nil
}

func realSearchPagesQuery(db *sql.DB, searchTerm, language string) ([]Page, error) {
	query := "SELECT title, url, language, last_updated, content FROM pages WHERE language = ?"
	args := []interface{}{language}

	if searchTerm != "" {
		query += " AND content LIKE ?"
		args = append(args, "%"+searchTerm+"%")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("rows.Close failed: %v", err)
		}
	}()

	var pages []Page
	for rows.Next() {
		var page Page
		if err := rows.Scan(&page.Title, &page.URL, &page.Language, &page.LastUpdated, &page.Content); err != nil {
			log.Printf("SearchPagesQuery row scan error: %v", err)
			continue
		}
		pages = append(pages, page)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return pages, nil
}

// ---- Assign real implementations ----

func init() {
	InsertUserQuery = realInsertUserQuery
	GetUserIDQuery = realGetUserIDQuery
	GetUserByIDQuery = realGetUserByIDQuery
	GetUserByUsernameQuery = realGetUserByUsernameQuery
	SearchPagesQuery = realSearchPagesQuery
}
