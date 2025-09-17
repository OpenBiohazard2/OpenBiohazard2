package state

import (
	"log"

	"github.com/OpenBiohazard2/OpenBiohazard2/client"
	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/OpenBiohazard2/OpenBiohazard2/game"
	"github.com/OpenBiohazard2/OpenBiohazard2/render"
	"github.com/OpenBiohazard2/OpenBiohazard2/ui"
	"github.com/OpenBiohazard2/OpenBiohazard2/ui_render"
)

type MainMenuStateInput struct {
	RenderDef           *render.RenderDef
	UIRenderer          *ui_render.UIRenderer
	MenuBackgroundImage *render.Image16Bit
	MenuTextImages      []*render.Image16Bit
	Menu                *ui.Menu
}

type InventoryStateInput struct {
	RenderDef           *render.RenderDef
	UIRenderer          *ui_render.UIRenderer
	InventoryMenuImages []*render.Image16Bit
	InventoryItemImages []*render.Image16Bit
	InventoryMenu       *ui.InventoryMenu
	HealthDisplay       *ui.HealthDisplay
	InventoryManager    *ui.InventoryManager
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
		UIRenderer:          ui_render.NewUIRenderer(renderDef),
		InventoryMenuImages: inventoryMenuImages,
		InventoryItemImages: inventoryItemImages,
		InventoryMenu:       ui.NewInventoryMenu(),
		HealthDisplay:       ui.NewHealthDisplay(),
		InventoryManager:    ui.NewInventoryManager(),
	}
}

func HandleInventory(inventoryStateInput *InventoryStateInput, gameStateManager *GameStateManager, windowHandler *client.WindowHandler) {
	renderDef := inventoryStateInput.RenderDef
	inventoryMenuImages := inventoryStateInput.InventoryMenuImages
	inventoryItemImages := inventoryStateInput.InventoryItemImages
	inventoryMenu := inventoryStateInput.InventoryMenu
	healthDisplay := inventoryStateInput.HealthDisplay
	inventoryManager := inventoryStateInput.InventoryManager

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
	inventoryStateInput.UIRenderer.GenerateInventoryImage(inventoryMenuImages, inventoryItemImages, inventoryMenu, healthDisplay, inventoryManager, timeElapsedSeconds)
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
		mainMenuStateInput.UIRenderer.UpdateMainMenu(mainMenuStateInput.MenuBackgroundImage, mainMenuStateInput.MenuTextImages,
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
			mainMenuStateInput.UIRenderer.UpdateMainMenu(mainMenuStateInput.MenuBackgroundImage, mainMenuStateInput.MenuTextImages,
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
		uiRenderer := ui_render.NewUIRenderer(renderDef)
		uiRenderer.GenerateSaveScreenImage(saveScreenImageRender)

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
		specialMenuStateInput.UIRenderer.UpdateSpecialMenu(specialMenuStateInput.MenuBackgroundImage, specialMenuStateInput.MenuTextImages,
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
			specialMenuStateInput.UIRenderer.UpdateSpecialMenu(specialMenuStateInput.MenuBackgroundImage, specialMenuStateInput.MenuTextImages,
				specialMenuStateInput.Menu.CurrentOption)
			gameStateManager.UpdateLastTimeChangeState(windowHandler)

			specialMenuStateInput.Menu.IsNewOption = false
		}
	}
}
