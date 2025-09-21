package state

import (
	"github.com/OpenBiohazard2/OpenBiohazard2/client"
	"github.com/OpenBiohazard2/OpenBiohazard2/render"
	"github.com/OpenBiohazard2/OpenBiohazard2/resource"
	"github.com/OpenBiohazard2/OpenBiohazard2/ui_render"
)

func HandleLoadSave(renderDef *render.RenderDef, gameStateManager *GameStateManager, windowHandler *client.WindowHandler) {
	if !gameStateManager.ImageResourcesLoaded {
		// Initialize load save screen
		saveScreenImage := resource.LoadADTImage(resource.SAVE_SCREEN_FILE)
		uiRenderer := ui_render.NewUIRenderer(renderDef)
		uiRenderer.GenerateSaveScreenImage(saveScreenImage)

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
