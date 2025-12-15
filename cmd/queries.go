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

type SearchResult struct {
	Title       string     `json:"title"`
	URL         string     `json:"url"`
	Language    string     `json:"language"`
	LastUpdated *time.Time `json:"last_updated,omitempty"`
	Snippet     string     `json:"snippet"`
	Rank        float64    `json:"-"`
}

// ---- Function variables (can be replaced in tests) ----

var (
	InsertUserQuery        func(db *sql.DB, username, email, password string) (int64, error)
	GetUserIDQuery         func(db *sql.DB, username string) (int, error)
	GetUserByIDQuery       func(db *sql.DB, userID string) (int, string, string, string, error)
	GetUserByUsernameQuery func(db *sql.DB, username string) (int, string, string, string, error)
	SearchPagesQuery       func(db *sql.DB, searchTerm, language string, limit int) ([]SearchResult, error)
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

func realSearchPagesQuery(db *sql.DB, searchTerm, language string, limit int) ([]SearchResult, error) {
	cappedLimit := clampLimit(limit)
	languageCode := "en"
	regConfig := "english"
	if language == "da" {
		languageCode = "da"
		regConfig = "danish"
	}

	query := `
WITH fts AS (
    SELECT
        p.title,
        p.url,
        p.language,
        p.last_updated,
        ts_headline(
            $2::regconfig,
            p.content,
            plainto_tsquery($2::regconfig, $1),
            'MaxFragments=2, MinWords=5, MaxWords=18, StartSel=<b>, StopSel=</b>'
        ) AS snippet,
        ts_rank(p.tsv_document, plainto_tsquery($2::regconfig, $1)) +
        COALESCE(EXTRACT(EPOCH FROM (p.last_updated - NOW())) * 1e-8, 0) AS rank
    FROM pages p
    WHERE p.language = $4
      AND p.tsv_document @@ plainto_tsquery($2::regconfig, $1)
    ORDER BY rank DESC, p.last_updated DESC NULLS LAST
    LIMIT $3
),
fallback AS (
    SELECT
        p.title,
        p.url,
        p.language,
        p.last_updated,
        ts_headline(
            $2::regconfig,
            p.content,
            plainto_tsquery($2::regconfig, $1),
            'MaxFragments=2, MinWords=5, MaxWords=18, StartSel=<b>, StopSel=</b>'
        ) AS snippet,
        similarity(p.title, $1) * 1.5 + similarity(p.content, $1) AS rank
    FROM pages p
    WHERE p.language = $4
      AND (
        p.title ILIKE '%' || $1 || '%'
        OR p.content ILIKE '%' || $1 || '%'
        OR p.title % $1
        OR p.content % $1
      )
      AND NOT EXISTS (SELECT 1 FROM fts f WHERE f.url = p.url)
    ORDER BY rank DESC, p.last_updated DESC NULLS LAST
    LIMIT $3
)
SELECT
    title,
    url,
    language,
    last_updated,
    snippet,
    rank
FROM (
    SELECT * FROM fts
    UNION ALL
    SELECT * FROM fallback WHERE (SELECT COUNT(*) FROM fts) < $3
) AS combined
ORDER BY rank DESC, last_updated DESC NULLS LAST
LIMIT $3;
`

	rows, err := db.Query(query, searchTerm, regConfig, cappedLimit, languageCode)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("rows.Close failed: %v", err)
		}
	}()

	var pages []SearchResult
	for rows.Next() {
		var page SearchResult
		var snippet sql.NullString
		var lastUpdated sql.NullTime

		if err := rows.Scan(&page.Title, &page.URL, &page.Language, &lastUpdated, &snippet, &page.Rank); err != nil {
			log.Printf("SearchPagesQuery row scan error: %v", err)
			continue
		}
		if snippet.Valid {
			page.Snippet = snippet.String
		}
		if lastUpdated.Valid {
			page.LastUpdated = &lastUpdated.Time
		}
		pages = append(pages, page)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return pages, nil
}

func clampLimit(limit int) int {
	if limit <= 0 {
		return 10
	}
	if limit > 50 {
		return 50
	}
	return limit
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
