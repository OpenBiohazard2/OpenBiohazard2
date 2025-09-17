package ui_render

import (
	"image"

	"github.com/OpenBiohazard2/OpenBiohazard2/render"
)

// UIRenderer handles UI-specific rendering operations
type UIRenderer struct {
	renderDef *render.RenderDef
}

// NewUIRenderer creates a new UI renderer
func NewUIRenderer(renderDef *render.RenderDef) *UIRenderer {
	return &UIRenderer{
		renderDef: renderDef,
	}
}

// Helper functions for UI rendering
func (r *UIRenderer) ClearScreen() {
	r.renderDef.ScreenImageManager.Clear()
}

func (r *UIRenderer) GetScreenImage() *render.Image16Bit {
	return r.renderDef.ScreenImageManager.GetScreenImage()
}

func (r *UIRenderer) UpdateVideoBuffer(screenImage *render.Image16Bit) {
	r.renderDef.VideoBuffer.UpdateSurface(screenImage)
}

// UpdateMainMenu renders the main menu
func (r *UIRenderer) UpdateMainMenu(
	menuBackgroundImage *render.Image16Bit,
	menuTextImages []*render.Image16Bit,
	mainMenuOption int,
) {
	r.ClearScreen()
	screenImage := r.GetScreenImage()
	buildMainMenuBackground(screenImage, menuBackgroundImage)
	buildMainMenuText(screenImage, menuTextImages, mainMenuOption)
	r.UpdateVideoBuffer(screenImage)
}

// UpdateSpecialMenu renders the special menu
func (r *UIRenderer) UpdateSpecialMenu(
	menuBackgroundImage *render.Image16Bit,
	menuTextImages []*render.Image16Bit,
	mainMenuOption int,
) {
	r.ClearScreen()
	screenImage := r.GetScreenImage()
	buildMainMenuBackground(screenImage, menuBackgroundImage)
	buildSpecialMenuText(screenImage, menuTextImages, mainMenuOption)
	r.UpdateVideoBuffer(screenImage)
}

func buildMainMenuBackground(screenImage *render.Image16Bit, menuBackgroundImage *render.Image16Bit) {
	screenImage.WriteSubImage(image.Point{0, 0}, menuBackgroundImage, image.Rect(0, 0, render.BACKGROUND_IMAGE_WIDTH, render.BACKGROUND_IMAGE_HEIGHT))
}

func buildMainMenuText(screenImage *render.Image16Bit, menuTextImages []*render.Image16Bit, mainMenuOption int) {
	buildTitleText(screenImage, menuTextImages)
	buildMainMenuOptions(screenImage, menuTextImages, mainMenuOption)
}

func buildTitleText(screenImage *render.Image16Bit, menuTextImages []*render.Image16Bit) {
	screenImage.WriteSubImage(image.Point{18, 30}, menuTextImages[1], image.Rect(0, 0, 128, 81))
	screenImage.WriteSubImage(image.Point{146, 31}, menuTextImages[1], image.Rect(0, 81, 128, 81+47))
	screenImage.WriteSubImage(image.Point{146, 78}, menuTextImages[2], image.Rect(0, 0, 128, 34))
	screenImage.WriteSubImage(image.Point{274, 31}, menuTextImages[2], image.Rect(0, 34, 46, 34+82))
}

func buildMainMenuOptions(screenImage *render.Image16Bit, menuTextImages []*render.Image16Bit, mainMenuOption int) {
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

func buildSpecialMenuText(screenImage *render.Image16Bit, menuTextImages []*render.Image16Bit, mainMenuOption int) {
	buildTitleText(screenImage, menuTextImages)
	buildSpecialMenuOptions(screenImage, menuTextImages, mainMenuOption)
}

func buildSpecialMenuOptions(screenImage *render.Image16Bit, menuTextImages []*render.Image16Bit, specialMenuOption int) {
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
