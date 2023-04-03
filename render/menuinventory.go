package render

import (
	"image"
	"image/color"

	"github.com/OpenBiohazard2/OpenBiohazard2/gui"
)

const (
	ITEMLIST_POS_X = 220
	ITEMLIST_POS_Y = 70
)

var (
	totalInventoryTime        = float64(0)
	updateInventoryCursorTime = float64(30) // milliseconds

	playerInventoryItems = InitializeInventoryItems()
)

type InventoryItem struct {
	Id   int
	Num  int
	Size int
}

func InitializeInventoryItems() []InventoryItem {
	playerInventoryItems := make([]InventoryItem, 11)
	playerInventoryItems[0] = InventoryItem{Id: 2, Num: 18, Size: 1}                      // hand gun
	playerInventoryItems[1] = InventoryItem{Id: 1, Num: 1, Size: 0}                       // knife
	playerInventoryItems[gui.RESERVED_ITEM_SLOT] = InventoryItem{Id: 47, Num: 1, Size: 0} // lighter
	return playerInventoryItems
}

func (renderDef *RenderDef) GenerateInventoryImage(
	inventoryMenuImages []*Image16Bit,
	inventoryItemImages []*Image16Bit,
	inventoryMenu *gui.InventoryMenu,
	timeElapsedSeconds float64,
) {
	screenImage.Clear()
	totalInventoryTime += timeElapsedSeconds * 1000
	totalHealthTime += timeElapsedSeconds * 1000
	buildBackground(inventoryMenuImages, inventoryMenu)
	buildItems(inventoryMenuImages, inventoryItemImages, inventoryMenu)
	renderDef.VideoBuffer.UpdateSurface(screenImage)
}

func buildItems(
	inventoryMenuImages []*Image16Bit,
	inventoryItemImages []*Image16Bit,
	inventoryMenu *gui.InventoryMenu,
) {
	// Item in top right corner
	reservedItemX := (playerInventoryItems[gui.RESERVED_ITEM_SLOT].Id % 6) * 40
	reservedItemY := (playerInventoryItems[gui.RESERVED_ITEM_SLOT].Id / 6) * 30
	screenImage.WriteSubImage(image.Point{ITEMLIST_POS_X + 45, ITEMLIST_POS_Y - 35},
		inventoryItemImages[0], image.Rect(reservedItemX, reservedItemY, reservedItemX+40, reservedItemY+30))

	// Empty inventory slots
	for row := 0; row < 4; row++ {
		// left slot
		leftItemIndex := (2 * row)
		leftItemX := (playerInventoryItems[leftItemIndex].Id % 6) * 40
		leftItemY := (playerInventoryItems[leftItemIndex].Id / 6) * 30
		screenImage.WriteSubImage(image.Point{ITEMLIST_POS_X + 5, ITEMLIST_POS_Y + 3 + 30*row},
			inventoryItemImages[0], image.Rect(leftItemX, leftItemY, leftItemX+40, leftItemY+30))
		// right slot
		rightItemIndex := (2 * row) + 1
		rightItemX := (playerInventoryItems[rightItemIndex].Id % 6) * 40
		rightItemY := (playerInventoryItems[rightItemIndex].Id / 6) * 30
		screenImage.WriteSubImage(image.Point{ITEMLIST_POS_X + 45, ITEMLIST_POS_Y + 3 + 30*row},
			inventoryItemImages[0], image.Rect(rightItemX, rightItemY, rightItemX+40, rightItemY+30))
	}

	// Equipped item
	screenImage.WriteSubImage(image.Point{172, 35}, inventoryItemImages[2], image.Rect(40, 90, 40+80, 90+30))

	// Item cursor surrounding item
	if inventoryMenu.IsEditingItemScreen() {
		displayInventoryMainCursor(inventoryMenuImages, inventoryMenu)
	}
}

func displayInventoryMainCursor(inventoryMenuImages []*Image16Bit, inventoryMenu *gui.InventoryMenu) {
	var cursorX, cursorY int
	cursorFrameOffsetX := 3
	cursorFrameOffsetY := 1

	if totalInventoryTime >= updateInventoryCursorTime {
		inventoryMenu.UpdateCursorBlink()
		totalInventoryTime = 0
	}
	brightnessFactor := inventoryMenu.GetCursorBlinkBrightnessFactor()

	// Special item in top right corner
	if inventoryMenu.IsCursorOnReservedItem() {
		cursorX = ITEMLIST_POS_X + cursorFrameOffsetX + 40
		cursorY = ITEMLIST_POS_Y + cursorFrameOffsetY - 38
	} else {
		cursorX = ITEMLIST_POS_X + cursorFrameOffsetX + inventoryMenu.GetMainCursorColumn()*40
		cursorY = ITEMLIST_POS_Y + cursorFrameOffsetY + inventoryMenu.GetMainCursorRow()*30
	}
	screenImage.WriteSubImageUniformBrightness(image.Point{cursorX, cursorY},
		inventoryMenuImages[3], image.Rect(0, 30, 44, 30+34), brightnessFactor)
}

