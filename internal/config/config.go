package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

var serviceName string

func GetServiceName() string {
	return serviceName
}

type Debug struct {
	Logger   bool `env:"DEBUG_LOGGER"`
	Database bool `env:"DEBUG_DATABASE"`
	Telegram bool `env:"DEBUG_TELEGRAM"`
}

type Server struct {
	ServerName string `env:"SERVICE_NAME" env-default:"tg-bot-clubbing"`
}

type Config struct {
	Debug    Debug
	Database Database
	Server   Server
	Telegram Telegram
}

type Telegram struct {
	Token string `env:"TELEGRAM_BOT_TOKEN"`
}

type PartsOut struct {
	Debug
	Database
	Server
	Telegram
}

func (cfg *Config) ToParts() PartsOut {
	return PartsOut{
		Debug:    cfg.Debug,
		Database: cfg.Database,
		Server:   cfg.Server,
		Telegram: cfg.Telegram,
	}
}

func NewConfig() (*Config, error) {
	cfg := new(Config)

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	serviceName = cfg.Server.ServerName

	return cfg, nil
}
