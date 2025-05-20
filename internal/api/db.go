package api

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var roomsTable = `
CREATE TABLE IF NOT EXISTS rooms (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	webex_room_id TEXT NOT NULL
);`

var boardsTable = `
CREATE TABLE IF NOT EXISTS boards (
	id INTEGER PRIMARY KEY NOT NULL UNIQUE,
	board_name TEXT NOT NULL,
	room_id INTEGER,
	notify_new_tickets INTEGER NOT NULL DEFAULT 0,
	notify_stale_tickets INTEGER NOT NULL DEFAULT 0,
	notify_sla_breach INTEGER NOT NULL DEFAULT 0,
	FOREIGN KEY(room_id) REFERENCES rooms(id)
);`

type DB struct {
	conn *sql.DB
}

type Room struct {
	ID          int
	WebexRoomID string
}

type Board struct {
	ID        int
	RoomID    int
	BoardName string
}

func newDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	return &DB{conn: db}, nil
}

func (db *DB) InitSchema() error {
	if _, err := db.conn.Exec(roomsTable); err != nil {
		return fmt.Errorf("creating rooms table: %w", err)
	}

	if _, err := db.conn.Exec(boardsTable); err != nil {
		return fmt.Errorf("creating boards table: %w", err)
	}

	return nil
}
