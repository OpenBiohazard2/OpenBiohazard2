package ui_render

import (
	"image"
	"image/color"

	"github.com/OpenBiohazard2/OpenBiohazard2/geometry"
	"github.com/OpenBiohazard2/OpenBiohazard2/resource"
	"github.com/OpenBiohazard2/OpenBiohazard2/ui"
)

const (
	ITEMLIST_POS_X = 220
	ITEMLIST_POS_Y = 70
	HEALTH_POS_X   = 58
	HEALTH_POS_Y   = 29
)

// GenerateInventoryImage renders the inventory menu
func (r *UIRenderer) GenerateInventoryImage(
	inventoryMenuImages []*resource.Image16Bit,
	inventoryItemImages []*resource.Image16Bit,
	inventoryMenu *ui.InventoryMenu,
	healthDisplay *ui.HealthDisplay,
	inventoryManager *ui.InventoryManager,
	timeElapsedSeconds float64,
) {
	r.ClearScreen()
	screenImage := r.GetScreenImage()
	inventoryManager.UpdateInventoryTime(timeElapsedSeconds)
	healthDisplay.UpdateHealthDisplay(timeElapsedSeconds)
	buildBackground(screenImage, inventoryMenuImages, inventoryMenu, healthDisplay)
	buildItems(screenImage, inventoryMenuImages, inventoryItemImages, inventoryMenu, inventoryManager)
	r.UpdateVideoBuffer(screenImage)
}

func buildItems(
	screenImage *resource.Image16Bit,
	inventoryMenuImages []*resource.Image16Bit,
	inventoryItemImages []*resource.Image16Bit,
	inventoryMenu *ui.InventoryMenu,
	inventoryManager *ui.InventoryManager,
) {
	// Item in top right corner
	playerInventoryItems := inventoryManager.GetPlayerInventoryItems()
	reservedItemX := (playerInventoryItems[ui.RESERVED_ITEM_SLOT].Id % 6) * 40
	reservedItemY := (playerInventoryItems[ui.RESERVED_ITEM_SLOT].Id / 6) * 30
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
		displayInventoryMainCursor(screenImage, inventoryMenuImages, inventoryMenu, inventoryManager)
	}
}

