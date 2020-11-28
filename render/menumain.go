package render

import (
	"image"
)

func (renderDef *RenderDef) UpdateMainMenu(
	menuBackgroundImage *Image16Bit,
	menuTextImages []*Image16Bit,
	mainMenuOption int,
) {
	screenImage.Clear()
	buildMainMenuBackground(menuBackgroundImage)
	buildMainMenuText(menuTextImages, mainMenuOption)
	renderDef.VideoBuffer.UpdateSurface(screenImage.GetPixelsForRendering())
}

func buildMainMenuBackground(menuBackgroundImage *Image16Bit) {
	screenImage.WriteSubImage(image.Point{0, 0}, menuBackgroundImage, image.Rect(0, 0, 320, 240))
}

func buildMainMenuText(menuTextImages []*Image16Bit, mainMenuOption int) {
	buildTitleText(menuTextImages)
	buildMainMenuOptions(menuTextImages, mainMenuOption)
}

func buildTitleText(menuTextImages []*Image16Bit) {
	screenImage.WriteSubImage(image.Point{18, 30}, menuTextImages[1], image.Rect(0, 0, 128, 81))
	screenImage.WriteSubImage(image.Point{146, 31}, menuTextImages[1], image.Rect(0, 81, 128, 81+47))
	screenImage.WriteSubImage(image.Point{146, 78}, menuTextImages[2], image.Rect(0, 0, 128, 34))
	screenImage.WriteSubImage(image.Point{274, 31}, menuTextImages[2], image.Rect(0, 34, 46, 34+82))
}

func buildMainMenuOptions(menuTextImages []*Image16Bit, mainMenuOption int) {
	selectedOption := 1.0
	otherOption := 0.3

	optionsBrightness := [4]float64{otherOption, otherOption, otherOption, otherOption}
	optionsBrightness[mainMenuOption] = selectedOption

	// Load Game
	screenImage.WriteSubImageUniformBrightness(image.Point{114, 134}, menuTextImages[0], image.Rect(70, 29, 70+106, 29+13),
		optionsBrightness[0])

	// New Game
	screenImage.WriteSubImageUniformBrightness(image.Point{95, 154}, menuTextImages[0], image.Rect(54, 17, 54+147, 17+12),
		optionsBrightness[1])

	// Special
	screenImage.WriteSubImageUniformBrightness(image.Point{130, 174}, menuTextImages[0], image.Rect(94, 96, 94+74, 96+14),
		optionsBrightness[2])

	// Option
	screenImage.WriteSubImageUniformBrightness(image.Point{125, 194}, menuTextImages[0], image.Rect(88, 43, 88+74, 43+14),
		optionsBrightness[3])
}

func (renderDef *RenderDef) UpdateSpecialMenu(
	menuBackgroundImage *Image16Bit,
	menuTextImages []*Image16Bit,
	mainMenuOption int,
) {
	screenImage.Clear()
	buildMainMenuBackground(menuBackgroundImage)
	buildSpecialMenuText(menuTextImages, mainMenuOption)
	renderDef.VideoBuffer.UpdateSurface(screenImage.GetPixelsForRendering())
}

func buildSpecialMenuText(menuTextImages []*Image16Bit, mainMenuOption int) {
	buildTitleText(menuTextImages)
	buildSpecialMenuOptions(menuTextImages, mainMenuOption)
}

func buildSpecialMenuOptions(menuTextImages []*Image16Bit, specialMenuOption int) {
	selectedOption := 1.0
	otherOption := 0.3

	optionsBrightness := [2]float64{otherOption, otherOption}
	optionsBrightness[specialMenuOption] = selectedOption

	// Special title
	screenImage.WriteSubImageUniformBrightness(image.Point{125, 134}, menuTextImages[0], image.Rect(94, 96, 94+74, 96+14),
		otherOption)

	// Gallery
	screenImage.WriteSubImageUniformBrightness(image.Point{120, 154}, menuTextImages[0], image.Rect(169, 96, 169+75, 96+14),
		optionsBrightness[0])

	// Exit
	screenImage.WriteSubImageUniformBrightness(image.Point{135, 174}, menuTextImages[0], image.Rect(105, 124, 105+45, 124+14),
		optionsBrightness[1])
}
