# mod-wide binaries
CREATE_BIN := bin/create_imdb
QUERY_BIN  := bin/query_imdb

# Default target
.PHONY: all
all: test

# ------- testing --------------

.PHONY: test
test:
	go test ./...

.PHONY: create-test
create-test:
	go test ./cmd/create/...

.PHONY: query-test
query-test:
	go test ./cmd/query/...

# ------- build and run ---------------------

# build only create
$(CREATE_BIN): cmd/create/*.go internal/*.go
	mkdir -p bin
	go build -o $(CREATE_BIN) ./cmd/create

# build only query
$(QUERY_BIN): cmd/query/*.go internal/*.go
	mkdir -p bin
	go build -o $(QUERY_BIN) ./cmd/query

.PHONY: create
create: $(CREATE_BIN)
	$(CREATE_BIN)

# custom flags get passed onto the go file
# make query ARGS="--sql 'SELECT * FROM movies LIMIT 5;'"
.PHONY: query
query: $(QUERY_BIN)
	$(QUERY_BIN) $(ARGS)

# --------------- utility ----------------------

.PHONY: clean
clean:
	rm -rf bin imdb.db imdb.db-shm imdb.db-wal
