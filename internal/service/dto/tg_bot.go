package dto

import (
	"github.com/mike7109/tg-bot-clubbing/internal/apperrors"
	"github.com/mike7109/tg-bot-clubbing/internal/entity"
	"github.com/mike7109/tg-bot-clubbing/pkg/utls"
)

type SavePage struct {
	Url         string
	UserName    string
	Description *string
	Title       *string
	Category    *string
}

func NewSavePage(url, userName string) *SavePage {
	return &SavePage{
		Url:      url,
		UserName: userName,
	}
}

func (p *SavePage) SetDescription(description string) {
	p.Description = &description
}

func (p *SavePage) SetTitle(name string) {
	p.Title = &name
}

func (p *SavePage) SetCategory(category string) {
	p.Category = &category
}

func (p *SavePage) ToEntity() *entity.UrlPage {
	return &entity.UrlPage{
		URL:         p.Url,
		UserName:    p.UserName,
		Description: p.Description,
		Title:       p.Title,
		Category:    p.Category,
	}
}

func (p *SavePage) Validate() error {
	if p.Url == "" {
		return apperrors.ErrNoURL
	}

	if p.UserName == "" {
		return apperrors.ErrNoUserName
	}

	if !utls.IsURL(p.Url) {
		return apperrors.ErrInvalidURL
	}

	return nil
}

type ListPage struct {
	SavePage     []*entity.UrlPage
	NumPage      int
	LastPage     int
	HaveNextPage bool
	HavePrevPage bool
	IsNoPage     bool
}

func NewListPage(pages []*entity.UrlPage, currentMaxUrl, totalUrl, currentPage, totalPages int) *ListPage {
	listPage := &ListPage{
		SavePage: pages,
		NumPage:  currentPage,
		LastPage: totalPages,
	}

	if currentPage < totalPages-1 {
		listPage.HaveNextPage = true
	}

	if currentPage > 0 {
		listPage.HavePrevPage = true
	}

	if currentMaxUrl >= totalUrl {
		listPage.HaveNextPage = false
	}

	return listPage
}
