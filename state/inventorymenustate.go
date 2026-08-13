package state

import (
	"fmt"

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

func NewInventoryStateInput(renderDef *render.RenderDef) (*InventoryStateInput, error) {
	inventoryMenuImages, err := resource.LoadTIMImages(resource.INVENTORY_FILE)
	if err != nil {
		return nil, fmt.Errorf("failed to load inventory menu images: %w", err)
	}
	inventoryItemImages, err := resource.LoadTIMImages(resource.ITEMALL_FILE)
	if err != nil {
		return nil, fmt.Errorf("failed to load inventory item images: %w", err)
	}

	return &InventoryStateInput{
		RenderDef:           renderDef,
		UIRenderer:          ui_render.NewUIRenderer(renderDef),
		InventoryMenuImages: inventoryMenuImages,
		InventoryItemImages: inventoryItemImages,
		InventoryMenu:       ui.NewInventoryMenu(),
		HealthDisplay:       ui.NewHealthDisplay(),
		InventoryManager:    ui.NewInventoryManager(),
	}, nil
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
	}

	if windowHandler.InputHandler.IsActiveOnce(client.PLAYER_VIEW_INVENTORY) {
		gameStateManager.UpdateGameState(GAME_STATE_MAIN_GAME)
	}

	if windowHandler.InputHandler.IsActiveOnce(client.ACTION_BUTTON) {
		if inventoryMenu.IsCursorOnTopMenu() {
			if inventoryMenu.IsTopMenuExit() {
				gameStateManager.UpdateGameState(GAME_STATE_MAIN_GAME)
			} else if inventoryMenu.IsTopMenuCursorOnItems() {
				inventoryMenu.SetEditItemScreen()
			}
		}
	}

	inventoryMenu.HandleSwitchMenuOption(windowHandler)

	timeElapsedSeconds := windowHandler.GetTimeSinceLastFrame()
	inventoryStateInput.UIRenderer.GenerateInventoryImage(inventoryMenuImages, inventoryItemImages, inventoryMenu, healthDisplay, inventoryManager, timeElapsedSeconds)
	renderDef.RenderSolidVideoBuffer()
}
