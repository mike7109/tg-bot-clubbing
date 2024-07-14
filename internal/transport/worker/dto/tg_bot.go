package dto

import (
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mike7109/tg-bot-clubbing/internal/service/dto"
)

func ParseMessageForPage(msg *tgApi.Message) (*dto.SavePage, error) {
	page := dto.NewSavePage(msg.Text, msg.From.UserName)
	if err := page.Validate(); err != nil {
		return nil, err
	}

	return page, nil
}
