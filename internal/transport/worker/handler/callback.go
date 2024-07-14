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

		page, exist := button.GetDataValue(buttonTarget, "p")
		if !exist {
			return fmt.Errorf("page not found in data")
		}

		pages, err := tgBotService.GetPageHandler(ctx, userName, int(page.(float64)), 10)
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

		if countPage == 0 {
			return messages.SendNoSavedPagesMessage(tgBot, chatID)
		}

		msg := editListMsg(pages, chatID, messageID, int(page.(float64)), countPage, 1)

		_, err = tgBot.Send(msg)

		return err
	}
}

func WantToDeleteUrl(ctx context.Context, tgBot *tgApi.BotAPI, tgBotService service.ITgBotService) CallbackHandlerFunc {
	return func(ctx context.Context, callbackUpdate *tgApi.CallbackQuery, buttonTarget *button.Button) error {
		callback := tgApi.NewCallback(callbackUpdate.ID, callbackUpdate.Data)
		if _, err := tgBot.Request(callback); err != nil {
			return err
		}

		userName := callbackUpdate.Message.Chat.UserName
		chatID := callbackUpdate.Message.Chat.ID
		messageID := callbackUpdate.Message.MessageID

		page, exist := button.GetDataValue(buttonTarget, "p")
		if !exist {
			return fmt.Errorf("page not found in data")
		}

		pages, err := tgBotService.GetPageHandler(ctx, userName, int(page.(float64)), 10)
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

		if countPage == 0 {
			return messages.SendNoSavedPagesMessage(tgBot, chatID)
		}

		msg := editListMsg(pages, chatID, messageID, int(page.(float64)), countPage, 1)

		_, err = tgBot.Send(msg)

		return err
	}
}

func CancelWantToDeleteUrl(ctx context.Context, tgBot *tgApi.BotAPI, tgBotService service.ITgBotService) CallbackHandlerFunc {
	return func(ctx context.Context, callbackUpdate *tgApi.CallbackQuery, buttonTarget *button.Button) error {
		callback := tgApi.NewCallback(callbackUpdate.ID, callbackUpdate.Data)
		if _, err := tgBot.Request(callback); err != nil {
			return err
		}

		userName := callbackUpdate.Message.Chat.UserName
		chatID := callbackUpdate.Message.Chat.ID
		messageID := callbackUpdate.Message.MessageID

		page, exist := button.GetDataValue(buttonTarget, "p")
		if !exist {
			return fmt.Errorf("page not found in data")
		}

		pages, err := tgBotService.GetPageHandler(ctx, userName, int(page.(float64)), 10)
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

		if countPage == 0 {
			return messages.SendNoSavedPagesMessage(tgBot, chatID)
		}

		msg := cancelEditListMsgForNotPAge(pages, chatID, messageID, int(page.(float64)), countPage)

		_, err = tgBot.Send(msg)

		return err
	}
}

func NextPageCommand(ctx context.Context, tgBot *tgApi.BotAPI, tgBotService service.ITgBotService) CallbackHandlerFunc {
	return func(ctx context.Context, callbackUpdate *tgApi.CallbackQuery, buttonTarget *button.Button) error {
		callback := tgApi.NewCallback(callbackUpdate.ID, callbackUpdate.Data)
		if _, err := tgBot.Request(callback); err != nil {
			return err
		}

		userName := callbackUpdate.Message.Chat.UserName
		chatID := callbackUpdate.Message.Chat.ID
		messageID := callbackUpdate.Message.MessageID

		page, exist := button.GetDataValue(buttonTarget, "p")
		if !exist {
			return fmt.Errorf("page not found in data")
		}

		wantToDelete, exist := button.GetDataValue(buttonTarget, "d")
		if !exist {
			return fmt.Errorf("wantToDelete not found in data")
		}

		pages, err := tgBotService.GetPageHandler(ctx, userName, int(page.(float64)), 10)
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

		if countPage == 0 {
			return messages.SendNoSavedPagesMessage(tgBot, chatID)
		}

		var msg tgApi.Chattable

		if int(wantToDelete.(float64)) == 1 {
			msg = editListMsg(pages, chatID, messageID, int(page.(float64)), countPage, int(wantToDelete.(float64)))
		} else {
			msg = creatListMsg(pages, chatID, messageID, int(page.(float64)), countPage)
		}

		_, err = tgBot.Send(msg)

		return err
	}
}