func buildBackground(inventoryMenuImages []*Image16Bit, inventoryMenu *gui.InventoryMenu) {
	// The inventory image is split up into many small components
	// Combine them manually back into a single image
	// source image is 256x256
	// dest image is 320x240
	backgroundColor := color.RGBA{5, 5, 31, 255}
	screenImage.FillPixels(image.Point{0, 0}, image.Rect(0, 0, 320, 240), backgroundColor)

	buildPlayerFace(inventoryMenuImages)
	buildHealthECG(inventoryMenuImages, backgroundColor)

	// Equipped item
	screenImage.WriteSubImage(image.Point{161, 29}, inventoryMenuImages[0], image.Rect(50, 211, 50+11, 211+39)) // left
	screenImage.WriteSubImage(image.Point{172, 29}, inventoryMenuImages[0], image.Rect(0, 158, 80, 158+6))      // top
	screenImage.WriteSubImage(image.Point{252, 29}, inventoryMenuImages[0], image.Rect(91, 164, 91+5, 164+39))  // right
	screenImage.WriteSubImage(image.Point{172, 65}, inventoryMenuImages[0], image.Rect(0, 155, 80, 155+3))      // bottom

	// Extra item
	screenImage.WriteSubImage(image.Point{260, 29}, inventoryMenuImages[0], image.Rect(0, 211, 50, 211+41))

	buildMenuTabs(inventoryMenuImages, inventoryMenu)

	// Item slots
	screenImage.WriteSubImage(image.Point{ITEMLIST_POS_X, ITEMLIST_POS_Y + 3},
		inventoryMenuImages[0], image.Rect(114, 92, 114+5, 92+120)) // left
	screenImage.WriteSubImage(image.Point{ITEMLIST_POS_X, ITEMLIST_POS_Y},
		inventoryMenuImages[0], image.Rect(0, 140, 90, 140+3)) // top
	screenImage.WriteSubImage(image.Point{ITEMLIST_POS_X + 85, ITEMLIST_POS_Y + 3},
		inventoryMenuImages[0], image.Rect(114, 92, 114+5, 92+120)) // right
	screenImage.WriteSubImage(image.Point{ITEMLIST_POS_X, ITEMLIST_POS_Y + 123},
		inventoryMenuImages[0], image.Rect(0, 140, 90, 140+4)) // bottom

	buildDescription(inventoryMenuImages)
}

func buildPlayerFace(inventoryMenuImages []*Image16Bit) {
	// Player
	screenImage.WriteSubImage(image.Point{7, 16}, inventoryMenuImages[0], image.Rect(106, 152, 106+4, 152+60))  // left
	screenImage.WriteSubImage(image.Point{11, 16}, inventoryMenuImages[0], image.Rect(0, 140, 39, 140+4))       // top
	screenImage.WriteSubImage(image.Point{49, 16}, inventoryMenuImages[0], image.Rect(109, 152, 109+4, 152+60)) // right
	screenImage.WriteSubImage(image.Point{11, 72}, inventoryMenuImages[0], image.Rect(0, 140, 39, 140+4))       // bottom
	screenImage.WriteSubImage(image.Point{11, 21}, inventoryMenuImages[1], image.Rect(1, 73, 1+37, 73+8))       // player name
	screenImage.WriteSubImage(image.Point{11, 31}, inventoryMenuImages[1], image.Rect(0, 85, 38, 85+42))        // player image
	screenImage.WriteSubImage(image.Point{11, 30}, inventoryMenuImages[0], image.Rect(56, 164, 56+38, 164+1))   // line between name and image

	// Pipes to the left of player image
	screenImage.WriteSubImage(image.Point{0, 17}, inventoryMenuImages[0], image.Rect(107, 242, 107+7, 242+14))
	screenImage.WriteSubImage(image.Point{0, 33}, inventoryMenuImages[0], image.Rect(107, 242, 107+7, 242+14))
	screenImage.WriteSubImage(image.Point{0, 49}, inventoryMenuImages[0], image.Rect(107, 242, 107+7, 242+14))

	// Pipes to the right of player image
	screenImage.WriteSubImage(image.Point{53, 32}, inventoryMenuImages[0], image.Rect(56, 186, 56+7, 186+7))
	screenImage.WriteSubImage(image.Point{53, HEALTH_POS_X + 2}, inventoryMenuImages[0], image.Rect(56, 186, 56+7, 186+7))
}

