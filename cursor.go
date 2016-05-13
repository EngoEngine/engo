package engo

// Cursor is a reference to standard cursors, to be used in conjunction with `SetCursor`. What they look like, is
// different for each platform.
type Cursor uint8

const (
	// CursorNone can be used to reset the cursor.
	CursorNone Cursor = iota
	CursorArrow
	CursorCrosshair
	CursorHand
	CursorIBeam
	CursorHResize
	CursorVResize
)
