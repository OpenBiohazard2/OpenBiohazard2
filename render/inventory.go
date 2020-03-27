package render

import (
	"github.com/samuelyuan/openbiohazard2/fileio"
)

const (
	HEALTH_FINE           = 0
	HEALTH_YELLOW_CAUTION = 1
	HEALTH_ORANGE_CAUTION = 2
	HEALTH_DANGER         = 3
	HEALTH_POISON         = 4
)

var (
	totalInventoryTime = float64(0)
	updateHealthTime   = float64(30) // milliseconds
	ecgOffsetX         = 0
	healthECGViews     = [5]HealthECGView{
		NewHealthECGFine(),
		NewHealthECGYellowCaution(),
		NewHealthECGOrangeCaution(),
		NewHealthECGDanger(),
		NewHealthECGPoison(),
	}
)

type HealthECGView struct {
	Color    [3]int
	Gradient [3]int
	Lines    [80][2]int
}

func (renderDef *RenderDef) GenerateInventoryImage(
	inventoryImages []*fileio.TIMOutput,
	inventoryItemImages []*fileio.TIMOutput,
	timeElapsedSeconds float64) {
	renderDef.VideoBuffer.ClearSurface()
	newImageColors := renderDef.VideoBuffer.ImagePixels
	totalInventoryTime += timeElapsedSeconds * 1000
	buildBackground(inventoryImages, newImageColors)
	buildItems(inventoryItemImages, newImageColors)
	renderDef.VideoBuffer.UpdateSurface(newImageColors)
}

func buildItems(inventoryItemImages []*fileio.TIMOutput, newImageColors []uint16) {
	// Item in top right corner
	copyPixels(inventoryItemImages[0].PixelData, 0, 0, 40, 30, newImageColors, 265, 35)

	// Empty inventory slots
	copyPixels(inventoryItemImages[0].PixelData, 0, 0, 40, 30, newImageColors, 225, 73)
	copyPixels(inventoryItemImages[0].PixelData, 0, 0, 40, 30, newImageColors, 265, 73)
	copyPixels(inventoryItemImages[0].PixelData, 0, 0, 40, 30, newImageColors, 225, 103)
	copyPixels(inventoryItemImages[0].PixelData, 0, 0, 40, 30, newImageColors, 265, 103)
	copyPixels(inventoryItemImages[0].PixelData, 0, 0, 40, 30, newImageColors, 225, 133)
	copyPixels(inventoryItemImages[0].PixelData, 0, 0, 40, 30, newImageColors, 265, 133)
	copyPixels(inventoryItemImages[0].PixelData, 0, 0, 40, 30, newImageColors, 225, 163)
	copyPixels(inventoryItemImages[0].PixelData, 0, 0, 40, 30, newImageColors, 265, 163)

	// Equipped item
	copyPixels(inventoryItemImages[2].PixelData, 40, 90, 80, 30, newImageColors, 172, 35)
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
	copyPixels(inventoryImages[0].PixelData, 114, 92, 5, 120, newImageColors, 220, 73) // left
	copyPixels(inventoryImages[0].PixelData, 0, 140, 90, 3, newImageColors, 220, 70)   // top
	copyPixels(inventoryImages[0].PixelData, 114, 92, 5, 120, newImageColors, 305, 73) // right
	copyPixels(inventoryImages[0].PixelData, 0, 140, 90, 4, newImageColors, 220, 193)  // bottom

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
	copyPixels(inventoryImages[0].PixelData, 56, 186, 7, 7, newImageColors, 53, 60)
}

func buildHealthECG(inventoryImages []*fileio.TIMOutput, newImageColors []uint16, backgroundColor [3]int) {
	// Draw health background
	copyPixels(inventoryImages[0].PixelData, 0, 92, 99, 47, newImageColors, 60, 29)
	// Sloped line to the right of Condition
	for i := 0; i < 8; i++ {
		fillPixels(newImageColors, 129-i, 68+i, 30+i, 1, backgroundColor[0], backgroundColor[1], backgroundColor[2])
	}

	// Draw ECG lines
	healthStatus := HEALTH_FINE
	ecgView := healthECGViews[healthStatus]

	if totalInventoryTime >= updateHealthTime {
		ecgOffsetX = (ecgOffsetX + 1) % 128
		totalInventoryTime = 0
	}

	for columnNum := 0; columnNum < 32; columnNum++ {
		startX := ecgOffsetX - columnNum
		if startX < 0 || startX >= len(ecgView.Lines) {
			continue
		}
		// Draw a vertical line
		destX := startX + 70
		destY := ecgView.Lines[startX][0] + 31
		width := 1
		height := ecgView.Lines[startX][1] + 1

		// lines to the left will have a darker color
		lineColor := ecgView.Color
		gradientColor := ecgView.Gradient
		red := lineColor[0] - (gradientColor[0] * columnNum)
		if red < 0 {
			red = 0
		}
		green := lineColor[1] - (gradientColor[1] * columnNum)
		if green < 0 {
			green = 0
		}
		blue := lineColor[2] - (gradientColor[2] * columnNum)
		if blue < 0 {
			blue = 0
		}
		fillPixels(newImageColors, destX, destY, width, height, red, green, blue)
	}
}

