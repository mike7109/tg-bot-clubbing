package button

import (
	"fmt"
	"github.com/mailru/easyjson"
)

//go:generate easyjson -all button.go

type Command string

const (
	DeleteCommand             Command = "1"
	WantToDeleteCommand       Command = "2"
	CancelWantToDeleteCommand Command = "3"
	NextPageCommand           Command = "4"
	ListCommand               Command = "5"
)

type CommandButton string

const (
	SwitchPageCommandButton            CommandButton = "1"
	DeleteURLCommandButton             CommandButton = "2"
	WantToDeleteURLCommandButton       CommandButton = "3"
	CancelWantToDeleteURLCommandButton CommandButton = "4"
)

type Button struct {
	Text    string                 `json:"-"`
	Data    map[string]interface{} `json:"d,omitempty"`
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

func (b *Button) ToListButton() (*ListButton, error) {
	out := &ListButton{}

	cmd, exist := GetDataValue(b, "c")
	if !exist {
		return nil, fmt.Errorf("cmd not found in data")
	}
	out.Cmd = CommandButton(cmd.(string))

	wantToDelete, exist := GetDataValue(b, "d")
	if !exist {
		return nil, fmt.Errorf("wantToDelete not found in data")
	}

	out.WithDelete = int(wantToDelete.(float64))

	page, exist := GetDataValue(b, "p")
	if !exist {
		return nil, fmt.Errorf("page not found in data")
	}
	out.CurrentPage = int(page.(float64))

	if id, exist := GetDataValue(b, "id"); exist {
		out.ID = int(id.(float64))
	}

	return out, nil
}
