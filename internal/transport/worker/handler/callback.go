package handler

import (
	"context"
	"errors"
	"fmt"
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mike7109/tg-bot-clubbing/internal/apperrors"
	"github.com/mike7109/tg-bot-clubbing/internal/entity"
	"github.com/mike7109/tg-bot-clubbing/internal/service"
	"github.com/mike7109/tg-bot-clubbing/internal/transport/worker/update_entity/button"
	"github.com/mike7109/tg-bot-clubbing/pkg/messages"
)

type CallbackHandlerFunc func(ctx context.Context, callback *tgApi.CallbackQuery, buttonTarget *button.Button) error

func DeleteUrl(ctx context.Context, tgBot *tgApi.BotAPI, tgBotService service.ITgBotService) CallbackHandlerFunc {
	return func(ctx context.Context, callbackUpdate *tgApi.CallbackQuery, buttonTarget *button.Button) error {
		callback := tgApi.NewCallback(callbackUpdate.ID, callbackUpdate.Data)
		if _, err := tgBot.Request(callback); err != nil {
			return err
		}

		userName := callbackUpdate.Message.Chat.UserName
		chatID := callbackUpdate.Message.Chat.ID
		messageID := callbackUpdate.Message.MessageID

		buttonID, exist := button.GetDataValue(buttonTarget, "id")
		if !exist {
			return fmt.Errorf("buttonID not found in data")
		}

		err := tgBotService.DeleteHandler(ctx, int(buttonID.(float64)), userName)
		if err != nil {
			return err
		}

		pages, err := tgBotService.ListHandler(ctx, userName)
		if err != nil {
			switch {
			case errors.Is(err, apperrors.ErrNoPages):
				msgZeroPage := editListMsgForNotPAge(chatID, messageID)
				_, err = tgBot.Send(msgZeroPage)
				return err
			case errors.Is(err, apperrors.ErrNoSave):
				return messages.SendErrorHandler(tgBot, chatID)
			default:
				return messages.SendErrorHandler(tgBot, chatID)
			}
		}

		msg := editListMsg(pages, chatID, messageID)

		_, err = tgBot.Send(msg)

		return err
	}
}

func editListMsg(pages []*entity.Page, chatID int64, MessageID int) tgApi.Chattable {
	var msg string

	builder := button.NewBuilder()

	for _, page := range pages {
		msg += page.String()
		but := page.ToButton(button.DeleteCommand)
		builder.AddButton(but)
	}

	keyboard := builder.Build()

	msgConfig := tgApi.NewEditMessageText(chatID, MessageID, msg)
	msgConfig.DisableWebPagePreview = true
	msgConfig.ReplyMarkup = &keyboard

	return msgConfig
}

func editListMsgForNotPAge(chatID int64, MessageID int) tgApi.Chattable {
	msgConfig := tgApi.NewEditMessageText(chatID, MessageID, messages.MsgNoSavedPages)
	msgConfig.DisableWebPagePreview = true
	msgConfig.ReplyMarkup = nil

	return msgConfig
}
