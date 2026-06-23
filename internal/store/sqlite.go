package store

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

type Store struct {
	DB *sql.DB
}

// -------------------- INIT --------------------

func New(dbPath string) *Store {

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	s := &Store{DB: db}
	s.init()

	return s
}

// -------------------- SCHEMA --------------------

func (s *Store) init() {

	// agents table
	_, _ = s.DB.Exec(`
	CREATE TABLE IF NOT EXISTS agents (
		agent_id TEXT PRIMARY KEY,
		hostname TEXT,
		last_seen TEXT,
		status TEXT
	)`)

	// events table
	_, _ = s.DB.Exec(`
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		agent_id TEXT,
		hostname TEXT,
		event_type TEXT,
		payload TEXT,
		timestamp TEXT
	)`)

	// alerts table
	_, _ = s.DB.Exec(`
	CREATE TABLE IF NOT EXISTS alerts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		agent_id TEXT,
		title TEXT,
		severity TEXT,
		description TEXT,
		timestamp TEXT
	)`)
}
