package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// Article represents a news article found on the consulate's website.
type Article struct {
	Title string
	URL   string
}

// CheckResult holds the outcome of a single monitoring cycle.
type CheckResult struct {
	CheckedAt     time.Time
	ArticlesFound int
	NewArticles   int
	AlertsSent    int
	Error         string
}

// DB wraps the SQLite database connection.
type DB struct {
	conn *sql.DB
}

// Open opens (or creates) the SQLite database at path and runs migrations.
func Open(path string) (*DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("opening sqlite db at %q: %w", path, err)
	}

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	return db, nil
}

// Close closes the underlying database connection.
func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) migrate() error {
	_, err := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS seen_articles (
			url       TEXT PRIMARY KEY,
			title     TEXT NOT NULL,
			notified  INTEGER NOT NULL DEFAULT 0,
			seen_at   DATETIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS check_log (
			id              INTEGER PRIMARY KEY AUTOINCREMENT,
			checked_at      DATETIME NOT NULL,
			articles_found  INTEGER NOT NULL DEFAULT 0,
			new_articles    INTEGER NOT NULL DEFAULT 0,
			alerts_sent     INTEGER NOT NULL DEFAULT 0,
			error           TEXT
		);
	`)
	return err
}

// IsNew returns true if the article URL has not been seen before.
// The URL is the canonical unique identifier — titles are reused across publications.
func (db *DB) IsNew(url string) (bool, error) {
	var count int
	err := db.conn.QueryRow(
		`SELECT COUNT(1) FROM seen_articles WHERE url = ?`, url,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("checking if article is new: %w", err)
	}
	return count == 0, nil
}

// MarkAsSeen inserts the article into seen_articles.
// notified indicates whether a Telegram alert was sent.
func (db *DB) MarkAsSeen(article Article, notified bool) error {
	notifiedInt := 0
	if notified {
		notifiedInt = 1
	}
	_, err := db.conn.Exec(
		`INSERT OR IGNORE INTO seen_articles (url, title, notified, seen_at) VALUES (?, ?, ?, ?)`,
		article.URL, article.Title, notifiedInt, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("marking article as seen: %w", err)
	}
	return nil
}

// LogCheck inserts a record into check_log for observability.
func (db *DB) LogCheck(result CheckResult) error {
	errStr := sql.NullString{}
	if result.Error != "" {
		errStr = sql.NullString{String: result.Error, Valid: true}
	}
	_, err := db.conn.Exec(
		`INSERT INTO check_log (checked_at, articles_found, new_articles, alerts_sent, error)
		 VALUES (?, ?, ?, ?, ?)`,
		result.CheckedAt.UTC(),
		result.ArticlesFound,
		result.NewArticles,
		result.AlertsSent,
		errStr,
	)
	if err != nil {
		return fmt.Errorf("logging check result: %w", err)
	}
	return nil
}
