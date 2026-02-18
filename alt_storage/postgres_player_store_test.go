package main

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func TestPostgresPlayerStore(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://localhost:5432/playerdb_test?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db.Exec("CREATE TABLE IF NOT EXISTS players (name TEXT PRIMARY KEY, wins INTEGER NOT NULL DEFAULT 0)")
	t.Cleanup(func() {
		db.Exec("DROP TABLE players")
		db.Close()
	})

	store := NewPostgresPlayerStore(db)

	// clean slate before each run
	db.Exec("DELETE FROM players")

	t.Run("get score for unknown player returns 0", func(t *testing.T) {
		got := store.GetPlayerScore("Unknown")
		if got != 0 {
			t.Errorf("got %d want 0", got)
		}
	})

	t.Run("record win for a new player then get score", func(t *testing.T) {
		store.RecordWin("Pepper")
		got := store.GetPlayerScore("Pepper")
		if got != 1 {
			t.Errorf("got %d want 1", got)
		}
	})

	t.Run("record multiple wins and get score", func(t *testing.T) {
		store.RecordWin("Floyd")
		store.RecordWin("Floyd")
		store.RecordWin("Floyd")
		got := store.GetPlayerScore("Floyd")
		if got != 3 {
			t.Errorf("got %d want 3", got)
		}
	})

	t.Run("scores are separate per player", func(t *testing.T) {
		got := store.GetPlayerScore("Pepper")
		if got != 1 {
			t.Errorf("Pepper: got %d want 1", got)
		}
		got = store.GetPlayerScore("Floyd")
		if got != 3 {
			t.Errorf("Floyd: got %d want 3", got)
		}
	})
}
