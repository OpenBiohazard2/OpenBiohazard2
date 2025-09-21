package render

import (
	"github.com/OpenBiohazard2/OpenBiohazard2/geometry"
	"github.com/OpenBiohazard2/OpenBiohazard2/resource"
)

// ScreenImageManager manages the screen image for menu rendering
type ScreenImageManager struct {
	screenImage *resource.Image16Bit
}

// NewScreenImageManager creates a new screen image manager
func NewScreenImageManager() *ScreenImageManager {
	return &ScreenImageManager{
		screenImage: resource.NewImage16Bit(0, 0, geometry.BACKGROUND_IMAGE_WIDTH, geometry.BACKGROUND_IMAGE_HEIGHT),
	}
}

// GetScreenImage returns the screen image for rendering
func (sim *ScreenImageManager) GetScreenImage() *resource.Image16Bit {
	return sim.screenImage
}

// Clear clears the screen image
func (sim *ScreenImageManager) Clear() {
	sim.screenImage.Clear()
}
