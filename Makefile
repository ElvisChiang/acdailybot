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
	@echo ".dump turnip" | sqlite3 acbot.db

# Example for reset db
# record it in crontab
resethl:
	@echo "delete from highlight where channelid = -436800666" | sqlite3 acbot.db
resetturnip:
	@echo "delete from turnip where channelid = -436800666" | sqlite3 acbot.db
