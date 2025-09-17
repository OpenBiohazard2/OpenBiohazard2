package ui

import (
	"testing"
)

func TestNewHealthDisplay(t *testing.T) {
	hd := NewHealthDisplay()

	// Test initial state
	if hd.totalHealthTime != 0 {
		t.Errorf("Expected totalHealthTime to be 0, got %f", hd.totalHealthTime)
	}

	if hd.updateHealthTimeMs != 30 {
		t.Errorf("Expected updateHealthTimeMs to be 30, got %f", hd.updateHealthTimeMs)
	}

	if hd.ecgOffsetX != 0 {
		t.Errorf("Expected ecgOffsetX to be 0, got %d", hd.ecgOffsetX)
	}

	// Test that all health views are initialized
	if len(hd.healthECGViews) != 5 {
		t.Errorf("Expected 5 health ECG views, got %d", len(hd.healthECGViews))
	}
}

func TestUpdateHealthDisplay(t *testing.T) {
	hd := NewHealthDisplay()

	// Test time accumulation
	hd.UpdateHealthDisplay(0.1) // 100ms
	if hd.totalHealthTime != 100 {
		t.Errorf("Expected totalHealthTime to be 100, got %f", hd.totalHealthTime)
	}

	hd.UpdateHealthDisplay(0.05) // 50ms more
	if hd.totalHealthTime != 150 {
		t.Errorf("Expected totalHealthTime to be 150, got %f", hd.totalHealthTime)
	}
}

func TestUpdateECGAnimation(t *testing.T) {
	hd := NewHealthDisplay()

	// Test that animation doesn't update before threshold
	hd.totalHealthTime = 25 // Below 30ms threshold
	hd.UpdateECGAnimation()

	if hd.ecgOffsetX != 0 {
		t.Errorf("Expected ecgOffsetX to remain 0, got %d", hd.ecgOffsetX)
	}

	// Test that animation updates when threshold is reached
	hd.totalHealthTime = 30 // At threshold
	hd.UpdateECGAnimation()

	if hd.ecgOffsetX != 1 {
		t.Errorf("Expected ecgOffsetX to be 1, got %d", hd.ecgOffsetX)
	}

	if hd.totalHealthTime != 0 {
		t.Errorf("Expected totalHealthTime to be reset to 0, got %f", hd.totalHealthTime)
	}
}

func TestUpdateECGAnimation_WrapAround(t *testing.T) {
	hd := NewHealthDisplay()
	hd.ecgOffsetX = 127 // Near end of cycle

	hd.totalHealthTime = 30
	hd.UpdateECGAnimation()

	if hd.ecgOffsetX != 0 {
		t.Errorf("Expected ecgOffsetX to wrap around to 0, got %d", hd.ecgOffsetX)
	}
}

func TestMultipleHealthDisplays(t *testing.T) {
	// This was impossible with global state!
	hd1 := NewHealthDisplay()
	hd2 := NewHealthDisplay()

	// Update them independently
	hd1.UpdateHealthDisplay(0.1)
	hd2.UpdateHealthDisplay(0.2)

	if hd1.totalHealthTime != 100 {
		t.Errorf("Expected hd1.totalHealthTime to be 100, got %f", hd1.totalHealthTime)
	}

	if hd2.totalHealthTime != 200 {
		t.Errorf("Expected hd2.totalHealthTime to be 200, got %f", hd2.totalHealthTime)
	}

	// They should be independent
	if hd1.totalHealthTime == hd2.totalHealthTime {
		t.Error("Health displays should be independent")
	}
}

func TestHealthDisplay_Isolation(t *testing.T) {
	// Test that multiple instances don't interfere
	hd1 := NewHealthDisplay()
	hd2 := NewHealthDisplay()

	// Set up different states
	hd1.totalHealthTime = 50
	hd1.ecgOffsetX = 10

	hd2.totalHealthTime = 25
	hd2.ecgOffsetX = 5

	// Update one
	hd1.UpdateECGAnimation()

	// The other should be unchanged
	if hd2.totalHealthTime != 25 {
		t.Errorf("Expected hd2.totalHealthTime to remain 25, got %f", hd2.totalHealthTime)
	}

	if hd2.ecgOffsetX != 5 {
		t.Errorf("Expected hd2.ecgOffsetX to remain 5, got %d", hd2.ecgOffsetX)
	}
}

// Benchmark tests
func BenchmarkUpdateHealthDisplay(b *testing.B) {
	hd := NewHealthDisplay()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hd.UpdateHealthDisplay(0.016) // ~60 FPS
	}
}

func BenchmarkUpdateECGAnimation(b *testing.B) {
	hd := NewHealthDisplay()
	hd.totalHealthTime = 30 // Ready to animate

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hd.UpdateECGAnimation()
		hd.totalHealthTime = 30 // Reset for next iteration
	}
}
