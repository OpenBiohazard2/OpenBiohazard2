package render

import (
	"image"
)

func (renderDef *RenderDef) GenerateSaveScreenImage(saveScreenImage *Image16Bit) {
	renderDef.ScreenImageManager.Clear()
	screenImage := renderDef.ScreenImageManager.GetScreenImage()
	buildSaveScreenBackground(screenImage, saveScreenImage)
	renderDef.VideoBuffer.UpdateSurface(screenImage)
}

func buildSaveScreenBackground(screenImage *Image16Bit, saveScreenImage *Image16Bit) {
	screenImage.WriteSubImage(image.Point{0, 0}, saveScreenImage, image.Rect(0, 0, BACKGROUND_IMAGE_WIDTH, BACKGROUND_IMAGE_HEIGHT))
}
