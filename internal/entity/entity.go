package entity

import (
	"fmt"
	"github.com/mike7109/tg-bot-clubbing/internal/transport/worker/update_entity/button"
	"strconv"
)

type Page struct {
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

func (p *Page) SetDescription(description string) {
	p.Description = &description
}

func (p *Page) SetTitle(name string) {
	p.Title = &name
}

func (p *Page) SetCategory(category string) {
	p.Category = &category
}

func (p *Page) String() string {
	return fmt.Sprintf("%d. %s\n", p.Metadata.Number, p.URL)
}

func (p *Page) ToButton(cmd button.Command) *button.Button {
	b := button.NewButton(strconv.Itoa(p.Metadata.Number), cmd)
	button.SetDataValue(b, "id", p.ID)
	return b
}