func editListMsg(pages []*entity.Page, chatID int64, MessageID int, numPage int, countPage int, wantToDelete int) tgApi.Chattable {
	var msg string

	builder := button.NewBuilder()

	if numPage > 0 {
		butPrev := button.NewButton("<", button.NextPageCommand)
		button.SetDataValue(butPrev, "p", numPage-1)
		button.SetDataValue(butPrev, "d", wantToDelete)
		butFirst := button.NewButton("<<", button.NextPageCommand)
		button.SetDataValue(butFirst, "p", 0)
		button.SetDataValue(butFirst, "d", wantToDelete)
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
		butNext := button.NewButton(">", button.NextPageCommand)
		button.SetDataValue(butNext, "p", numPage+1)
		button.SetDataValue(butNext, "d", wantToDelete)
		butEnd := button.NewButton(">>", button.NextPageCommand)
		button.SetDataValue(butEnd, "p", lastPage)
		button.SetDataValue(butEnd, "d", wantToDelete)
		builder.AddButtonTopRows(butNext, butEnd)
	}

	for _, page := range pages {
		msg += page.String()
		but := page.ToButton(button.DeleteCommand)
		button.SetDataValue(but, "p", numPage)
		builder.AddButton(but)
	}

	wantToDeleteCommand := button.NewButton("Назад", button.CancelWantToDeleteCommand)
	button.SetDataValue(wantToDeleteCommand, "p", numPage)
	button.SetDataValue(wantToDeleteCommand, "d", 0)

	builder.AddButtonBottomRow(wantToDeleteCommand)

	keyboard := builder.Build()

	msgConfig := tgApi.NewEditMessageText(chatID, MessageID, msg)
	msgConfig.DisableWebPagePreview = true
	msgConfig.ReplyMarkup = &keyboard

	return msgConfig
}

func creatListMsg(pages []*entity.Page, chatID int64, MessageID int, numPage int, countPage int) tgApi.Chattable {
	var msg string

	builder := button.NewBuilder()

	if numPage > 0 {
		butPrev := button.NewButton("<", button.NextPageCommand)
		button.SetDataValue(butPrev, "p", numPage-1)
		button.SetDataValue(butPrev, "d", 0)
		butFirst := button.NewButton("<<", button.NextPageCommand)
		button.SetDataValue(butFirst, "p", 0)
		button.SetDataValue(butFirst, "d", 0)
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
		butNext := button.NewButton(">", button.NextPageCommand)
		button.SetDataValue(butNext, "p", numPage+1)
		button.SetDataValue(butNext, "d", 0)
		butEnd := button.NewButton(">>", button.NextPageCommand)
		button.SetDataValue(butEnd, "p", lastPage)
		button.SetDataValue(butEnd, "d", 0)

		builder.AddButtonTopRows(butNext, butEnd)
	}

	for _, page := range pages {
		msg += page.String()
	}

	wantToDeleteCommand := button.NewButton("Удалить по номерам", button.WantToDeleteCommand)
	button.SetDataValue(wantToDeleteCommand, "p", numPage)

	builder.AddButtonBottomRow(wantToDeleteCommand)

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

func cancelEditListMsgForNotPAge(pages []*entity.Page, chatID int64, MessageID int, numPage int, countPage int) tgApi.Chattable {
	var msg string

	builder := button.NewBuilder()

	if numPage > 0 {
		butPrev := button.NewButton("<", button.NextPageCommand)
		button.SetDataValue(butPrev, "p", numPage-1)
		button.SetDataValue(butPrev, "d", 0)
		butFirst := button.NewButton("<<", button.NextPageCommand)
		button.SetDataValue(butFirst, "p", 0)
		button.SetDataValue(butFirst, "d", 0)
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
		butNext := button.NewButton(">", button.NextPageCommand)
		button.SetDataValue(butNext, "p", numPage+1)
		button.SetDataValue(butNext, "d", 0)
		butEnd := button.NewButton(">>", button.NextPageCommand)
		button.SetDataValue(butEnd, "p", lastPage)
		button.SetDataValue(butEnd, "d", 0)

		builder.AddButtonTopRows(butNext, butEnd)
	}

	for _, page := range pages {
		msg += page.String()
	}

	wantToDeleteCommand := button.NewButton("Удалить по номерам", button.WantToDeleteCommand)
	button.SetDataValue(wantToDeleteCommand, "p", numPage)

	builder.AddButtonBottomRow(wantToDeleteCommand)

	keyboard := builder.Build()

	msgConfig := tgApi.NewEditMessageText(chatID, MessageID, msg)
	msgConfig.DisableWebPagePreview = true
	msgConfig.ReplyMarkup = &keyboard

	return msgConfig
}
