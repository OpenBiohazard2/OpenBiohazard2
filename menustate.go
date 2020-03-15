package main

import (
	"./client"
	"./fileio"
	"./game"
	"./render"
)

const (
	GAME_STATE_MAIN_MENU    = 0
	GAME_STATE_MAIN_GAME    = 1
	GAME_STATE_INVENTORY    = 2
	GAME_STATE_LOAD_SAVE    = 3
	GAME_STATE_SPECIAL_MENU = 4

	STATE_CHANGE_DELAY = 0.5 // in seconds
)

type GameStateManager struct {
	GameState            int
	MainMenuOption       int
	SpecialMenuOption    int
	ImageResourcesLoaded bool
	LastTimeChangeState  float64
}

type MainMenuStateInput struct {
	RenderDef                 *render.RenderDef
	MenuBackgroundImageOutput *fileio.ADTOutput
	MenuBackgroundTextImages  []*fileio.TIMOutput
}

type InventoryStateInput struct {
	RenderDef           *render.RenderDef
	InventoryImages     []*fileio.TIMOutput
	InventoryItemImages []*fileio.TIMOutput
}

func NewGameStateManager() *GameStateManager {
	return &GameStateManager{
		GameState:            GAME_STATE_MAIN_MENU,
		MainMenuOption:       0,
		SpecialMenuOption:    0,
		ImageResourcesLoaded: false,
		LastTimeChangeState:  windowHandler.GetCurrentTime(),
	}
}

func (gameStateManager *GameStateManager) CanUpdateGameState() bool {
	return windowHandler.GetCurrentTime()-gameStateManager.LastTimeChangeState >= STATE_CHANGE_DELAY
}

func (gameStateManager *GameStateManager) UpdateGameState(newGameState int) {
	gameStateManager.GameState = newGameState
	gameStateManager.ImageResourcesLoaded = false
}

func (gameStateManager *GameStateManager) UpdateLastTimeChangeState() {
	gameStateManager.LastTimeChangeState = windowHandler.GetCurrentTime()
}

func NewInventoryStateInput(renderDef *render.RenderDef) *InventoryStateInput {
	inventoryImages, _ := fileio.LoadTIMImages(game.INVENTORY_FILE)
	inventoryItemImages, _ := fileio.LoadTIMImages(game.ITEMALL_FILE)
	return &InventoryStateInput{
		RenderDef:           renderDef,
		InventoryImages:     inventoryImages,
		InventoryItemImages: inventoryItemImages,
	}
}

func handleInventory(inventoryStateInput *InventoryStateInput, gameStateManager *GameStateManager) {
	renderDef := inventoryStateInput.RenderDef
	inventoryImages := inventoryStateInput.InventoryImages
	inventoryItemImages := inventoryStateInput.InventoryItemImages

	if gameStateManager.ImageResourcesLoaded == false {
		gameStateManager.ImageResourcesLoaded = true
		gameStateManager.UpdateLastTimeChangeState()
	}

	if windowHandler.InputHandler.IsActive(client.PLAYER_VIEW_INVENTORY) {
		if gameStateManager.CanUpdateGameState() {
			gameStateManager.UpdateGameState(GAME_STATE_MAIN_GAME)
			gameStateManager.UpdateLastTimeChangeState()
		}
	}

	timeElapsedSeconds := windowHandler.GetTimeSinceLastFrame()
	renderDef.GenerateInventoryImage(inventoryImages, inventoryItemImages, timeElapsedSeconds)
	renderDef.RenderSolidVideoBuffer()
}

