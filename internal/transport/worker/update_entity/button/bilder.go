package button

import tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Builder struct {
	buttons                   []*Button
	buttonsKeyboardTopRow     []tgApi.InlineKeyboardButton
	buttonsKeyboardMiddleRows [][]tgApi.InlineKeyboardButton
	buttonsKeyboardBottomRow  []tgApi.InlineKeyboardButton
	buttonsKeyboardRows       []tgApi.InlineKeyboardButton
	InlineKeyboardMarkup      tgApi.InlineKeyboardMarkup
}

func NewBuilder() *Builder {
	return &Builder{
		buttons:                   make([]*Button, 0),
		buttonsKeyboardTopRow:     make([]tgApi.InlineKeyboardButton, 0),
		buttonsKeyboardMiddleRows: make([][]tgApi.InlineKeyboardButton, 0),
		buttonsKeyboardBottomRow:  make([]tgApi.InlineKeyboardButton, 0),
		InlineKeyboardMarkup:      tgApi.InlineKeyboardMarkup{},
	}
}

func (b *Builder) AddButton(button *Button) {
	b.buttons = append(b.buttons, button)
}

func (b *Builder) AddButtons(buttons ...*Button) {
	b.buttons = append(b.buttons, buttons...)
}

func (b *Builder) AddButtonTopRow(button *Button) {
	b.buttonsKeyboardTopRow = append(b.buttonsKeyboardTopRow, tgApi.NewInlineKeyboardButtonData(button.Text, Marshal(button)))
}

func (b *Builder) AddButtonTopRows(buttons ...*Button) {
	for _, button := range buttons {
		b.buttonsKeyboardTopRow = append(b.buttonsKeyboardTopRow, tgApi.NewInlineKeyboardButtonData(button.Text, Marshal(button)))
	}
}

func (b *Builder) AddButtonMiddleRow(button *Button) {
	b.buttonsKeyboardMiddleRows = append(b.buttonsKeyboardMiddleRows, []tgApi.InlineKeyboardButton{tgApi.NewInlineKeyboardButtonData(button.Text, Marshal(button))})
}

func (b *Builder) AddButtonBottomRow(button *Button) {
	b.buttonsKeyboardBottomRow = append(b.buttonsKeyboardBottomRow, tgApi.NewInlineKeyboardButtonData(button.Text, Marshal(button)))
}

func (b *Builder) Build() tgApi.InlineKeyboardMarkup {
	var row []tgApi.InlineKeyboardButton

	if len(b.buttonsKeyboardTopRow) > 0 {
		b.buttonsKeyboardMiddleRows = append([][]tgApi.InlineKeyboardButton{b.buttonsKeyboardTopRow}, b.buttonsKeyboardMiddleRows...)
	}

	for i, button := range b.buttons {
		row = append(row, tgApi.NewInlineKeyboardButtonData(button.Text, Marshal(button)))

		// Если набрали 5 кнопок в строке или это последняя кнопка
		if (i+1)%5 == 0 || i == len(b.buttons)-1 {
			b.buttonsKeyboardMiddleRows = append(b.buttonsKeyboardMiddleRows, row)
			row = []tgApi.InlineKeyboardButton{} // Очищаем строку для следующей партии кнопок
		}

	}

	if len(b.buttonsKeyboardBottomRow) > 0 {
		b.buttonsKeyboardMiddleRows = append(b.buttonsKeyboardMiddleRows, b.buttonsKeyboardBottomRow)
	}

	b.InlineKeyboardMarkup = tgApi.NewInlineKeyboardMarkup(b.buttonsKeyboardMiddleRows...)

	return b.InlineKeyboardMarkup
}
