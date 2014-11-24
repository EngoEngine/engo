package engi

type Text struct {
	Content string
	font    *Font
}

func NewText(content string, face *Font) *Text {
	return &Text{Content: content, font: face}
}

func (text Text) Width() float32 {
	return float32(text.font.CellWidth() * len(text.Content))
}

func (text Text) Height() float32 {
	return float32(text.font.CellHeight())
}

func (text Text) Draw(batch *Batch, position Point) {
	text.font.Print(batch, text.Content, position.X, position.Y, 0xffffff)
}
