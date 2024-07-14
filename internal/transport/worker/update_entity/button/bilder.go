package button

import tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Builder struct {
	buttons              []*Button
	buttonsKeyboard      [][]tgApi.InlineKeyboardButton
	InlineKeyboardMarkup tgApi.InlineKeyboardMarkup
}

func NewBuilder() *Builder {
	return &Builder{
		buttons:              make([]*Button, 0),
		buttonsKeyboard:      make([][]tgApi.InlineKeyboardButton, 0),
		InlineKeyboardMarkup: tgApi.InlineKeyboardMarkup{},
	}
}

func (b *Builder) AddButton(button *Button) {
	b.buttons = append(b.buttons, button)
}

func (b *Builder) Build() tgApi.InlineKeyboardMarkup {
	var row []tgApi.InlineKeyboardButton

	for i, button := range b.buttons {
		row = append(row, tgApi.NewInlineKeyboardButtonData(button.Text, Marshal(button)))

		// Если набрали 5 кнопок в строке или это последняя кнопка
		if (i+1)%5 == 0 || i == len(b.buttons)-1 {
			b.buttonsKeyboard = append(b.buttonsKeyboard, row)
			row = []tgApi.InlineKeyboardButton{} // Очищаем строку для следующей партии кнопок
		}
	}

	b.InlineKeyboardMarkup = tgApi.NewInlineKeyboardMarkup(b.buttonsKeyboard...)

	return b.InlineKeyboardMarkup
}