func buildMenuTabs(inventoryMenuImages []*Image16Bit, inventoryMenu *gui.InventoryMenu) {
	var selectedOption [3]float64
	if inventoryMenu.IsCursorOnTopMenu() {
		// Cursor is on this option, but it's not selected
		selectedOption = [3]float64{1.0, 1.0, 1.0}
	} else {
		// Highlight option in red if this option is selected
		selectedOption = [3]float64{1.0, 0.5, 0.5}
	}

	otherOption := [3]float64{0.4, 0.4, 0.4}

	optionsBrightness := [4][3]float64{otherOption, otherOption, otherOption, otherOption}
	optionsBrightness[inventoryMenu.GetTopMenuSelectedOption()] = selectedOption

	// File
	screenImage.WriteSubImage(image.Point{111, 16}, inventoryMenuImages[0], image.Rect(3, 164, 3+49, 164+12))
	screenImage.WriteSubImageVariableBrightness(image.Point{112, 17},
		inventoryMenuImages[5], image.Rect(0, 0, 47, 10), optionsBrightness[0])

	// Map
	screenImage.WriteSubImage(image.Point{160, 16}, inventoryMenuImages[0], image.Rect(3, 164, 3+49, 164+12))
	screenImage.WriteSubImageVariableBrightness(image.Point{161, 17},
		inventoryMenuImages[5], image.Rect(0, 10, 47, 10+10), optionsBrightness[1])

	// Item
	screenImage.WriteSubImage(image.Point{209, 16}, inventoryMenuImages[0], image.Rect(3, 164, 3+49, 164+12))
	screenImage.WriteSubImageVariableBrightness(image.Point{210, 17},
		inventoryMenuImages[5], image.Rect(0, 20, 47, 20+10), optionsBrightness[2])

	// Exit
	screenImage.WriteSubImage(image.Point{258, 16}, inventoryMenuImages[0], image.Rect(3, 164, 3+49, 164+12))
	screenImage.WriteSubImageVariableBrightness(image.Point{259, 17},
		inventoryMenuImages[5], image.Rect(0, 30, 47, 30+10), optionsBrightness[3])
}

func buildDescription(inventoryMenuImages []*Image16Bit) {
	descriptionColor := color.RGBA{6, 13, 23, 255}
	screenImage.FillPixels(image.Point{13, 174}, image.Rect(13, 174, 13+201, 174+49), descriptionColor)
	screenImage.WriteSubImage(image.Point{8, 174}, inventoryMenuImages[0], image.Rect(106, 163, 106+5, 163+49))   // left
	screenImage.WriteSubImage(image.Point{8, 170}, inventoryMenuImages[0], image.Rect(0, 147, 83, 147+4))         // top left
	screenImage.WriteSubImage(image.Point{91, 170}, inventoryMenuImages[0], image.Rect(0, 80, 128, 80+4))         // top right
	screenImage.WriteSubImage(image.Point{214, 174}, inventoryMenuImages[0], image.Rect(106, 163, 106+5, 163+49)) // right
	screenImage.WriteSubImage(image.Point{8, 223}, inventoryMenuImages[0], image.Rect(0, 147, 83, 147+4))         // bottom left
	screenImage.WriteSubImage(image.Point{91, 223}, inventoryMenuImages[0], image.Rect(0, 80, 128, 80+4))         // bottom right

	// Pipes to the right of description
	screenImage.WriteSubImage(image.Point{219, 212}, inventoryMenuImages[0], image.Rect(107, 242, 107+7, 242+14))
	screenImage.WriteSubImage(image.Point{226, 215}, inventoryMenuImages[0], image.Rect(56, 178, 56+35, 178+7))
	screenImage.WriteSubImage(image.Point{261, 215}, inventoryMenuImages[0], image.Rect(56, 178, 56+35, 178+7))
	screenImage.WriteSubImage(image.Point{296, 215}, inventoryMenuImages[0], image.Rect(56, 178, 56+24, 178+7))
}
