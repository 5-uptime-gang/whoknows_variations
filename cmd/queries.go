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

func InsertUserQuery(db *sql.DB, username string, email string, password string) (int64, error) {
	query := "INSERT INTO users (username, email, password) values (?, ?, ?)"
	res, err := db.Exec(query, username, email, password)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetUserIDQuery(db *sql.DB, username string) (int, error) {
	query := "SELECT id FROM users WHERE username = ?"
	var id int
	err := db.QueryRow(query, username).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetUserByIDQuery(db *sql.DB, userID string) (int, string, string, string, error) {
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

func GetUserByUsernameQuery(db *sql.DB, username string) (int, string, string, string, error) {
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

func SearchPagesQuery(db *sql.DB, searchTerm string, language string) ([]Page, error) {
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
    defer rows.Close()

    var pages []Page
    for rows.Next() {
        var page Page
        err := rows.Scan(&page.Title, &page.URL, &page.Language, &page.LastUpdated, &page.Content)
        if err != nil {
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

