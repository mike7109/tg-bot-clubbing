package entity

import (
	"fmt"
	"github.com/mike7109/tg-bot-clubbing/internal/transport/worker/update_entity/button"
	"strconv"
)

type Page struct {
	UrlPage     []*UrlPage
	currentPage int
	pageSize    int
}

type UrlPage struct {
	ID          int
	URL         string
	UserName    string
	Description *string
	Title       *string
	Category    *string
	Metadata    Metadata
}

type Metadata struct {
	Views  int
	Number int
}

func (p *UrlPage) SetDescription(description string) {
	p.Description = &description
}

func (p *UrlPage) SetTitle(name string) {
	p.Title = &name
}

func (p *UrlPage) SetCategory(category string) {
	p.Category = &category
}

func (p *UrlPage) String() string {
	return fmt.Sprintf("%d. %s\n", p.Metadata.Number, p.URL)
}

func (p *UrlPage) ToButton(cmd button.Command) *button.Button {
	b := button.NewButton(strconv.Itoa(p.Metadata.Number), cmd)
	button.SetDataValue(b, "id", p.ID)
	return b
}
