package handler

import (
	"context"
	"errors"
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mike7109/tg-bot-clubbing/internal/apperrors"
	"github.com/mike7109/tg-bot-clubbing/internal/service"
	schema "github.com/mike7109/tg-bot-clubbing/internal/service/dto"
	"github.com/mike7109/tg-bot-clubbing/internal/transport/worker/dto"
	"github.com/mike7109/tg-bot-clubbing/internal/transport/worker/update_entity/button"
	"github.com/mike7109/tg-bot-clubbing/pkg/messages"
)

type MessageHandlerFunc func(ctx context.Context, msg *tgApi.Message) error

func Start(ctx context.Context, tgBotApi *tgApi.BotAPI, tgBotService service.ITgBotService) MessageHandlerFunc {
	return func(ctx context.Context, msgUpdate *tgApi.Message) error {

		msg := tgBotService.StartHandler()

		_, err := tgBotApi.Send(tgApi.NewMessage(msgUpdate.Chat.ID, msg))
		if err != nil {
			return err
		}

		return nil
	}
}

func Help(ctx context.Context, tgBot *tgApi.BotAPI, tgBotService service.ITgBotService) MessageHandlerFunc {
	return func(ctx context.Context, msgUpdate *tgApi.Message) error {
		msg := tgBotService.HelpHandler()

		_, err := tgBot.Send(tgApi.NewMessage(msgUpdate.Chat.ID, msg))
		if err != nil {
			return err
		}

		return nil
	}
}

func Save(ctx context.Context, tgBot *tgApi.BotAPI, tgBotService service.ITgBotService) MessageHandlerFunc {
	return func(ctx context.Context, msgUpdate *tgApi.Message) error {
		chatID := msgUpdate.Chat.ID

		page, err := dto.ParseMessageForPage(msgUpdate)
		if err != nil {
			switch {
			case errors.Is(err, apperrors.ErrNoURL):
				return messages.SendInvalidUrlMessage(tgBot, chatID)
			case errors.Is(err, apperrors.ErrNoUserName):
				return messages.SendErrNoUserName(tgBot, chatID)
			case errors.Is(err, apperrors.ErrInvalidURL):
				return messages.SendInvalidUrlMessage(tgBot, chatID)
			}

			return err
		}

		msg, err := tgBotService.SaveHandler(ctx, page)
		if err != nil {
			if errors.Is(err, apperrors.ErrNoSave) {
				return messages.SendErrorHandler(tgBot, chatID)
			}
			return err
		}

		return messages.SendMessage(tgBot, msgUpdate.Chat.ID, msg)
	}
}

func ListUrl(ctx context.Context, tgBot *tgApi.BotAPI, tgBotService service.ITgBotService) MessageHandlerFunc {
	return func(ctx context.Context, msgUpdate *tgApi.Message) error {
		userName := msgUpdate.From.UserName
		chatID := msgUpdate.Chat.ID

		pages, err := tgBotService.GetPageHandler(ctx, userName, 0, 10)
		if err != nil {
			switch {
			case errors.Is(err, apperrors.ErrNoUrl):
				return messages.SendNoSavedPagesMessage(tgBot, chatID)
			case errors.Is(err, apperrors.ErrNoSave):
				return messages.SendErrorHandler(tgBot, chatID)
			default:
				return messages.SendErrorHandler(tgBot, chatID)
			}
		}

		msg := createListMsg(pages, msgUpdate.Chat.ID)

		_, err = tgBot.Send(msg)

		return err
	}
}

func createListMsg(listPage *schema.ListPage, chatID int64) tgApi.Chattable {
	var msg string

	builder := button.NewBuilder()

	if listPage.HaveNextPage {
		butNext := button.NewButton(">", button.ListCommand)
		button.SetDataValue(butNext, "p", 1)
		button.SetDataValue(butNext, "d", 0)
		button.SetDataValue(butNext, "c", button.SwitchPageCommandButton)

		butEnd := button.NewButton(">>", button.ListCommand)
		button.SetDataValue(butEnd, "p", listPage.LastPage)
		button.SetDataValue(butEnd, "d", 0)
		button.SetDataValue(butEnd, "c", button.SwitchPageCommandButton)
		builder.AddButtonTopRows(butNext, butEnd)
	}

	but := button.NewButton("Удалить по номерам", button.ListCommand)
	button.SetDataValue(but, "p", 0)
	button.SetDataValue(but, "c", button.WantToDeleteURLCommandButton)
	button.SetDataValue(but, "d", 0)
	builder.AddButton(but)

	for _, page := range listPage.SavePage {
		msg += page.String()
	}

	keyboard := builder.Build()

	msgConfig := tgApi.NewMessage(chatID, msg)
	msgConfig.DisableWebPagePreview = true
	msgConfig.ReplyMarkup = keyboard

	return msgConfig
}