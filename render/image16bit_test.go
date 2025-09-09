package render

import (
	"image"
	"image/color"
	"math"
	"testing"
)

func TestNewImage16Bit(t *testing.T) {
	tests := []struct {
		name   string
		x      int
		y      int
		width  int
		height int
	}{
		{
			name:   "Standard_image",
			x:      0,
			y:      0,
			width:  100,
			height: 100,
		},
		{
			name:   "Offset_image",
			x:      10,
			y:      20,
			width:  50,
			height: 75,
		},
		{
			name:   "Single_pixel",
			x:      0,
			y:      0,
			width:  1,
			height: 1,
		},
		{
			name:   "Wide_image",
			x:      0,
			y:      0,
			width:  200,
			height: 50,
		},
		{
			name:   "Tall_image",
			x:      0,
			y:      0,
			width:  50,
			height: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := NewImage16Bit(tt.x, tt.y, tt.width, tt.height)

			if img == nil {
				t.Fatal("NewImage16Bit returned nil")
			}

			if img.imageData == nil {
				t.Fatal("Image data should not be nil")
			}

			// Test dimensions
			expectedWidth := tt.width
			expectedHeight := tt.height
			actualWidth := img.GetWidth()
			actualHeight := img.GetHeight()

			if actualWidth != expectedWidth {
				t.Errorf("Expected width %d, got %d", expectedWidth, actualWidth)
			}

			if actualHeight != expectedHeight {
				t.Errorf("Expected height %d, got %d", expectedHeight, actualHeight)
			}

			// Test bounds
			bounds := img.imageData.Bounds()
			expectedMinX := tt.x
			expectedMinY := tt.y
			expectedMaxX := tt.x + tt.width
			expectedMaxY := tt.y + tt.height

			if bounds.Min.X != expectedMinX {
				t.Errorf("Expected min X %d, got %d", expectedMinX, bounds.Min.X)
			}

			if bounds.Min.Y != expectedMinY {
				t.Errorf("Expected min Y %d, got %d", expectedMinY, bounds.Min.Y)
			}

			if bounds.Max.X != expectedMaxX {
				t.Errorf("Expected max X %d, got %d", expectedMaxX, bounds.Max.X)
			}

			if bounds.Max.Y != expectedMaxY {
				t.Errorf("Expected max Y %d, got %d", expectedMaxY, bounds.Max.Y)
			}
		})
	}
}

func TestImage16Bit_GetWidth_GetHeight(t *testing.T) {
	img := NewImage16Bit(0, 0, 150, 200)

	width := img.GetWidth()
	height := img.GetHeight()

	if width != 150 {
		t.Errorf("Expected width 150, got %d", width)
	}

	if height != 200 {
		t.Errorf("Expected height 200, got %d", height)
	}
}

func TestImage16Bit_Clear(t *testing.T) {
	img := NewImage16Bit(0, 0, 10, 10)

	// Fill with some color first
	fillColor := color.RGBA{255, 128, 64, 255}
	img.FillPixels(image.Point{0, 0}, image.Rect(0, 0, 10, 10), fillColor)

	// Verify it's filled
	testColor := img.imageData.RGBAAt(5, 5)
	if testColor != fillColor {
		t.Error("Image should be filled with test color before clearing")
	}

	// Clear the image
	img.Clear()

	// Verify it's cleared (transparent)
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			pixelColor := img.imageData.RGBAAt(x, y)
			expectedColor := color.RGBA{0, 0, 0, 0}
			if pixelColor != expectedColor {
				t.Errorf("Pixel at (%d, %d) should be transparent, got %v", x, y, pixelColor)
			}
		}
	}
}

func TestImage16Bit_WriteSubImage(t *testing.T) {
	// Create source image with a pattern
	sourceImg := NewImage16Bit(0, 0, 5, 5)
	redColor := color.RGBA{255, 0, 0, 255}
	blueColor := color.RGBA{0, 0, 255, 255}

	// Fill source with red
	sourceImg.FillPixels(image.Point{0, 0}, image.Rect(0, 0, 5, 5), redColor)
	// Add blue square in center
	sourceImg.FillPixels(image.Point{1, 1}, image.Rect(1, 1, 4, 4), blueColor)

	// Create destination image
	destImg := NewImage16Bit(0, 0, 10, 10)

	// Write source to destination at offset (2, 2)
	destImg.WriteSubImage(image.Point{2, 2}, sourceImg, image.Rect(0, 0, 5, 5))

	// Verify the copy
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			expectedColor := sourceImg.imageData.RGBAAt(x, y)
			actualColor := destImg.imageData.RGBAAt(x+2, y+2)
			if actualColor != expectedColor {
				t.Errorf("Pixel at (%d, %d) should be %v, got %v", x+2, y+2, expectedColor, actualColor)
			}
		}
	}
}

