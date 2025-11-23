package main

import (
	"database/sql"
	"imdbsqlite/internal"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

// create temp csv
func writeTempCSV(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed writing temp CSV %s: %v", name, err)
	}
	return path
}

// create temp db cxn
func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	// in memory sqlite to avoid messing w/ imdb.sql
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}

	_, err = db.Exec(internal.Schema) // from create_imdb.go (schema const)
	if err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}
	return db
}

func TestCreateAndLoadMinimalCSV(t *testing.T) {
	tmpDir := t.TempDir()

	// tiny datasets, 1 row each
	moviesCSV := writeTempCSV(t, tmpDir, "movies.csv",
		"movie_id,name,year,rank\n1,Test Movie,2000,7.5\n")

	actorsCSV := writeTempCSV(t, tmpDir, "actors.csv",
		"actor_id,first_name,last_name,gender\n10,Jane,Doe,F\n")

	genresCSV := writeTempCSV(t, tmpDir, "genres.csv",
		"movie_id,genre\n1,Drama\n")

	rolesCSV := writeTempCSV(t, tmpDir, "roles.csv",
		"actor_id,movie_id,role\n10,1,Lead\n")

	// read CSVs using readCSV() from create_imdb.go
	moviesData, err := readCSV(moviesCSV)
	if err != nil {
		t.Fatalf("readCSV movies: %v", err)
	}

	actorsData, err := readCSV(actorsCSV)
	if err != nil {
		t.Fatalf("readCSV actors: %v", err)
	}

	genresData, err := readCSV(genresCSV)
	if err != nil {
		t.Fatalf("readCSV genres: %v", err)
	}

	rolesData, err := readCSV(rolesCSV)
	if err != nil {
		t.Fatalf("readCSV roles: %v", err)
	}

	// open temp DB
	db := openTestDB(t)
	defer db.Close()

	// insert using write functions from create_imdb.go
	if err := writeMovies(db, moviesData); err != nil {
		t.Fatalf("writeMovies failed: %v", err)
	}
	if err := writeActors(db, actorsData); err != nil {
		t.Fatalf("writeActors failed: %v", err)
	}
	if err := writeGenres(db, genresData); err != nil {
		t.Fatalf("writeGenres failed: %v", err)
	}
	if err := writeRoles(db, rolesData); err != nil {
		t.Fatalf("writeRoles failed: %v", err)
	}

	// validate row counts
	checkCount := func(table string, expected int) {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count)
		if err != nil {
			t.Fatalf("count query failed on %s: %v", table, err)
		}
		if count != expected {
			t.Fatalf("%s: expected %d rows, got %d", table, expected, count)
		}
	}

	checkCount("movies", 1)
	checkCount("actors", 1)
	checkCount("genres", 1)
	checkCount("roles", 1)
}
