package button

import (
	"github.com/mailru/easyjson"
)

//go:generate easyjson -all button.go

type Command string

const (
	DeleteCommand Command = "1"
)

type Button struct {
	Text    string                 `json:"-"`
	Data    map[string]interface{} `json:"d"`
	Command Command                `json:"c"`
}

func NewButton(text string, cmd Command) *Button {
	return &Button{
		Text:    text,
		Data:    make(map[string]interface{}),
		Command: cmd,
	}
}

func Marshal(b *Button) string {
	data, err := easyjson.Marshal(b)
	if err != nil {
		panic(err)
	}

	return string(data)
}

func UnmarshalButton(data string) (*Button, error) {
	var b Button
	err := easyjson.Unmarshal([]byte(data), &b)
	if err != nil {
		return nil, err
	}

	return &b, nil
}

// GetDataValue function to extract value from Data field
func GetDataValue(b *Button, key string) (interface{}, bool) {
	value, exists := b.Data[key]
	return value, exists
}

// SetDataValue function to set value to Data field
func SetDataValue(b *Button, key string, value interface{}) {
	b.Data[key] = value
}
