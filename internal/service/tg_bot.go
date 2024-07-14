package service

import (
	"context"
	"github.com/mike7109/tg-bot-clubbing/internal/apperrors"
	"github.com/mike7109/tg-bot-clubbing/internal/entity"
	"github.com/mike7109/tg-bot-clubbing/internal/service/dto"
	"github.com/mike7109/tg-bot-clubbing/pkg/messages"
)

type ITgBotService interface {
	StartHandler() string
	HelpHandler() string
	SaveHandler(ctx context.Context, page *dto.SavePage) (string, error)
	ListHandler(ctx context.Context, userName string, offset int) ([]*entity.Page, error)
	DeleteHandler(ctx context.Context, id int, userName string) error
	GetPageHandler(ctx context.Context, userName string, page, pageSize int) ([]*entity.Page, error)
	CountHandler(ctx context.Context, userName string) (int, error)
}

type TgBotService struct {
	storage IStorage
	//cache   ICache
}

func NewTgBotService(storage IStorage) TgBotService {
	return TgBotService{
		storage: storage,
	}
}

func (t TgBotService) StartHandler() string {
	return messages.MsgHelp
}

func (t TgBotService) HelpHandler() string {
	return messages.MsgHelp
}

func (t TgBotService) SaveHandler(ctx context.Context, page *dto.SavePage) (string, error) {
	if err := t.storage.Save(ctx, page.ToEntity()); err != nil {
		return "", apperrors.ErrNoSave
	}

	return messages.MsgSaved, nil
}

func (t TgBotService) ListHandler(ctx context.Context, userName string, offset int) ([]*entity.Page, error) {
	pages, err := t.storage.ListUrl(ctx, userName, offset, 10)
	if err != nil {
		return nil, apperrors.ErrNoPages
	}

	if len(pages) == 0 {
		return nil, apperrors.ErrNoPages
	}

	return pages, nil
}

func (t TgBotService) DeleteHandler(ctx context.Context, id int, userName string) error {
	return t.storage.DeleteUrl(ctx, id, userName)
}

func (t TgBotService) GetPageHandler(ctx context.Context, userName string, page, pageSize int) ([]*entity.Page, error) {
	offset := (page) * pageSize

	pages, err := t.storage.ListUrl(ctx, userName, offset, pageSize)

	if err != nil {
		return nil, apperrors.ErrNoPages
	}

	if len(pages) == 0 {
		return nil, apperrors.ErrNoPages
	}

	return pages, nil
}

func (t TgBotService) CountHandler(ctx context.Context, userName string) (int, error) {
	return t.storage.CountPage(ctx, userName)
}
