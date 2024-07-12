package main

import (
	"github.com/mike7109/tg-bot-clubbing/internal/core"
	"log"
)

func main() {
	err := core.New()
	if err != nil {
		log.Fatal("Failed to create bot: ", err)
	}
}
