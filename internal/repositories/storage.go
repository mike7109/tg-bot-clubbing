package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mike7109/tg-bot-clubbing/internal/apperrors"
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
func (s *Storage) Save(ctx context.Context, p *entity.UrlPage) error {
	q := `INSERT OR REPLACE INTO pages (url, user_name, name, description, category) VALUES (?, ?, ?, ?, ?);`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName, p.Title, p.Description, p.Category); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	return nil
}

// PickRandom picks random page from storage.
func (s *Storage) PickRandom(ctx context.Context, userName string) (*entity.UrlPage, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.ErrNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick random page: %w", err)
	}

	return &entity.UrlPage{
		URL:      url,
		UserName: userName,
	}, nil
}

func (s *Storage) DeleteUrl(ctx context.Context, id int, userName string) error {
	q := `DELETE FROM pages WHERE id = ? AND user_name = ?`
	if _, err := s.db.ExecContext(ctx, q, id, userName); err != nil {
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
func (s *Storage) IsExists(ctx context.Context, page *entity.UrlPage) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, page.URL, page.UserName).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}

	return count > 0, nil
}

// ListUrl returns list of saved pages.
func (s *Storage) ListUrl(ctx context.Context, userName string, offset int, limit int) ([]*entity.UrlPage, error) {
	q := `
			SELECT id, url, name, description, category FROM pages WHERE user_name = ? ORDER BY created_at ASC LIMIT ? OFFSET ?;
		`

	rows, err := s.db.QueryContext(ctx, q, userName, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNoSavedPages
		}

		return nil, fmt.Errorf("can't list pages: %w", err)
	}
	defer rows.Close()

	number := 1

	var pages []*entity.UrlPage
	for rows.Next() {
		page := &entity.UrlPage{UserName: userName}
		if err = rows.Scan(&page.ID, &page.URL, &page.Title, &page.Description, &page.Category); err != nil {
			return nil, fmt.Errorf("can't scan page: %w", err)
		}

		page.Metadata.Number = number

		pages = append(pages, page)
		number++
	}

	if len(pages) == 0 {
		return nil, apperrors.ErrNoSavedPages
	}

	return pages, nil
}

func (s *Storage) CountPage(ctx context.Context, userName string) (int, error) {
	q := `SELECT COUNT(*) FROM pages WHERE user_name = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, userName).Scan(&count); err != nil {
		return 0, fmt.Errorf("can't count pages: %w", err)
	}

	return count, nil
}