func handleMainMenu(mainMenuStateInput *MainMenuStateInput, gameStateManager *GameStateManager) {
	maxOptions := 4
	renderDef := mainMenuStateInput.RenderDef
	if gameStateManager.ImageResourcesLoaded == false {
		menuBackgroundImageOutput := fileio.LoadADTFile(game.MENU_IMAGE_FILE)
		menuBackgroundTextImages, _ := fileio.LoadTIMImages(game.MENU_TEXT_FILE)
		renderDef.GenerateMainMenuImage(menuBackgroundImageOutput, menuBackgroundTextImages)

		mainMenuStateInput.MenuBackgroundImageOutput = menuBackgroundImageOutput
		mainMenuStateInput.MenuBackgroundTextImages = menuBackgroundTextImages
		gameStateManager.ImageResourcesLoaded = true
		gameStateManager.MainMenuOption = 0
		gameStateManager.UpdateLastTimeChangeState()
	}

	renderDef.RenderTransparentVideoBuffer()

	if windowHandler.InputHandler.IsActive(client.ACTION_BUTTON) {
		if gameStateManager.CanUpdateGameState() {
			if gameStateManager.MainMenuOption == 0 {
				gameStateManager.UpdateGameState(GAME_STATE_LOAD_SAVE)
				gameStateManager.UpdateLastTimeChangeState()
			} else if gameStateManager.MainMenuOption == 1 {
				gameStateManager.UpdateGameState(GAME_STATE_MAIN_GAME)
				gameStateManager.UpdateLastTimeChangeState()
			} else if gameStateManager.MainMenuOption == 2 {
				gameStateManager.UpdateGameState(GAME_STATE_SPECIAL_MENU)
				gameStateManager.UpdateLastTimeChangeState()
			}
		}
	}

	if windowHandler.InputHandler.IsActive(client.MENU_UP_BUTTON) {
		if gameStateManager.CanUpdateGameState() {
			gameStateManager.MainMenuOption--
			if gameStateManager.MainMenuOption < 0 {
				gameStateManager.MainMenuOption = 0
			}
			renderDef.UpdateMainMenu(mainMenuStateInput.MenuBackgroundImageOutput, mainMenuStateInput.MenuBackgroundTextImages, gameStateManager.MainMenuOption)
			gameStateManager.UpdateLastTimeChangeState()
		}
	}

	if windowHandler.InputHandler.IsActive(client.MENU_DOWN_BUTTON) {
		if gameStateManager.CanUpdateGameState() {
			gameStateManager.MainMenuOption++
			if gameStateManager.MainMenuOption >= maxOptions {
				gameStateManager.MainMenuOption = maxOptions - 1
			}
			renderDef.UpdateMainMenu(mainMenuStateInput.MenuBackgroundImageOutput, mainMenuStateInput.MenuBackgroundTextImages, gameStateManager.MainMenuOption)
			gameStateManager.UpdateLastTimeChangeState()
		}
	}
}

func handleLoadSave(renderDef *render.RenderDef, gameStateManager *GameStateManager) {
	if gameStateManager.ImageResourcesLoaded == false {
		// Initialize load save screen
		saveScreenImage := fileio.LoadADTFile(game.SAVE_SCREEN_FILE)
		renderDef.GenerateSaveScreenImage(saveScreenImage)

		gameStateManager.ImageResourcesLoaded = true
		gameStateManager.UpdateLastTimeChangeState()
	}

	renderDef.RenderTransparentVideoBuffer()
	if windowHandler.InputHandler.IsActive(client.ACTION_BUTTON) {
		if gameStateManager.CanUpdateGameState() {
			gameStateManager.UpdateGameState(GAME_STATE_MAIN_MENU)
			gameStateManager.UpdateLastTimeChangeState()
		}
	}
}

func handleSpecialMenu(mainMenuStateInput *MainMenuStateInput, gameStateManager *GameStateManager) {
	maxOptions := 2
	renderDef := mainMenuStateInput.RenderDef
	if gameStateManager.ImageResourcesLoaded == false {
		menuBackgroundImageOutput := fileio.LoadADTFile(game.MENU_IMAGE_FILE)
		menuBackgroundTextImages, _ := fileio.LoadTIMImages(game.MENU_TEXT_FILE)
		renderDef.GenerateSpecialMenuImage(menuBackgroundImageOutput, menuBackgroundTextImages)

		mainMenuStateInput.MenuBackgroundImageOutput = menuBackgroundImageOutput
		mainMenuStateInput.MenuBackgroundTextImages = menuBackgroundTextImages
		gameStateManager.ImageResourcesLoaded = true
		gameStateManager.SpecialMenuOption = 0
		gameStateManager.UpdateLastTimeChangeState()
	}

	renderDef.RenderTransparentVideoBuffer()
	if windowHandler.InputHandler.IsActive(client.ACTION_BUTTON) {
		if gameStateManager.CanUpdateGameState() {
			if gameStateManager.SpecialMenuOption == 0 {

			} else if gameStateManager.SpecialMenuOption == 1 {
				// Exit
				gameStateManager.UpdateGameState(GAME_STATE_MAIN_MENU)
				gameStateManager.UpdateLastTimeChangeState()
			}
		}
	}

	if windowHandler.InputHandler.IsActive(client.MENU_UP_BUTTON) {
		if gameStateManager.CanUpdateGameState() {
			gameStateManager.SpecialMenuOption--
			if gameStateManager.SpecialMenuOption < 0 {
				gameStateManager.SpecialMenuOption = 0
			}
			renderDef.UpdateSpecialMenu(mainMenuStateInput.MenuBackgroundImageOutput, mainMenuStateInput.MenuBackgroundTextImages, gameStateManager.SpecialMenuOption)
			gameStateManager.UpdateLastTimeChangeState()
		}
	}

	if windowHandler.InputHandler.IsActive(client.MENU_DOWN_BUTTON) {
		if gameStateManager.CanUpdateGameState() {
			gameStateManager.SpecialMenuOption++
			if gameStateManager.SpecialMenuOption >= maxOptions {
				gameStateManager.SpecialMenuOption = maxOptions - 1
			}
			renderDef.UpdateSpecialMenu(mainMenuStateInput.MenuBackgroundImageOutput, mainMenuStateInput.MenuBackgroundTextImages, gameStateManager.SpecialMenuOption)
			gameStateManager.UpdateLastTimeChangeState()
		}
	}
}
