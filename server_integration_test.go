package poker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// server_integration_test.go
func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	// db, err := sql.Open("postgres", "postgres://localhost:5432/playerdb_test?sslmode=disable")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// db.Exec("CREATE TABLE IF NOT EXISTS players (name TEXT PRIMARY KEY, wins INTEGER NOT NULL DEFAULT 0)")
	// t.Cleanup(func() {
	// 	db.Exec("DROP TABLE players")
	// 	db.Close()
	// })

	// store := NewPostgresPlayerStore(db)

	// // clean slate before each run
	// db.Exec("DELETE FROM players")
	database, cleanDatabase := createTempFile(t, "")
	defer cleanDatabase()
	store, err := NewFileSystemPlayerStore(database)
	assertNoError(t, err)
	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		assertStatus(t, response.Code, http.StatusOK)

		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		assertStatus(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := []Player{
			{"Pepper", 3},
		}
		assertLeague(t, got, want)
	})
}
