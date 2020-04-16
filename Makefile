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

# while sleep 300 ; do make dump ; done
dump:
	@mkdir -p backup
	@echo ".dump highlight" | sqlite3 acbot.db | tee backup/highlight_`date "+%Y%m%d-%H_%M_%S"`.sql
	@echo ".dump turnip" | sqlite3 acbot.db | tee backup/turnip_`date "+%Y%m%d-%H_%M_%S"`.sql

# Example for reset db
resethl:
	@echo "delete from highlight where channelid = -436800666" | sqlite3 acbot.db
resetturnip:
	@echo "delete from turnip where channelid = -436800666" | sqlite3 acbot.db
