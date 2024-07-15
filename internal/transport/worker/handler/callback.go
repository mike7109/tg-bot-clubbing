package handler

import (
	"context"
	"errors"
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mike7109/tg-bot-clubbing/internal/apperrors"
	"github.com/mike7109/tg-bot-clubbing/internal/entity"
	"github.com/mike7109/tg-bot-clubbing/internal/service"
	"github.com/mike7109/tg-bot-clubbing/internal/transport/worker/update_entity/button"
	"github.com/mike7109/tg-bot-clubbing/pkg/messages"
)

type CallbackHandlerFunc func(ctx context.Context, callback *tgApi.CallbackQuery, buttonTarget *button.Button) error

func ListCallback(ctx context.Context, tgBot *tgApi.BotAPI, tgBotService service.ITgBotService) CallbackHandlerFunc {
	return func(ctx context.Context, callbackUpdate *tgApi.CallbackQuery, buttonTarget *button.Button) error {
		callback := tgApi.NewCallback(callbackUpdate.ID, callbackUpdate.Data)
		if _, err := tgBot.Request(callback); err != nil {
			return err
		}

		userName := callbackUpdate.Message.Chat.UserName
		chatID := callbackUpdate.Message.Chat.ID
		messageID := callbackUpdate.Message.MessageID

		listButton, err := buttonTarget.ToListButton()
		if err != nil {
			return err
		}

		WithDeleteButton := listButton.WithDelete

		switch listButton.Cmd {
		case button.WantToDeleteURLCommandButton:
			WithDeleteButton = 1
		case button.CancelWantToDeleteURLCommandButton:
			WithDeleteButton = 0
		case button.SwitchPageCommandButton:
		case button.DeleteURLCommandButton:
			err = tgBotService.DeleteHandler(ctx, listButton.ID, userName)
			if err != nil {
				return err
			}
		}

		pages, err := tgBotService.GetPageHandler(ctx, userName, listButton.CurrentPage, 10)
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

		countPage, err := tgBotService.CountHandler(ctx, userName)
		if err != nil {
			return err
		}

		msg := CreateListPages(pages, chatID, messageID, listButton.CurrentPage, countPage, WithDeleteButton)

		_, err = tgBot.Send(msg)

		return nil
	}
}

func CreateListPages(pages []*entity.UrlPage, chatID int64, MessageID int, numPage int, countPage int, withDelete int) tgApi.Chattable {
	var msg string

	builder := button.NewBuilder()

	if numPage > 0 {
		butPrev := button.NewButton("<", button.ListCommand)
		button.SetDataValue(butPrev, "p", numPage-1)
		button.SetDataValue(butPrev, "d", withDelete)
		button.SetDataValue(butPrev, "c", button.SwitchPageCommandButton)
		butFirst := button.NewButton("<<", button.ListCommand)
		button.SetDataValue(butFirst, "p", 0)
		button.SetDataValue(butFirst, "d", withDelete)
		button.SetDataValue(butFirst, "c", button.SwitchPageCommandButton)
		builder.AddButtonTopRows(butFirst, butPrev)
	}

	var lastPage int

	coinPageDiv := countPage % 10
	if coinPageDiv == 0 {
		lastPage = (countPage - 1) / 10
	} else {
		lastPage = countPage / 10
	}

	if countPage > (numPage+1)*10 {
		butNext := button.NewButton(">", button.ListCommand)
		button.SetDataValue(butNext, "p", numPage+1)
		button.SetDataValue(butNext, "d", withDelete)
		button.SetDataValue(butNext, "c", button.SwitchPageCommandButton)
		butEnd := button.NewButton(">>", button.ListCommand)
		button.SetDataValue(butEnd, "p", lastPage)
		button.SetDataValue(butEnd, "d", withDelete)
		button.SetDataValue(butEnd, "c", button.SwitchPageCommandButton)

		builder.AddButtonTopRows(butNext, butEnd)
	}

	for _, page := range pages {
		msg += page.String()
		if withDelete == 1 {
			but := page.ToButton(button.ListCommand)
			button.SetDataValue(but, "c", button.DeleteURLCommandButton)
			button.SetDataValue(but, "p", numPage)
			button.SetDataValue(but, "d", withDelete)
			builder.AddButton(but)
		}
	}

	if withDelete == 0 {
		wantToDeleteCommand := button.NewButton("Удалить по номерам", button.ListCommand)
		button.SetDataValue(wantToDeleteCommand, "p", numPage)
		button.SetDataValue(wantToDeleteCommand, "c", button.WantToDeleteURLCommandButton)
		button.SetDataValue(wantToDeleteCommand, "d", 1)
		builder.AddButtonBottomRow(wantToDeleteCommand)
	}

	if withDelete == 1 {
		wantToDeleteCommand := button.NewButton("Назад", button.ListCommand)
		button.SetDataValue(wantToDeleteCommand, "p", numPage)
		button.SetDataValue(wantToDeleteCommand, "c", button.CancelWantToDeleteURLCommandButton)
		button.SetDataValue(wantToDeleteCommand, "d", 0)
		builder.AddButtonBottomRow(wantToDeleteCommand)
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
