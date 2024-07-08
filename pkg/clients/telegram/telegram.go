package telegram

import (
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func NewTelegramClient(token string, debug bool) *tgApi.BotAPI {
	bot, err := tgApi.NewBotAPI(token)
	if err != nil {
		log.Fatal("Failed to create bot: ", err)
	} else {
		log.Printf("Bot created: %v", bot.Self.UserName)
	}

	bot.Debug = debug

	return bot
}
