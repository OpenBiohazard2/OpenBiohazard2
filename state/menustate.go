package state

import (
	"log"

	"github.com/OpenBiohazard2/OpenBiohazard2/client"
	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/OpenBiohazard2/OpenBiohazard2/game"
	"github.com/OpenBiohazard2/OpenBiohazard2/gui"
	"github.com/OpenBiohazard2/OpenBiohazard2/render"
)

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

func HandleInventory(inventoryStateInput *InventoryStateInput, gameStateManager *GameStateManager, windowHandler *client.WindowHandler) {
	renderDef := inventoryStateInput.RenderDef
	inventoryMenuImages := inventoryStateInput.InventoryMenuImages
	inventoryItemImages := inventoryStateInput.InventoryItemImages
	inventoryMenu := inventoryStateInput.InventoryMenu

	if !gameStateManager.ImageResourcesLoaded {
		inventoryMenu.Reset()
		gameStateManager.ImageResourcesLoaded = true
		gameStateManager.UpdateLastTimeChangeState(windowHandler)
	}

	if windowHandler.InputHandler.IsActive(client.PLAYER_VIEW_INVENTORY) {
		if gameStateManager.CanUpdateGameState(windowHandler) {
			gameStateManager.UpdateGameState(GAME_STATE_MAIN_GAME)
			gameStateManager.UpdateLastTimeChangeState(windowHandler)
		}
	}

	if windowHandler.InputHandler.IsActive(client.ACTION_BUTTON) {
		if gameStateManager.CanUpdateGameState(windowHandler) {
			if inventoryMenu.IsCursorOnTopMenu() {
				if inventoryMenu.IsTopMenuExit() {
					gameStateManager.UpdateGameState(GAME_STATE_MAIN_GAME)
				} else if inventoryMenu.IsTopMenuCursorOnItems() {
					inventoryMenu.SetEditItemScreen()
				}
			}
			gameStateManager.UpdateLastTimeChangeState(windowHandler)
		}
	}

	if gameStateManager.CanUpdateGameState(windowHandler) {
		inventoryMenu.HandleSwitchMenuOption(windowHandler)
		gameStateManager.UpdateLastTimeChangeState(windowHandler)
	}

	timeElapsedSeconds := windowHandler.GetTimeSinceLastFrame()
	renderDef.GenerateInventoryImage(inventoryMenuImages, inventoryItemImages, inventoryMenu, timeElapsedSeconds)
	renderDef.RenderSolidVideoBuffer()
}

func HandleMainMenu(mainMenuStateInput *MainMenuStateInput, gameStateManager *GameStateManager, windowHandler *client.WindowHandler) {
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
		gameStateManager.UpdateLastTimeChangeState(windowHandler)
	}

	renderDef.RenderTransparentVideoBuffer()

	if gameStateManager.CanUpdateGameState(windowHandler) {
		mainMenuStateInput.Menu.HandleMenuEvent(windowHandler)

		if mainMenuStateInput.Menu.IsOptionSelected {
			switch mainMenuStateInput.Menu.CurrentOption {
			case 0:
				gameStateManager.UpdateGameState(GAME_STATE_LOAD_SAVE)
				gameStateManager.UpdateLastTimeChangeState(windowHandler)
			case 1:
				gameStateManager.UpdateGameState(GAME_STATE_MAIN_GAME)
				gameStateManager.UpdateLastTimeChangeState(windowHandler)
			case 2:
				gameStateManager.UpdateGameState(GAME_STATE_SPECIAL_MENU)
				gameStateManager.UpdateLastTimeChangeState(windowHandler)
			}

			mainMenuStateInput.Menu.IsOptionSelected = false
		} else if mainMenuStateInput.Menu.IsNewOption {
			renderDef.UpdateMainMenu(mainMenuStateInput.MenuBackgroundImage, mainMenuStateInput.MenuTextImages,
				mainMenuStateInput.Menu.CurrentOption)
			gameStateManager.UpdateLastTimeChangeState(windowHandler)

			mainMenuStateInput.Menu.IsNewOption = false
		}
	}
}

func HandleLoadSave(renderDef *render.RenderDef, gameStateManager *GameStateManager, windowHandler *client.WindowHandler) {
	if !gameStateManager.ImageResourcesLoaded {
		// Initialize load save screen
		saveScreenImageADTOutput, err := fileio.LoadADTFile(game.SAVE_SCREEN_FILE)
		if err != nil {
			log.Fatal("Error loading save screen image: ", err)
		}
		saveScreenImageRender := render.ConvertPixelsToImage16Bit(saveScreenImageADTOutput.PixelData)
		renderDef.GenerateSaveScreenImage(saveScreenImageRender)

		gameStateManager.ImageResourcesLoaded = true
		gameStateManager.UpdateLastTimeChangeState(windowHandler)
	}

	renderDef.RenderTransparentVideoBuffer()
	if windowHandler.InputHandler.IsActive(client.ACTION_BUTTON) {
		if gameStateManager.CanUpdateGameState(windowHandler) {
			gameStateManager.UpdateGameState(GAME_STATE_MAIN_MENU)
			gameStateManager.UpdateLastTimeChangeState(windowHandler)
		}
	}
}

func HandleSpecialMenu(specialMenuStateInput *MainMenuStateInput, gameStateManager *GameStateManager, windowHandler *client.WindowHandler) {
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
		gameStateManager.UpdateLastTimeChangeState(windowHandler)
	}

	renderDef.RenderTransparentVideoBuffer()

	if gameStateManager.CanUpdateGameState(windowHandler) {
		specialMenuStateInput.Menu.HandleMenuEvent(windowHandler)

		if specialMenuStateInput.Menu.IsOptionSelected {
			switch specialMenuStateInput.Menu.CurrentOption {
			case 0:
				// TODO: Load gallery
			case 1:
				// Exit
				gameStateManager.UpdateGameState(GAME_STATE_MAIN_MENU)
				gameStateManager.UpdateLastTimeChangeState(windowHandler)
			}

			specialMenuStateInput.Menu.IsOptionSelected = false
		} else if specialMenuStateInput.Menu.IsNewOption {
			renderDef.UpdateSpecialMenu(specialMenuStateInput.MenuBackgroundImage, specialMenuStateInput.MenuTextImages,
				specialMenuStateInput.Menu.CurrentOption)
			gameStateManager.UpdateLastTimeChangeState(windowHandler)

			specialMenuStateInput.Menu.IsNewOption = false
		}
	}
}
