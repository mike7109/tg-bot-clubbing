package service

import (
	"context"
	"github.com/mike7109/tg-bot-clubbing/internal/entity"
)

type IStorage interface {
	Save(ctx context.Context, p *entity.UrlPage) error
	PickRandom(ctx context.Context, userName string) (*entity.UrlPage, error)
	DeleteUrl(ctx context.Context, id int, userName string) error
	DeleteAll(ctx context.Context, userName string) error
	IsExists(ctx context.Context, p *entity.UrlPage) (bool, error)
	ListUrl(ctx context.Context, userName string, offset int, limit int) ([]*entity.UrlPage, error)
	CountPage(ctx context.Context, userName string) (int, error)
}
