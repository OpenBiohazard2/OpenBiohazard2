package client

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type WindowHandler struct {
	glfwWindow   *glfw.Window
	InputHandler *InputHandler

	firstFrame    bool
	deltaTime     float64
	lastFrameTime float64
}

func NewWindowHandler(width, height int, title string) *WindowHandler {
	// Initialize and create window
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	glfwWindow, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(fmt.Errorf("Could not create OpenGL renderer: %v", err))
	}
	glfwWindow.MakeContextCurrent()

	// Check for resize
	glfwWindow.SetSizeCallback(resizeCallback)
	glfwWindow.GetSize()

	inputHandler := NewInputHandler()

	// Keyboard callback
	glfwWindow.SetKeyCallback(inputHandler.keyCallback)
	// Mouse callback
	glfwWindow.SetCursorPosCallback(inputHandler.mouseCallback)

	return &WindowHandler{
		glfwWindow:   glfwWindow,
		InputHandler: inputHandler,
		firstFrame:   true,
	}
}

// Resize the screen
func resizeCallback(w *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

func (windowHandler *WindowHandler) StartFrame() {
	windowHandler.glfwWindow.SwapBuffers()

	// Window events for keyboard and mouse
	glfw.PollEvents()

	if windowHandler.InputHandler.IsActive(PROGRAM_QUIT) {
		windowHandler.glfwWindow.SetShouldClose(true)
	}

	// Set frame time
	currentFrameTime := glfw.GetTime()

	if windowHandler.firstFrame {
		windowHandler.lastFrameTime = currentFrameTime
		windowHandler.firstFrame = false
	}

	windowHandler.deltaTime = currentFrameTime - windowHandler.lastFrameTime
	windowHandler.lastFrameTime = currentFrameTime

	windowHandler.InputHandler.updateCursor()
}

func (windowHandler *WindowHandler) ShouldClose() bool {
	return windowHandler.glfwWindow.ShouldClose()
}

func (windowHandler *WindowHandler) GetTimeSinceLastFrame() float64 {
	return windowHandler.deltaTime
}

func (windowHandler *WindowHandler) GetCurrentTime() float64 {
	return glfw.GetTime()
}
