package engine

import (
	"fmt"
	"time"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type CursorPosHandler func(float64, float64, float64, float64) bool
type MouseButtonHandler func(bool, bool) bool
type MouseScrollHandler func(float64, float64) bool
type KeyPressHandler func(int, int, int) bool

type WindowManager struct {
	Window *glfw.Window
	Width  int
	Height int

	fpsLock float64
	lastFps float64

	cursorPosHandlers   []CursorPosHandler
	mouseButtonHandlers []MouseButtonHandler
	mouseScrollHandlers []MouseScrollHandler
	keyPressHandlers    []KeyPressHandler

	prevPosX, prevPosY float64
	posInit            bool
	leftPressed        bool
	rightPressed       bool

	loopCursor bool
}

func NewWindowManager(title string, width, height int) (*WindowManager, error) {
	// init glfw
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	// set glfw window hints
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	//glfw.WindowHint(glfw.Samples, 4)

	// create window
	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, err
	}
	// actually creating the OpenGL context
	window.MakeContextCurrent()

	// init OpenGL
	if err := gl.Init(); err != nil {
		return nil, err
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// set clear color
	gl.Enable(gl.DEPTH_TEST)
	gl.FrontFace(gl.CCW)
	gl.CullFace(gl.BACK)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	// set default values
	windowManager := WindowManager{
		Window:  window,
		Width:   width,
		Height:  height,
		fpsLock: -1.0,

		prevPosX:     0.0,
		prevPosY:     0.0,
		posInit:      false,
		leftPressed:  false,
		rightPressed: false,

		loopCursor: false,
	}

	// add handlers
	windowManager.Window.SetCursorPosCallback(windowManager.onCursorPos)
	windowManager.Window.SetCursorEnterCallback(windowManager.onCursorEnter)
	windowManager.Window.SetMouseButtonCallback(windowManager.onMouseButton)
	windowManager.Window.SetScrollCallback(windowManager.onMouseScroll)
	windowManager.Window.SetKeyCallback(windowManager.onKeyPress)

	return &windowManager, nil
}
func (windowManager *WindowManager) Close() {
	glfw.Terminate()
}
func (windowManager *WindowManager) RunMainLoop(render func()) {
	for !windowManager.Window.ShouldClose() {
		// set frame start
		frameStart := time.Now()
		// reset gl states
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// render user defined function
		render()
		// swap front with back buffer
		windowManager.Window.SwapBuffers()
		// get inputs
		glfw.PollEvents()
		// get the time after the rendering
		frameEnd := time.Now()

		// frame lock if specified
		deltaTime := frameEnd.Sub(frameStart).Seconds() * 1000.0
		timeToWait := (1000.0 / windowManager.fpsLock) - deltaTime
		if timeToWait > 0.0 && windowManager.fpsLock > 0.0 {
			time.Sleep(time.Duration(timeToWait/1000) * time.Second)
			deltaTime = deltaTime + timeToWait
		}
		windowManager.lastFps = 1000.0 / deltaTime
	}
}

func (windowManager *WindowManager) LockFPS(fps float64) {
	windowManager.fpsLock = fps
}
func (windowManager *WindowManager) GetFPS() float64 {
	return windowManager.lastFps
}

func (windowManager *WindowManager) EnableCursorLoop() {
	windowManager.loopCursor = true
	windowManager.Window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
}

func (windowManager *WindowManager) SetTitle(title string) {
	windowManager.Window.SetTitle(title)
}
func (windowManager *WindowManager) SetClearColor(r, g, b float32) {
	gl.ClearColor(r, g, b, 1.0)
}

func (windowManager *WindowManager) onCursorPos(w *glfw.Window, x float64, y float64) {
	if !windowManager.posInit {
		windowManager.posInit = true
		windowManager.prevPosX = x
		windowManager.prevPosY = y
	}
	deltaX := x - windowManager.prevPosX
	deltaY := y - windowManager.prevPosY
	for _, handler := range windowManager.cursorPosHandlers {
		if handler(x, y, deltaX, deltaY) {
			break
		}
	}
	windowManager.prevPosX = x
	windowManager.prevPosY = y
}
func (windowManager *WindowManager) onMouseButton(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	// save pressed button
	if button == glfw.MouseButtonLeft {
		if action == glfw.Press {
			windowManager.leftPressed = true
		} else if action == glfw.Release {
			windowManager.leftPressed = false
		}
	} else if button == glfw.MouseButtonRight {
		if action == glfw.Press {
			windowManager.rightPressed = true
		} else if action == glfw.Release {
			windowManager.rightPressed = false
		}
	}

	// inform all handlers
	for _, handler := range windowManager.mouseButtonHandlers {
		if handler(windowManager.leftPressed, windowManager.rightPressed) {
			break
		}
	}
}
func (windowManager *WindowManager) onMouseScroll(w *glfw.Window, x float64, y float64) {
	for _, handler := range windowManager.mouseScrollHandlers {
		if handler(x, y) {
			break
		}
	}
}
func (windowManager *WindowManager) onCursorEnter(w *glfw.Window, entered bool) {
	if !entered {
		windowManager.posInit = false

		// loop
		if windowManager.loopCursor {
			x := windowManager.prevPosX
			y := windowManager.prevPosY
			w := float64(windowManager.Width)
			h := float64(windowManager.Height)
			var border float64 = 20

			if x < border {
				x = w - 1
			} else if x > w-border {
				x = 1
			}

			if y < border {
				y = h - 1
			} else if y > h-border {
				y = 1
			}

			windowManager.Window.SetCursorPos(x, y)
		}
	}
}
func (windowManager *WindowManager) onKeyPress(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	for _, handler := range windowManager.keyPressHandlers {
		if handler(int(key), int(action), int(mods)) {
			break
		}
	}
}

type Interactable interface {
	OnCursorPosMove(x, y, dx, dy float64) bool
	OnMouseButtonPress(leftPressed, rightPressed bool) bool
	OnMouseScroll(x, y float64) bool
	OnKeyPress(key, action, mods int) bool
}

func (windowManager *WindowManager) AddInteractable(interactable Interactable) {
	windowManager.AddCursorPosHandler(interactable.OnCursorPosMove)
	windowManager.AddMouseButtonHandler(interactable.OnMouseButtonPress)
	windowManager.AddMouseScrollHandler(interactable.OnMouseScroll)
	windowManager.AddKeyPressHandler(interactable.OnKeyPress)
}
func (windowManager *WindowManager) AddCursorPosHandler(handler CursorPosHandler) {
	windowManager.cursorPosHandlers = append(windowManager.cursorPosHandlers, handler)
}
func (windowManager *WindowManager) AddMouseButtonHandler(handler MouseButtonHandler) {
	windowManager.mouseButtonHandlers = append(windowManager.mouseButtonHandlers, handler)
}
func (windowManager *WindowManager) AddMouseScrollHandler(handler MouseScrollHandler) {
	windowManager.mouseScrollHandlers = append(windowManager.mouseScrollHandlers, handler)
}
func (windowManager *WindowManager) AddKeyPressHandler(handler KeyPressHandler) {
	windowManager.keyPressHandlers = append(windowManager.keyPressHandlers, handler)
}
