package state

import (
	"fmt"

	"github.com/OpenBiohazard2/OpenBiohazard2/client"
	"github.com/OpenBiohazard2/OpenBiohazard2/render"
	"github.com/OpenBiohazard2/OpenBiohazard2/resource"
	"github.com/OpenBiohazard2/OpenBiohazard2/ui_render"
)

func HandleLoadSave(renderDef *render.RenderDef, gameStateManager *GameStateManager, windowHandler *client.WindowHandler) error {
	if !gameStateManager.ImageResourcesLoaded {
		// Initialize load save screen
		saveScreenImage, err := resource.LoadADTImage(resource.SAVE_SCREEN_FILE)
		if err != nil {
			return fmt.Errorf("failed to load save screen image: %w", err)
		}
		uiRenderer := ui_render.NewUIRenderer(renderDef)
		uiRenderer.GenerateSaveScreenImage(saveScreenImage)

		gameStateManager.ImageResourcesLoaded = true
	}

	renderDef.RenderTransparentVideoBuffer()
	if windowHandler.InputHandler.IsActiveOnce(client.ACTION_BUTTON) {
		gameStateManager.UpdateGameState(GAME_STATE_MAIN_MENU)
	}
	return nil
}