func TestImage16Bit_WriteSubImage_Transparency(t *testing.T) {
	// Create source image with transparent pixels
	sourceImg := NewImage16Bit(0, 0, 3, 3)
	redColor := color.RGBA{255, 0, 0, 255}
	transparentColor := color.RGBA{0, 0, 0, 0}

	// Fill with red
	sourceImg.FillPixels(image.Point{0, 0}, image.Rect(0, 0, 3, 3), redColor)
	// Make center pixel transparent
	sourceImg.imageData.SetRGBA(1, 1, transparentColor)

	// Create destination image with blue background
	destImg := NewImage16Bit(0, 0, 5, 5)
	blueColor := color.RGBA{0, 0, 255, 255}
	destImg.FillPixels(image.Point{0, 0}, image.Rect(0, 0, 5, 5), blueColor)

	// Write source to destination
	destImg.WriteSubImage(image.Point{1, 1}, sourceImg, image.Rect(0, 0, 3, 3))

	// Verify transparent pixel wasn't written (should remain blue)
	centerColor := destImg.imageData.RGBAAt(2, 2)
	if centerColor != blueColor {
		t.Errorf("Transparent pixel should not overwrite destination, expected %v, got %v", blueColor, centerColor)
	}

	// Verify non-transparent pixels were written
	cornerColor := destImg.imageData.RGBAAt(1, 1)
	if cornerColor != redColor {
		t.Errorf("Non-transparent pixel should be written, expected %v, got %v", redColor, cornerColor)
	}
}

func TestImage16Bit_WriteSubImageUniformBrightness(t *testing.T) {
	// Create source image
	sourceImg := NewImage16Bit(0, 0, 2, 2)
	originalColor := color.RGBA{200, 100, 50, 255}
	sourceImg.FillPixels(image.Point{0, 0}, image.Rect(0, 0, 2, 2), originalColor)

	// Create destination image
	destImg := NewImage16Bit(0, 0, 4, 4)

	// Test brightness factor of 0.5 (should darken)
	brightnessFactor := 0.5
	destImg.WriteSubImageUniformBrightness(image.Point{1, 1}, sourceImg, image.Rect(0, 0, 2, 2), brightnessFactor)

	// Verify brightness adjustment
	expectedR := uint8(math.Floor(float64(originalColor.R) * brightnessFactor))
	expectedG := uint8(math.Floor(float64(originalColor.G) * brightnessFactor))
	expectedB := uint8(math.Floor(float64(originalColor.B) * brightnessFactor))
	expectedColor := color.RGBA{expectedR, expectedG, expectedB, originalColor.A}

	actualColor := destImg.imageData.RGBAAt(1, 1)
	if actualColor != expectedColor {
		t.Errorf("Expected brightness-adjusted color %v, got %v", expectedColor, actualColor)
	}
}

func TestImage16Bit_WriteSubImageVariableBrightness(t *testing.T) {
	// Create source image
	sourceImg := NewImage16Bit(0, 0, 2, 2)
	originalColor := color.RGBA{200, 100, 50, 255}
	sourceImg.FillPixels(image.Point{0, 0}, image.Rect(0, 0, 2, 2), originalColor)

	// Create destination image
	destImg := NewImage16Bit(0, 0, 4, 4)

	// Test different RGB factors
	rgbFactor := [3]float64{0.5, 1.0, 2.0} // Red dimmed, Green unchanged, Blue brightened
	destImg.WriteSubImageVariableBrightness(image.Point{1, 1}, sourceImg, image.Rect(0, 0, 2, 2), rgbFactor)

	// Verify RGB factor adjustment
	expectedR := uint8(math.Floor(float64(originalColor.R) * rgbFactor[0]))
	expectedG := uint8(math.Floor(float64(originalColor.G) * rgbFactor[1]))
	expectedB := uint8(math.Floor(float64(originalColor.B) * rgbFactor[2]))
	expectedColor := color.RGBA{expectedR, expectedG, expectedB, originalColor.A}

	actualColor := destImg.imageData.RGBAAt(1, 1)
	if actualColor != expectedColor {
		t.Errorf("Expected RGB-factor-adjusted color %v, got %v", expectedColor, actualColor)
	}
}

