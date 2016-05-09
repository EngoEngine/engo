package engo

type Cursor uint8

const (
	CursorNone = iota
	CursorArrow
	CursorCrosshair
	CursorHand
	CursorIBeam
	CursorHResize
	CursorVResize
)
