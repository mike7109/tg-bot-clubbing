package task

import (
	"context"
	"errors"
	"fmt"
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mike7109/tg-bot-clubbing/internal/entity"
	"github.com/mike7109/tg-bot-clubbing/internal/repositories"
	"github.com/mike7109/tg-bot-clubbing/internal/service/processor"
	"github.com/mike7109/tg-bot-clubbing/pkg/messages"
	"github.com/mike7109/tg-bot-clubbing/pkg/utls"
	"log"
	"regexp"
)

func Start(ctx context.Context, tgBot *tgApi.BotAPI) processor.ProcessingFunc {
	return func(ctx context.Context, update tgApi.Update, msg *tgApi.Message) error {
		msgConfig := tgApi.NewMessage(msg.Chat.ID, messages.MsgHello)
		_, err := tgBot.Send(msgConfig)
		if err != nil {
			return err
		}

		return nil
	}
}

func Help(ctx context.Context, tgBot *tgApi.BotAPI) processor.ProcessingFunc {
	return func(ctx context.Context, update tgApi.Update, msg *tgApi.Message) error {
		msgConfig := tgApi.NewMessage(msg.Chat.ID, messages.MsgHelp)
		_, err := tgBot.Send(msgConfig)
		if err != nil {
			return err
		}

		return nil
	}
}

func Rnd(ctx context.Context, tgBot *tgApi.BotAPI, storage *repositories.Storage) processor.ProcessingFunc {
	return func(ctx context.Context, update tgApi.Update, msg *tgApi.Message) error {
		page, err := storage.PickRandom(ctx, msg.From.UserName)
		if err != nil && !errors.Is(err, entity.ErrNoSavedPages) {
			return err
		}

		if errors.Is(err, entity.ErrNoSavedPages) {
			msgConfig := tgApi.NewMessage(msg.Chat.ID, messages.MsgNoSavedPages)
			_, err = tgBot.Send(msgConfig)
			if err != nil {
				return err
			}

			return nil
		}

		msgConfig := tgApi.NewMessage(msg.Chat.ID, page.URL)
		_, err = tgBot.Send(msgConfig)
		if err != nil {
			return err
		}

		if err = storage.Remove(context.Background(), page); err != nil {
			log.Println("Failed to remove page: ", err)
			return nil
		}

		return nil
	}
}

func Save(ctx context.Context, tgBot *tgApi.BotAPI, storage *repositories.Storage) processor.ProcessingFunc {
	return func(ctx context.Context, update tgApi.Update, msg *tgApi.Message) error {

		re := regexp.MustCompile(`^/add\s+(\S+)(?:\s+(.+?))?(?:\s+(.+))?(?:\s+(.+?))?$`)

		matches := re.FindStringSubmatch(msg.Text)

		if len(matches) == 0 {
			msgConfig := tgApi.NewMessage(msg.Chat.ID, messages.MsgInvalidAddCommand)
			_, err := tgBot.Send(msgConfig)
			if err != nil {
				return err
			}

			return nil
		}

		url := matches[1]
		if !utls.IsURL(url) {
			msgConfig := tgApi.NewMessage(msg.Chat.ID, messages.MsgInvalidUrl)
			_, err := tgBot.Send(msgConfig)
			if err != nil {
				return err
			}

			return nil
		}

		var description, title, category *string

		if len(matches) > 2 {
			category = &matches[2]
		}
		if len(matches) > 3 {
			title = &matches[3]
		}
		if len(matches) > 4 {
			description = &matches[4]
		}

		page := &entity.Page{
			UserName:    msg.From.UserName,
			URL:         url,
			Title:       title,
			Category:    category,
			Description: description,
		}

		isExists, err := storage.IsExists(ctx, page)
		if err != nil {
			return err
		}
		if isExists {
			msgConfig := tgApi.NewMessage(msg.Chat.ID, messages.MsgAlreadyExists)
			_, err = tgBot.Send(msgConfig)
			if err != nil {
				return err
			}

			return nil
		}

		if err := storage.Save(ctx, page); err != nil {
			return err
		}

		msgConfig := tgApi.NewMessage(msg.Chat.ID, messages.MsgSaved)
		_, err = tgBot.Send(msgConfig)
		if err != nil {
			return err
		}

		return nil
	}
}

func SaveSimple(ctx context.Context, tgBot *tgApi.BotAPI, storage *repositories.Storage) processor.ProcessingFunc {
	return func(ctx context.Context, update tgApi.Update, msg *tgApi.Message) error {
		page := &entity.Page{
			UserName: msg.From.UserName,
			URL:      msg.Text,
		}

		isExists, err := storage.IsExists(ctx, page)
		if err != nil {
			return err
		}
		if isExists {
			msgConfig := tgApi.NewMessage(msg.Chat.ID, messages.MsgAlreadyExists)
			_, err = tgBot.Send(msgConfig)
			if err != nil {
				return err
			}

			return nil
		}

		if err := storage.Save(ctx, page); err != nil {
			return err
		}

		msgConfig := tgApi.NewMessage(msg.Chat.ID, messages.MsgSaved)
		_, err = tgBot.Send(msgConfig)
		if err != nil {
			return err
		}

		return nil
	}
}

func ListUrl(ctx context.Context, tgBot *tgApi.BotAPI, storage *repositories.Storage) processor.ProcessingFunc {
	return func(ctx context.Context, update tgApi.Update, msg *tgApi.Message) error {
		pages, err := storage.ListUrl(ctx, msg.From.UserName)
		if err != nil && !errors.Is(err, entity.ErrNoSavedPages) {
			return err
		}

		if errors.Is(err, entity.ErrNoSavedPages) {
			msgConfig := tgApi.NewMessage(msg.Chat.ID, messages.MsgNoSavedPages)
			_, err = tgBot.Send(msgConfig)
			if err != nil {
				return err
			}

			return nil
		}

		var urlList string
		for i, page := range pages {
			urlList += fmt.Sprintf("%d. %s ", i+1, page.URL)
			if page.Category != nil {
				urlList += fmt.Sprintf("%s ", *page.Category)
			}
			if page.Title != nil {
				urlList += fmt.Sprintf("%s ", *page.Title)
			}
			if page.Description != nil {
				urlList += fmt.Sprintf("%s ", *page.Description)
			}
			urlList += "\n"
		}

		msgConfig := tgApi.NewMessage(msg.Chat.ID, urlList)
		msgConfig.DisableWebPagePreview = true // Отключаем веб-превью
		_, err = tgBot.Send(msgConfig)
		if err != nil {
			return err
		}

		return nil
	}
}

func DeleteAll(ctx context.Context, tgBot *tgApi.BotAPI, storage *repositories.Storage) processor.ProcessingFunc {
	return func(ctx context.Context, update tgApi.Update, msg *tgApi.Message) error {
		err := storage.DeleteAll(ctx, msg.From.UserName)
		if err != nil {
			return err
		}

		msgConfig := tgApi.NewMessage(msg.Chat.ID, messages.MsgDeletedAll)
		_, err = tgBot.Send(msgConfig)
		if err != nil {
			return err
		}

		return nil
	}
}
