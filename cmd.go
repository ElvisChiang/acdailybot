package main

import (
	"database/sql"
	"errors"
	"log"
	"strings"
)

// Command main entry for parsing command
func Command(db *sql.DB, channel int64, who string, isAdmin bool, msg string) (result string, err error) {
	err = errors.New("execuse me?")
	result = ""

	cmdstring := strings.Fields(msg)
	if len(cmdstring) == 0 { // should be never happen
		err = nil
		return
	}
	command := cmdstring[0]
	lowerCmd := strings.ToLower(command)
	param := strings.TrimSpace(msg[len(lowerCmd):])
	adminMark := ""
	if isAdmin {
		adminMark = "*"
	}
	log.Printf("command [%s%s] `%s` `%s`\n", who, adminMark, lowerCmd, param)

	switch lowerCmd {
	case "/hl":
		err = Highlight(db, channel, who, param)
	case "/hl_remove":
		err = Remove(db, channel, who, param, isAdmin)
	case "/hl_reset":
		err = ResetAll(db, channel, isAdmin)
	case "/hl_list":
		result, err = HighlightList(db, channel)
	}
	return
}

// Highlight my data /hl
func Highlight(db *sql.DB, channel int64, who string, msg string) (err error) {
	err = nil

	return
}

// Remove my data /remove
func Remove(db *sql.DB, channel int64, who string, user string, isAdmin bool) (err error) {
	err = nil

	if user != who && isAdmin == false {
		return errors.New("normal user can only remove his/her own data")
	}

	return
}

// HighlightList all the result in DB
func HighlightList(db *sql.DB, channel int64) (result string, err error) {
	err = nil
	result = ""

	return
}

// ResetAll in channel
func ResetAll(db *sql.DB, channel int64, isAdmin bool) (err error) {
	err = nil
	if isAdmin == false {
		return errors.New("you are not admin")
	}

	return
}
