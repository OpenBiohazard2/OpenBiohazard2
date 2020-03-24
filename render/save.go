package render

import (
	"github.com/samuelyuan/openbiohazard2/fileio"
)

func (renderDef *RenderDef) GenerateSaveScreenImage(saveScreenImageOutput *fileio.ADTOutput) {
	renderDef.VideoBuffer.ClearSurface()
	newImageColors := renderDef.VideoBuffer.ImagePixels
	buildSaveScreenBackground(saveScreenImageOutput, newImageColors)
	renderDef.VideoBuffer.UpdateSurface(newImageColors)
}

func buildSaveScreenBackground(saveScreenImageOutput *fileio.ADTOutput, newImageColors []uint16) {
	copyPixelsTransparent(saveScreenImageOutput.PixelData, 0, 0, 320, 240, newImageColors, 0, 0)
}
