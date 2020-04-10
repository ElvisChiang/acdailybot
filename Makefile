SRC = $(shell find . -type f -name '*.go' | grep -v test | grep -v vendor)
LDFLAG = # -ldflags "-linkmode external -extldflags -static"

all: build

build: clean
	go build ${LDFLAG} -o acdailybot ${SRC}

run:
	@echo "source code: ${SRC}"
	@go run ${SRC}

clean:
	@rm -f acdailybot

dump:
	@echo ".dump highlight" | sqlite3 acbot.db
