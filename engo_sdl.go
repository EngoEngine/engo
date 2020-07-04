// +build sdl

package engo

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/EngoEngine/gl"

	"github.com/Noofbiz/sdlMojaveFix"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	// Window is the sdl Window used for engo
	Window *sdl.Window

	cursorNone      *sdl.Cursor
	cursorArrow     *sdl.Cursor
	cursorIBeam     *sdl.Cursor
	cursorCrosshair *sdl.Cursor
	cursorHand      *sdl.Cursor
	cursorHResize   *sdl.Cursor
	cursorVResize   *sdl.Cursor

	Gl           *gl.Context
	sdlGLContext sdl.GLContext

	scale = float32(1)
)

func init() {
	runtime.LockOSThread()
}

// fatalErr calls log.Fatal with the given error if it is non-nil.
func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// CreateWindow opens the window and gets a GL surface for rendering
func CreateWindow(title string, width, height int, fullscreen bool, msaa int) {
	CurrentBackEnd = BackEndSDL

	err := sdl.Init(sdl.INIT_EVERYTHING)
	fatalErr(err)

	if !opts.HeadlessMode {
		cursorNone = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_NO)
		cursorArrow = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_ARROW)
		cursorIBeam = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_IBEAM)
		cursorCrosshair = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_CROSSHAIR)
		cursorHand = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_HAND)
		cursorHResize = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_SIZENS)
		cursorVResize = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_SIZEWE)
	}

	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 2)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)

	if msaa > 0 {
		sdl.GLSetAttribute(sdl.GL_MULTISAMPLEBUFFERS, 1)
		sdl.GLSetAttribute(sdl.GL_MULTISAMPLESAMPLES, msaa)
	}

	SetVSync(opts.VSync)

	Window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, int32(width), int32(height), sdl.WINDOW_OPENGL)
	fatalErr(err)

	sdlGLContext, err = Window.GLCreateContext()
	fatalErr(err)

	Gl = gl.NewContext()

	if fullscreen {
		Window.SetFullscreen(sdl.WINDOW_FULLSCREEN)
	}
	if opts.NotResizable {
		Window.SetResizable(false)
	} else {
		Window.SetResizable(true)
	}

	gameWidth, gameHeight = float32(width), float32(height)

	w, h := Window.GetSize()
	windowWidth, windowHeight = float32(w), float32(h)

	fw, fh := Window.GLGetDrawableSize()
	canvasWidth, canvasHeight = float32(fw), float32(fh)

	if windowWidth <= canvasWidth && windowHeight <= canvasHeight {
		scale = canvasWidth / windowWidth
	}
}

// DestroyWindow handles the termination of windows
func DestroyWindow() {
	sdl.GLDeleteContext(sdlGLContext)
	Window.Destroy()
	sdl.Quit()
}

// SetTitle sets the title of the window
func SetTitle(title string) {
	if opts.HeadlessMode {
		log.Println("Title set to:", title)
	} else {
		Window.SetTitle(title)
	}
}

