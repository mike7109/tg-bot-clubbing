package service

import (
	"context"
	"fmt"
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Process struct {
	tgBot    *tgApi.BotAPI
	fCommand IFactoryCommand
}

func NewProcess(tgBot *tgApi.BotAPI, factoryCommand IFactoryCommand) *Process {
	return &Process{
		tgBot:    tgBot,
		fCommand: factoryCommand,
	}
}

func (p Process) process(ctx context.Context, update tgApi.Update) error {
	if update.Message == nil {
		return nil
	}

	msg, err := p.fCommand.CreateCommand(update).Execute()
	if err != nil {
		return err
	}

	_, err = p.tgBot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (p Process) Start() error {
	u := tgApi.NewUpdate(0)
	u.Timeout = 60

	ctx := context.TODO()

	for update := range p.tgBot.GetUpdatesChan(u) {
		if err := p.process(ctx, update); err != nil {
			return fmt.Errorf("failed to process update: %v", err)
		}
	}

	return nil
}
