package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"

	_ "modernc.org/sqlite"
)

// premade queries
var Queries = map[string]struct {
	Description string
	SQL         string
}{
	"avg-genres": {
		Description: "Highest rating genres",
		SQL: `
            SELECT g.genre, COUNT(*) AS movie_count, ROUND(AVG(m.rank), 2) AS avg_rank
            FROM genres g
            JOIN movies m ON g.movie_id = m.movie_id
            GROUP BY g.genre
            ORDER BY avg_rank DESC;
        `,
	},

	"top-actors": {
		Description: "Actors with highest avg-rated movies (min 3 movies)",
		SQL: `
            SELECT a.first_name || ' ' || a.last_name AS actor_name,
            COUNT(*) AS movie_count,
            ROUND(AVG(m.rank), 2) AS avg_rank
            FROM actors a
            JOIN roles r ON a.actor_id = r.actor_id
            JOIN movies m ON r.movie_id = m.movie_id
            GROUP BY a.actor_id
            HAVING movie_count >= 3
            ORDER BY avg_rank DESC
            LIMIT 20;
        `,
	},

	"co-stars": {
		Description: "Most popular actor pairs/co-stars",
		SQL: `
            SELECT a1.first_name || ' ' || a1.last_name AS actor1,
            a2.first_name || ' ' || a2.last_name AS actor2,
            COUNT(*) AS shared_movies
            FROM roles r1
            JOIN roles r2 ON r1.movie_id = r2.movie_id AND r1.actor_id < r2.actor_id
            JOIN actors a1 ON a1.actor_id = r1.actor_id
            JOIN actors a2 ON a2.actor_id = r2.actor_id
            GROUP BY actor1, actor2
            ORDER BY shared_movies DESC
            LIMIT 20;
        `,
	},

	"above-genre-average": {
		Description: "Movies with ratings above their genre average",
		SQL: `
            WITH genre_avg AS (
                SELECT genre, AVG(rank) AS avg_rank
                FROM genres JOIN movies USING(movie_id)
                GROUP BY genre
            )
            SELECT m.name, g.genre, m.rank, ROUND(ga.avg_rank,2) AS genre_avg
            FROM movies m
            JOIN genres g ON m.movie_id = g.movie_id
            JOIN genre_avg ga ON g.genre = ga.genre
            WHERE m.rank > ga.avg_rank
            ORDER BY (m.rank - ga.avg_rank) DESC
            LIMIT 30;
        `,
	},
}

func main() {

	// open imdb sqlite db
	db, err := sql.Open("sqlite", "imdb.db")
	if err != nil {
		log.Fatalf("Failed to open imdb.db: %v\n"+
			"If the database has not been created yet, run:\n\n"+
			"    go run ./cmd/create\n", err)
	}
	defer db.Close()

	fmt.Println("Connected to imdb.db")

	// Parse flags
	listFlag := flag.Bool("list", false, "List all predefined query options")
	queryFlag := flag.String("query", "", "Run a predefined query by name")
	sqlFlag := flag.String("sql", "", "Run a custom SQL query")
	flag.Parse()

	if *listFlag {
		fmt.Println("Available queries:")
		for name, q := range Queries {
			fmt.Printf("  %-20s  %s\n", name, q.Description)
		}
		return
	}

	if *sqlFlag != "" {
		runQuery(db, *sqlFlag)
		return
	}

	if *queryFlag != "" {
		q, exists := Queries[*queryFlag]
		if !exists {
			log.Fatalf("Unknown query: %s. Use --list to see options.", *queryFlag)
		}
		runQuery(db, q.SQL)
		return
	}

	fmt.Println("No query provided. Use --list or --query <name> or --sql <query-string>")
}

// runs query in sqlite db
func runQuery(db *sql.DB, sql string) {
	rows, err := db.Query(sql)
	if err != nil {
		log.Fatalf("Query error: %v\nSQL: %s", err, sql)
	}
	defer rows.Close()

	printTable(rows)
	fmt.Println("\nQuery completed.")
}

// pads the columns
func pad(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

// prints a border like: +-----+--------+----+
func printBorder(colWidths []int) {
	fmt.Print("+")
	for _, w := range colWidths {
		fmt.Print(strings.Repeat("-", w+2) + "+")
	}
	fmt.Println()
}

// pretty prints tables of query results
func printTable(rows *sql.Rows) {
	cols, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	// read all rows into memory as [][]string
	values := make([]interface{}, len(cols))
	for i := range values {
		var v interface{}
		values[i] = &v
	}

	var data [][]string

	for rows.Next() {
		err := rows.Scan(values...)
		if err != nil {
			log.Fatal(err)
		}

		row := make([]string, len(cols))
		for i, ptr := range values {
			v := *(ptr.(*interface{}))
			row[i] = fmt.Sprint(v)
		}
		data = append(data, row)
	}

	// calc column widths
	colWidths := make([]int, len(cols))
	for i, col := range cols {
		colWidths[i] = len(col)
	}

	for _, row := range data {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// print header
	printBorder(colWidths)
	fmt.Print("|")
	for i, col := range cols {
		fmt.Printf(" %s |", pad(col, colWidths[i]))
	}
	fmt.Println()
	printBorder(colWidths)

	// print rows
	for _, row := range data {
		fmt.Print("|")
		for i, cell := range row {
			fmt.Printf(" %s |", pad(cell, colWidths[i]))
		}
		fmt.Println()
	}
	printBorder(colWidths)
}
