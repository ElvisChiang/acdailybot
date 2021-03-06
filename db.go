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
	create table %s (id integer not null primary key, channelid integer, name text, message text);
	delete from %s;
	create table %s (id integer not null primary key, channelid integer, name text, buy integer,
		sell1_am integer, sell1_pm integer,
		sell2_am integer, sell2_pm integer,
		sell3_am integer, sell3_pm integer,
		sell4_am integer, sell4_pm integer,
		sell5_am integer, sell5_pm integer,
		sell6_am integer, sell6_pm integer
		);
	delete from %s;
	`
	sqlStmt := fmt.Sprintf(template,
		getDBTableName(), getDBTableName(),
		getTurnipDBTableName(), getTurnipDBTableName(),
	)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	return
}

func initDB(isReset bool) (db *sql.DB, err error) {
	err = nil

	_, err = os.Stat(getDBFilename())
	if os.IsNotExist(err) {
		log.Println("DB doesn't exist, init it")
		isReset = true
	}

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

func removeHLEntry(db *sql.DB, channelid int64, who string) (err error) {
	err = nil

	template := `
	DELETE from %s where channelid = %d and name = ?;
	`
	sqlStmt := fmt.Sprintf(template, getDBTableName(), channelid)

	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
	_, err = stmt.Exec(who)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	return
}

func insertHLEntry(db *sql.DB, channelid int64, who string, message string) (err error) {
	err = nil

	template := `
	INSERT INTO %s(channelid, name, message) VALUES (?, ?, ?);
	`
	sqlStmt := fmt.Sprintf(template, getDBTableName())
	stmt, err := db.Prepare(sqlStmt)
	defer stmt.Close()
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
	_, err = stmt.Exec(channelid, who, message)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	return
}

func queryHLEntry(db *sql.DB, channelid int64, who string) (message string, err error) {
	err = nil
	message = ""
	// TODO: query entry in db

	return
}

func replaceHLEntry(db *sql.DB, channelid int64, who string, message string) (err error) {
	err = nil

	err = removeHLEntry(db, channelid, who)
	if err != nil {
		return
	}

	err = insertHLEntry(db, channelid, who, message)
	if err != nil {
		return
	}

	return
}

func queryAllHLEntry(db *sql.DB, channelid int64) (message string, err error) {
	err = nil
	message = ""

	template := `
	SELECT name, message from %s where channelid = %d
	`
	sqlStmt := fmt.Sprintf(template, getDBTableName(), channelid)

	rows, err := db.Query(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var msg string
		err = rows.Scan(&name, &msg)
		if err != nil {
			log.Printf(err.Error())
			continue
		}
		message = fmt.Sprintf("%s```%s:```\t%s \n", message, name, msg)
		fmt.Println(name, msg)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err.Error())
		return
	}

	return
}

func resetAllHLEntry(db *sql.DB, channelid int64) (err error) {
	err = nil

	template := `
	DELETE from %s where channelid = %d;
	`
	sqlStmt := fmt.Sprintf(template, getDBTableName(), channelid)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	return
}
