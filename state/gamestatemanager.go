package state

const (
	GAME_STATE_MAIN_MENU    = 0
	GAME_STATE_MAIN_GAME    = 1
	GAME_STATE_INVENTORY    = 2
	GAME_STATE_LOAD_SAVE    = 3
	GAME_STATE_SPECIAL_MENU = 4
)

type GameStateManager struct {
	GameState            int
	ImageResourcesLoaded bool
}

func NewGameStateManager() *GameStateManager {
	return &GameStateManager{
		GameState:            GAME_STATE_MAIN_MENU,
		ImageResourcesLoaded: false,
	}
}

func (gameStateManager *GameStateManager) UpdateGameState(newGameState int) {
	gameStateManager.GameState = newGameState
	gameStateManager.ImageResourcesLoaded = false
}
