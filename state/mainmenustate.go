package state

import (
	"github.com/OpenBiohazard2/OpenBiohazard2/client"
	"github.com/OpenBiohazard2/OpenBiohazard2/render"
	"github.com/OpenBiohazard2/OpenBiohazard2/resource"
	"github.com/OpenBiohazard2/OpenBiohazard2/ui"
	"github.com/OpenBiohazard2/OpenBiohazard2/ui_render"
)

type MainMenuStateInput struct {
	RenderDef           *render.RenderDef
	UIRenderer          *ui_render.UIRenderer
	MenuBackgroundImage *resource.Image16Bit
	MenuTextImages      []*resource.Image16Bit
	Menu                *ui.Menu
}

func HandleMainMenu(mainMenuStateInput *MainMenuStateInput, gameStateManager *GameStateManager, windowHandler *client.WindowHandler) {
	renderDef := mainMenuStateInput.RenderDef
	if !gameStateManager.ImageResourcesLoaded {
		mainMenuStateInput.MenuBackgroundImage = resource.LoadADTImage(resource.MENU_IMAGE_FILE)
		mainMenuStateInput.MenuTextImages = resource.LoadTIMImages(resource.MENU_TEXT_FILE)
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