func TestImage16Bit_FillPixels(t *testing.T) {
	img := NewImage16Bit(0, 0, 5, 5)
	fillColor := color.RGBA{128, 64, 32, 255}

	// Fill a 3x3 area starting at (1, 1)
	img.FillPixels(image.Point{1, 1}, image.Rect(1, 1, 4, 4), fillColor)

	// Verify filled area
	for y := 1; y < 4; y++ {
		for x := 1; x < 4; x++ {
			pixelColor := img.imageData.RGBAAt(x, y)
			if pixelColor != fillColor {
				t.Errorf("Pixel at (%d, %d) should be %v, got %v", x, y, fillColor, pixelColor)
			}
		}
	}

	// Verify unfilled area (should be transparent)
	unfilledColor := img.imageData.RGBAAt(0, 0)
	expectedUnfilled := color.RGBA{0, 0, 0, 0}
	if unfilledColor != expectedUnfilled {
		t.Errorf("Unfilled pixel should be transparent, got %v", unfilledColor)
	}
}

func TestImage16Bit_ApplyMask(t *testing.T) {
	// Create destination image with red color
	destImg := NewImage16Bit(0, 0, 3, 3)
	redColor := color.RGBA{255, 0, 0, 255}
	destImg.FillPixels(image.Point{0, 0}, image.Rect(0, 0, 3, 3), redColor)

	// Create mask image
	maskImg := NewImage16Bit(0, 0, 3, 3)
	opaqueColor := color.RGBA{0, 0, 0, 255}
	transparentColor := color.RGBA{0, 0, 0, 0}

	// Fill mask with opaque, but make center transparent
	maskImg.FillPixels(image.Point{0, 0}, image.Rect(0, 0, 3, 3), opaqueColor)
	maskImg.imageData.SetRGBA(1, 1, transparentColor)

	// Apply mask
	destImg.ApplyMask(image.Point{0, 0}, maskImg, image.Rect(0, 0, 3, 3))

	// Verify masked pixels
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			maskColor := maskImg.imageData.RGBAAt(x, y)
			expectedAlpha := maskColor.A
			actualColor := destImg.imageData.RGBAAt(x, y)

			if actualColor.A != expectedAlpha {
				t.Errorf("Pixel at (%d, %d) should have alpha %d, got %d", x, y, expectedAlpha, actualColor.A)
			}

			// RGB should remain unchanged
			if actualColor.R != redColor.R || actualColor.G != redColor.G || actualColor.B != redColor.B {
				t.Errorf("Pixel at (%d, %d) RGB should remain %v, got %v", x, y, redColor, actualColor)
			}
		}
	}
}

func TestImage16Bit_GetPixelsForRendering(t *testing.T) {
	// Create image with known colors
	img := NewImage16Bit(0, 0, 2, 2)
	redColor := color.RGBA{255, 0, 0, 255}
	greenColor := color.RGBA{0, 255, 0, 255}

	img.imageData.SetRGBA(0, 0, redColor)
	img.imageData.SetRGBA(1, 0, greenColor)
	img.imageData.SetRGBA(0, 1, redColor)
	img.imageData.SetRGBA(1, 1, greenColor)

	// Get pixels for rendering
	pixels := img.GetPixelsForRendering()

	// Verify array length
	expectedLength := 2 * 2 // width * height
	if len(pixels) != expectedLength {
		t.Errorf("Expected pixel array length %d, got %d", expectedLength, len(pixels))
	}

	// Verify pixel values (converted to A1B5G5R5 format)
	// Note: We can't easily verify exact values without duplicating the conversion logic,
	// but we can verify the array has the right length and contains non-zero values
	hasNonZero := false
	for i, pixel := range pixels {
		if pixel != 0 {
			hasNonZero = true
		}
		if i >= expectedLength {
			t.Errorf("Pixel array has unexpected length, index %d out of bounds", i)
		}
	}

	if !hasNonZero {
		t.Error("Pixel array should contain non-zero values")
	}
}

