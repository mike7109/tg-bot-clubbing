package config

type Database struct {
	Path string `env:"SQLITE_PATH" env-default:"data/sqlite/storage.db"`
}
