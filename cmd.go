package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

// TODO: store in db, used to be a global motd
var motd = ""

// Command main entry for parsing command
func Command(db *sql.DB, channel int64, who string, isAdmin bool, msg string) (result string, err error) {
	err = errors.New("execuse me?")
	result = ""

	msg2 := strings.ReplaceAll(msg, "@ac_dailybot", "")
	cmdstring := strings.Fields(msg2)
	if len(cmdstring) == 0 { // should be never happen
		err = nil
		return
	}
	command := cmdstring[0]
	lowerCmd := strings.ToLower(command)
	param := strings.TrimSpace(msg2[len(lowerCmd):])
	adminMark := ""
	if isAdmin {
		adminMark = "*"
	}
	log.Printf("command [%s%s] `%s` `%s`", who, adminMark, lowerCmd, param)

	var replyAllList = false
	switch lowerCmd {
	case "/motd":
		if who == backdoorUser {
			log.Printf("set motd = `%s`", param)
			motd = param
			err = nil
			result = ""
		}
	case "/hl":
		err = Highlight(db, channel, who, param)
		if err == nil {
			replyAllList = true
		}
	case "/hl_remove":
		err = Remove(db, channel, who, param, isAdmin)
	case "/hl_reset":
		log.Printf("reset `%s` `%s`", lowerCmd, param)
		if who == param {
			err = ResetAll(db, channel, isAdmin)
		} else {
			err = errors.New("簽名以重設資料")
		}
	case "/hl_list":
		// result, err = HighlightList(db, channel)
		replyAllList = true
		err = nil
	default:
		// err = errors.New("what's run?")
		log.Printf("receive unknown message: %s", msg)
		return
	}
	if err != nil {
		return
	}

	if replyAllList {
		result, err = HighlightList(db, channel)

		if err == nil {
			fmt.Printf("DEBUG: motd `%s` result `%s`", motd, result)
			result = motd + "\n" + "=== #動森高光 ===\n" + result
			fmt.Printf("DEBUG: result `%s`\n", result)
		}
	}
	return
}

// Highlight my data /hl
func Highlight(db *sql.DB, channelid int64, who string, msg string) (err error) {
	err = nil

	if msg == "" {
		err = errors.New("say something")
		return
	}

	return replaceHLEntry(db, channelid, who, msg)
}

// Remove my data /remove
func Remove(db *sql.DB, channelid int64, who string, user string, isAdmin bool) (err error) {
	err = nil

	if user == "" {
		user = who
	}

	if user != who && isAdmin == false {
		log.Printf("normal user can only remove his/her own data")
		return errors.New("normal user can only remove his/her own data")
	}

	return removeHLEntry(db, channelid, user)
}

// HighlightList all the result in DB
func HighlightList(db *sql.DB, channelid int64) (result string, err error) {
	err = nil
	result = ""

	return queryAllHLEntry(db, channelid)
}

// ResetAll in channel
func ResetAll(db *sql.DB, channelid int64, isAdmin bool) (err error) {
	err = nil

	if isAdmin == false {
		return errors.New("只有管理員能重設資料")
	}

	return resetAllHLEntry(db, channelid)
}
