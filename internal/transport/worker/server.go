package worker

import (
	"context"
	"fmt"
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mike7109/tg-bot-clubbing/internal/transport/worker/handler"
	"github.com/mike7109/tg-bot-clubbing/internal/transport/worker/update_entity/button"
	"github.com/mike7109/tg-bot-clubbing/internal/transport/worker/update_entity/message"
	"github.com/mike7109/tg-bot-clubbing/pkg/messages"
	"github.com/mike7109/tg-bot-clubbing/pkg/utls"
	"go.uber.org/zap"
)

type TaskProcessor struct {
	messagesHandler          map[message.Command]handler.MessageHandlerFunc
	callbackHandler          map[button.Command]handler.CallbackHandlerFunc
	NoErrorHandler           func(ctx context.Context, update tgApi.Update)
	ErrorHandler             func(ctx context.Context, err error, text string, chatID int64)
	ProcessorNotFoundHandler func(ctx context.Context, text string, chatID int64)
}

func NewTaskProcessor(tgBot *tgApi.BotAPI) TaskProcessor {
	defaultNoErrorHandler := func(ctx context.Context, update tgApi.Update) {
		fmt.Println("Processed update, acking")
	}

	defaultErrorHandler := func(ctx context.Context, err error, text string, chatID int64) {
		fmt.Println("Failed to process update, rejecting", zap.String("Text", text), err)
		_, _ = tgBot.Send(tgApi.NewMessage(chatID, messages.ErrorHandler))
	}

	defaultProcessorNotFoundHandler := func(ctx context.Context, text string, chatID int64) {
		fmt.Println("No processor found for update type", zap.String("Text", text))
		_, _ = tgBot.Send(tgApi.NewMessage(chatID, messages.MsgUnknownCommand))
	}

	return TaskProcessor{
		messagesHandler:          make(map[message.Command]handler.MessageHandlerFunc),
		callbackHandler:          make(map[button.Command]handler.CallbackHandlerFunc),
		NoErrorHandler:           defaultNoErrorHandler,
		ErrorHandler:             defaultErrorHandler,
		ProcessorNotFoundHandler: defaultProcessorNotFoundHandler,
	}
}

func (c TaskProcessor) AddHandlerMessages(routingKey message.Command, fn handler.MessageHandlerFunc) {
	c.messagesHandler[routingKey] = fn
}

func (c TaskProcessor) AddHandlerCallback(routingKey button.Command, fn handler.CallbackHandlerFunc) {
	c.callbackHandler[routingKey] = fn
}

func (c TaskProcessor) Consume(ctx context.Context, update tgApi.Update) {
	c.Process(ctx, update)
}

func (c TaskProcessor) Process(ctx context.Context, update tgApi.Update) {
	var err error
	var chatID int64
	var data string

	if update.Message != nil {
		data = update.Message.Text
		chatID = update.Message.Chat.ID

		cmd := c.parseCommand(update.Message)

		msgHandler, ok := c.messagesHandler[cmd]
		if !ok {
			c.ProcessorNotFoundHandler(ctx, data, chatID)
			return
		}

		err = msgHandler(ctx, update.Message)
	}

	if update.CallbackQuery != nil {
		data = update.CallbackQuery.Data
		chatID = update.CallbackQuery.Message.Chat.ID

		var buttonTarget *button.Button

		buttonTarget, err = button.UnmarshalButton(data)
		if err == nil {
			callbackHandler, ok := c.callbackHandler[buttonTarget.Command]
			if !ok {
				c.ProcessorNotFoundHandler(ctx, data, chatID)
				return
			}

			err = callbackHandler(ctx, update.CallbackQuery, buttonTarget)
		}
	}

	if err != nil {
		c.ErrorHandler(ctx, err, data, chatID)
		return
	}

	c.NoErrorHandler(ctx, update)

}

func (c TaskProcessor) parseCommand(msg *tgApi.Message) message.Command {
	text := msg.Text

	if utls.IsAddCmd(text) {
		return message.AddCommand
	}

	return message.Command(text)
}
