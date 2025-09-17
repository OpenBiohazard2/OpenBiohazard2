package ui_render

import (
	"image"

	"github.com/OpenBiohazard2/OpenBiohazard2/render"
)

// GenerateSaveScreenImage renders the save screen
func (r *UIRenderer) GenerateSaveScreenImage(saveScreenImage *render.Image16Bit) {
	r.ClearScreen()
	screenImage := r.GetScreenImage()
	buildSaveScreenBackground(screenImage, saveScreenImage)
	r.UpdateVideoBuffer(screenImage)
}

func buildSaveScreenBackground(screenImage *render.Image16Bit, saveScreenImage *render.Image16Bit) {
	screenImage.WriteSubImage(image.Point{0, 0}, saveScreenImage, image.Rect(0, 0, render.BACKGROUND_IMAGE_WIDTH, render.BACKGROUND_IMAGE_HEIGHT))
}
