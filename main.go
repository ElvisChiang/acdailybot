package main

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// should only for debug
const resetDBAtStartup = false
const backdoorUser = "elvisfb"

var lastHLMsg map[int64]int
var lastTurnipMsg map[int64]int

const typeMessageHighlight = 1
const typeMessageTurnip = 2

var sendMessageType = 0

func main() {

	lastHLMsg = make(map[int64]int)
	lastTurnipMsg = make(map[int64]int)

	db, err := initDB(resetDBAtStartup)

	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if !(update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup") {
			log.Printf("Non group message from %s, ignore", update.Message.Chat.UserName)
			continue
		}

		user := tgbotapi.ChatConfigWithUser{
			ChatID:             update.Message.Chat.ID,
			SuperGroupUsername: update.Message.Chat.UserName,
			UserID:             update.Message.From.ID,
		}

		talker, err := bot.GetChatMember(user)

		adminMark := ""
		if talker.IsAdministrator() || talker.IsCreator() {
			adminMark = "*"
		}
		log.Printf("[%d/%s%s] `%s`",
			update.Message.Chat.ID,
			update.Message.From.UserName,
			adminMark,
			update.Message.Text)

		if len(strings.TrimSpace(update.Message.Text)) == 0 {
			continue
		}

		username := update.Message.From.UserName
		if len(username) == 0 {
			username = update.Message.From.FirstName
		}
		if len(username) == 0 {
			username = update.Message.From.LastName
		}
		if len(username) == 0 {
			continue
		}
		result, err := Command(db, update.Message.Chat.ID, username, talker.IsAdministrator() || talker.IsCreator(), update.Message.Text)

		if err != nil {
			log.Printf("failed to process [%s] `%s`: `%s`", update.Message.From.UserName, update.Message.Text, err.Error())
			result = err.Error()
		}

		if result == "" {
			// result = "ok"
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, result)
		msg.ParseMode = "Markdown"
		msg.DisableWebPagePreview = true
		msg.ReplyToMessageID = update.Message.MessageID
		ret, err := bot.Send(msg)

		if err == nil {
			var lastMsg = lastHLMsg
			if sendMessageType == 0 {
				continue
			} else if sendMessageType == typeMessageTurnip {
				lastMsg = lastTurnipMsg
			}
			log.Printf("messageID %d -> %d in chat %d\n", ret.MessageID, lastMsg[ret.Chat.ID], ret.Chat.ID)
			if lastMsg[ret.Chat.ID] != 0 {
				bot.DeleteMessage(tgbotapi.DeleteMessageConfig{ChatID: ret.Chat.ID, MessageID: lastMsg[ret.Chat.ID]})
			}
			lastMsg[ret.Chat.ID] = ret.MessageID
		}
	}
}
