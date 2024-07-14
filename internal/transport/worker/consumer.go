package worker

import (
	"context"
	"fmt"
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mike7109/tg-bot-clubbing/internal/service"
	"github.com/mike7109/tg-bot-clubbing/internal/transport/worker/handler"
	"github.com/mike7109/tg-bot-clubbing/internal/transport/worker/update_entity/button"
	"github.com/mike7109/tg-bot-clubbing/internal/transport/worker/update_entity/message"

	"sync"
)

type TgProcessor struct {
	chDone chan struct{}
}

func NewTgProcessor(ctx context.Context, tgBot *tgApi.BotAPI, tgBotService service.ITgBotService) (*TgProcessor, error) {
	p := &TgProcessor{
		chDone: make(chan struct{}),
	}
	err := p.newTgProcessor(ctx, tgBot, tgBotService)
	if err != nil {
		return nil, fmt.Errorf("failed to create new tg processor: %w", err)
	}

	return p, nil
}

func (p *TgProcessor) newTgProcessor(ctx context.Context, tgBot *tgApi.BotAPI, tgBotService service.ITgBotService) error {
	u := tgApi.NewUpdate(0)
	u.Timeout = 60

	taskProcessor := NewTaskProcessor(tgBot)

	// AddCmd task processor
	taskProcessor.AddHandlerMessages(message.StartCommand, handler.Start(ctx, tgBot, tgBotService))
	taskProcessor.AddHandlerMessages(message.HelpCommand, handler.Help(ctx, tgBot, tgBotService))
	taskProcessor.AddHandlerMessages(message.AddCommand, handler.Save(ctx, tgBot, tgBotService))
	taskProcessor.AddHandlerMessages(message.ListCommand, handler.ListUrl(ctx, tgBot, tgBotService))

	// AddCallback task processor
	taskProcessor.AddHandlerCallback(button.ListCommand, handler.ListCallback(ctx, tgBot, tgBotService))

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
