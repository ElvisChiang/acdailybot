package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// should only for debug
const resetDBAtStartup = false

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

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.Chat.Type != "group" {
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
		log.Printf("[%d/%s%s] %s",
			update.Message.Chat.ID,
			update.Message.From.UserName,
			adminMark,
			update.Message.Text)

		result, err := Command(db, update.Message.Chat.ID, update.Message.From.UserName, talker.IsAdministrator(), update.Message.Text)

		if err != nil {
			log.Printf("failed to process [%s] `%s`: `%s`", update.Message.From.UserName, update.Message.Text, err.Error())
		}

		if result == "" {
			return
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ParseMode = "Markdown"
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}
