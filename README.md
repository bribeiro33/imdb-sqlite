# MDb SQLite Builder & Query Tool

A Go-based system for constructing, querying, and extending an IMDb-style relational movie database.

## Overview

This project provides a fully working SQLite relational database constructed from four CSV files:

* movies – movie metadata

* actors – actor info

* genres – movie–genre relationships

* roles – actor roles in movies

The repository contains two Go CLI programs:

1. create_imdb

    * Reads all CSV files in parallel

    * Creates a clean SQLite schema

    * Loads 300k+ rows

    * Includes robust error handling and unit tests

2. query_imdb

    * Runs predefined or custom SQL queries

    * Prints results using a custom table formatter

    * Supports CLI flags and can accept full SQL manually

    * Unit tests

The project also includes a Makefile that automates builds, tests, and execution.

## Usage

### Build & populate the database
```{bash}
make create
```
This deletes any existing imdb.db, re-creates the schema, and loads all CSVs.

### Run queries
```{bash}
make query ARGS="--list"
```
to list the flags for the prewritten queries e.g.
```{bash}
make query ARGS="--top-genres"
```
or 
```{bash}
make query ARGS="--sql 'SELECT * FROM movies LIMIT 5;'"
```
run direct sql

### Run tests
```{bash}
make test        # all tests
make create-test # tests for the create command
make query-test  # tests for the query command
```

### Schema
The database schema (stored in internal/schema.go) writes four tables:
* movies(movie_id, name, year, rank)
* actors(actor_id, first_name, last_name, gender)
* genres(movie_id, genre)
* roles(actor_id, movie_id, role)

## Extending to a Personal Movie Collection
```{bash}
CREATE TABLE my_movies (
    my_id INTEGER PRIMARY KEY AUTOINCREMENT,
    movie_id INTEGER,          -- link to imbd data
    location TEXT,             -- streaming location or physical shelf
    personal_rating REAL,      
    notes TEXT,
    FOREIGN KEY(movie_id) REFERENCES movies(movie_id)
);
```
* Tracks which movies you own
* Tracks where you can find it
* Tracks your personal score (can then compare to imdb)
* Allows for notes like "best soundtrack" "rewatch" for some customizability/tagging

### Advantages
I could join my movies with the imdb movies and find out things like my collections breakdown by genre or decade.
I could also compare how different on average my ratings are from the imdb avgs, so I could see what genres I like more than the avg population.
I could create a UI using go to make it easy to update and naviagte my personal collection and then also run those predefined aggregate functions.
Able to modify how I'd like with new tables and indexes rather than interfacing with imdb directly. 

## Possible Future Enhancements
* Adding a reviews table (ref movie, ref reviewer) and reviewer table
* ML movie recommender givin the initial training dataset of my_movies and its data
* Especially for my_movies, a GUI of some sort to make seeing and updating my collection easier, as well as a dashboard with interesting aggregated stats.
