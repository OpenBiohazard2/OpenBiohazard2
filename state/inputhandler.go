package state

import (
	"github.com/OpenBiohazard2/OpenBiohazard2/client"
	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/OpenBiohazard2/OpenBiohazard2/game"
	"github.com/OpenBiohazard2/OpenBiohazard2/world"
)

type InputHandler struct {
	windowHandler    *client.WindowHandler
	gameStateManager *GameStateManager
}

func NewInputHandler(windowHandler *client.WindowHandler, gameStateManager *GameStateManager) *InputHandler {
	return &InputHandler{
		windowHandler:    windowHandler,
		gameStateManager: gameStateManager,
	}
}

func (h *InputHandler) HandleTankMovement(gameDef *game.GameDef, timeElapsedSeconds float64, collisionEntities []fileio.CollisionEntity) {
	if h.windowHandler.InputHandler.IsActive(client.PLAYER_FORWARD) {
		gameDef.Player.HandlePlayerInputForward(collisionEntities, timeElapsedSeconds)
	}

	if h.windowHandler.InputHandler.IsActive(client.PLAYER_BACKWARD) {
		gameDef.Player.HandlePlayerInputBackward(collisionEntities, timeElapsedSeconds)
	}

	if !h.windowHandler.InputHandler.IsActive(client.PLAYER_FORWARD) &&
		!h.windowHandler.InputHandler.IsActive(client.PLAYER_BACKWARD) {
		gameDef.Player.PoseNumber = -1
	}
}

func (h *InputHandler) HandleTankRotation(gameDef *game.GameDef, timeElapsedSeconds float64) {
	if h.windowHandler.InputHandler.IsActive(client.PLAYER_ROTATE_LEFT) {
		gameDef.Player.RotatePlayerLeft(timeElapsedSeconds)
	}

	if h.windowHandler.InputHandler.IsActive(client.PLAYER_ROTATE_RIGHT) {
		gameDef.Player.RotatePlayerRight(timeElapsedSeconds)
	}
}

func (h *InputHandler) HandleActionButton(gameDef *game.GameDef, collisionEntities []fileio.CollisionEntity) {
	if h.windowHandler.InputHandler.IsActive(client.ACTION_BUTTON) {
		if h.gameStateManager.CanUpdateGameState(h.windowHandler) {
			gameDef.HandlePlayerActionButton(collisionEntities)
			h.gameStateManager.UpdateLastTimeChangeState(h.windowHandler)
		}
	}
}

func (h *InputHandler) HandleInventoryToggle() {
	if h.windowHandler.InputHandler.IsActive(client.PLAYER_VIEW_INVENTORY) {
		if h.gameStateManager.CanUpdateGameState(h.windowHandler) {
			h.gameStateManager.UpdateGameState(GAME_STATE_INVENTORY)
			h.gameStateManager.UpdateLastTimeChangeState(h.windowHandler)
		}
	}
}

func (h *InputHandler) HandleAllInput(gameDef *game.GameDef, timeElapsedSeconds float64, gameWorld *world.GameWorld) {
	collisionEntities := gameWorld.GameRoom.CollisionEntities

	h.HandleTankMovement(gameDef, timeElapsedSeconds, collisionEntities)
	h.HandleTankRotation(gameDef, timeElapsedSeconds)
	h.HandleActionButton(gameDef, collisionEntities)
	h.HandleInventoryToggle()
}
