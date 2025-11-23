package internal

const Schema = `
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS genres;
DROP TABLE IF EXISTS actors;
DROP TABLE IF EXISTS movies;

CREATE TABLE movies (
    movie_id      INTEGER PRIMARY KEY,
    name          TEXT NOT NULL,
    year          INTEGER,
    rank          REAL
);

CREATE TABLE actors (
    actor_id     INTEGER PRIMARY KEY,
    first_name   TEXT,
    last_name    TEXT,
    gender       TEXT
);

CREATE TABLE genres (
    movie_id    INTEGER,
    genre       TEXT,
    FOREIGN KEY (movie_id) REFERENCES movies(movie_id)
);

CREATE INDEX idx_genres_movie_id ON genres(movie_id);

CREATE TABLE roles (
    actor_id    INTEGER,
    movie_id    INTEGER,
    role        TEXT,
    FOREIGN KEY (actor_id) REFERENCES actors(actor_id),
    FOREIGN KEY (movie_id) REFERENCES movies(movie_id)
);

CREATE INDEX idx_roles_actor_id ON roles(actor_id);
CREATE INDEX idx_roles_movie_id ON roles(movie_id);
`
