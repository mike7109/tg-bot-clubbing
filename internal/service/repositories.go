package service

import (
	"context"
	"github.com/mike7109/tg-bot-clubbing/internal/entity"
)

type IStorage interface {
	Save(ctx context.Context, p *entity.Page) error
	PickRandom(ctx context.Context, userName string) (*entity.Page, error)
	Remove(ctx context.Context, p *entity.Page) error
	DeleteAll(ctx context.Context, userName string) error
	IsExists(ctx context.Context, p *entity.Page) (bool, error)
	ListUrl(ctx context.Context, userName string) ([]*entity.Page, error)
}