// RunIteration runs one iteration per frame
func RunIteration() {
	Time.Tick()

	// First check for new keypresses
	if !opts.HeadlessMode {
		Input.update()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				Exit()
			case *sdl.KeyboardEvent:
				key := Key(e.Keysym.Sym)
				if e.GetType() == sdl.KEYUP {
					Input.keys.Set(key, false)
				} else if e.GetType() == sdl.KEYDOWN {
					Input.keys.Set(key, true)
				}
				Input.Modifier = Modifier(sdl.GetModState())

				// SDL supports codes for both left/right mods, to keep the API similar
				// with the GLFW implementation, either will trigger the mod for that key
				if Input.Modifier&sdl.KMOD_SHIFT != 0 {
					Input.Modifier = Shift
				} else if Input.Modifier&sdl.KMOD_CTRL != 0 {
					Input.Modifier = Control
				} else if Input.Modifier&sdl.KMOD_ALT != 0 {
					Input.Modifier = Alt
				} else if Input.Modifier&sdl.KMOD_GUI != 0 {
					Input.Modifier = Super
				}

			case *sdl.MouseWheelEvent:
				Input.Mouse.ScrollX = float32(e.X)
				Input.Mouse.ScrollY = float32(e.Y)
			case *sdl.MouseButtonEvent:
				Input.Mouse.X, Input.Mouse.Y = float32(e.X)/(opts.GlobalScale.X), float32(e.Y)/(opts.GlobalScale.Y)

				switch e.Button {
				case sdl.BUTTON_LEFT:
					Input.Mouse.Button = MouseButtonLeft
				case sdl.BUTTON_MIDDLE:
					Input.Mouse.Button = MouseButtonMiddle
				case sdl.BUTTON_RIGHT:
					Input.Mouse.Button = MouseButtonRight
				case sdl.BUTTON_X1:
					Input.Mouse.Button = MouseButton4
				case sdl.BUTTON_X2:
					Input.Mouse.Button = MouseButton5
				}

				Input.Mouse.Modifer = Input.Modifier

				if e.State == sdl.PRESSED {
					Input.Mouse.Action = Press
				} else {
					Input.Mouse.Action = Release
				}
			case *sdl.MouseMotionEvent:
				Input.Mouse.X, Input.Mouse.Y = float32(e.X)/opts.GlobalScale.X, float32(e.Y)/opts.GlobalScale.Y
				if Input.Mouse.Action != Release && Input.Mouse.Action != Press {
					Input.Mouse.Action = Move
				}
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_RESIZED {

					w, h := Window.GetSize()
					fw, fh := Window.GLGetDrawableSize()

					message := WindowResizeMessage{
						OldWidth:  int(windowWidth),
						OldHeight: int(windowHeight),
						NewWidth:  int(w),
						NewHeight: int(h),
					}

					Gl.Viewport(0, 0, int(fw), int(fh))
					windowWidth, windowHeight = float32(w), float32(h)

					oldCanvasW, oldCanvasH := canvasWidth, canvasHeight

					canvasWidth, canvasHeight = float32(fw), float32(fh)

					ResizeXOffset += oldCanvasW - canvasWidth
					ResizeYOffset += oldCanvasH - canvasHeight

					if !opts.ScaleOnResize {
						gameWidth, gameHeight = float32(w), float32(h)
					}

					if windowWidth <= canvasWidth && windowHeight <= canvasHeight {
						scale = canvasWidth / windowWidth
					}

					Mailbox.Dispatch(message)
				}
			case *sdl.TextInputEvent:
				n := bytes.IndexByte(e.Text[:], 0)
				s := string(e.Text[:n])
				if len(s) == 1 {
					Mailbox.Dispatch(TextMessage{[]rune(s)[0]})
				}
			}
		}
	}

	// Then update the world and all Systems
	currentUpdater.Update(Time.Delta())

	// Lastly, forget keypresses and swap buffers
	if !opts.HeadlessMode {
		// reset values to avoid catching the same "signal" twice
		Input.Mouse.ScrollX, Input.Mouse.ScrollY = 0, 0
		Input.Mouse.Action = Neutral
		sdlMojaveFix.UpdateNSGLContext(sdlGLContext)
		Window.GLSwap()
	}
}

// RunPreparation is called automatically when calling Open. It should only be called once.
func RunPreparation(defaultScene Scene) {
	Time = NewClock()
	SetScene(defaultScene, false)
}

func runLoop(defaultScene Scene, headless bool) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		closeEvent()
	}()

	RunPreparation(defaultScene)
	ticker := time.NewTicker(time.Duration(int(time.Second) / opts.FPSLimit))

	// Start tick, minimize the delta
	Time.Tick()

	for {
		select {
		case <-ticker.C:
			RunIteration()
		case <-resetLoopTicker:
			ticker.Stop()
			ticker = time.NewTicker(time.Duration(int(time.Second) / opts.FPSLimit))
		case <-closeGame:
			ticker.Stop()
			closeEvent()
			return
		}
	}
}

