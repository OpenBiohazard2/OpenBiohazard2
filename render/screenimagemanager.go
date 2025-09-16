package render

// ScreenImageManager manages the screen image for menu rendering
type ScreenImageManager struct {
	screenImage *Image16Bit
}

// NewScreenImageManager creates a new screen image manager
func NewScreenImageManager() *ScreenImageManager {
	return &ScreenImageManager{
		screenImage: NewImage16Bit(0, 0, BACKGROUND_IMAGE_WIDTH, BACKGROUND_IMAGE_HEIGHT),
	}
}

// GetScreenImage returns the screen image for rendering
func (sim *ScreenImageManager) GetScreenImage() *Image16Bit {
	return sim.screenImage
}

// Clear clears the screen image
func (sim *ScreenImageManager) Clear() {
	sim.screenImage.Clear()
}
