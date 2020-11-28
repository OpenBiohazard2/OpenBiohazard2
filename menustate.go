package main

import (
	"github.com/samuelyuan/openbiohazard2/client"
	"github.com/samuelyuan/openbiohazard2/fileio"
	"github.com/samuelyuan/openbiohazard2/game"
	"github.com/samuelyuan/openbiohazard2/render"
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
	RenderDef           *render.RenderDef
	MenuBackgroundImage *render.Image16Bit
	MenuTextImages      []*render.Image16Bit
}

type InventoryStateInput struct {
	RenderDef           *render.RenderDef
	InventoryMenuImages []*render.Image16Bit
	InventoryItemImages []*render.Image16Bit
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
	inventoryMenuImagesTIMOutput, _ := fileio.LoadTIMImages(game.INVENTORY_FILE)
	inventoryMenuImages := make([]*render.Image16Bit, len(inventoryMenuImagesTIMOutput))
	for i := 0; i < len(inventoryMenuImages); i++ {
		inventoryMenuImages[i] = render.ConvertPixelsToImage16Bit(inventoryMenuImagesTIMOutput[i].PixelData)
	}

	inventoryItemImagesTIMOutput, _ := fileio.LoadTIMImages(game.ITEMALL_FILE)
	inventoryItemImages := make([]*render.Image16Bit, len(inventoryItemImagesTIMOutput))
	for i := 0; i < len(inventoryItemImages); i++ {
		inventoryItemImages[i] = render.ConvertPixelsToImage16Bit(inventoryItemImagesTIMOutput[i].PixelData)
	}
	return &InventoryStateInput{
		RenderDef:           renderDef,
		InventoryMenuImages: inventoryMenuImages,
		InventoryItemImages: inventoryItemImages,
	}
}

func handleInventory(inventoryStateInput *InventoryStateInput, gameStateManager *GameStateManager) {
	renderDef := inventoryStateInput.RenderDef
	inventoryMenuImages := inventoryStateInput.InventoryMenuImages
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

	if windowHandler.InputHandler.IsActive(client.MENU_LEFT_BUTTON) {
		if gameStateManager.CanUpdateGameState() {
			if render.IsCursorOnTopMenu() {
				render.PrevTopMenuOption()
			} else if render.IsEditingItemScreen() {
				render.PrevItemInList()
			}
			gameStateManager.UpdateLastTimeChangeState()
		}
	} else if windowHandler.InputHandler.IsActive(client.MENU_RIGHT_BUTTON) {
		if gameStateManager.CanUpdateGameState() {
			if render.IsCursorOnTopMenu() {
				render.NextTopMenuOption()
			} else if render.IsEditingItemScreen() {
				render.NextItemInList()
			}
			gameStateManager.UpdateLastTimeChangeState()
		}
	} else if windowHandler.InputHandler.IsActive(client.MENU_UP_BUTTON) {
		if gameStateManager.CanUpdateGameState() {
			if render.IsEditingItemScreen() {
				render.PrevRowInItemList()
			}
			gameStateManager.UpdateLastTimeChangeState()
		}
	} else if windowHandler.InputHandler.IsActive(client.MENU_DOWN_BUTTON) {
		if gameStateManager.CanUpdateGameState() {
			if render.IsEditingItemScreen() {
				render.NextRowInItemList()
			}
			gameStateManager.UpdateLastTimeChangeState()
		}
	} else if windowHandler.InputHandler.IsActive(client.ACTION_BUTTON) {
		if gameStateManager.CanUpdateGameState() {
			if render.IsCursorOnTopMenu() {
				if render.IsTopMenuExit() {
					gameStateManager.UpdateGameState(GAME_STATE_MAIN_GAME)
				} else if render.IsTopMenuCursorOnItems() {
					render.SetEditItemScreen()
				}
			}
			gameStateManager.UpdateLastTimeChangeState()
		}
	}

	timeElapsedSeconds := windowHandler.GetTimeSinceLastFrame()
	renderDef.GenerateInventoryImage(inventoryMenuImages, inventoryItemImages, timeElapsedSeconds)
	renderDef.RenderSolidVideoBuffer()
}

