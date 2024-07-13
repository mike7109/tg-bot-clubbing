package service

import (
	"context"
	"fmt"
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mike7109/tg-bot-clubbing/internal/repositories"
	"github.com/mike7109/tg-bot-clubbing/internal/service/processor"
	"github.com/mike7109/tg-bot-clubbing/internal/service/processor/task"
	"github.com/mike7109/tg-bot-clubbing/pkg/commands"
	"sync"
)

type TgProcessor struct {
	chDone chan struct{}
}

func NewTgProcessor(ctx context.Context, tgBot *tgApi.BotAPI, storage *repositories.Storage) (*TgProcessor, error) {
	p := &TgProcessor{
		chDone: make(chan struct{}),
	}
	err := p.newTgProcessor(ctx, tgBot, storage)
	if err != nil {
		return nil, fmt.Errorf("failed to create new tg processor: %w", err)
	}

	return p, nil
}

func (p *TgProcessor) newTgProcessor(ctx context.Context, tgBot *tgApi.BotAPI, storage *repositories.Storage) error {
	u := tgApi.NewUpdate(0)
	u.Timeout = 60

	taskProcessor := processor.NewTaskProcessor(tgBot)

	// AddCmd task processor
	taskProcessor.AddTaskProcessor(commands.StartCmd, task.Start(ctx, tgBot))
	taskProcessor.AddTaskProcessor(commands.HelpCmd, task.Help(ctx, tgBot))
	taskProcessor.AddTaskProcessor(commands.RndCmd, task.Rnd(ctx, tgBot, storage))
	taskProcessor.AddTaskProcessor(commands.AddCmd, task.Save(ctx, tgBot, storage))
	taskProcessor.AddTaskProcessor(commands.AddSimpleCmd, task.SaveSimple(ctx, tgBot, storage))
	taskProcessor.AddTaskProcessor(commands.ListUrl, task.ListUrl(ctx, tgBot, storage))

	var wg sync.WaitGroup

	consume := func(updates <-chan tgApi.Update) {
		for update := range updates {
			wg.Add(1)
			go func() {
				defer wg.Done()
				taskProcessor.Consume(ctx, update)
			}()
		}
	}

	go consume(tgBot.GetUpdatesChan(u))
	fmt.Println("Worker started")

	go func() {
		select {
		case <-ctx.Done():
			tgBot.StopReceivingUpdates()
			wg.Wait()
			close(p.chDone)
			fmt.Println("Worker stopped")
		}
	}()

	return nil
}

func (p *TgProcessor) WaitClose() {
	<-p.chDone
}
