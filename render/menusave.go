package render

import (
	"image"
)

func (renderDef *RenderDef) GenerateSaveScreenImage(saveScreenImage *Image16Bit) {
	screenImage.Clear()
	buildSaveScreenBackground(saveScreenImage)
	renderDef.VideoBuffer.UpdateSurface(screenImage)
}

func buildSaveScreenBackground(saveScreenImage *Image16Bit) {
	screenImage.WriteSubImage(image.Point{0, 0}, saveScreenImage, image.Rect(0, 0, 320, 240))
}
