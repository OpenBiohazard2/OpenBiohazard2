package ui_render

import (
	"image"

	"github.com/OpenBiohazard2/OpenBiohazard2/render"
	"github.com/OpenBiohazard2/OpenBiohazard2/resource"
)

// GenerateSaveScreenImage renders the save screen
func (r *UIRenderer) GenerateSaveScreenImage(saveScreenImage *resource.Image16Bit) {
	r.ClearScreen()
	screenImage := r.GetScreenImage()
	buildSaveScreenBackground(screenImage, saveScreenImage)
	r.UpdateVideoBuffer(screenImage)
}

func buildSaveScreenBackground(screenImage *resource.Image16Bit, saveScreenImage *resource.Image16Bit) {
	screenImage.WriteSubImage(image.Point{0, 0}, saveScreenImage, image.Rect(0, 0, render.BACKGROUND_IMAGE_WIDTH, render.BACKGROUND_IMAGE_HEIGHT))
}