// CursorPos returns the current cursor position
func CursorPos() (x, y float32) {
	w, h, _ := sdl.GetMouseState()
	return float32(w), float32(h)
}

// WindowSize gets the current window size
func WindowSize() (w, h int) {
	width, height := Window.GetSize()
	return int(width), int(height)
}

// WindowWidth gets the current window width
func WindowWidth() float32 {
	return windowWidth
}

// WindowHeight gets the current window height
func WindowHeight() float32 {
	return windowHeight
}

// CanvasWidth gets the width of the current OpenGL Framebuffer
func CanvasWidth() float32 {
	return canvasWidth
}

// CanvasHeight gets the height of the current OpenGL Framebuffer
func CanvasHeight() float32 {
	return canvasHeight
}

// CanvasScale gets the ratio of the canvas to the window sizes
func CanvasScale() float32 {
	return scale
}

// SetCursor sets the pointer of the mouse to the defined standard cursor
func SetCursor(c Cursor) {
	var cur *sdl.Cursor
	switch c {
	case CursorNone:
		cur = cursorNone
	case CursorArrow:
		cur = cursorArrow
	case CursorCrosshair:
		cur = cursorCrosshair
	case CursorHand:
		cur = cursorHand
	case CursorIBeam:
		cur = cursorIBeam
	case CursorHResize:
		cur = cursorHResize
	case CursorVResize:
		cur = cursorVResize
	}
	sdl.SetCursor(cur)
}

// SetVSync sets whether or not to use VSync
func SetVSync(enabled bool) {
	opts.VSync = enabled
	if opts.VSync {
		err := sdl.GLSetSwapInterval(-1)
		if err != nil {
			sdl.GLSetSwapInterval(1)
		}
	} else {
		sdl.GLSetSwapInterval(0)
	}
}

//SetCursorVisibility sets the visibility of the cursor.
//If true the cursor is visible, if false the cursor is not.
func SetCursorVisibility(visible bool) {
	if visible {
		sdl.ShowCursor(sdl.ENABLE)
	} else {
		sdl.ShowCursor(sdl.DISABLE)
	}
}

// openFile is the desktop-specific way of opening a file
func openFile(url string) (io.ReadCloser, error) {
	return os.Open(url)
}

// IsAndroidChrome tells if the browser is Chrome for android
func IsAndroidChrome() bool {
	return false
}