func handleMainMenu(mainMenuStateInput *MainMenuStateInput, gameStateManager *GameStateManager) {
	maxOptions := 4
	renderDef := mainMenuStateInput.RenderDef
	if gameStateManager.ImageResourcesLoaded == false {
		menuBackgroundImageADTOutput := fileio.LoadADTFile(game.MENU_IMAGE_FILE)
		menuBackgroundImage := render.ConvertPixelsToImage16Bit(menuBackgroundImageADTOutput.PixelData)

		menuBackgroundTextImagesTIMOutput, _ := fileio.LoadTIMImages(game.MENU_TEXT_FILE)
		menuTextImages := make([]*render.Image16Bit, len(menuBackgroundTextImagesTIMOutput))
		for i := 0; i < len(menuBackgroundTextImagesTIMOutput); i++ {
			menuTextImages[i] = render.ConvertPixelsToImage16Bit(menuBackgroundTextImagesTIMOutput[i].PixelData)
		}

		renderDef.UpdateMainMenu(menuBackgroundImage, menuTextImages, 0)

		mainMenuStateInput.MenuBackgroundImage = menuBackgroundImage
		mainMenuStateInput.MenuTextImages = menuTextImages
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
			renderDef.UpdateMainMenu(mainMenuStateInput.MenuBackgroundImage, mainMenuStateInput.MenuTextImages, gameStateManager.MainMenuOption)
			gameStateManager.UpdateLastTimeChangeState()
		}
	}

	if windowHandler.InputHandler.IsActive(client.MENU_DOWN_BUTTON) {
		if gameStateManager.CanUpdateGameState() {
			gameStateManager.MainMenuOption++
			if gameStateManager.MainMenuOption >= maxOptions {
				gameStateManager.MainMenuOption = maxOptions - 1
			}
			renderDef.UpdateMainMenu(mainMenuStateInput.MenuBackgroundImage, mainMenuStateInput.MenuTextImages, gameStateManager.MainMenuOption)
			gameStateManager.UpdateLastTimeChangeState()
		}
	}
}

func handleLoadSave(renderDef *render.RenderDef, gameStateManager *GameStateManager) {
	if gameStateManager.ImageResourcesLoaded == false {
		// Initialize load save screen
		saveScreenImageADTOutput := fileio.LoadADTFile(game.SAVE_SCREEN_FILE)
		saveScreenImageRender := render.ConvertPixelsToImage16Bit(saveScreenImageADTOutput.PixelData)
		renderDef.GenerateSaveScreenImage(saveScreenImageRender)

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
		menuBackgroundImageADTOutput := fileio.LoadADTFile(game.MENU_IMAGE_FILE)
		menuBackgroundImage := render.ConvertPixelsToImage16Bit(menuBackgroundImageADTOutput.PixelData)

		menuBackgroundTextImagesTIMOutput, _ := fileio.LoadTIMImages(game.MENU_TEXT_FILE)
		menuTextImages := make([]*render.Image16Bit, len(menuBackgroundTextImagesTIMOutput))
		for i := 0; i < len(menuBackgroundTextImagesTIMOutput); i++ {
			menuTextImages[i] = render.ConvertPixelsToImage16Bit(menuBackgroundTextImagesTIMOutput[i].PixelData)
		}

		renderDef.UpdateSpecialMenu(menuBackgroundImage, menuTextImages, 0)

		mainMenuStateInput.MenuBackgroundImage = menuBackgroundImage
		mainMenuStateInput.MenuTextImages = menuTextImages
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
			renderDef.UpdateSpecialMenu(mainMenuStateInput.MenuBackgroundImage, mainMenuStateInput.MenuTextImages, gameStateManager.SpecialMenuOption)
			gameStateManager.UpdateLastTimeChangeState()
		}
	}

	if windowHandler.InputHandler.IsActive(client.MENU_DOWN_BUTTON) {
		if gameStateManager.CanUpdateGameState() {
			gameStateManager.SpecialMenuOption++
			if gameStateManager.SpecialMenuOption >= maxOptions {
				gameStateManager.SpecialMenuOption = maxOptions - 1
			}
			renderDef.UpdateSpecialMenu(mainMenuStateInput.MenuBackgroundImage, mainMenuStateInput.MenuTextImages, gameStateManager.SpecialMenuOption)
			gameStateManager.UpdateLastTimeChangeState()
		}
	}
}
