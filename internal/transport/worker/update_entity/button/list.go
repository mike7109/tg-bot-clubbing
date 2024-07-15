package button

type ListButton struct {
	Cmd         CommandButton
	CurrentPage int
	WithDelete  int
	ID          int
}

func CreateListButton(cmd CommandButton, currentPage, withDelete int) *Button {
	b := NewButton("", ListCommand)
	SetDataValue(b, "p", currentPage)
	SetDataValue(b, "d", withDelete)
	SetDataValue(b, "c", cmd)
	return b
}

func CreateListButtonWithID(cmd CommandButton, currentPage, withDelete int, id int) *Button {
	b := CreateListButton(cmd, currentPage, withDelete)
	SetDataValue(b, "id", id)
	return b
}
