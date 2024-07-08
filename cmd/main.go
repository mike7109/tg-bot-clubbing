package main

import (
	"github.com/mike7109/tg-bot-clubbing/internal/core"
	"log"
)

func main() {
	tgBot, err := core.New()
	if err != nil {
		log.Fatal("Failed to create bot: ", err)
	}

	if err = tgBot.Start(); err != nil {
		log.Fatal("Failed to start bot: ", err)
	}
}
