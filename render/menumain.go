package render

import (
	"github.com/samuelyuan/openbiohazard2/fileio"
)

func (renderDef *RenderDef) GenerateMainMenuImage(
	menuBackgroundImageOutput *fileio.ADTOutput,
	menuBackgroundTextOutput []*fileio.TIMOutput) {
	renderDef.VideoBuffer.ClearSurface()
	newImageColors := renderDef.VideoBuffer.ImagePixels
	buildMainMenuBackground(menuBackgroundImageOutput, newImageColors)
	buildMainMenuText(menuBackgroundTextOutput, newImageColors, 0)
	renderDef.VideoBuffer.UpdateSurface(newImageColors)
}

func buildMainMenuBackground(backgroundImageOutput *fileio.ADTOutput, newImageColors []uint16) {
	copyPixelsTransparent(backgroundImageOutput.PixelData, 0, 0, 320, 240, newImageColors, 0, 0)
}

func buildMainMenuText(menuBackgroundTextOutput []*fileio.TIMOutput, newImageColors []uint16, mainMenuOption int) {
	buildTitleText(menuBackgroundTextOutput, newImageColors)
	buildMainMenuOptions(menuBackgroundTextOutput, newImageColors, mainMenuOption)
}

func buildTitleText(menuBackgroundTextOutput []*fileio.TIMOutput, newImageColors []uint16) {
	copyPixelsTransparent(menuBackgroundTextOutput[1].PixelData, 0, 0, 128, 81, newImageColors, 18, 30)
	copyPixelsTransparent(menuBackgroundTextOutput[1].PixelData, 0, 81, 128, 47, newImageColors, 146, 31)
	copyPixelsTransparent(menuBackgroundTextOutput[2].PixelData, 0, 0, 128, 34, newImageColors, 146, 78)
	copyPixelsTransparent(menuBackgroundTextOutput[2].PixelData, 0, 34, 46, 82, newImageColors, 274, 31)
}

func buildMainMenuOptions(menuBackgroundTextOutput []*fileio.TIMOutput, newImageColors []uint16, mainMenuOption int) {
	selectedOption := 1.0
	otherOption := 0.3

	optionsBrightness := [4]float64{otherOption, otherOption, otherOption, otherOption}
	optionsBrightness[mainMenuOption] = selectedOption

	// Load Game
	copyPixelsBrightness(menuBackgroundTextOutput[0].PixelData, 70, 29, 106, 13, newImageColors, 114, 134, optionsBrightness[0])

	// New Game
	copyPixelsBrightness(menuBackgroundTextOutput[0].PixelData, 54, 17, 147, 12, newImageColors, 95, 154, optionsBrightness[1])

	// Special
	copyPixelsBrightness(menuBackgroundTextOutput[0].PixelData, 94, 96, 74, 14, newImageColors, 130, 174, optionsBrightness[2])

	// Option
	copyPixelsBrightness(menuBackgroundTextOutput[0].PixelData, 88, 43, 74, 14, newImageColors, 125, 194, optionsBrightness[3])
}

func (renderDef *RenderDef) UpdateMainMenu(
	menuBackgroundImageOutput *fileio.ADTOutput,
	menuBackgroundTextOutput []*fileio.TIMOutput,
	mainMenuOption int) {
	renderDef.VideoBuffer.ClearSurface()
	newImageColors := renderDef.VideoBuffer.ImagePixels
	buildMainMenuBackground(menuBackgroundImageOutput, newImageColors)
	buildMainMenuText(menuBackgroundTextOutput, newImageColors, mainMenuOption)
	renderDef.VideoBuffer.UpdateSurface(newImageColors)
}

func (renderDef *RenderDef) GenerateSpecialMenuImage(
	menuBackgroundImageOutput *fileio.ADTOutput,
	menuBackgroundTextOutput []*fileio.TIMOutput) {
	renderDef.VideoBuffer.ClearSurface()
	newImageColors := renderDef.VideoBuffer.ImagePixels
	buildMainMenuBackground(menuBackgroundImageOutput, newImageColors)
	buildSpecialMenuText(menuBackgroundTextOutput, newImageColors, 0)
	renderDef.VideoBuffer.UpdateSurface(newImageColors)
}

func buildSpecialMenuText(menuBackgroundTextOutput []*fileio.TIMOutput, newImageColors []uint16, mainMenuOption int) {
	buildTitleText(menuBackgroundTextOutput, newImageColors)
	buildSpecialMenuOptions(menuBackgroundTextOutput, newImageColors, mainMenuOption)
}

func buildSpecialMenuOptions(menuBackgroundTextOutput []*fileio.TIMOutput, newImageColors []uint16, specialMenuOption int) {
	selectedOption := 1.0
	otherOption := 0.3

	optionsBrightness := [2]float64{otherOption, otherOption}
	optionsBrightness[specialMenuOption] = selectedOption

	// Special title
	copyPixelsBrightness(menuBackgroundTextOutput[0].PixelData, 94, 96, 74, 14, newImageColors, 125, 134, otherOption)

	// Gallery
	copyPixelsBrightness(menuBackgroundTextOutput[0].PixelData, 169, 96, 75, 14, newImageColors, 120, 154, optionsBrightness[0])

	// Exit
	copyPixelsBrightness(menuBackgroundTextOutput[0].PixelData, 105, 124, 45, 14, newImageColors, 135, 174, optionsBrightness[1])
}

func (renderDef *RenderDef) UpdateSpecialMenu(
	menuBackgroundImageOutput *fileio.ADTOutput,
	menuBackgroundTextOutput []*fileio.TIMOutput,
	mainMenuOption int) {
	renderDef.VideoBuffer.ClearSurface()
	newImageColors := renderDef.VideoBuffer.ImagePixels
	buildMainMenuBackground(menuBackgroundImageOutput, newImageColors)
	buildSpecialMenuText(menuBackgroundTextOutput, newImageColors, mainMenuOption)
	renderDef.VideoBuffer.UpdateSurface(newImageColors)
}
