package main

import (
	"database/sql"
	"imdbsqlite/internal"
	"testing"

	_ "modernc.org/sqlite"
)

// build a tiny test DB with known contents.
func buildTestQueryDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	// load schema from create_imdb.go
	if _, err := db.Exec(internal.Schema); err != nil {
		t.Fatalf("schema error: %v", err)
	}

	// insert a little bit of data
	db.Exec(`INSERT INTO movies VALUES (1, 'Test Movie', 2000, 8.0)`)
	db.Exec(`INSERT INTO movies VALUES (2, 'Bad Movie', 2001, 3.0)`)

	db.Exec(`INSERT INTO actors VALUES (10, 'Jane', 'Doe', 'F')`)
	db.Exec(`INSERT INTO actors VALUES (20, 'John', 'Smith', 'M')`)

	db.Exec(`INSERT INTO genres VALUES (1, 'Drama')`)
	db.Exec(`INSERT INTO genres VALUES (2, 'Comedy')`)

	db.Exec(`INSERT INTO roles VALUES (10, 1, 'Lead')`)
	db.Exec(`INSERT INTO roles VALUES (20, 2, 'Sidekick')`)

	return db
}

func TestRunQuerySimpleSelect(t *testing.T) {
	db := buildTestQueryDB(t)
	defer db.Close()

	rows, err := db.Query(`SELECT name, rank FROM movies ORDER BY rank DESC`)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}

	defer rows.Close()

	// Just print to stdout (prints during test run)
	printTable(rows)
}

func TestPredefinedQueriesRun(t *testing.T) {
	db := buildTestQueryDB(t)
	defer db.Close()

	for key, q := range Queries {
		rows, err := db.Query(q.SQL)
		if err != nil {
			t.Fatalf("predefined query %s failed: %v", key, err)
		}
		rows.Close() // does it run is all that matter rn
	}
}
