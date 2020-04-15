package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// TODO: store in db, used to be a global motd
var motd = ""

// Command main entry for parsing command
func Command(db *sql.DB, channel int64, who string, isAdmin bool, msg string) (result string, err error) {
	err = errors.New("execuse me?")
	result = ""

	msg2 := strings.ReplaceAll(msg, "@ac_dailybot", "")
	msg2 = strings.ReplaceAll(msg2, "@testturnipbot", "")

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
	var replyAllTurnipList = false
	sendMessageType = 0
	switch lowerCmd {
	case "/motd":
		if who == backdoorUser {
			log.Printf("set motd = `%s`", param)
			motd = param
			err = nil
			result = ""
		}
	case "/hl", "/hi":
		err = Highlight(db, channel, who, param)
		if err == nil {
			replyAllList = true
		}
	case "/hl_remove":
		err = Remove(db, channel, who, param, isAdmin, false)
	case "/hl_reset":
		log.Printf("reset `%s` `%s`", lowerCmd, param)
		if who == param {
			err = ResetAll(db, channel, isAdmin, false)
		} else {
			err = errors.New("簽名以重設資料")
		}
	case "/hl_list":
		replyAllList = true
		err = nil
	case "/turnip", "/posh":
		err = Turnip(db, channel, who, strings.ToLower(msg))
		if err == nil {
			replyAllTurnipList = true
		}
	case "/turnip_remove":
		err = Remove(db, channel, who, param, isAdmin, true)
	case "/turnip_reset":
		log.Printf("reset `%s` `%s`", lowerCmd, param)
		if who == param {
			err = ResetAll(db, channel, isAdmin, true)
		} else {
			err = errors.New("簽名以重設資料")
		}
	case "/turnip_list":
		replyAllTurnipList = true
		err = nil
	default:
		// err = errors.New("what's run?")
		log.Printf("receive unknown message: %s", msg)
		err = nil
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
		}
		sendMessageType = typeMessageHighlight
		return
	}
	if replyAllTurnipList {
		result, err = TurnipList(db, channel)

		if err == nil {
			fmt.Printf("DEBUG: motd `%s` result `%s`", motd, result)
			result = motd + "\n" + "=== #動森大頭菜 ===\n" + result
		}
		sendMessageType = typeMessageTurnip
		return
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
func Remove(db *sql.DB, channelid int64, who string, user string, isAdmin bool, isTurnip bool) (err error) {
	err = nil

	if user == "" {
		user = who
	}

	if user != who && isAdmin == false {
		log.Printf("normal user can only remove his/her own data")
		return errors.New("normal user can only remove his/her own data")
	}

	if isTurnip {
		removeTurnipEntry(db, channelid, user)
		return
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
func ResetAll(db *sql.DB, channelid int64, isAdmin bool, isTurnip bool) (err error) {
	err = nil

	if isAdmin == false {
		return errors.New("只有管理員能重設資料")
	}

	if isTurnip {
		return resetAllTurnipEntry(db, channelid)
	}
	return resetAllHLEntry(db, channelid)
}

// --- Turnip ---

// Turnip my data /turnip
func Turnip(db *sql.DB, channelid int64, who string, msg string) (err error) {
	err = nil
	errReturn := errors.New("範例語法: \n/turnip buy 123 sell 101,-,103,")

	if msg == "" {
		err = errors.New("say something")
		return
	}

	price := new(Price)
	price.buy = 0
	r := regexp.MustCompile("\\s*buy\\s*(\\d+)\\s*sell(.*)$")

	regexResult := r.FindStringSubmatch(msg)

	if len(regexResult) == 0 {
		return errReturn
	}

	price.buy, err = strconv.Atoi(regexResult[1])

	if err != nil {
		return errReturn
	}

	trim := strings.ReplaceAll(regexResult[2], " ", "")
	sell := strings.Split(trim, ",")

	for i, s := range sell {
		if s == "" {
			break
		}
		if s == "-" {
			price.sell[i] = 0
			continue
		}
		p, err := strconv.Atoi(s)
		if err != nil {
			fmt.Printf("atoi error `%s`\n", s)
			return errReturn
		}
		if p > MaxOfTurnip {
			return errReturn
		}
		price.sell[i] = p
	}

	fmt.Printf("input: `%s` regex: buy `%d` sell `%d %d %d %d %d %d %d %d %d %d %d %d`\n", msg, price.buy,
		price.sell[0], price.sell[1],
		price.sell[2], price.sell[3],
		price.sell[4], price.sell[5],
		price.sell[6], price.sell[7],
		price.sell[8], price.sell[9],
		price.sell[10], price.sell[11],
	)

	return replaceTurnipEntry(db, channelid, who, price)
}

// TurnipList all the result in DB
func TurnipList(db *sql.DB, channelid int64) (result string, err error) {
	err = nil
	result = ""

	return queryAllTurnipEntry(db, channelid)
}
