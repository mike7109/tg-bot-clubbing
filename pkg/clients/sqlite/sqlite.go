package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"time"
)

// NewSqliteClient creates new SQLite storage.
func NewSqliteClient(path string) (*sql.DB, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("can't create directory: %w", err)
	}

	// Check if file exists, create if not
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("can't create database file: %w", err)
		}
		file.Close()
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	if err := createTable(db); err != nil {
		return nil, fmt.Errorf("can't create table: %w", err)
	}

	return db, nil
}

func createTable(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	q := `
			CREATE TABLE IF NOT EXISTS pages (
			id integer primary key,
			url TEXT NOT NULL,
			user_name TEXT NOT NULL,
			name TEXT,
			description TEXT,
			category TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (url, user_name) -- Добавление уникального индекса
		);
	`

	_, err := db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	q = `CREATE UNIQUE INDEX idx_pages_url_user_name ON pages (url, user_name);`

	_, err = db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}
