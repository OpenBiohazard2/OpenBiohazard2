package state

import (
	"fmt"

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

func HandleMainMenu(mainMenuStateInput *MainMenuStateInput, gameStateManager *GameStateManager, windowHandler *client.WindowHandler) error {
	renderDef := mainMenuStateInput.RenderDef
	if !gameStateManager.ImageResourcesLoaded {
		menuBackgroundImage, err := resource.LoadADTImage(resource.MENU_IMAGE_FILE)
		if err != nil {
			return fmt.Errorf("failed to load main menu background image: %w", err)
		}
		menuTextImages, err := resource.LoadTIMImages(resource.MENU_TEXT_FILE)
		if err != nil {
			return fmt.Errorf("failed to load main menu text images: %w", err)
		}
		mainMenuStateInput.MenuBackgroundImage = menuBackgroundImage
		mainMenuStateInput.MenuTextImages = menuTextImages
		mainMenuStateInput.Menu.CurrentOption = 0
		mainMenuStateInput.UIRenderer.UpdateMainMenu(mainMenuStateInput.MenuBackgroundImage, mainMenuStateInput.MenuTextImages,
			mainMenuStateInput.Menu.CurrentOption)

		gameStateManager.ImageResourcesLoaded = true
	}

	renderDef.RenderTransparentVideoBuffer()

	mainMenuStateInput.Menu.HandleMenuEvent(windowHandler)

	if mainMenuStateInput.Menu.IsOptionSelected {
		switch mainMenuStateInput.Menu.CurrentOption {
		case 0:
			gameStateManager.UpdateGameState(GAME_STATE_LOAD_SAVE)
		case 1:
			gameStateManager.UpdateGameState(GAME_STATE_MAIN_GAME)
		case 2:
			gameStateManager.UpdateGameState(GAME_STATE_SPECIAL_MENU)
		}

		mainMenuStateInput.Menu.IsOptionSelected = false
	} else if mainMenuStateInput.Menu.IsNewOption {
		mainMenuStateInput.UIRenderer.UpdateMainMenu(mainMenuStateInput.MenuBackgroundImage, mainMenuStateInput.MenuTextImages,
			mainMenuStateInput.Menu.CurrentOption)

		mainMenuStateInput.Menu.IsNewOption = false
	}
	return nil
}
