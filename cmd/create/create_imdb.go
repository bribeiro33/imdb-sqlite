package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"imdbsqlite/internal"
	"log"
	"os"
	"strconv"
	"sync"

	_ "modernc.org/sqlite"
)

func must(context string, err error) {
	if err != nil {
		log.Fatalf("\nError: %s\nDetails: %v\n", context, err)
	}
}

func main() {

	// delete db if it exists
	if _, err := os.Stat("imdb.db"); err == nil {
		fmt.Println("Removing existing imdb.db ...")
		must("unable to remove old imdb.db", os.Remove("imdb.db"))
	}

	// open new db
	db, err := sql.Open("sqlite", "imdb.db")
	must("failed to open sqlite database", err)
	defer db.Close()

	// create schema
	fmt.Println("Creating schema...")
	_, err = db.Exec(internal.Schema)
	must("failed to create schema", err)
	fmt.Println("Schema created")

	// load all csvs in parallel
	fmt.Println("Reading CSVs...")

	var moviesData, actorsData, genresData, rolesData [][]string
	var readErr error

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		moviesData, readErr = readCSV("data/aan520_w5a_movies2000After.csv")
	}()

	go func() {
		defer wg.Done()
		actorsData, readErr = readCSV("data/aan520_w5a_actors2000After.csv")
	}()

	go func() {
		defer wg.Done()
		genresData, readErr = readCSV("data/aan520_w5a_genres2000After.csv")
	}()

	go func() {
		defer wg.Done()
		rolesData, readErr = readCSV("data/aan520_w5a_roles2000After.csv")
	}()

	wg.Wait()

	if readErr != nil {
		must("failed while reading CSV files", readErr)
	}

	fmt.Println("Finished reading all CSV files")

	// db writes, sequential
	fmt.Println("Writing to SQLite...")

	if err := writeMovies(db, moviesData); err != nil {
		must("failed inserting movies", err)
	}
	if err := writeActors(db, actorsData); err != nil {
		must("failed inserting actors", err)
	}
	if err := writeGenres(db, genresData); err != nil {
		must("failed inserting genres", err)
	}
	if err := writeRoles(db, rolesData); err != nil {
		must("failed inserting roles", err)
	}

	fmt.Println("All data loaded successfully into imdb.db")
}

// csv reader helper
func readCSV(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %w", path, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not parse CSV %s: %w", path, err)
	}

	fmt.Println("Read:", len(records)-1, "rows from", path)
	return records, nil
}

// writer helpers for each table
// if increase db complexity, could create a general purpose write function
// that takes in the necessary table structure info

func writeMovies(db *sql.DB, data [][]string) error {
	fmt.Println("Inserting movies...")

	tx, err := db.Begin()
	must("begin tx for movies", err)
	stmt, err := tx.Prepare(`INSERT INTO movies(movie_id, name, year, rank) VALUES (?, ?, ?, ?)`)
	must("prepare movies insert", err)

	for i, row := range data {
		if i == 0 {
			continue
		} // skip headers

		// convert to db type
		id, _ := strconv.Atoi(row[0])
		year, _ := strconv.Atoi(row[2])
		rank, _ := strconv.ParseFloat(row[3], 64)

		_, err := stmt.Exec(id, row[1], year, rank)
		if err != nil {
			return fmt.Errorf("insert movie row %d: %w", i+1, err)
		}
	}

	stmt.Close()
	fmt.Println("Movies inserted:", len(data)-1)
	return tx.Commit()
}

func writeActors(db *sql.DB, data [][]string) error {
	fmt.Println("Inserting actors...")

	tx, err := db.Begin()
	must("begin tx for actors", err)
	stmt, err := tx.Prepare(`INSERT INTO actors(actor_id, first_name, last_name, gender) VALUES (?, ?, ?, ?)`)
	must("prep movie insert", err)

	for i, row := range data {
		if i == 0 {
			continue
		}

		id, _ := strconv.Atoi(row[0])

		_, err := stmt.Exec(id, row[1], row[2], row[3])
		if err != nil {
			return fmt.Errorf("insert actor row %d: %w", i+1, err)
		}
	}

	stmt.Close()
	fmt.Println("Actors inserted:", len(data)-1)
	return tx.Commit()
}

func writeGenres(db *sql.DB, data [][]string) error {
	fmt.Println("Inserting genres...")

	tx, err := db.Begin()
	must("begin tx for genres", err)
	stmt, err := tx.Prepare(`INSERT INTO genres(movie_id, genre) VALUES (?, ?)`)
	must("prepare genres insert", err)

	for i, row := range data {
		if i == 0 {
			continue
		}

		id, _ := strconv.Atoi(row[0])

		_, err := stmt.Exec(id, row[1])
		if err != nil {
			return fmt.Errorf("insert genre row %d: %w", i+1, err)
		}
	}

	stmt.Close()
	fmt.Println("Genres inserted:", len(data)-1)
	return tx.Commit()
}

func writeRoles(db *sql.DB, data [][]string) error {
	fmt.Println("Inserting roles...")

	tx, err := db.Begin()
	must("begin tx for roles", err)
	stmt, err := tx.Prepare(`INSERT INTO roles(actor_id, movie_id, role) VALUES (?, ?, ?)`)
	must("prepare roles insert", err)

	for i, row := range data {
		if i == 0 {
			continue
		}

		aid, _ := strconv.Atoi(row[0])
		mid, _ := strconv.Atoi(row[1])

		_, err := stmt.Exec(aid, mid, row[2])
		if err != nil {
			return fmt.Errorf("insert role row %d: %w", i+1, err)
		}
	}

	stmt.Close()
	fmt.Println("Roles inserted:", len(data)-1)
	return tx.Commit()
}
