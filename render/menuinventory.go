package render

import (
	"github.com/samuelyuan/openbiohazard2/fileio"
)

const (
	ITEMLIST_POS_X      = 220
	ITEMLIST_POS_Y      = 70
	MAX_TOP_MENU_SLOTS  = 4
	MAX_INVENTORY_SLOTS = 8
	RESERVED_ITEM_SLOT  = 10
)

var (
	totalInventoryTime        = float64(0)
	updateInventoryCursorTime = float64(30) // milliseconds

	Status_Function0           = 3
	Status_Function1           = 0
	Status_MenuCursor0         = 2
	Status_InventoryMainCursor = 0
	Status_BlinkSwitch0        = false
	Status_BlinkTimer0         = 50

	playerInventoryItems = InitializeInventoryItems()
)

type InventoryItem struct {
	Id   int
	Num  int
	Size int
}

func InitializeInventoryCursor() {
	Status_Function0 = 3
	Status_Function1 = 0
	Status_MenuCursor0 = 2
	Status_InventoryMainCursor = 0
	Status_BlinkSwitch0 = false
	Status_BlinkTimer0 = 50
}

func InitializeInventoryItems() []InventoryItem {
	playerInventoryItems := make([]InventoryItem, 11)
	playerInventoryItems[0] = InventoryItem{Id: 2, Num: 18, Size: 1}                  // hand gun
	playerInventoryItems[1] = InventoryItem{Id: 1, Num: 1, Size: 0}                   // knife
	playerInventoryItems[RESERVED_ITEM_SLOT] = InventoryItem{Id: 47, Num: 1, Size: 0} // lighter
	return playerInventoryItems
}

func (renderDef *RenderDef) GenerateInventoryImage(
	inventoryImages []*fileio.TIMOutput,
	inventoryItemImages []*fileio.TIMOutput,
	timeElapsedSeconds float64) {
	renderDef.VideoBuffer.ClearSurface()
	newImageColors := renderDef.VideoBuffer.ImagePixels
	totalInventoryTime += timeElapsedSeconds * 1000
	totalHealthTime += timeElapsedSeconds * 1000
	buildBackground(inventoryImages, newImageColors)
	buildItems(inventoryImages, inventoryItemImages, newImageColors)
	renderDef.VideoBuffer.UpdateSurface(newImageColors)
}

func buildItems(inventoryImages []*fileio.TIMOutput, inventoryItemImages []*fileio.TIMOutput, newImageColors []uint16) {
	// Item in top right corner
	reservedItemX := (playerInventoryItems[RESERVED_ITEM_SLOT].Id % 6) * 40
	reservedItemY := (playerInventoryItems[RESERVED_ITEM_SLOT].Id / 6) * 30
	copyPixels(inventoryItemImages[0].PixelData, reservedItemX, reservedItemY, 40, 30, newImageColors, ITEMLIST_POS_X+45, ITEMLIST_POS_Y-35)

	// Empty inventory slots
	for row := 0; row < 4; row++ {
		// left slot
		leftItemIndex := (2 * row)
		leftItemX := (playerInventoryItems[leftItemIndex].Id % 6) * 40
		leftItemY := (playerInventoryItems[leftItemIndex].Id / 6) * 30
		copyPixels(inventoryItemImages[0].PixelData, leftItemX, leftItemY, 40, 30, newImageColors, ITEMLIST_POS_X+5, ITEMLIST_POS_Y+3+30*row)
		// right slot
		rightItemIndex := (2 * row) + 1
		rightItemX := (playerInventoryItems[rightItemIndex].Id % 6) * 40
		rightItemY := (playerInventoryItems[rightItemIndex].Id / 6) * 30
		copyPixels(inventoryItemImages[0].PixelData, rightItemX, rightItemY, 40, 30, newImageColors, ITEMLIST_POS_X+45, ITEMLIST_POS_Y+3+30*row)
	}

	// Equipped item
	copyPixels(inventoryItemImages[2].PixelData, 40, 90, 80, 30, newImageColors, 172, 35)

	// Item cursor surrounding item
	if IsEditingItemScreen() {
		displayInventoryMainCursor(inventoryImages, newImageColors)
	}
}

