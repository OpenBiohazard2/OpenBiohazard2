package render

import (
	"math"
)

const (
	IMAGE_SURFACE_WIDTH  = 320
	IMAGE_SURFACE_HEIGHT = 240
)

func NewSurface2D() []uint16 {
	return make([]uint16, IMAGE_SURFACE_WIDTH*IMAGE_SURFACE_HEIGHT)
}

func copyPixels(sourcePixels [][]uint16, startX int, startY int, width int, height int,
	destPixels []uint16, destX int, destY int) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			destPixels[((destY+y)*IMAGE_SURFACE_WIDTH)+(destX+x)] = sourcePixels[startY+y][startX+x]
		}
	}
}

func copyPixelsTransparent(sourcePixels [][]uint16, startX int, startY int, width int, height int,
	destPixels []uint16, destX int, destY int) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := sourcePixels[startY+y][startX+x]
			newPixel := uint16(0)
			if pixel > 0 {
				newPixel = (1 << 15) | pixel
			}

			// Overwrite pixel if it's not transparent
			if newPixel > 0 {
				destPixels[((destY+y)*IMAGE_SURFACE_WIDTH)+(destX+x)] = newPixel
			}
		}
	}
}

// Multiply pixels by a brightness factor
// Less than 1.0 will darken it
func copyPixelsBrightness(sourcePixels [][]uint16, startX int, startY int, width int, height int,
	destPixels []uint16, destX int, destY int, factor float64) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := sourcePixels[startY+y][startX+x]
			newPixel := int(0)
			if pixel > 0 {
				pixelR := (pixel % 32)
				pixelG := ((pixel >> 5) % 32)
				pixelB := ((pixel >> 10) % 32)
				newR := int(math.Floor(float64(pixelR) * factor))
				newG := int(math.Floor(float64(pixelG) * factor))
				newB := int(math.Floor(float64(pixelB) * factor))
				newPixel = (1 << 15) | (newB << 10) | (newG << 5) | newR
			}

			// Overwrite pixel if it's not transparent
			if newPixel > 0 {
				destPixels[((destY+y)*IMAGE_SURFACE_WIDTH)+(destX+x)] = uint16(newPixel)
			}
		}
	}
}

func fillPixels(newImageColors []uint16, destX int, destY int, width int, height int,
	r int, g int, b int) {
	// Convert color to A1R5G5B5 format
	newR := uint16(math.Round(float64(r) / 8.0))
	newG := uint16(math.Round(float64(g) / 8.0))
	newB := uint16(math.Round(float64(b) / 8.0))
	color := (1 << 15) | (newB << 10) | (newG << 5) | newR

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			newImageColors[((destY+y)*IMAGE_SURFACE_WIDTH)+(destX+x)] = color
		}
	}
}

func buildSurface2DVertexBuffer() []float32 {
	z := float32(0.999)
	return []float32{
		// (-1, 1, z)
		-1.0, 1.0, z, 0.0, 0.0,
		// (-1, -1, z)
		-1.0, -1.0, z, 0.0, 1.0,
		// (1, -1, z)
		1.0, -1.0, z, 1.0, 1.0,

		// (1, -1, z)
		1.0, -1.0, z, 1.0, 1.0,
		// (1, 1, z)
		1.0, 1.0, z, 1.0, 0.0,
		// (-1, 1, z)
		-1.0, 1.0, z, 0.0, 0.0,
	}
}
