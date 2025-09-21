package geometry

import "image"

const (
	BACKGROUND_IMAGE_WIDTH  = 320
	BACKGROUND_IMAGE_HEIGHT = 240
	
	// Depth constants for 2D rendering
	BACKGROUND_DEPTH = 0.999  // Background elements drawn in the back
)

// BACKGROUND_IMAGE_RECT represents the full background image rectangle (0,0 to width,height)
var BACKGROUND_IMAGE_RECT = image.Rect(0, 0, BACKGROUND_IMAGE_WIDTH, BACKGROUND_IMAGE_HEIGHT)

// ConvertToScreenX normalizes x coordinates between -1 and 1
func ConvertToScreenX(x float32) float32 {
	return 2.0*(x/float32(BACKGROUND_IMAGE_WIDTH)) - 1.0
}

// ConvertToScreenY normalizes y coordinates between -1 and 1
func ConvertToScreenY(y float32) float32 {
	return -1.0 * (2.0*(y/float32(BACKGROUND_IMAGE_HEIGHT)) - 1.0)
}

// ConvertToTextureU normalizes u coordinates between 0 and 1
func ConvertToTextureU(u float32) float32 {
	return u / float32(BACKGROUND_IMAGE_WIDTH)
}

// ConvertToTextureV normalizes v coordinates between 0 and 1
func ConvertToTextureV(v float32) float32 {
	return v / float32(BACKGROUND_IMAGE_HEIGHT)
}
