package service

import (
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/url"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

type IFactoryCommand interface {
	CreateCommand(update tgApi.Update) ICommand
}

type FactoryCommand struct {
	storage IStorage
}

func NewFactoryCommand(storage IStorage) FactoryCommand {
	return FactoryCommand{storage}

}

func (f FactoryCommand) CreateCommand(update tgApi.Update) ICommand {
	if isAddCmd(update.Message.Text) {
		return NewSavePageCommand(f.storage, update.Message.Text, update.Message.From.UserName, update.Message.Chat.ID)
	}

	switch update.Message.Text {
	case StartCmd:
		return StartCommand{update}

	case HelpCmd:
		return HelpCommand{update}
	case RndCmd:
		return NewRndCommand(f.storage, update.Message.From.UserName, update.Message.Chat.ID)
	default:
		return UnknownCommand{update}
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
