package state

import (
	"github.com/OpenBiohazard2/OpenBiohazard2/client"
)

const (
	GAME_STATE_MAIN_MENU    = 0
	GAME_STATE_MAIN_GAME    = 1
	GAME_STATE_INVENTORY    = 2
	GAME_STATE_LOAD_SAVE    = 3
	GAME_STATE_SPECIAL_MENU = 4

	STATE_CHANGE_DELAY = 0.2 // in seconds
)

type GameStateManager struct {
	GameState            int
	ImageResourcesLoaded bool
	LastTimeChangeState  float64
}

func NewGameStateManager() *GameStateManager {
	return &GameStateManager{
		GameState:            GAME_STATE_MAIN_MENU,
		ImageResourcesLoaded: false,
		LastTimeChangeState:  0, // Will be set when windowHandler is available
	}
}

func (gameStateManager *GameStateManager) CanUpdateGameState(windowHandler *client.WindowHandler) bool {
	return windowHandler.GetCurrentTime()-gameStateManager.LastTimeChangeState >= STATE_CHANGE_DELAY
}

func (gameStateManager *GameStateManager) UpdateGameState(newGameState int) {
	gameStateManager.GameState = newGameState
	gameStateManager.ImageResourcesLoaded = false
}

func (gameStateManager *GameStateManager) UpdateLastTimeChangeState(windowHandler *client.WindowHandler) {
	gameStateManager.LastTimeChangeState = windowHandler.GetCurrentTime()
}
