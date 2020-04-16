package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func removeTurnipEntry(db *sql.DB, channelid int64, who string) (err error) {
	err = nil

	template := `
	DELETE from %s where channelid = %d and name = ?;
	`
	sqlStmt := fmt.Sprintf(template, getTurnipDBTableName(), channelid)

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

func insertTurnipEntry(db *sql.DB, channelid int64, who string, price *Price) (err error) {
	err = nil

	template := `
	INSERT INTO %s(channelid, name, buy,
		sell1_am, sell1_pm,
		sell2_am, sell2_pm,
		sell3_am, sell3_pm,
		sell4_am, sell4_pm,
		sell5_am, sell5_pm,
		sell6_am, sell6_pm
		) VALUES (%d, "%s", %d,
		%d, %d,
		%d, %d,
		%d, %d,
		%d, %d,
		%d, %d,
		%d, %d
		);
	`
	sqlStmt := fmt.Sprintf(template, getTurnipDBTableName(), channelid, who, price.buy,
		price.sell[0], price.sell[1],
		price.sell[2], price.sell[3],
		price.sell[4], price.sell[5],
		price.sell[6], price.sell[7],
		price.sell[8], price.sell[9],
		price.sell[10], price.sell[11],
	)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	return
}

func queryTurnipEntry(db *sql.DB, channelid int64, who string) (message string, err error) {
	err = nil
	message = ""
	// TOOD: query turnip entry in db

	return
}

func replaceTurnipEntry(db *sql.DB, channelid int64, who string, price *Price) (err error) {
	err = nil

	err = removeTurnipEntry(db, channelid, who)
	if err != nil {
		return
	}

	err = insertTurnipEntry(db, channelid, who, price)
	if err != nil {
		return
	}

	return
}

func queryAllTurnipEntry(db *sql.DB, channelid int64) (message string, err error) {
	err = nil
	message = ""
	/* 	SELECT name, buy,
	   	sell1_am, sell1_pm,
	   	sell2_am, sell2_pm
	   	sell3_am, sell3_pm
	   	sell4_am, sell4_pm
	   	sell5_am, sell5_pm
	   	sell6_am, sell6_pm
	*/
	template := `select * from %s where channelid = %d
	`
	sqlStmt := fmt.Sprintf(template, getTurnipDBTableName(), channelid)

	rows, err := db.Query(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var msg string
		price := new(Price)
		var tmp int
		err = rows.Scan(&tmp, &tmp, &name, &price.buy,
			&price.sell[0], &price.sell[1],
			&price.sell[2], &price.sell[3],
			&price.sell[4], &price.sell[5],
			&price.sell[6], &price.sell[7],
			&price.sell[8], &price.sell[9],
			&price.sell[10], &price.sell[11],
		)
		if err != nil {
			log.Printf(err.Error())
			continue
		}
		predictURL := fmt.Sprintf("[預測](https://turnipprophet.io/index.html?prices=%d.%d.%d.%d.%d.%d.%d.%d.%d.%d.%d.%d.%d)",
			price.buy,
			price.sell[0], price.sell[1],
			price.sell[2], price.sell[3],
			price.sell[4], price.sell[5],
			price.sell[6], price.sell[7],
			price.sell[8], price.sell[9],
			price.sell[10], price.sell[11],
		)
		message = fmt.Sprintf("%s```%s:``` 買 %d 賣 %d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d %s\n",
			message, name, price.buy,
			price.sell[0], price.sell[1],
			price.sell[2], price.sell[3],
			price.sell[4], price.sell[5],
			price.sell[6], price.sell[7],
			price.sell[8], price.sell[9],
			price.sell[10], price.sell[11],
			predictURL,
		)
		fmt.Println(name, msg)
	}
	for true {
		newMsg := strings.ReplaceAll(message, ",0 ", " ")
		if newMsg == message {
			break
		}
		message = newMsg
	}

	fmt.Printf("DBUG:\n%s\n", message)
	err = rows.Err()
	if err != nil {
		log.Println(err.Error())
		return
	}

	return
}

func resetAllTurnipEntry(db *sql.DB, channelid int64) (err error) {
	err = nil

	template := `
	DELETE from %s where channelid = %d;
	`
	sqlStmt := fmt.Sprintf(template, getTurnipDBTableName(), channelid)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	return
}
