package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func NewSQLiteConnection(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		session_id TEXT,
		trace_id TEXT,
		name TEXT NOT NULL,
		timestamp INTEGER NOT NULL,
		received_at INTEGER NOT NULL,
		type TEXT NOT NULL,
		duration REAL,
		metadata TEXT
	);`

	if _, err := db.Exec(query); err != nil {
		return nil, err
	}

	return db, nil
}