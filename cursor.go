package engo

// Cursor is a reference to standard GLFW-cursors, to be used in conjunction with `SetCursor`
type Cursor uint8

const (
	// CursorNone can be used to reset the cursor
	CursorNone = iota
	CursorArrow
	CursorCrosshair
	CursorHand
	CursorIBeam
	CursorHResize
	CursorVResize
)