func displayInventoryMainCursor(inventoryImages []*fileio.TIMOutput, newImageColors []uint16) {
	var cursorX, cursorY int
	cursorFrameOffsetX := 3
	cursorFrameOffsetY := 1

	if totalInventoryTime >= updateInventoryCursorTime {
		if Status_Function1 == 0 {
			if !Status_BlinkSwitch0 {
				if Status_BlinkTimer0 > 120 {
					Status_BlinkSwitch0 = true
				}
				Status_BlinkTimer0 += 3
			} else {
				if Status_BlinkTimer0 < 50 {
					Status_BlinkSwitch0 = false
				}
				Status_BlinkTimer0 -= 3
			}
		} else {
			Status_BlinkTimer0 = 60
		}
		totalInventoryTime = 0
	}
	brightnessFactor := float64(Status_BlinkTimer0) / 128.0

	// Special item in top right corner
	if Status_InventoryMainCursor == RESERVED_ITEM_SLOT {
		cursorX = ITEMLIST_POS_X + cursorFrameOffsetX + 40
		cursorY = ITEMLIST_POS_Y + cursorFrameOffsetY - 38
	} else {
		cursorX = ITEMLIST_POS_X + cursorFrameOffsetX + (Status_InventoryMainCursor%2)*40
		cursorY = ITEMLIST_POS_Y + cursorFrameOffsetY + (Status_InventoryMainCursor/2)*30
	}
	copyPixelsBrightness(inventoryImages[3].PixelData, 0, 30, 44, 34, newImageColors, cursorX, cursorY, brightnessFactor)
}

func buildBackground(inventoryImages []*fileio.TIMOutput, newImageColors []uint16) {
	// The inventory image is split up into many small components
	// Combine them manually back into a single image
	// source image is 256x256
	// dest image is 320x240
	backgroundColor := [3]int{5, 5, 31}
	fillPixels(newImageColors, 0, 0, 320, 240, backgroundColor[0], backgroundColor[1], backgroundColor[2])

	buildPlayerFace(inventoryImages, newImageColors)
	buildHealthECG(inventoryImages, newImageColors, backgroundColor)

	// Equipped item
	copyPixels(inventoryImages[0].PixelData, 50, 211, 11, 39, newImageColors, 161, 29) // left
	copyPixels(inventoryImages[0].PixelData, 0, 158, 80, 6, newImageColors, 172, 29)   // top
	copyPixels(inventoryImages[0].PixelData, 91, 164, 5, 39, newImageColors, 252, 29)  // right
	copyPixels(inventoryImages[0].PixelData, 0, 155, 80, 3, newImageColors, 172, 65)   // bottom

	// Extra item
	copyPixels(inventoryImages[0].PixelData, 0, 211, 50, 41, newImageColors, 260, 29)

	buildMenuTabs(inventoryImages, newImageColors)

	// Item slots
	copyPixels(inventoryImages[0].PixelData, 114, 92, 5, 120, newImageColors, ITEMLIST_POS_X, ITEMLIST_POS_Y+3)    // left
	copyPixels(inventoryImages[0].PixelData, 0, 140, 90, 3, newImageColors, ITEMLIST_POS_X, ITEMLIST_POS_Y)        // top
	copyPixels(inventoryImages[0].PixelData, 114, 92, 5, 120, newImageColors, ITEMLIST_POS_X+85, ITEMLIST_POS_Y+3) // right
	copyPixels(inventoryImages[0].PixelData, 0, 140, 90, 4, newImageColors, ITEMLIST_POS_X, ITEMLIST_POS_Y+123)    // bottom

	buildDescription(inventoryImages, newImageColors)
}