func NewHealthECGFine() HealthECGView {
	lines := [80][2]int{
		{15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0},
		{15, 0}, {15, 0}, {15, 0}, {14, 0}, {13, 0}, {12, 0}, {12, 0}, {13, 2}, {15, 3}, {18, 2},
		{20, 0}, {16, 4}, {8, 8}, {5, 3}, {4, 0}, {5, 3}, {8, 7}, {15, 4}, {19, 5}, {24, 3},
		{27, 0}, {25, 2}, {21, 4}, {16, 5}, {14, 2}, {13, 0}, {14, 2}, {16, 3}, {19, 0}, {19, 0},
		{18, 0}, {16, 2}, {14, 2}, {13, 0}, {12, 0}, {13, 0}, {14, 1}, {15, 0}, {15, 0}, {15, 0},
		{15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0},
		{15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0},
		{15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0},
	}

	return HealthECGView{
		Color:    [3]int{20, 255, 20}, // green,
		Gradient: [3]int{1, 8, 1},
		Lines:    lines,
	}
}

func NewHealthECGYellowCaution() HealthECGView {
	lines := [80][2]int{
		{15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0},
		{15, 0}, {15, 0}, {14, 0}, {13, 0}, {12, 0}, {12, 0}, {14, 0}, {13, 2}, {15, 3}, {18, 2},
		{20, 0}, {16, 4}, {8, 8}, {6, 2}, {5, 0}, {6, 2}, {8, 5}, {13, 2}, {15, 0}, {16, 4},
		{20, 2}, {22, 0}, {21, 0}, {16, 5}, {15, 0}, {14, 0}, {14, 0}, {13, 0}, {13, 0}, {14, 0},
		{15, 0}, {15, 0}, {15, 0}, {14, 0}, {14, 0}, {14, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0},
		{15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0},
		{15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0},
		{15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0},
	}

	return HealthECGView{
		Color:    [3]int{255, 255, 20}, // yellow
		Gradient: [3]int{8, 8, 1},
		Lines:    lines,
	}
}

func NewHealthECGOrangeCaution() HealthECGView {
	lines := [80][2]int{
		{15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0},
		{15, 0}, {15, 0}, {16, 0}, {16, 0}, {17, 0}, {17, 0}, {17, 0}, {16, 0}, {15, 0}, {14, 0},
		{14, 0}, {14, 0}, {15, 0}, {15, 0}, {15, 0}, {14, 0}, {11, 3}, {10, 0}, {9, 0}, {10, 4},
		{13, 3}, {16, 3}, {19, 0}, {20, 0}, {19, 0}, {18, 0}, {16, 2}, {14, 2}, {13, 0}, {12, 0},
		{12, 0}, {12, 0}, {13, 0}, {14, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0},
		{15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0},
		{15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0},
		{15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0},
	}

	return HealthECGView{
		Color:    [3]int{255, 80, 20}, // orange
		Gradient: [3]int{8, 4, 1},
		Lines:    lines,
	}
}

func NewHealthECGDanger() HealthECGView {
	lines := [80][2]int{
		{15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0},
		{15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {14, 0},
		{14, 0}, {13, 0}, {13, 0}, {14, 0}, {15, 0}, {15, 0}, {15, 0}, {16, 0}, {17, 2}, {17, 0},
		{17, 0}, {14, 3}, {10, 4}, {9, 0}, {10, 2}, {12, 3}, {15, 0}, {16, 0}, {16, 0}, {16, 0},
		{16, 0}, {15, 0}, {14, 0}, {13, 0}, {13, 0}, {14, 0}, {14, 0}, {15, 0}, {15, 0}, {15, 0},
		{15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0},
		{15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0},
		{15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0},
	}

	return HealthECGView{
		Color:    [3]int{255, 20, 20}, // red
		Gradient: [3]int{8, 1, 1},
		Lines:    lines,
	}
}

func NewHealthECGPoison() HealthECGView {
	lines := [80][2]int{
		{15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0},
		{15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {14, 1}, {13, 1}, {12, 1}, {12, 1},
		{13, 2}, {15, 3}, {18, 2}, {20, 1}, {16, 4}, {8, 8}, {5, 3}, {4, 1}, {5, 3}, {8, 7},
		{15, 4}, {19, 5}, {24, 2}, {26, 1}, {25, 1}, {15, 10}, {14, 1}, {15, 1}, {16, 2}, {18, 1},
		{17, 1}, {10, 7}, {9, 1}, {10, 2}, {12, 4}, {16, 1}, {17, 1}, {18, 1}, {18, 1}, {18, 1},
		{17, 1}, {16, 1}, {15, 1}, {15, 1}, {15, 1}, {15, 1}, {15, 1}, {12, 3}, {10, 2}, {9, 1},
		{10, 6}, {16, 3}, {19, 1}, {19, 1}, {17, 2}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0},
		{15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0}, {15, 0},
	}

	return HealthECGView{
		Color:    [3]int{255, 20, 255}, // purple
		Gradient: [3]int{8, 1, 8},
		Lines:    lines,
	}
}

func buildMenuTabs(inventoryImages []*fileio.TIMOutput, newImageColors []uint16) {
	// File
	copyPixels(inventoryImages[0].PixelData, 3, 164, 49, 12, newImageColors, 111, 16)

	// Map
	copyPixels(inventoryImages[0].PixelData, 3, 164, 49, 12, newImageColors, 160, 16)

	// Item
	copyPixels(inventoryImages[0].PixelData, 3, 164, 49, 12, newImageColors, 209, 16)

	// Exit
	copyPixels(inventoryImages[0].PixelData, 3, 164, 49, 12, newImageColors, 258, 16)
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
