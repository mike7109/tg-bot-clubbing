package service

import (
	"context"
	"errors"
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mike7109/tg-bot-clubbing/internal/entity"
	"log"
)

type ICommand interface {
	Execute() (tgApi.MessageConfig, error)
}

type SavePageCommand struct {
	storage  IStorage
	url      string
	userName string
	chatID   int64
}

func NewSavePageCommand(storage IStorage, url string, userName string, chatID int64) *SavePageCommand {
	return &SavePageCommand{
		storage:  storage,
		url:      url,
		userName: userName,
		chatID:   chatID,
	}
}

func (c SavePageCommand) Execute() (tgApi.MessageConfig, error) {
	page := &entity.Page{
		URL:      c.url,
		UserName: c.userName,
	}

	isExists, err := c.storage.IsExists(context.Background(), page)
	if err != nil {
		return tgApi.MessageConfig{}, err
	}
	if isExists {
		return tgApi.NewMessage(c.chatID, msgAlreadyExists), nil
	}

	if err := c.storage.Save(context.Background(), page); err != nil {
		return tgApi.MessageConfig{}, err
	}

	msg := tgApi.NewMessage(c.chatID, msgSaved)

	return msg, nil
}

type RndCommand struct {
	storage  IStorage
	userName string
	chatID   int64
}

func NewRndCommand(storage IStorage, userName string, chatID int64) *RndCommand {
	return &RndCommand{
		storage:  storage,
		userName: userName,
		chatID:   chatID,
	}
}

func (c RndCommand) Execute() (tgApi.MessageConfig, error) {
	var msg tgApi.MessageConfig
	page, err := c.storage.PickRandom(context.Background(), c.userName)
	if err != nil && !errors.Is(err, entity.ErrNoSavedPages) {
		return msg, err
	}
	if errors.Is(err, entity.ErrNoSavedPages) {
		return tgApi.NewMessage(c.chatID, msgNoSavedPages), nil
	}

	msg = tgApi.NewMessage(c.chatID, page.URL)

	if err = c.storage.Remove(context.Background(), page); err != nil {
		log.Println("Failed to remove page: ", err)
		return msg, nil
	}

	return msg, nil
}

type HelpCommand struct {
	update tgApi.Update
}

func (c HelpCommand) Execute() (tgApi.MessageConfig, error) {
	msg := tgApi.NewMessage(c.update.Message.Chat.ID, msgHelp)

	return msg, nil

}

type StartCommand struct {
	update tgApi.Update
}

func (c StartCommand) Execute() (tgApi.MessageConfig, error) {
	msg := tgApi.NewMessage(c.update.Message.Chat.ID, msgHello)

	return msg, nil
}

type UnknownCommand struct {
	update tgApi.Update
}

func (c UnknownCommand) Execute() (tgApi.MessageConfig, error) {
	msg := tgApi.NewMessage(c.update.Message.Chat.ID, msgUnknownCommand)

	return msg, nil

}