func buildPlayerFace(inventoryImages []*fileio.TIMOutput, newImageColors []uint16) {
	// Player
	copyPixels(inventoryImages[0].PixelData, 106, 152, 4, 60, newImageColors, 7, 16)  // left
	copyPixels(inventoryImages[0].PixelData, 0, 140, 39, 4, newImageColors, 11, 16)   // top
	copyPixels(inventoryImages[0].PixelData, 109, 152, 4, 60, newImageColors, 49, 16) // right
	copyPixels(inventoryImages[0].PixelData, 0, 140, 39, 4, newImageColors, 11, 72)   // bottom
	copyPixels(inventoryImages[1].PixelData, 1, 73, 37, 8, newImageColors, 11, 21)    // player name
	copyPixels(inventoryImages[1].PixelData, 0, 85, 38, 42, newImageColors, 11, 31)   // player image
	copyPixels(inventoryImages[0].PixelData, 56, 164, 38, 1, newImageColors, 11, 30)  // line between name and image

	// Pipes to the left of player image
	copyPixels(inventoryImages[0].PixelData, 107, 242, 7, 14, newImageColors, 0, 17)
	copyPixels(inventoryImages[0].PixelData, 107, 242, 7, 14, newImageColors, 0, 33)
	copyPixels(inventoryImages[0].PixelData, 107, 242, 7, 14, newImageColors, 0, 49)

	// Pipes to the right of player image
	copyPixels(inventoryImages[0].PixelData, 56, 186, 7, 7, newImageColors, 53, 32)
	copyPixels(inventoryImages[0].PixelData, 56, 186, 7, 7, newImageColors, 53, HEALTH_POS_X+2)
}

func buildMenuTabs(inventoryImages []*fileio.TIMOutput, newImageColors []uint16) {
	var selectedOption [3]float64
	if IsCursorOnTopMenu() {
		// Cursor is on this option, but it's not selected
		selectedOption = [3]float64{1.0, 1.0, 1.0}
	} else {
		// Highlight option in red if this option is selected
		selectedOption = [3]float64{1.0, 0.5, 0.5}
	}

	otherOption := [3]float64{0.4, 0.4, 0.4}

	optionsBrightness := [4][3]float64{otherOption, otherOption, otherOption, otherOption}
	optionsBrightness[Status_MenuCursor0] = selectedOption

	// File
	copyPixels(inventoryImages[0].PixelData, 3, 164, 49, 12, newImageColors, 111, 16)
	copyPixelsBrightnessColor(inventoryImages[5].PixelData, 0, 0, 47, 10, newImageColors, 112, 17,
		optionsBrightness[0][0], optionsBrightness[0][1], optionsBrightness[0][2])

	// Map
	copyPixels(inventoryImages[0].PixelData, 3, 164, 49, 12, newImageColors, 160, 16)
	copyPixelsBrightnessColor(inventoryImages[5].PixelData, 0, 10, 47, 10, newImageColors, 161, 17,
		optionsBrightness[1][0], optionsBrightness[1][1], optionsBrightness[1][2])

	// Item
	copyPixels(inventoryImages[0].PixelData, 3, 164, 49, 12, newImageColors, 209, 16)
	copyPixelsBrightnessColor(inventoryImages[5].PixelData, 0, 20, 47, 10, newImageColors, 210, 17,
		optionsBrightness[2][0], optionsBrightness[2][1], optionsBrightness[2][2])

	// Exit
	copyPixels(inventoryImages[0].PixelData, 3, 164, 49, 12, newImageColors, 258, 16)
	copyPixelsBrightnessColor(inventoryImages[5].PixelData, 0, 30, 47, 10, newImageColors, 259, 17,
		optionsBrightness[3][0], optionsBrightness[3][1], optionsBrightness[3][2])
}

