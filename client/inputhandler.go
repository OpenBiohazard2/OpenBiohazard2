package client

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl64"
)

type Action int

const (
	ACTION_BUTTON         Action = iota
	MENU_UP_BUTTON        Action = iota
	MENU_DOWN_BUTTON      Action = iota
	MENU_LEFT_BUTTON      Action = iota
	MENU_RIGHT_BUTTON     Action = iota
	PLAYER_FORWARD        Action = iota
	PLAYER_BACKWARD       Action = iota
	PLAYER_ROTATE_LEFT    Action = iota
	PLAYER_ROTATE_RIGHT   Action = iota
	PLAYER_VIEW_INVENTORY Action = iota
	DEBUG_DUMP            Action = iota
	PROGRAM_QUIT          Action = iota
)

type InputHandler struct {
	actionToKeyMap map[Action]glfw.Key
	keysPressed    [glfw.KeyLast]bool
	// keysPressedLast is the key state from the previous frame, for edge detection.
	keysPressedLast [glfw.KeyLast]bool

	firstCursorAction    bool
	cursor               mgl64.Vec2
	cursorChange         mgl64.Vec2
	cursorLast           mgl64.Vec2
	bufferedCursorChange mgl64.Vec2
}

func NewInputHandler() *InputHandler {
	actionToKeyMap := map[Action]glfw.Key{
		ACTION_BUTTON:         glfw.KeyEnter,
		MENU_UP_BUTTON:        glfw.KeyUp,
		MENU_DOWN_BUTTON:      glfw.KeyDown,
		MENU_LEFT_BUTTON:      glfw.KeyLeft,
		MENU_RIGHT_BUTTON:     glfw.KeyRight,
		PLAYER_FORWARD:        glfw.KeyW,
		PLAYER_BACKWARD:       glfw.KeyS,
		PLAYER_ROTATE_LEFT:    glfw.KeyA,
		PLAYER_ROTATE_RIGHT:   glfw.KeyD,
		PLAYER_VIEW_INVENTORY: glfw.KeyTab,
		DEBUG_DUMP:            glfw.KeyBackslash,
		PROGRAM_QUIT:          glfw.KeyEscape,
	}

	return &InputHandler{
		actionToKeyMap:    actionToKeyMap,
		firstCursorAction: false,
	}
}

// IsActive returns true while the mapped key is held down. Use for continuous actions
// like movement.
func (handler *InputHandler) IsActive(a Action) bool {
	return handler.keysPressed[handler.actionToKeyMap[a]]
}

// IsActiveOnce returns true only on the frame the mapped key is first pressed. Use for
// discrete actions like menu navigation or confirming a selection.
func (handler *InputHandler) IsActiveOnce(a Action) bool {
	key := handler.actionToKeyMap[a]
	return handler.keysPressed[key] && !handler.keysPressedLast[key]
}

// snapshotKeyState records the current key state as last frame's state. Call once per
// frame, before polling new events.
func (handler *InputHandler) snapshotKeyState() {
	handler.keysPressedLast = handler.keysPressed
}

func (handler *InputHandler) keyCallback(window *glfw.Window, key glfw.Key, scancode int,
	action glfw.Action, mods glfw.ModifierKey) {

	switch action {
	case glfw.Press:
		handler.keysPressed[key] = true
	case glfw.Release:
		handler.keysPressed[key] = false
	}
}

func (handler *InputHandler) getCursorChange() mgl64.Vec2 {
	return handler.cursorChange
}

func (handler *InputHandler) updateCursor() {
	handler.cursorChange[0] = handler.bufferedCursorChange[0]
	handler.cursorChange[1] = handler.bufferedCursorChange[1]
	handler.cursor[0] = handler.cursorLast[0]
	handler.cursor[1] = handler.cursorLast[1]

	handler.bufferedCursorChange[0] = 0
	handler.bufferedCursorChange[1] = 0
}

func (handler *InputHandler) mouseCallback(window *glfw.Window, xPos float64, yPos float64) {
	if handler.firstCursorAction {
		handler.cursorLast[0] = xPos
		handler.cursorLast[1] = yPos
		handler.firstCursorAction = false
	}

	handler.bufferedCursorChange[0] += xPos - handler.cursorLast[0]
	handler.bufferedCursorChange[1] += handler.cursorLast[1] - yPos

	handler.cursorLast[0] = xPos
	handler.cursorLast[1] = yPos
}
