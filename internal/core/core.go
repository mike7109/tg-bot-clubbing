package core

import (
	"github.com/joho/godotenv"
	"github.com/mike7109/tg-bot-clubbing/internal/config"
	"github.com/mike7109/tg-bot-clubbing/internal/repositories"
	"github.com/mike7109/tg-bot-clubbing/internal/service"
	"github.com/mike7109/tg-bot-clubbing/pkg/clients/sqlite"
	"github.com/mike7109/tg-bot-clubbing/pkg/clients/telegram"
	"log"
)

type IProcess interface {
	Start() error
}

type Core struct {
	process IProcess
}

func New() (*Core, error) {
	_ = godotenv.Load()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("can't get config: ", err)
	}

	db, err := sqlite.NewSqliteClient(cfg.Database.Path)
	if err != nil {
		log.Fatal("can't init database: ", err)
	}

	storage := repositories.NewStorage(db)

	factory := service.NewFactoryCommand(storage)

	tgClient := telegram.NewTelegramClient(cfg.Telegram.Token, cfg.Debug.Telegram)
	process := service.NewProcess(tgClient, factory)

	return &Core{
		process: process,
	}, nil
}

func (c *Core) Start() error {
	return c.process.Start()
}