func buildDescription(inventoryImages []*fileio.TIMOutput, newImageColors []uint16) {
	descriptionColor := [3]int{6, 13, 23}
	fillPixels(newImageColors, 13, 174, 201, 49, descriptionColor[0], descriptionColor[1], descriptionColor[2])
	copyPixels(inventoryImages[0].PixelData, 106, 163, 5, 49, newImageColors, 8, 174)   // left
	copyPixels(inventoryImages[0].PixelData, 0, 147, 83, 4, newImageColors, 8, 170)     // top left
	copyPixels(inventoryImages[0].PixelData, 0, 80, 128, 4, newImageColors, 91, 170)    // top right
	copyPixels(inventoryImages[0].PixelData, 106, 163, 5, 49, newImageColors, 214, 174) // right
	copyPixels(inventoryImages[0].PixelData, 0, 147, 83, 4, newImageColors, 8, 223)     // bottom left
	copyPixels(inventoryImages[0].PixelData, 0, 80, 128, 4, newImageColors, 91, 223)    // bottom right

	// Pipes to the right of description
	copyPixels(inventoryImages[0].PixelData, 107, 242, 7, 14, newImageColors, 219, 212)
	copyPixels(inventoryImages[0].PixelData, 56, 178, 35, 7, newImageColors, 226, 215)
	copyPixels(inventoryImages[0].PixelData, 56, 178, 35, 7, newImageColors, 261, 215)
	copyPixels(inventoryImages[0].PixelData, 56, 178, 24, 7, newImageColors, 296, 215)
}

func NextTopMenuOption() {
	Status_MenuCursor0++
	if Status_MenuCursor0 >= MAX_TOP_MENU_SLOTS {
		Status_MenuCursor0 = MAX_TOP_MENU_SLOTS - 1
	}
}

func PrevTopMenuOption() {
	Status_MenuCursor0--
	if Status_MenuCursor0 < 0 {
		Status_MenuCursor0 = 0
	}
}

func NextItemInList() {
	if Status_InventoryMainCursor == RESERVED_ITEM_SLOT {
		return
	}

	Status_InventoryMainCursor++
	if Status_InventoryMainCursor >= MAX_INVENTORY_SLOTS {
		Status_InventoryMainCursor = MAX_INVENTORY_SLOTS - 1
	}
}

func PrevItemInList() {
	if Status_InventoryMainCursor == RESERVED_ITEM_SLOT {
		return
	}

	Status_InventoryMainCursor--
	if Status_InventoryMainCursor < 0 {
		Status_InventoryMainCursor = 0
	}
}

func NextRowInItemList() {
	if Status_InventoryMainCursor == RESERVED_ITEM_SLOT {
		Status_InventoryMainCursor = 1
		return
	}

	if Status_InventoryMainCursor+2 < MAX_INVENTORY_SLOTS {
		Status_InventoryMainCursor += 2
	}
}

func PrevRowInItemList() {
	// Return to top menu
	if Status_InventoryMainCursor == RESERVED_ITEM_SLOT {
		Status_InventoryMainCursor = 0
		SetCursorTopMenu()
		return
	}

	if Status_InventoryMainCursor-2 >= 0 {
		Status_InventoryMainCursor -= 2
	} else if Status_InventoryMainCursor == 1 {
		Status_InventoryMainCursor = RESERVED_ITEM_SLOT
	}
}

func IsCursorOnTopMenu() bool {
	return Status_Function0 < 3
}

func IsEditingItemScreen() bool {
	return Status_Function0 == 3 && Status_MenuCursor0 == 2
}

func IsTopMenuCursorOnItems() bool {
	return Status_Function0 < 3 && Status_MenuCursor0 == 2
}

func IsTopMenuExit() bool {
	return Status_MenuCursor0 == 3
}

func SetCursorTopMenu() {
	// Can only naviagate top menu with cursor
	Status_Function0 = 2
}

func SetEditItemScreen() {
	Status_Function0 = 3
}
