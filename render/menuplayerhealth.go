package render

import (
	"image"
	"image/color"
)

const (
	HEALTH_FINE           = 0
	HEALTH_YELLOW_CAUTION = 1
	HEALTH_ORANGE_CAUTION = 2
	HEALTH_DANGER         = 3
	HEALTH_POISON         = 4
	HEALTH_POS_X          = 58
	HEALTH_POS_Y          = 29
)

var (
	totalHealthTime  = float64(0)
	updateHealthTime = float64(30) // milliseconds
	ecgOffsetX       = 0
	healthECGViews   = [5]HealthECGView{
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

func buildHealthECG(inventoryMenuImages []*Image16Bit, backgroundColor color.RGBA) {
	// Draw health background
	screenImage.WriteSubImage(image.Point{HEALTH_POS_X + 2, HEALTH_POS_Y}, inventoryMenuImages[0], image.Rect(0, 92, 99, 92+47))
	// Sloped line to the right of Condition
	for i := 0; i < 8; i++ {
		screenImage.FillPixels(image.Point{129 - i, 68 + i}, image.Rect(129-i, 68+i, 159, 69+i), backgroundColor)
	}

	// Draw ECG lines
	healthStatus := HEALTH_FINE
	ecgView := healthECGViews[healthStatus]

	if totalHealthTime >= updateHealthTime {
		ecgOffsetX = (ecgOffsetX + 1) % 128
		totalHealthTime = 0
	}

	for columnNum := 0; columnNum < 32; columnNum++ {
		startX := ecgOffsetX - columnNum
		if startX < 0 || startX >= len(ecgView.Lines) {
			continue
		}
		// Draw a vertical line
		destX := startX + HEALTH_POS_X + 12
		destY := ecgView.Lines[startX][0] + HEALTH_POS_Y + 2
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
		finalColor := color.RGBA{uint8(red), uint8(green), uint8(blue), 255}
		screenImage.FillPixels(image.Point{destX, destY}, image.Rect(destX, destY, destX+width, destY+height), finalColor)
	}

	drawPlayerCondition(inventoryMenuImages, healthStatus)
}

func drawPlayerCondition(inventoryMenuImages []*Image16Bit, healthStatus int) {
	screenImage.WriteSubImage(image.Point{HEALTH_POS_X + 47, HEALTH_POS_Y + 25},
		inventoryMenuImages[4], image.Rect(0, healthStatus*11, 44, (healthStatus+1)*11))
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
