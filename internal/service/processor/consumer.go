package processor

import (
	"context"
	"fmt"
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"net/url"
	"strings"
)

const (
	ErrorHandler    = "Произошла ошибка"
	NotFoundHandler = "Я не знаю такой команды"
)

type ProcessingFunc func(ctx context.Context, update tgApi.Update, msg *tgApi.Message) error

type TaskProcessor struct {
	routes                   map[string]ProcessingFunc
	NoErrorHandler           func(ctx context.Context, update tgApi.Update)
	ErrorHandler             func(ctx context.Context, err error, update tgApi.Update)
	ProcessorNotFoundHandler func(ctx context.Context, update tgApi.Update)
}

func NewTaskProcessor(tgBot *tgApi.BotAPI) TaskProcessor {
	defaultNoErrorHandler := func(ctx context.Context, update tgApi.Update) {
		fmt.Println("Processed update, acking")
	}

	defaultErrorHandler := func(ctx context.Context, err error, update tgApi.Update) {
		fmt.Println("Failed to process update, rejecting", zap.String("Text", update.Message.Text), err)
		_, _ = tgBot.Send(tgApi.NewMessage(update.Message.Chat.ID, ErrorHandler))
	}

	defaultProcessorNotFoundHandler := func(ctx context.Context, updates tgApi.Update) {
		fmt.Println("No processor found for update type", zap.String("Text", updates.Message.Text))
		_, _ = tgBot.Send(tgApi.NewMessage(updates.Message.Chat.ID, NotFoundHandler))
	}

	return TaskProcessor{
		routes:                   make(map[string]ProcessingFunc),
		NoErrorHandler:           defaultNoErrorHandler,
		ErrorHandler:             defaultErrorHandler,
		ProcessorNotFoundHandler: defaultProcessorNotFoundHandler,
	}
}

func (c TaskProcessor) AddTaskProcessor(routingKey string, fn ProcessingFunc) {
	c.routes[routingKey] = fn
}

func (c TaskProcessor) Consume(ctx context.Context, update tgApi.Update) {
	c.Process(ctx, update)
}

func (c TaskProcessor) Process(ctx context.Context, update tgApi.Update) {
	cmd := c.parseCommand(update)

	processor, ok := c.routes[cmd]
	if !ok {
		c.ProcessorNotFoundHandler(ctx, update)
		return
	}

	err := processor(ctx, update, update.Message)
	if err != nil {
		c.ErrorHandler(ctx, err, update)
		return
	}

	c.NoErrorHandler(ctx, update)
}

func (c TaskProcessor) parseCommand(update tgApi.Update) string {
	text := update.Message.Text

	// Разделяем текст на команду и остальную часть
	parts := strings.SplitN(text, " ", 2)
	if len(parts) > 1 {
		return parts[0]
	}

	if isAddCmd(update.Message.Text) {
		return "/add"
	}

	return update.Message.Text
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
