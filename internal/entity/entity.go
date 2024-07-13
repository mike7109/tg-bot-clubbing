package entity

import (
	"errors"
)

var ErrNoSavedPages = errors.New("no saved pages")

type Page struct {
	URL         string
	UserName    string
	Description *string
	Title       *string
	Category    *string
}

func (p *Page) SetDescription(description string) {
	p.Description = &description
}

func (p *Page) SetTitle(name string) {
	p.Title = &name
}

func (p *Page) SetCategory(category string) {
	p.Category = &category
}
