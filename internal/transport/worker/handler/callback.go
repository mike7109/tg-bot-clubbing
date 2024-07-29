package handler

import (
	"context"
	"errors"
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mike7109/tg-bot-clubbing/internal/apperrors"
	"github.com/mike7109/tg-bot-clubbing/internal/entity"
	"github.com/mike7109/tg-bot-clubbing/internal/service"
	"github.com/mike7109/tg-bot-clubbing/internal/service/dto"
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
			case errors.Is(err, apperrors.ErrNoUrl):
				msgZeroPage := editListMsgForNotPage(chatID, messageID)
				_, err = tgBot.Send(msgZeroPage)
				return err
			case errors.Is(err, apperrors.ErrNoSave):
				return messages.SendErrorHandler(tgBot, chatID)
			default:
				return messages.SendErrorHandler(tgBot, chatID)
			}
		}

		msg := CreateListPages(pages, chatID, messageID, WithDeleteButton)

		_, err = tgBot.Send(msg)

		return nil
	}
}

func CreateListPages(listPage *dto.ListPage, chatID int64, MessageID int, withDelete int) tgApi.Chattable {
	var msg string
	builder := button.NewBuilder()

	addNavigationButtons := func() {
		if listPage.HavePrevPage {
			builder.AddButtonTopRows(
				createNavButton("<<", 0, withDelete),
				createNavButton("<", listPage.NumPage-1, withDelete),
			)
		}

		if listPage.HaveNextPage {
			builder.AddButtonTopRows(
				createNavButton(">", listPage.NumPage+1, withDelete),
				createNavButton(">>", listPage.LastPage, withDelete),
			)
		}
	}

	addPageButtons := func() {
		for _, page := range listPage.SavePage {
			msg += page.String()
			if withDelete == 1 {
				builder.AddButton(createDeleteButton(page, listPage.NumPage, withDelete))
			}
		}
	}

	addBottomButton := func() {
		var bottomButton *button.Button
		if withDelete == 0 {
			bottomButton = createActionButton("Удалить по номерам", listPage.NumPage, 1, button.WantToDeleteURLCommandButton)
		} else {
			bottomButton = createActionButton("Назад", listPage.NumPage, 0, button.CancelWantToDeleteURLCommandButton)
		}
		builder.AddButtonBottomRow(bottomButton)
	}

	addNavigationButtons()
	addPageButtons()
	addBottomButton()

	keyboard := builder.Build()
	msgConfig := tgApi.NewEditMessageText(chatID, MessageID, msg)
	msgConfig.DisableWebPagePreview = true
	msgConfig.ReplyMarkup = &keyboard

	return msgConfig
}

func createNavButton(text string, page, withDelete int) *button.Button {
	btn := button.NewButton(text, button.ListCommand)
	button.SetDataValue(btn, "p", page)
	button.SetDataValue(btn, "d", withDelete)
	button.SetDataValue(btn, "c", button.SwitchPageCommandButton)
	return btn
}

func createDeleteButton(page *entity.UrlPage, numPage, withDelete int) *button.Button {
	btn := page.ToButton(button.ListCommand)
	button.SetDataValue(btn, "c", button.DeleteURLCommandButton)
	button.SetDataValue(btn, "p", numPage)
	button.SetDataValue(btn, "d", withDelete)
	return btn
}

func createActionButton(text string, numPage, withDelete int, commandButton button.CommandButton) *button.Button {
	btn := button.NewButton(text, button.ListCommand)
	button.SetDataValue(btn, "p", numPage)
	button.SetDataValue(btn, "c", commandButton)
	button.SetDataValue(btn, "d", withDelete)
	return btn
}

func editListMsgForNotPage(chatID int64, MessageID int) tgApi.Chattable {
	msgConfig := tgApi.NewEditMessageText(chatID, MessageID, messages.MsgNoSavedPages)
	msgConfig.DisableWebPagePreview = true
	msgConfig.ReplyMarkup = nil

	return msgConfig
}
