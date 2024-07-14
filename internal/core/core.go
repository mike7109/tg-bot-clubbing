package core

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mike7109/tg-bot-clubbing/internal/config"
	"github.com/mike7109/tg-bot-clubbing/internal/repositories"
	"github.com/mike7109/tg-bot-clubbing/internal/service"
	"github.com/mike7109/tg-bot-clubbing/internal/transport/worker"
	"github.com/mike7109/tg-bot-clubbing/pkg/clients/sqlite"
	"github.com/mike7109/tg-bot-clubbing/pkg/clients/telegram"
	"log"
	"os"
	"os/signal"
)

func New() error {
	_ = godotenv.Load()

	fmt.Println("Starting bot...")

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("can't get config: ", err)
	}

	fmt.Println("Config loaded...")

	db, err := sqlite.NewSqliteClient(cfg.Database.Path)
	if err != nil {
		log.Fatal("can't init database: ", err)
	}

	fmt.Println("Database connected...")

	storage := repositories.NewStorage(db)

	tgBotService := service.NewTgBotService(storage)

	tgClient := telegram.NewTelegramClient(cfg.Telegram.Token, cfg.Debug.Telegram)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	p, err := worker.NewTgProcessor(ctx, tgClient, tgBotService)
	if err != nil {
		log.Fatal("can't create telegram processor: ", err)
	}

	fmt.Println("Telegram processor created...")

	defer func() {
		fmt.Println("Shutting down")
		p.WaitClose()
		db.Close()
		fmt.Println("Exiting")
	}()

	// Ожидание событий в канале ошибок или сигналов
	select {
	case <-ctx.Done():
	}

	return err
}
