package resource

import (
	"image"
	"image/color"
	"math"
)

type Image16Bit struct {
	imageData *image.RGBA
}

func NewImage16Bit(x int, y int, width int, height int) *Image16Bit {
	return &Image16Bit{
		imageData: image.NewRGBA(image.Rect(x, y, x+width, y+height)),
	}
}

func ConvertPixelsToImage16Bit(sourcePixels [][]uint16) *Image16Bit {
	width := len(sourcePixels[0])
	height := len(sourcePixels)
	imageData := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := convert16BitColor(sourcePixels[y][x])
			imageData.SetRGBA(x, y, c)
		}
	}

	return &Image16Bit{
		imageData: imageData,
	}
}

func (image16bit *Image16Bit) GetWidth() int {
	return image16bit.imageData.Bounds().Dx()
}

func (image16bit *Image16Bit) GetHeight() int {
	return image16bit.imageData.Bounds().Dy()
}

func (image16bit *Image16Bit) Clear() {
	minPoint := image16bit.imageData.Bounds().Min
	maxPoint := image16bit.imageData.Bounds().Max
	for y := minPoint.Y; y < maxPoint.Y; y++ {
		for x := minPoint.X; x < maxPoint.X; x++ {
			image16bit.imageData.SetRGBA(x, y, color.RGBA{0, 0, 0, 0})
		}
	}
}

func (image16bit *Image16Bit) WriteSubImage(
	destOrigin image.Point,
	sourceImage *Image16Bit,
	sourceBounds image.Rectangle,
) {
	subImageWidth := sourceBounds.Dx()
	subImageHeight := sourceBounds.Dy()
	for offsetY := 0; offsetY < subImageHeight; offsetY++ {
		for offsetX := 0; offsetX < subImageWidth; offsetX++ {
			sourceColor := sourceImage.imageData.RGBAAt(sourceBounds.Min.X+offsetX, sourceBounds.Min.Y+offsetY)
			// Skip writing transparent pixels
			if sourceColor.A > 0 {
				image16bit.imageData.SetRGBA(destOrigin.X+offsetX, destOrigin.Y+offsetY, sourceColor)
			}
		}
	}
}

// Multiply pixels by a brightness factor
// Less than 1.0 will darken it
func (image16bit *Image16Bit) WriteSubImageUniformBrightness(
	destOrigin image.Point,
	sourceImage *Image16Bit,
	sourceBounds image.Rectangle,
	factor float64,
) {
	subImageWidth := sourceBounds.Dx()
	subImageHeight := sourceBounds.Dy()
	for offsetY := 0; offsetY < subImageHeight; offsetY++ {
		for offsetX := 0; offsetX < subImageWidth; offsetX++ {
			sourceColor := sourceImage.imageData.RGBAAt(sourceBounds.Min.X+offsetX, sourceBounds.Min.Y+offsetY)
			newR := int(math.Floor(float64(sourceColor.R) * factor))
			newG := int(math.Floor(float64(sourceColor.G) * factor))
			newB := int(math.Floor(float64(sourceColor.B) * factor))
			modifiedSourceColor := color.RGBA{uint8(newR), uint8(newG), uint8(newB), sourceColor.A}
			// Skip writing transparent pixels
			if modifiedSourceColor.A > 0 {
				image16bit.imageData.SetRGBA(destOrigin.X+offsetX, destOrigin.Y+offsetY, modifiedSourceColor)
			}
		}
	}
}

