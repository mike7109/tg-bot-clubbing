package message

type Command string

const (
	HelpCommand      Command = "/help"
	StartCommand     Command = "/start"
	AddCommand       Command = "/add"
	AddSimpleCommand Command = "/add_simple"
	ListCommand      Command = "/list"
	DeleteAllCommand Command = "/delete_all"
	DeleteCommand    Command = "/delete"
)
