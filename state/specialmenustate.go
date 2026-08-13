package state

import (
	"fmt"

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

func HandleSpecialMenu(specialMenuStateInput *SpecialMenuStateInput, gameStateManager *GameStateManager, windowHandler *client.WindowHandler) error {
	renderDef := specialMenuStateInput.RenderDef
	if !gameStateManager.ImageResourcesLoaded {
		menuBackgroundImage, err := resource.LoadADTImage(resource.MENU_IMAGE_FILE)
		if err != nil {
			return fmt.Errorf("failed to load special menu background image: %w", err)
		}
		menuTextImages, err := resource.LoadTIMImages(resource.MENU_TEXT_FILE)
		if err != nil {
			return fmt.Errorf("failed to load special menu text images: %w", err)
		}
		specialMenuStateInput.MenuBackgroundImage = menuBackgroundImage
		specialMenuStateInput.MenuTextImages = menuTextImages
		specialMenuStateInput.Menu.CurrentOption = 0
		specialMenuStateInput.UIRenderer.UpdateSpecialMenu(specialMenuStateInput.MenuBackgroundImage, specialMenuStateInput.MenuTextImages,
			specialMenuStateInput.Menu.CurrentOption)

		gameStateManager.ImageResourcesLoaded = true
	}

	renderDef.RenderTransparentVideoBuffer()

	specialMenuStateInput.Menu.HandleMenuEvent(windowHandler)

	if specialMenuStateInput.Menu.IsOptionSelected {
		switch specialMenuStateInput.Menu.CurrentOption {
		case 0:
			// TODO: Load gallery
		case 1:
			// Exit
			gameStateManager.UpdateGameState(GAME_STATE_MAIN_MENU)
		}

		specialMenuStateInput.Menu.IsOptionSelected = false
	} else if specialMenuStateInput.Menu.IsNewOption {
		specialMenuStateInput.UIRenderer.UpdateSpecialMenu(specialMenuStateInput.MenuBackgroundImage, specialMenuStateInput.MenuTextImages,
			specialMenuStateInput.Menu.CurrentOption)

		specialMenuStateInput.Menu.IsNewOption = false
	}
	return nil
}
