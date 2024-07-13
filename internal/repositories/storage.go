package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mike7109/tg-bot-clubbing/internal/entity"
)

type Storage struct {
	db *sql.DB
}

// New creates new SQLite storage.
func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

// Save saves page to storage.
func (s *Storage) Save(ctx context.Context, p *entity.Page) error {
	q := `INSERT INTO pages (url, user_name, name, description, category) VALUES (?, ?, ?, ?, ?)`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName, p.Title, p.Description, p.Category); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	return nil
}

// PickRandom picks random page from storage.
func (s *Storage) PickRandom(ctx context.Context, userName string) (*entity.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, entity.ErrNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick random page: %w", err)
	}

	return &entity.Page{
		URL:      url,
		UserName: userName,
	}, nil
}

// Remove removes page from storage.
func (s *Storage) Remove(ctx context.Context, page *entity.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name = ?`
	if _, err := s.db.ExecContext(ctx, q, page.URL, page.UserName); err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}

	return nil
}

func (s *Storage) DeleteAll(ctx context.Context, userName string) error {
	q := `DELETE FROM pages WHERE user_name = ?`
	if _, err := s.db.ExecContext(ctx, q, userName); err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}

	return nil
}

// IsExists checks if page exists in storage.
func (s *Storage) IsExists(ctx context.Context, page *entity.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, page.URL, page.UserName).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}

	return count > 0, nil
}

// ListUrl returns list of saved pages.
func (s *Storage) ListUrl(ctx context.Context, userName string) ([]*entity.Page, error) {
	q := `SELECT url, name, description, category FROM pages WHERE user_name = ? ORDER BY created_at ASC`

	rows, err := s.db.QueryContext(ctx, q, userName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNoSavedPages
		}
	}

	defer rows.Close()

	var pages []*entity.Page
	for rows.Next() {
		page := &entity.Page{UserName: userName}
		if err = rows.Scan(&page.URL, &page.Title, &page.Description, &page.Category); err != nil {
			return nil, fmt.Errorf("can't scan page: %w", err)
		}
		pages = append(pages, page)
	}

	return pages, nil
}
