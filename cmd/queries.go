package main

import (
	"database/sql"
	"log"
	"time"
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
	GetUserCountQuery      func(db *sql.DB) (float64, error)
)

// ---- Real implementations ----

func realInsertUserQuery(db *sql.DB, username, email, password string) (int64, error) {
	// PostgreSQL does not support LastInsertId() reliably via database/sql,
	// so we use RETURNING.
	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id"

	var id int64
	if err := db.QueryRow(query, username, email, password).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func realGetUserIDQuery(db *sql.DB, username string) (int, error) {
	query := "SELECT id FROM users WHERE username = $1"
	var id int
	err := db.QueryRow(query, username).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func realGetUserByIDQuery(db *sql.DB, userID string) (int, string, string, string, error) {
	// Be explicit about columns (more robust than SELECT *)
	query := "SELECT id, username, email, password FROM users WHERE id = $1"
	row := db.QueryRow(query, userID)

	var id int
	var username, email, password string
	if err := row.Scan(&id, &username, &email, &password); err != nil {
		return 0, "", "", "", err
	}
	return id, username, email, password, nil
}

func realGetUserByUsernameQuery(db *sql.DB, username string) (int, string, string, string, error) {
	query := "SELECT id, username, email, password FROM users WHERE username = $1"
	row := db.QueryRow(query, username)

	var id int
	var dbUsername, email, password string
	if err := row.Scan(&id, &dbUsername, &email, &password); err != nil {
		return 0, "", "", "", err
	}
	return id, dbUsername, email, password, nil
}

func realSearchPagesQuery(db *sql.DB, searchTerm, language string) ([]Page, error) {
	// Use ILIKE for case-insensitive search in PostgreSQL (optional but nice).
	// If you want case-sensitive, swap back to LIKE.
	query := "SELECT title, url, language, last_updated, content FROM pages WHERE language = $1"
	args := []interface{}{language}

	if searchTerm != "" {
		query += " AND content ILIKE $2"
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

func realGetUserCountQuery(db *sql.DB) (float64, error) {
	var count float64
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ---- Assign real implementations ----

func init() {
	InsertUserQuery = realInsertUserQuery
	GetUserIDQuery = realGetUserIDQuery
	GetUserByIDQuery = realGetUserByIDQuery
	GetUserByUsernameQuery = realGetUserByUsernameQuery
	SearchPagesQuery = realSearchPagesQuery
	GetUserCountQuery = realGetUserCountQuery
}
