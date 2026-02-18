package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewPostgresPlayerStore(db *sql.DB) *PostgresPlayerStore {
	return &PostgresPlayerStore{db}
}

type PostgresPlayerStore struct {
	store *sql.DB
}

func (i *PostgresPlayerStore) RecordWin(name string) {
	i.store.Exec(
		"INSERT INTO players (name, wins) VALUES ($1, 1) ON CONFLICT (name) DO UPDATE SET wins = players.wins + 1",
		name,
	)
}

func (i *PostgresPlayerStore) GetPlayerScore(name string) int {
	var wins int
	i.store.QueryRow("SELECT wins FROM players WHERE name = $1", name).Scan(&wins)
	return wins
}

func (i *PostgresPlayerStore) GetLeague() []Player {
	var league []Player
	rows, err := i.store.Query("SELECT name, wins FROM players")
	if err != nil {
		return league
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var wins int
		if err := rows.Scan(&name, &wins); err != nil {
			return league
		}
		league = append(league, Player{name, wins})
	}
	return league
}
