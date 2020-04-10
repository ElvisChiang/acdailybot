package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func resetDB() (err error) {
	err = nil
	os.Remove(getDBFilename())

	db, err := sql.Open("sqlite3", getDBFilename())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	template := `
	create table %s (channelid integer not null primary key, name text, message text);
	delete from %s;
	`
	sqlStmt := fmt.Sprintf(template, getDBTableName(), getDBTableName())

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	return
}

func initDB(isReset bool) (db *sql.DB, err error) {
	err = nil

	if isReset {
		ret := resetDB()

		if ret != nil {
			log.Printf("DB reset failed")
			return
		}
	}

	db, err = sql.Open("sqlite3", getDBFilename())
	if err != nil {
		log.Fatal(err)
	}

	return
}

func removeHLEntry(channel string, who string) (err error) {
	err = nil
	// TODO: remove entry in db
	return
}

func insertHLEntry(channel string, who string, message string) (err error) {
	err = nil
	// TODO: add entry in db
	return
}

func queryHLEntry(channel string, who string) (message string, err error) {
	err = nil
	message = ""
	// TOOD: query entry in db

	return
}

func replaceHLEntry(channel string, who string) (err error) {
	err = nil
	removeHLEntry(channel, who) // ignore

	// TODO: insert
	return
}

func queryAllHLEntry(channel string) (message string, err error) {
	err = nil
	message = ""

	// TODO: query all
	return
}