func displayInventoryMainCursor(screenImage *resource.Image16Bit, inventoryMenuImages []*resource.Image16Bit, inventoryMenu *ui.InventoryMenu, inventoryManager *ui.InventoryManager) {
	var cursorX, cursorY int
	cursorFrameOffsetX := 3
	cursorFrameOffsetY := 1

	if inventoryManager.ShouldUpdateCursor() {
		inventoryMenu.UpdateCursorBlink()
		inventoryManager.ResetInventoryTime()
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

func buildBackground(screenImage *resource.Image16Bit, inventoryMenuImages []*resource.Image16Bit, inventoryMenu *ui.InventoryMenu, healthDisplay *ui.HealthDisplay) {
	// The inventory image is split up into many small components
	// Combine them manually back into a single image
	// source image is 256x256
	// dest image is 320x240
	backgroundColor := color.RGBA{5, 5, 31, 255}
	screenImage.FillPixels(image.Point{0, 0}, geometry.BACKGROUND_IMAGE_RECT, backgroundColor)

	buildPlayerFace(screenImage, inventoryMenuImages)
	buildHealthECG(screenImage, healthDisplay, inventoryMenuImages, backgroundColor)

	// Equipped item
	screenImage.WriteSubImage(image.Point{161, 29}, inventoryMenuImages[0], image.Rect(50, 211, 50+11, 211+39)) // left
	screenImage.WriteSubImage(image.Point{172, 29}, inventoryMenuImages[0], image.Rect(0, 158, 80, 158+6))      // top
	screenImage.WriteSubImage(image.Point{252, 29}, inventoryMenuImages[0], image.Rect(91, 164, 91+5, 164+39))  // right
	screenImage.WriteSubImage(image.Point{172, 65}, inventoryMenuImages[0], image.Rect(0, 155, 80, 155+3))      // bottom

	// Extra item
	screenImage.WriteSubImage(image.Point{260, 29}, inventoryMenuImages[0], image.Rect(0, 211, 50, 211+41))

	buildMenuTabs(screenImage, inventoryMenuImages, inventoryMenu)

	// Item slots
	screenImage.WriteSubImage(image.Point{ITEMLIST_POS_X, ITEMLIST_POS_Y + 3},
		inventoryMenuImages[0], image.Rect(114, 92, 114+5, 92+120)) // left
	screenImage.WriteSubImage(image.Point{ITEMLIST_POS_X, ITEMLIST_POS_Y},
		inventoryMenuImages[0], image.Rect(0, 140, 90, 140+3)) // top
	screenImage.WriteSubImage(image.Point{ITEMLIST_POS_X + 85, ITEMLIST_POS_Y + 3},
		inventoryMenuImages[0], image.Rect(114, 92, 114+5, 92+120)) // right
	screenImage.WriteSubImage(image.Point{ITEMLIST_POS_X, ITEMLIST_POS_Y + 123},
		inventoryMenuImages[0], image.Rect(0, 140, 90, 140+4)) // bottom

	buildDescription(screenImage, inventoryMenuImages)
}

func buildPlayerFace(screenImage *resource.Image16Bit, inventoryMenuImages []*resource.Image16Bit) {
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

func buildMenuTabs(screenImage *resource.Image16Bit, inventoryMenuImages []*resource.Image16Bit, inventoryMenu *ui.InventoryMenu) {
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

func buildDescription(screenImage *resource.Image16Bit, inventoryMenuImages []*resource.Image16Bit) {
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

// Health rendering functions

func buildHealthECG(screenImage *resource.Image16Bit, healthDisplay *ui.HealthDisplay, inventoryMenuImages []*resource.Image16Bit, backgroundColor color.RGBA) {
	healthStatus := ui.HEALTH_FINE

	drawHealthBackground(screenImage, inventoryMenuImages, backgroundColor)
	healthDisplay.UpdateECGAnimation()
	drawECGLines(screenImage, healthDisplay, healthStatus)
	drawPlayerCondition(screenImage, inventoryMenuImages, healthStatus)
}

// drawHealthBackground draws the health background and sloped line
func drawHealthBackground(screenImage *resource.Image16Bit, inventoryMenuImages []*resource.Image16Bit, backgroundColor color.RGBA) {
	// Draw health background
	screenImage.WriteSubImage(image.Point{HEALTH_POS_X + 2, HEALTH_POS_Y}, inventoryMenuImages[0], image.Rect(0, 92, 99, 92+47))

	// Sloped line to the right of Condition
	for i := 0; i < 8; i++ {
		screenImage.FillPixels(image.Point{129 - i, 68 + i}, image.Rect(129-i, 68+i, 159, 69+i), backgroundColor)
	}
}

// drawECGLines draws the animated ECG lines with gradient colors
func drawECGLines(screenImage *resource.Image16Bit, healthDisplay *ui.HealthDisplay, healthStatus int) {
	ecgView := healthDisplay.GetHealthECGView(healthStatus)

	for columnNum := 0; columnNum < 32; columnNum++ {
		startX := healthDisplay.GetECGOffsetX() - columnNum
		if startX < 0 || startX >= 80 {
			continue
		}

		// Calculate line position and size
		destX := startX + HEALTH_POS_X + 12
		destY := ecgView.Lines[startX][0] + HEALTH_POS_Y + 2
		width := 1
		height := ecgView.Lines[startX][1] + 1

		// Calculate gradient color
		finalColor := ui.CalculateECGLineColor(ecgView, columnNum)

		// Draw the line
		screenImage.FillPixels(image.Point{destX, destY}, image.Rect(destX, destY, destX+width, destY+height), finalColor)
	}
}

func drawPlayerCondition(screenImage *resource.Image16Bit, inventoryMenuImages []*resource.Image16Bit, healthStatus int) {
	screenImage.WriteSubImage(image.Point{HEALTH_POS_X + 47, HEALTH_POS_Y + 25},
		inventoryMenuImages[4], image.Rect(0, healthStatus*11, 44, (healthStatus+1)*11))
}
