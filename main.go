package main

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// should only for debug
const resetDBAtStartup = false
const backdoorUser = "elvisfb"

// Debug for detail api message
const Debug = false

func main() {

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

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.From.UserName != backdoorUser {
			if !(update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup") {
				log.Printf("Non group message from %s, ignore", update.Message.Chat.UserName)
				continue
			}
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

		result, err := Command(db, update.Message.Chat.ID, update.Message.From.UserName, talker.IsAdministrator() || talker.IsCreator(), update.Message.Text)

		if err != nil {
			log.Printf("failed to process [%s] `%s`: `%s`", update.Message.From.UserName, update.Message.Text, err.Error())
			result = err.Error()
		}

		if result == "" {
			result = "ok"
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, result)
		msg.ParseMode = "Markdown"
		msg.DisableWebPagePreview = true
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}