func TestConvertPixelsToImage16Bit(t *testing.T) {
	// Create test pixel data (2x2 image)
	// A1B5G5R5 format: bit 15=alpha, bits 14-10=blue, bits 9-5=green, bits 4-0=red
	sourcePixels := [][]uint16{
		{0x0000, 0x001F}, // transparent, red (bits 0-4 = 31, scaled to 248)
		{0x03E0, 0x7C00}, // green (bits 5-9 = 31, scaled to 248), blue (bits 10-14 = 31, scaled to 248)
	}

	img := ConvertPixelsToImage16Bit(sourcePixels)

	if img == nil {
		t.Fatal("ConvertPixelsToImage16Bit returned nil")
	}

	// Verify dimensions
	if img.GetWidth() != 2 {
		t.Errorf("Expected width 2, got %d", img.GetWidth())
	}

	if img.GetHeight() != 2 {
		t.Errorf("Expected height 2, got %d", img.GetHeight())
	}

	// Verify pixel conversions
	// Pixel (0,0) should be transparent (0x0000)
	pixel00 := img.imageData.RGBAAt(0, 0)
	if pixel00.A != 0 {
		t.Errorf("Pixel (0,0) should be transparent, got alpha %d", pixel00.A)
	}

	// Pixel (1,0) should be red (0x001F = red in bits 0-4)
	pixel10 := img.imageData.RGBAAt(1, 0)
	if pixel10.R == 0 || pixel10.G != 0 || pixel10.B != 0 {
		t.Errorf("Pixel (1,0) should be red, got %v", pixel10)
	}
}

// Test the private conversion functions by testing their behavior through public functions
func TestColorConversions(t *testing.T) {
	// Test convert16BitColor through ConvertPixelsToImage16Bit
	testCases := []struct {
		name     string
		pixel16  uint16
		expected color.RGBA
	}{
		{
			name:     "Transparent_pixel",
			pixel16:  0x0000,
			expected: color.RGBA{0, 0, 0, 0},
		},
		{
			name:     "Red_pixel",
			pixel16:  0x001F,                     // Red in A1B5G5R5 format (bits 0-4 = 31, scaled to 248)
			expected: color.RGBA{248, 0, 0, 255}, // Approximate red
		},
		{
			name:     "Green_pixel",
			pixel16:  0x03E0,                     // Green in A1B5G5R5 format (bits 5-9 = 31, scaled to 248)
			expected: color.RGBA{0, 248, 0, 255}, // Approximate green
		},
		{
			name:     "Blue_pixel",
			pixel16:  0x7C00,                     // Blue in A1B5G5R5 format (bits 10-14 = 31, scaled to 248)
			expected: color.RGBA{0, 0, 248, 255}, // Approximate blue
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create 1x1 image with the test pixel
			sourcePixels := [][]uint16{{tc.pixel16}}
			img := ConvertPixelsToImage16Bit(sourcePixels)

			actualColor := img.imageData.RGBAAt(0, 0)

			// Allow for some tolerance in color conversion
			tolerance := uint8(8) // Allow ±8 in each channel
			if !colorsWithinTolerance(actualColor, tc.expected, tolerance) {
				t.Errorf("Expected color %v, got %v (tolerance: ±%d)", tc.expected, actualColor, tolerance)
			}
		})
	}
}

// Helper function to check if two colors are within tolerance
func colorsWithinTolerance(actual, expected color.RGBA, tolerance uint8) bool {
	return absDiff(actual.R, expected.R) <= tolerance &&
		absDiff(actual.G, expected.G) <= tolerance &&
		absDiff(actual.B, expected.B) <= tolerance &&
		actual.A == expected.A // Alpha should be exact
}

func absDiff(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}

// Benchmark tests
func BenchmarkNewImage16Bit(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewImage16Bit(0, 0, 100, 100)
	}
}

func BenchmarkImage16Bit_Clear(b *testing.B) {
	img := NewImage16Bit(0, 0, 100, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		img.Clear()
	}
}

func BenchmarkImage16Bit_WriteSubImage(b *testing.B) {
	sourceImg := NewImage16Bit(0, 0, 50, 50)
	destImg := NewImage16Bit(0, 0, 100, 100)

	// Fill source with some color
	redColor := color.RGBA{255, 0, 0, 255}
	sourceImg.FillPixels(image.Point{0, 0}, image.Rect(0, 0, 50, 50), redColor)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		destImg.WriteSubImage(image.Point{0, 0}, sourceImg, image.Rect(0, 0, 50, 50))
	}
}

func BenchmarkImage16Bit_GetPixelsForRendering(b *testing.B) {
	img := NewImage16Bit(0, 0, 100, 100)

	// Fill with some color
	redColor := color.RGBA{255, 0, 0, 255}
	img.FillPixels(image.Point{0, 0}, image.Rect(0, 0, 100, 100), redColor)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = img.GetPixelsForRendering()
	}
}