// GetKeyName returns the string returned from the given Key. So you can write "Press W to move forward"
// and get a W for QWERTY and Z for AZERTY
func GetKeyName(k Key) string {
	m := map[Key]sdl.Keycode{
		KeyGrave:        sdl.K_BACKQUOTE,
		KeyDash:         sdl.K_MINUS,
		KeyApostrophe:   sdl.K_QUOTE,
		KeySemicolon:    sdl.K_SEMICOLON,
		KeyEquals:       sdl.K_EQUALS,
		KeyComma:        sdl.K_COMMA,
		KeyPeriod:       sdl.K_PERIOD,
		KeySlash:        sdl.K_SLASH,
		KeyBackslash:    sdl.K_BACKSLASH,
		KeyBackspace:    sdl.K_BACKSPACE,
		KeyTab:          sdl.K_TAB,
		KeyCapsLock:     sdl.K_CAPSLOCK,
		KeySpace:        sdl.K_SPACE,
		KeyEnter:        sdl.K_RETURN,
		KeyEscape:       sdl.K_ESCAPE,
		KeyInsert:       sdl.K_INSERT,
		KeyPrintScreen:  sdl.K_PRINTSCREEN,
		KeyDelete:       sdl.K_DELETE,
		KeyPageUp:       sdl.K_PAGEUP,
		KeyPageDown:     sdl.K_PAGEDOWN,
		KeyHome:         sdl.K_HOME,
		KeyEnd:          sdl.K_END,
		KeyPause:        sdl.K_PAUSE,
		KeyScrollLock:   sdl.K_SCROLLLOCK,
		KeyArrowLeft:    sdl.K_LEFT,
		KeyArrowRight:   sdl.K_RIGHT,
		KeyArrowDown:    sdl.K_DOWN,
		KeyArrowUp:      sdl.K_UP,
		KeyLeftBracket:  sdl.K_LEFTBRACKET,
		KeyLeftShift:    sdl.K_LSHIFT,
		KeyLeftControl:  sdl.K_LCTRL,
		KeyLeftSuper:    sdl.K_LGUI,
		KeyLeftAlt:      sdl.K_LALT,
		KeyRightBracket: sdl.K_RIGHTBRACKET,
		KeyRightShift:   sdl.K_RSHIFT,
		KeyRightControl: sdl.K_RCTRL,
		KeyRightSuper:   sdl.K_RGUI,
		KeyRightAlt:     sdl.K_RALT,
		KeyZero:         sdl.K_0,
		KeyOne:          sdl.K_1,
		KeyTwo:          sdl.K_2,
		KeyThree:        sdl.K_3,
		KeyFour:         sdl.K_4,
		KeyFive:         sdl.K_5,
		KeySix:          sdl.K_6,
		KeySeven:        sdl.K_7,
		KeyEight:        sdl.K_8,
		KeyNine:         sdl.K_9,
		KeyF1:           sdl.K_F1,
		KeyF2:           sdl.K_F2,
		KeyF3:           sdl.K_F3,
		KeyF4:           sdl.K_F4,
		KeyF5:           sdl.K_F5,
		KeyF6:           sdl.K_F6,
		KeyF7:           sdl.K_F7,
		KeyF8:           sdl.K_F8,
		KeyF9:           sdl.K_F9,
		KeyF10:          sdl.K_F10,
		KeyF11:          sdl.K_F11,
		KeyF12:          sdl.K_F12,
		KeyA:            sdl.K_a,
		KeyB:            sdl.K_b,
		KeyC:            sdl.K_c,
		KeyD:            sdl.K_d,
		KeyE:            sdl.K_e,
		KeyF:            sdl.K_f,
		KeyG:            sdl.K_g,
		KeyH:            sdl.K_h,
		KeyI:            sdl.K_i,
		KeyJ:            sdl.K_j,
		KeyK:            sdl.K_k,
		KeyL:            sdl.K_l,
		KeyM:            sdl.K_m,
		KeyN:            sdl.K_n,
		KeyO:            sdl.K_o,
		KeyP:            sdl.K_p,
		KeyQ:            sdl.K_q,
		KeyR:            sdl.K_r,
		KeyS:            sdl.K_s,
		KeyT:            sdl.K_t,
		KeyU:            sdl.K_u,
		KeyV:            sdl.K_v,
		KeyW:            sdl.K_w,
		KeyX:            sdl.K_x,
		KeyY:            sdl.K_y,
		KeyZ:            sdl.K_z,
		KeyNumLock:      sdl.K_NUMLOCKCLEAR,
		KeyNumMultiply:  sdl.K_KP_MULTIPLY,
		KeyNumDivide:    sdl.K_KP_DIVIDE,
		KeyNumAdd:       sdl.K_KP_PLUS,
		KeyNumSubtract:  sdl.K_KP_MINUS,
		KeyNumZero:      sdl.K_KP_0,
		KeyNumOne:       sdl.K_KP_1,
		KeyNumTwo:       sdl.K_KP_2,
		KeyNumThree:     sdl.K_KP_3,
		KeyNumFour:      sdl.K_KP_4,
		KeyNumFive:      sdl.K_KP_5,
		KeyNumSix:       sdl.K_KP_6,
		KeyNumSeven:     sdl.K_KP_7,
		KeyNumEight:     sdl.K_KP_8,
		KeyNumNine:      sdl.K_KP_9,
		KeyNumDecimal:   sdl.K_KP_DECIMAL,
		KeyNumEnter:     sdl.K_KP_ENTER,
	}
	return sdl.GetKeyName(m[k])
}
