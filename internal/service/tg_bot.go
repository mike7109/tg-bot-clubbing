package service

import (
	"context"
	"github.com/mike7109/tg-bot-clubbing/internal/apperrors"
	"github.com/mike7109/tg-bot-clubbing/internal/entity"
	"github.com/mike7109/tg-bot-clubbing/internal/service/dto"
	"github.com/mike7109/tg-bot-clubbing/pkg/messages"
	"math"
)

type ITgBotService interface {
	StartHandler() string
	HelpHandler() string
	SaveHandler(ctx context.Context, page *dto.SavePage) (string, error)
	ListHandler(ctx context.Context, userName string, offset int) ([]*entity.UrlPage, error)
	DeleteHandler(ctx context.Context, id int, userName string) error
	GetPageHandler(ctx context.Context, userName string, page, pageSize int) (*dto.ListPage, error)
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

func (t TgBotService) ListHandler(ctx context.Context, userName string, offset int) ([]*entity.UrlPage, error) {
	pages, err := t.storage.ListUrl(ctx, userName, offset, 10)
	if err != nil {
		return nil, apperrors.ErrNoUrl
	}

	if len(pages) == 0 {
		return nil, apperrors.ErrNoUrl
	}

	return pages, nil
}

func (t TgBotService) DeleteHandler(ctx context.Context, id int, userName string) error {
	return t.storage.DeleteUrl(ctx, id, userName)
}

func (t TgBotService) GetPageHandler(ctx context.Context, userName string, page, pageSize int) (*dto.ListPage, error) {
	offset := (page) * pageSize

	totalUrl, err := t.storage.CountUrl(ctx, userName)
	if err != nil {
		return nil, apperrors.ErrNoUrl
	}

	if totalUrl == 0 {
		return nil, apperrors.ErrNoUrl
	}

	if offset >= totalUrl {
		if page > 0 {
			page--
			offset = page * pageSize
		} else {
			offset = 0
		}
	}

	pages, err := t.storage.ListUrl(ctx, userName, offset, pageSize)
	if err != nil {
		return nil, apperrors.ErrNoUrl
	}

	if len(pages) == 0 {
		return nil, apperrors.ErrNoUrl
	}

	currentMaxUrl := offset + len(pages)

	totalPages := int(math.Ceil(float64(totalUrl) / float64(pageSize)))

	pagesList := dto.NewListPage(pages, currentMaxUrl, totalUrl, page, totalPages)

	return pagesList, nil
}

func (t TgBotService) CountHandler(ctx context.Context, userName string) (int, error) {
	return t.storage.CountUrl(ctx, userName)
}