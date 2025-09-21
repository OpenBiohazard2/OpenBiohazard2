package state

import (
	"github.com/OpenBiohazard2/OpenBiohazard2/client"
	"github.com/OpenBiohazard2/OpenBiohazard2/render"
	"github.com/OpenBiohazard2/OpenBiohazard2/resource"
	"github.com/OpenBiohazard2/OpenBiohazard2/ui"
	"github.com/OpenBiohazard2/OpenBiohazard2/ui_render"
)

type InventoryStateInput struct {
	RenderDef           *render.RenderDef
	UIRenderer          *ui_render.UIRenderer
	InventoryMenuImages []*resource.Image16Bit
	InventoryItemImages []*resource.Image16Bit
	InventoryMenu       *ui.InventoryMenu
	HealthDisplay       *ui.HealthDisplay
	InventoryManager    *ui.InventoryManager
}

func NewInventoryStateInput(renderDef *render.RenderDef) *InventoryStateInput {
	inventoryMenuImages := resource.LoadTIMImages(resource.INVENTORY_FILE)
	inventoryItemImages := resource.LoadTIMImages(resource.ITEMALL_FILE)
	
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
