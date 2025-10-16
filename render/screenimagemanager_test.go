package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/OpenBiohazard2/OpenBiohazard2/geometry"
	"github.com/OpenBiohazard2/OpenBiohazard2/resource"
)

func TestNewScreenImageManager(t *testing.T) {
	manager := NewScreenImageManager()

	if manager == nil {
		t.Fatal("NewScreenImageManager() returned nil")
	}

	if manager.screenImage == nil {
		t.Fatal("Expected screenImage to be initialized")
	}

	// Test that the screen image has correct dimensions
	if manager.screenImage.GetWidth() != geometry.BACKGROUND_IMAGE_WIDTH {
		t.Errorf("Expected screen image width to be %d, got %d", geometry.BACKGROUND_IMAGE_WIDTH, manager.screenImage.GetWidth())
	}

	if manager.screenImage.GetHeight() != geometry.BACKGROUND_IMAGE_HEIGHT {
		t.Errorf("Expected screen image height to be %d, got %d", geometry.BACKGROUND_IMAGE_HEIGHT, manager.screenImage.GetHeight())
	}
}

func TestScreenImageManager_GetScreenImage(t *testing.T) {
	manager := NewScreenImageManager()
	screenImage := manager.GetScreenImage()

	if screenImage == nil {
		t.Fatal("GetScreenImage() returned nil")
	}

	// Test that we get the same instance
	if screenImage != manager.screenImage {
		t.Error("GetScreenImage() should return the same instance")
	}
}

func TestScreenImageManager_Clear(t *testing.T) {
	manager := NewScreenImageManager()
	screenImage := manager.GetScreenImage()

	// Fill the screen with some data first
	originalPixels := screenImage.GetPixelsForRendering()
	if len(originalPixels) == 0 {
		t.Fatal("Expected screen image to have pixels")
	}

	// Clear the screen
	manager.Clear()

	// Verify the screen was cleared
	clearedPixels := screenImage.GetPixelsForRendering()
	if len(clearedPixels) != len(originalPixels) {
		t.Error("Expected pixel count to remain the same after clear")
	}

	// Check that all pixels are now transparent (0)
	for i, pixel := range clearedPixels {
		if pixel != 0 {
			t.Errorf("Expected pixel %d to be 0 (transparent) after clear, got %d", i, pixel)
		}
	}
}

func TestScreenImageManager_Isolation(t *testing.T) {
	// Test that multiple managers are isolated
	manager1 := NewScreenImageManager()
	manager2 := NewScreenImageManager()

	screenImage1 := manager1.GetScreenImage()
	screenImage2 := manager2.GetScreenImage()

	// They should be different instances
	if screenImage1 == screenImage2 {
		t.Error("Expected different managers to have different screen images")
	}

	// Test that operations on one don't affect the other
	manager1.Clear()

	// manager2's screen image should still have its original state
	screenImage2Pixels := screenImage2.GetPixelsForRendering()
	if len(screenImage2Pixels) == 0 {
		t.Error("Expected manager2's screen image to still have pixels")
	}
}

func TestScreenImageManager_WriteOperations(t *testing.T) {
	manager := NewScreenImageManager()
	screenImage := manager.GetScreenImage()

	// Create a test image
	testImage := resource.NewImage16Bit(0, 0, 10, 10)
	redColor := color.RGBA{255, 0, 0, 255}
	testImage.FillPixels(image.Point{0, 0}, image.Rect(0, 0, 10, 10), redColor)

	// Write to screen image
	screenImage.WriteSubImage(image.Point{5, 5}, testImage, image.Rect(0, 0, 10, 10))

	// Verify the write operation
	resultPixels := screenImage.GetPixelsForRendering()
	if len(resultPixels) == 0 {
		t.Fatal("Expected screen image to have pixels after write operation")
	}

	// Check that some pixels were written (not all transparent)
	hasNonZeroPixels := false
	for _, pixel := range resultPixels {
		if pixel != 0 {
			hasNonZeroPixels = true
			break
		}
	}

	if !hasNonZeroPixels {
		t.Error("Expected screen image to have non-zero pixels after write operation")
	}
}

func TestScreenImageManager_ClearAfterWrite(t *testing.T) {
	manager := NewScreenImageManager()
	screenImage := manager.GetScreenImage()

	// Write some data
	testImage := resource.NewImage16Bit(0, 0, 5, 5)
	whiteColor := color.RGBA{255, 255, 255, 255}
	testImage.FillPixels(image.Point{0, 0}, image.Rect(0, 0, 5, 5), whiteColor)
	screenImage.WriteSubImage(image.Point{0, 0}, testImage, image.Rect(0, 0, 5, 5))

	// Verify data was written
	beforeClear := screenImage.GetPixelsForRendering()
	hasData := false
	for _, pixel := range beforeClear {
		if pixel != 0 {
			hasData = true
			break
		}
	}
	if !hasData {
		t.Fatal("Expected screen image to have data before clear")
	}

	// Clear and verify
	manager.Clear()
	afterClear := screenImage.GetPixelsForRendering()
	for i, pixel := range afterClear {
		if pixel != 0 {
			t.Errorf("Expected pixel %d to be 0 after clear, got %d", i, pixel)
		}
	}
}

func BenchmarkScreenImageManager_Clear(b *testing.B) {
	manager := NewScreenImageManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Clear()
	}
}

func BenchmarkScreenImageManager_GetScreenImage(b *testing.B) {
	manager := NewScreenImageManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.GetScreenImage()
	}
}
