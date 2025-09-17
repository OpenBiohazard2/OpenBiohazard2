package ui

import (
	"testing"
)

func TestNewInventoryManager(t *testing.T) {
	manager := NewInventoryManager()

	if manager == nil {
		t.Fatal("NewInventoryManager() returned nil")
	}

	if manager.totalInventoryTime != 0 {
		t.Errorf("Expected totalInventoryTime to be 0, got %f", manager.totalInventoryTime)
	}

	if manager.updateInventoryCursorTime != 30 {
		t.Errorf("Expected updateInventoryCursorTime to be 30, got %f", manager.updateInventoryCursorTime)
	}

	items := manager.GetPlayerInventoryItems()
	if len(items) != 11 {
		t.Errorf("Expected 11 inventory items, got %d", len(items))
	}

	// Check specific items
	if items[0].Id != 2 || items[0].Num != 18 || items[0].Size != 1 {
		t.Errorf("Expected first item to be hand gun (Id: 2, Num: 18, Size: 1), got (Id: %d, Num: %d, Size: %d)",
			items[0].Id, items[0].Num, items[0].Size)
	}
}

func TestInventoryManager_UpdateInventoryTime(t *testing.T) {
	manager := NewInventoryManager()

	// Test initial state
	if manager.totalInventoryTime != 0 {
		t.Errorf("Expected initial totalInventoryTime to be 0, got %f", manager.totalInventoryTime)
	}

	// Test time update
	manager.UpdateInventoryTime(0.1) // 100ms
	expected := 100.0
	if manager.totalInventoryTime != expected {
		t.Errorf("Expected totalInventoryTime to be %f after 0.1s, got %f", expected, manager.totalInventoryTime)
	}

	// Test multiple updates
	manager.UpdateInventoryTime(0.05) // 50ms
	expected = 150.0
	if manager.totalInventoryTime != expected {
		t.Errorf("Expected totalInventoryTime to be %f after additional 0.05s, got %f", expected, manager.totalInventoryTime)
	}
}

func TestInventoryManager_ShouldUpdateCursor(t *testing.T) {
	manager := NewInventoryManager()

	// Test initial state - should not update cursor
	if manager.ShouldUpdateCursor() {
		t.Error("Expected ShouldUpdateCursor to return false initially")
	}

	// Test after enough time has passed
	manager.UpdateInventoryTime(0.05) // 50ms - should be enough
	if !manager.ShouldUpdateCursor() {
		t.Error("Expected ShouldUpdateCursor to return true after 50ms")
	}
}

func TestInventoryManager_ResetInventoryTime(t *testing.T) {
	manager := NewInventoryManager()

	// Add some time
	manager.UpdateInventoryTime(0.1)
	if manager.totalInventoryTime == 0 {
		t.Error("Expected totalInventoryTime to be non-zero after update")
	}

	// Reset and verify
	manager.ResetInventoryTime()
	if manager.totalInventoryTime != 0 {
		t.Errorf("Expected totalInventoryTime to be 0 after reset, got %f", manager.totalInventoryTime)
	}
}

func TestInventoryManager_GetPlayerInventoryItems(t *testing.T) {
	manager := NewInventoryManager()

	items := manager.GetPlayerInventoryItems()

	// Test that we get the expected number of items
	if len(items) != 11 {
		t.Errorf("Expected 11 items, got %d", len(items))
	}

	// Test that modifications to returned slice don't affect internal state
	items[0].Id = 999
	originalItems := manager.GetPlayerInventoryItems()
	if originalItems[0].Id == 999 {
		t.Error("Modifying returned slice affected internal state")
	}
}

func TestInventoryManager_Isolation(t *testing.T) {
	// Test that multiple instances don't interfere
	manager1 := NewInventoryManager()
	manager2 := NewInventoryManager()

	// Update them independently
	manager1.UpdateInventoryTime(0.1)
	manager2.UpdateInventoryTime(0.2)

	// Verify they have different states
	if manager1.totalInventoryTime == manager2.totalInventoryTime {
		t.Error("Different managers should have different time values")
	}

	// Verify they have different inventory items (if modified)
	items1 := manager1.GetPlayerInventoryItems()
	_ = manager2.GetPlayerInventoryItems() // items2

	// Modify one
	items1[0].Id = 999

	// Verify the other is unaffected
	items2Again := manager2.GetPlayerInventoryItems()
	if items2Again[0].Id == 999 {
		t.Error("Modifying one manager's items affected another manager")
	}
}

// Benchmark tests
func BenchmarkUpdateInventoryTime(b *testing.B) {
	manager := NewInventoryManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.UpdateInventoryTime(0.016) // ~60 FPS
	}
}

func BenchmarkGetPlayerInventoryItems(b *testing.B) {
	manager := NewInventoryManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.GetPlayerInventoryItems()
	}
}
