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
	PLAYER_FORWARD        Action = iota
	PLAYER_BACKWARD       Action = iota
	PLAYER_ROTATE_LEFT    Action = iota
	PLAYER_ROTATE_RIGHT   Action = iota
	PLAYER_VIEW_INVENTORY Action = iota
	PROGRAM_QUIT          Action = iota
)

type InputHandler struct {
	actionToKeyMap  map[Action]glfw.Key
	keysPressed     [glfw.KeyLast]bool
	keysPressedOnce [glfw.KeyLast]bool

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
		PLAYER_FORWARD:        glfw.KeyW,
		PLAYER_BACKWARD:       glfw.KeyS,
		PLAYER_ROTATE_LEFT:    glfw.KeyA,
		PLAYER_ROTATE_RIGHT:   glfw.KeyD,
		PLAYER_VIEW_INVENTORY: glfw.KeyTab,
		PROGRAM_QUIT:          glfw.KeyEscape,
	}

	return &InputHandler{
		actionToKeyMap:    actionToKeyMap,
		firstCursorAction: false,
	}
}

func (handler *InputHandler) IsActive(a Action) bool {
	return handler.keysPressed[handler.actionToKeyMap[a]]
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
