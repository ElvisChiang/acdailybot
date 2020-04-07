SRC = $(shell find . -type f -name '*.go')

all: build

build:
	@GOOS=linux GOARCH=amd64 go build ${SRC}

run:
	@go run ${SRC}

