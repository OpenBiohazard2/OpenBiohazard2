package ui_render

import (
	"image"

	"github.com/OpenBiohazard2/OpenBiohazard2/geometry"
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
	screenImage.WriteSubImage(image.Point{0, 0}, saveScreenImage, geometry.BACKGROUND_IMAGE_RECT)
}
