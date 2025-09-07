package main

import (
	"log"

	"github.com/OpenBiohazard2/OpenBiohazard2/client"
	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/OpenBiohazard2/OpenBiohazard2/game"
	"github.com/OpenBiohazard2/OpenBiohazard2/gui"
	"github.com/OpenBiohazard2/OpenBiohazard2/render"
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

type MainMenuStateInput struct {
	RenderDef           *render.RenderDef
	MenuBackgroundImage *render.Image16Bit
	MenuTextImages      []*render.Image16Bit
	Menu                *gui.Menu
}

type InventoryStateInput struct {
	RenderDef           *render.RenderDef
	InventoryMenuImages []*render.Image16Bit
	InventoryItemImages []*render.Image16Bit
	InventoryMenu       *gui.InventoryMenu
}

func NewGameStateManager() *GameStateManager {
	return &GameStateManager{
		GameState:            GAME_STATE_MAIN_MENU,
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
		InventoryMenu:       gui.NewInventoryMenu(),
	}
}

func handleInventory(inventoryStateInput *InventoryStateInput, gameStateManager *GameStateManager) {
	renderDef := inventoryStateInput.RenderDef
	inventoryMenuImages := inventoryStateInput.InventoryMenuImages
	inventoryItemImages := inventoryStateInput.InventoryItemImages
	inventoryMenu := inventoryStateInput.InventoryMenu

	if !gameStateManager.ImageResourcesLoaded {
		inventoryMenu.Reset()
		gameStateManager.ImageResourcesLoaded = true
		gameStateManager.UpdateLastTimeChangeState()
	}

	if windowHandler.InputHandler.IsActive(client.PLAYER_VIEW_INVENTORY) {
		if gameStateManager.CanUpdateGameState() {
			gameStateManager.UpdateGameState(GAME_STATE_MAIN_GAME)
			gameStateManager.UpdateLastTimeChangeState()
		}
	}

	if windowHandler.InputHandler.IsActive(client.ACTION_BUTTON) {
		if gameStateManager.CanUpdateGameState() {
			if inventoryMenu.IsCursorOnTopMenu() {
				if inventoryMenu.IsTopMenuExit() {
					gameStateManager.UpdateGameState(GAME_STATE_MAIN_GAME)
				} else if inventoryMenu.IsTopMenuCursorOnItems() {
					inventoryMenu.SetEditItemScreen()
				}
			}
			gameStateManager.UpdateLastTimeChangeState()
		}
	}

	if gameStateManager.CanUpdateGameState() {
		inventoryMenu.HandleSwitchMenuOption(windowHandler)
		gameStateManager.UpdateLastTimeChangeState()
	}

	timeElapsedSeconds := windowHandler.GetTimeSinceLastFrame()
	renderDef.GenerateInventoryImage(inventoryMenuImages, inventoryItemImages, inventoryMenu, timeElapsedSeconds)
	renderDef.RenderSolidVideoBuffer()
}

func handleMainMenu(mainMenuStateInput *MainMenuStateInput, gameStateManager *GameStateManager) {
	renderDef := mainMenuStateInput.RenderDef
	if !gameStateManager.ImageResourcesLoaded {
		menuBackgroundImageADTOutput, err := fileio.LoadADTFile(game.MENU_IMAGE_FILE)
		if err != nil {
			log.Fatal("Error loading menu image: ", err)
		}
		menuBackgroundImage := render.ConvertPixelsToImage16Bit(menuBackgroundImageADTOutput.PixelData)

		menuBackgroundTextImagesTIMOutput, _ := fileio.LoadTIMImages(game.MENU_TEXT_FILE)
		menuTextImages := make([]*render.Image16Bit, len(menuBackgroundTextImagesTIMOutput))
		for i := 0; i < len(menuBackgroundTextImagesTIMOutput); i++ {
			menuTextImages[i] = render.ConvertPixelsToImage16Bit(menuBackgroundTextImagesTIMOutput[i].PixelData)
		}

		mainMenuStateInput.MenuBackgroundImage = menuBackgroundImage
		mainMenuStateInput.MenuTextImages = menuTextImages
		mainMenuStateInput.Menu.CurrentOption = 0
		renderDef.UpdateMainMenu(mainMenuStateInput.MenuBackgroundImage, mainMenuStateInput.MenuTextImages,
			mainMenuStateInput.Menu.CurrentOption)

		gameStateManager.ImageResourcesLoaded = true
		gameStateManager.UpdateLastTimeChangeState()
	}

	renderDef.RenderTransparentVideoBuffer()

	if gameStateManager.CanUpdateGameState() {
		mainMenuStateInput.Menu.HandleMenuEvent(windowHandler)

		if mainMenuStateInput.Menu.IsOptionSelected {
			if mainMenuStateInput.Menu.CurrentOption == 0 {
				gameStateManager.UpdateGameState(GAME_STATE_LOAD_SAVE)
				gameStateManager.UpdateLastTimeChangeState()
			} else if mainMenuStateInput.Menu.CurrentOption == 1 {
				gameStateManager.UpdateGameState(GAME_STATE_MAIN_GAME)
				gameStateManager.UpdateLastTimeChangeState()
			} else if mainMenuStateInput.Menu.CurrentOption == 2 {
				gameStateManager.UpdateGameState(GAME_STATE_SPECIAL_MENU)
				gameStateManager.UpdateLastTimeChangeState()
			}

			mainMenuStateInput.Menu.IsOptionSelected = false
		} else if mainMenuStateInput.Menu.IsNewOption {
			renderDef.UpdateMainMenu(mainMenuStateInput.MenuBackgroundImage, mainMenuStateInput.MenuTextImages,
				mainMenuStateInput.Menu.CurrentOption)
			gameStateManager.UpdateLastTimeChangeState()

			mainMenuStateInput.Menu.IsNewOption = false
		}
	}
}

func handleLoadSave(renderDef *render.RenderDef, gameStateManager *GameStateManager) {
	if !gameStateManager.ImageResourcesLoaded {
		// Initialize load save screen
		saveScreenImageADTOutput, err := fileio.LoadADTFile(game.SAVE_SCREEN_FILE)
		if err != nil {
			log.Fatal("Error loading save screen image: ", err)
		}
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

func handleSpecialMenu(specialMenuStateInput *MainMenuStateInput, gameStateManager *GameStateManager) {
	renderDef := specialMenuStateInput.RenderDef
	if !gameStateManager.ImageResourcesLoaded {
		menuBackgroundImageADTOutput, err := fileio.LoadADTFile(game.MENU_IMAGE_FILE)
		if err != nil {
			log.Fatal("Error loading menu image: ", err)
		}
		menuBackgroundImage := render.ConvertPixelsToImage16Bit(menuBackgroundImageADTOutput.PixelData)

		menuBackgroundTextImagesTIMOutput, _ := fileio.LoadTIMImages(game.MENU_TEXT_FILE)
		menuTextImages := make([]*render.Image16Bit, len(menuBackgroundTextImagesTIMOutput))
		for i := 0; i < len(menuBackgroundTextImagesTIMOutput); i++ {
			menuTextImages[i] = render.ConvertPixelsToImage16Bit(menuBackgroundTextImagesTIMOutput[i].PixelData)
		}

		specialMenuStateInput.MenuBackgroundImage = menuBackgroundImage
		specialMenuStateInput.MenuTextImages = menuTextImages
		specialMenuStateInput.Menu.CurrentOption = 0
		renderDef.UpdateSpecialMenu(specialMenuStateInput.MenuBackgroundImage, specialMenuStateInput.MenuTextImages,
			specialMenuStateInput.Menu.CurrentOption)

		gameStateManager.ImageResourcesLoaded = true
		gameStateManager.UpdateLastTimeChangeState()
	}

	renderDef.RenderTransparentVideoBuffer()

	if gameStateManager.CanUpdateGameState() {
		specialMenuStateInput.Menu.HandleMenuEvent(windowHandler)

		if specialMenuStateInput.Menu.IsOptionSelected {
			if specialMenuStateInput.Menu.CurrentOption == 0 {
				// TODO: Load gallery
			} else if specialMenuStateInput.Menu.CurrentOption == 1 {
				// Exit
				gameStateManager.UpdateGameState(GAME_STATE_MAIN_MENU)
				gameStateManager.UpdateLastTimeChangeState()
			}

			specialMenuStateInput.Menu.IsOptionSelected = false
		} else if specialMenuStateInput.Menu.IsNewOption {
			renderDef.UpdateSpecialMenu(specialMenuStateInput.MenuBackgroundImage, specialMenuStateInput.MenuTextImages,
				specialMenuStateInput.Menu.CurrentOption)
			gameStateManager.UpdateLastTimeChangeState()

			specialMenuStateInput.Menu.IsNewOption = false
		}
	}
}
