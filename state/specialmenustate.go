package state

import (
	"github.com/OpenBiohazard2/OpenBiohazard2/client"
	"github.com/OpenBiohazard2/OpenBiohazard2/render"
	"github.com/OpenBiohazard2/OpenBiohazard2/resource"
	"github.com/OpenBiohazard2/OpenBiohazard2/ui"
	"github.com/OpenBiohazard2/OpenBiohazard2/ui_render"
)

type SpecialMenuStateInput struct {
	RenderDef           *render.RenderDef
	UIRenderer          *ui_render.UIRenderer
	MenuBackgroundImage *resource.Image16Bit
	MenuTextImages      []*resource.Image16Bit
	Menu                *ui.Menu
}

func HandleSpecialMenu(specialMenuStateInput *SpecialMenuStateInput, gameStateManager *GameStateManager, windowHandler *client.WindowHandler) {
	renderDef := specialMenuStateInput.RenderDef
	if !gameStateManager.ImageResourcesLoaded {
		specialMenuStateInput.MenuBackgroundImage = resource.LoadADTImage(resource.MENU_IMAGE_FILE)
		specialMenuStateInput.MenuTextImages = resource.LoadTIMImages(resource.MENU_TEXT_FILE)
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
