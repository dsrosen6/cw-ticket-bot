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
	ID                 int
	BoardName          string
	RoomID             int
	NotifyNewTickets   bool
	NotifyStaleTickets bool
	NotifySlaBreach    bool
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

func (db *DB) GetBoardsWithNewTicketNoti() ([]Board, error) {
	return db.getBoardsWithNoti("notify_new_tickets", func(b *Board, v int) { b.NotifyNewTickets = v == 1 })
}

func (db *DB) GetBoardsWithStaleTicketNoti() ([]Board, error) {
	return db.getBoardsWithNoti("notify_stale_tickets", func(b *Board, v int) { b.NotifyStaleTickets = v == 1 })
}

func (db *DB) GetBoardsWithSlaBreachNoti() ([]Board, error) {
	return db.getBoardsWithNoti("notify_sla_breach", func(b *Board, v int) { b.NotifySlaBreach = v == 1 })
}

func (db *DB) getBoardsWithNoti(column string, setField func(*Board, int)) ([]Board, error) {
	query := fmt.Sprintf("SELECT id, board_name, room_id, %s FROM boards WHERE %s = 1", column, column)
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("querying boards with %s enabled: %w", column, err)
	}
	defer rows.Close()

	var boards []Board
	for rows.Next() {
		var b Board
		var flag int
		if err := rows.Scan(&b.ID, &b.BoardName, &b.RoomID, &flag); err != nil {
			return nil, fmt.Errorf("scanning board: %w", err)
		}
		setField(&b, flag)
		boards = append(boards, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating boards: %w", err)
	}

	return boards, nil
}