func (image16bit *Image16Bit) WriteSubImageVariableBrightness(
	destOrigin image.Point,
	sourceImage *Image16Bit,
	sourceBounds image.Rectangle,
	rgbFactor [3]float64,
) {
	subImageWidth := sourceBounds.Dx()
	subImageHeight := sourceBounds.Dy()
	for offsetY := 0; offsetY < subImageHeight; offsetY++ {
		for offsetX := 0; offsetX < subImageWidth; offsetX++ {
			sourceColor := sourceImage.imageData.RGBAAt(sourceBounds.Min.X+offsetX, sourceBounds.Min.Y+offsetY)
			newR := int(math.Floor(float64(sourceColor.R) * rgbFactor[0]))
			newG := int(math.Floor(float64(sourceColor.G) * rgbFactor[1]))
			newB := int(math.Floor(float64(sourceColor.B) * rgbFactor[2]))
			modifiedSourceColor := color.RGBA{uint8(newR), uint8(newG), uint8(newB), sourceColor.A}
			// Skip writing transparent pixels
			if modifiedSourceColor.A > 0 {
				image16bit.imageData.SetRGBA(destOrigin.X+offsetX, destOrigin.Y+offsetY, modifiedSourceColor)
			}
		}
	}
}

func (image16bit *Image16Bit) FillPixels(
	destOrigin image.Point,
	sourceBounds image.Rectangle,
	fillColor color.RGBA,
) {
	subImageWidth := sourceBounds.Dx()
	subImageHeight := sourceBounds.Dy()
	for offsetY := 0; offsetY < subImageHeight; offsetY++ {
		for offsetX := 0; offsetX < subImageWidth; offsetX++ {
			image16bit.imageData.SetRGBA(destOrigin.X+offsetX, destOrigin.Y+offsetY, fillColor)
		}
	}
}

func (image16bit *Image16Bit) ApplyMask(
	destOrigin image.Point,
	sourceImage *Image16Bit,
	sourceBounds image.Rectangle,
) {
	subImageWidth := sourceBounds.Dx()
	subImageHeight := sourceBounds.Dy()
	for offsetY := 0; offsetY < subImageHeight; offsetY++ {
		for offsetX := 0; offsetX < subImageWidth; offsetX++ {
			destOriginalColor := image16bit.imageData.RGBAAt(destOrigin.X+offsetX, destOrigin.Y+offsetY)
			sourceColor := sourceImage.imageData.RGBAAt(sourceBounds.Min.X+offsetX, sourceBounds.Min.Y+offsetY)

			// Set destination pixel to transparent if the mask pixel is transparent (pixel value is 0)
			destModifiedColor := color.RGBA{destOriginalColor.R, destOriginalColor.G, destOriginalColor.B, sourceColor.A}
			image16bit.imageData.SetRGBA(destOrigin.X+offsetX, destOrigin.Y+offsetY, destModifiedColor)
		}
	}
}

// OpenGL renderer uses a 1d array to render pixels to screen
func (image16bit *Image16Bit) GetPixelsForRendering() []uint16 {
	bounds := image16bit.imageData.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	pixelData1D := make([]uint16, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := (y * width) + x
			pixelData1D[index] = convertRGBAToA1B5G5R5(image16bit.imageData.RGBAAt(x, y))
		}
	}
	return pixelData1D
}

// Original pixel is 16 bits in the A1B5G5R5 format
// Extract RGBA values
func convert16BitColor(pixel uint16) color.RGBA {
	// r, g, b only has values in range 0 - 32
	// scale to 0 - 256
	pixelR := (pixel % 32) * 8
	pixelG := ((pixel >> 5) % 32) * 8
	pixelB := ((pixel >> 10) % 32) * 8

	if pixel > 0 {
		return color.RGBA{uint8(pixelR), uint8(pixelG), uint8(pixelB), uint8(255)}
	} else {
		return color.RGBA{0, 0, 0, 0}
	}
}

func convertRGBAToA1B5G5R5(pixelColor color.RGBA) uint16 {
	// Convert color to A1R5G5B5 format for OpenGL rendering
	newR := uint16(math.Round(float64(pixelColor.R) / 8.0))
	newG := uint16(math.Round(float64(pixelColor.G) / 8.0))
	newB := uint16(math.Round(float64(pixelColor.B) / 8.0))
	if newR >= 32 {
		newR = 31
	}
	if newG >= 32 {
		newG = 31
	}
	if newB >= 32 {
		newB = 31
	}
	newA := uint16(math.Floor(float64(pixelColor.A) / 255.0))
	return (newA << 15) | (newB << 10) | (newG << 5) | newR
}
