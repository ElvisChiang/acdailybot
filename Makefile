SRC = $(shell find . -type f -name '*.go' | grep -v test | grep -v vendor)

all: build

build:
	@GOOS=linux GOARCH=amd64 go build ${SRC}

run:
	@echo "source code: ${SRC}"
	@go run ${SRC}

dump:
	@echo ".dump highlight" | sqlite3 acbot.db
