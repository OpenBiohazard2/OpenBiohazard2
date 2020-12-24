package gui

import (
	"github.com/samuelyuan/openbiohazard2/client"
)

const (
	MAX_TOP_MENU_SLOTS  = 4
	MAX_INVENTORY_SLOTS = 8
	RESERVED_ITEM_SLOT  = 10
)

type InventoryMenu struct {
	Status_Function0           int
	Status_Function1           int
	Status_MenuCursor0         int
	Status_InventoryMainCursor int
	Status_BlinkSwitch0        bool
	Status_BlinkTimer0         int
}

func NewInventoryMenu() *InventoryMenu {
	inventoryMenu := &InventoryMenu{}
	inventoryMenu.Reset()
	return inventoryMenu
}

func (inventoryMenu *InventoryMenu) Reset() {
	inventoryMenu.Status_Function0 = 3
	inventoryMenu.Status_Function1 = 0
	inventoryMenu.Status_MenuCursor0 = 2
	inventoryMenu.Status_InventoryMainCursor = 0
	inventoryMenu.Status_BlinkSwitch0 = false
	inventoryMenu.Status_BlinkTimer0 = 50
}

func (inventoryMenu *InventoryMenu) HandleSwitchMenuOption(windowHandler *client.WindowHandler) {
	if inventoryMenu.IsCursorOnTopMenu() {
		if windowHandler.InputHandler.IsActive(client.MENU_LEFT_BUTTON) {
			inventoryMenu.PrevTopMenuOption()
		} else if windowHandler.InputHandler.IsActive(client.MENU_RIGHT_BUTTON) {
			inventoryMenu.NextTopMenuOption()
		}
		return
	}

	if inventoryMenu.IsEditingItemScreen() {
		if windowHandler.InputHandler.IsActive(client.MENU_LEFT_BUTTON) {
			inventoryMenu.PrevItemInList()
		} else if windowHandler.InputHandler.IsActive(client.MENU_RIGHT_BUTTON) {
			inventoryMenu.NextItemInList()
		} else if windowHandler.InputHandler.IsActive(client.MENU_UP_BUTTON) {
			inventoryMenu.PrevRowInItemList()
		} else if windowHandler.InputHandler.IsActive(client.MENU_DOWN_BUTTON) {
			inventoryMenu.NextRowInItemList()
		}
		return
	}
}

func (inventoryMenu *InventoryMenu) NextTopMenuOption() {
	inventoryMenu.Status_MenuCursor0++
	if inventoryMenu.Status_MenuCursor0 >= MAX_TOP_MENU_SLOTS {
		inventoryMenu.Status_MenuCursor0 = MAX_TOP_MENU_SLOTS - 1
	}
}

func (inventoryMenu *InventoryMenu) PrevTopMenuOption() {
	inventoryMenu.Status_MenuCursor0--
	if inventoryMenu.Status_MenuCursor0 < 0 {
		inventoryMenu.Status_MenuCursor0 = 0
	}
}

func (inventoryMenu *InventoryMenu) NextItemInList() {
	if inventoryMenu.Status_InventoryMainCursor == RESERVED_ITEM_SLOT {
		return
	}

	inventoryMenu.Status_InventoryMainCursor++
	if inventoryMenu.Status_InventoryMainCursor >= MAX_INVENTORY_SLOTS {
		inventoryMenu.Status_InventoryMainCursor = MAX_INVENTORY_SLOTS - 1
	}
}

func (inventoryMenu *InventoryMenu) PrevItemInList() {
	if inventoryMenu.Status_InventoryMainCursor == RESERVED_ITEM_SLOT {
		return
	}

	inventoryMenu.Status_InventoryMainCursor--
	if inventoryMenu.Status_InventoryMainCursor < 0 {
		inventoryMenu.Status_InventoryMainCursor = 0
	}
}

func (inventoryMenu *InventoryMenu) NextRowInItemList() {
	if inventoryMenu.Status_InventoryMainCursor == RESERVED_ITEM_SLOT {
		inventoryMenu.Status_InventoryMainCursor = 1
		return
	}

	if inventoryMenu.Status_InventoryMainCursor+2 < MAX_INVENTORY_SLOTS {
		inventoryMenu.Status_InventoryMainCursor += 2
	}
}

func (inventoryMenu *InventoryMenu) PrevRowInItemList() {
	// Return to top menu
	if inventoryMenu.Status_InventoryMainCursor == RESERVED_ITEM_SLOT {
		inventoryMenu.Status_InventoryMainCursor = 0
		inventoryMenu.SetCursorTopMenu()
		return
	}

	if inventoryMenu.Status_InventoryMainCursor-2 >= 0 {
		inventoryMenu.Status_InventoryMainCursor -= 2
	} else if inventoryMenu.Status_InventoryMainCursor == 1 {
		inventoryMenu.Status_InventoryMainCursor = RESERVED_ITEM_SLOT
	}
}

func (inventoryMenu *InventoryMenu) IsCursorOnTopMenu() bool {
	return inventoryMenu.Status_Function0 < 3
}

func (inventoryMenu *InventoryMenu) IsEditingItemScreen() bool {
	return inventoryMenu.Status_Function0 == 3 && inventoryMenu.GetTopMenuSelectedOption() == 2
}

func (inventoryMenu *InventoryMenu) IsTopMenuCursorOnItems() bool {
	return inventoryMenu.Status_Function0 < 3 && inventoryMenu.GetTopMenuSelectedOption() == 2
}

func (inventoryMenu *InventoryMenu) IsTopMenuExit() bool {
	return inventoryMenu.GetTopMenuSelectedOption() == 3
}

func (inventoryMenu *InventoryMenu) SetCursorTopMenu() {
	// Can only naviagate top menu with cursor
	inventoryMenu.Status_Function0 = 2
}

func (inventoryMenu *InventoryMenu) SetEditItemScreen() {
	inventoryMenu.Status_Function0 = 3
}

func (inventoryMenu *InventoryMenu) UpdateCursorBlink() {
	if inventoryMenu.Status_Function1 == 0 {
		if !inventoryMenu.Status_BlinkSwitch0 {
			if inventoryMenu.Status_BlinkTimer0 > 120 {
				inventoryMenu.Status_BlinkSwitch0 = true
			}
			inventoryMenu.Status_BlinkTimer0 += 3
		} else {
			if inventoryMenu.Status_BlinkTimer0 < 50 {
				inventoryMenu.Status_BlinkSwitch0 = false
			}
			inventoryMenu.Status_BlinkTimer0 -= 3
		}
	} else {
		inventoryMenu.Status_BlinkTimer0 = 60
	}
}

func (inventoryMenu *InventoryMenu) GetCursorBlinkBrightnessFactor() float64 {
	return float64(inventoryMenu.Status_BlinkTimer0) / 128.0
}

func (inventoryMenu *InventoryMenu) IsCursorOnReservedItem() bool {
	return inventoryMenu.Status_InventoryMainCursor == RESERVED_ITEM_SLOT
}

func (inventoryMenu *InventoryMenu) GetMainCursorRow() int {
	return inventoryMenu.Status_InventoryMainCursor / 2
}

func (inventoryMenu *InventoryMenu) GetMainCursorColumn() int {
	return inventoryMenu.Status_InventoryMainCursor % 2
}

func (inventoryMenu *InventoryMenu) GetTopMenuSelectedOption() int {
	return inventoryMenu.Status_MenuCursor0
}
