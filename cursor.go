package engo

// Cursor is a reference to standard cursors, to be used in conjunction with `SetCursor`. What they look like, is
// different for each platform.
type Cursor uint8

const (
	// CursorNone can be used to reset the cursor.
	CursorNone Cursor = iota
	// CursorArrow represents an arrow cursor
	CursorArrow
	// CursorCrosshair represents a crosshair cursor
	CursorCrosshair
	// CursorHand represents a hand cursor
	CursorHand
	// CursorIBeam represents an IBeam cursor
	CursorIBeam
	// CursorHResize represents a HResize cursor
	CursorHResize
	// CursorVResize represents a VResize cursor
	CursorVResize
)
