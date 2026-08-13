package state

import (
	"testing"
)

func TestNewGameStateManager(t *testing.T) {
	manager := NewGameStateManager()

	if manager == nil {
		t.Fatal("NewGameStateManager returned nil")
	}

	if manager.GameState != GAME_STATE_MAIN_MENU {
		t.Errorf("Expected initial state to be GAME_STATE_MAIN_MENU (%d), got %d", GAME_STATE_MAIN_MENU, manager.GameState)
	}

	if manager.ImageResourcesLoaded != false {
		t.Errorf("Expected ImageResourcesLoaded to be false, got %v", manager.ImageResourcesLoaded)
	}
}

func TestUpdateGameState(t *testing.T) {
	manager := NewGameStateManager()
	manager.ImageResourcesLoaded = true // Set to true to test it gets reset

	// Test updating to main game state
	manager.UpdateGameState(GAME_STATE_MAIN_GAME)

	if manager.GameState != GAME_STATE_MAIN_GAME {
		t.Errorf("Expected GameState to be GAME_STATE_MAIN_GAME (%d), got %d", GAME_STATE_MAIN_GAME, manager.GameState)
	}

	if manager.ImageResourcesLoaded != false {
		t.Errorf("Expected ImageResourcesLoaded to be reset to false, got %v", manager.ImageResourcesLoaded)
	}

	// Test updating to inventory state
	manager.UpdateGameState(GAME_STATE_INVENTORY)

	if manager.GameState != GAME_STATE_INVENTORY {
		t.Errorf("Expected GameState to be GAME_STATE_INVENTORY (%d), got %d", GAME_STATE_INVENTORY, manager.GameState)
	}

	if manager.ImageResourcesLoaded != false {
		t.Errorf("Expected ImageResourcesLoaded to be reset to false, got %v", manager.ImageResourcesLoaded)
	}
}

func TestGameStateConstants(t *testing.T) {
	// Test that all state constants have unique values
	states := []int{
		GAME_STATE_MAIN_MENU,
		GAME_STATE_MAIN_GAME,
		GAME_STATE_INVENTORY,
		GAME_STATE_LOAD_SAVE,
		GAME_STATE_SPECIAL_MENU,
	}

	// Check for duplicates
	seen := make(map[int]bool)
	for _, state := range states {
		if seen[state] {
			t.Errorf("Duplicate state value found: %d", state)
		}
		seen[state] = true
	}

	// Test specific expected values
	expectedStates := map[string]int{
		"GAME_STATE_MAIN_MENU":    0,
		"GAME_STATE_MAIN_GAME":    1,
		"GAME_STATE_INVENTORY":    2,
		"GAME_STATE_LOAD_SAVE":    3,
		"GAME_STATE_SPECIAL_MENU": 4,
	}

	if GAME_STATE_MAIN_MENU != expectedStates["GAME_STATE_MAIN_MENU"] {
		t.Errorf("GAME_STATE_MAIN_MENU should be %d, got %d", expectedStates["GAME_STATE_MAIN_MENU"], GAME_STATE_MAIN_MENU)
	}

	if GAME_STATE_MAIN_GAME != expectedStates["GAME_STATE_MAIN_GAME"] {
		t.Errorf("GAME_STATE_MAIN_GAME should be %d, got %d", expectedStates["GAME_STATE_MAIN_GAME"], GAME_STATE_MAIN_GAME)
	}

	if GAME_STATE_INVENTORY != expectedStates["GAME_STATE_INVENTORY"] {
		t.Errorf("GAME_STATE_INVENTORY should be %d, got %d", expectedStates["GAME_STATE_INVENTORY"], GAME_STATE_INVENTORY)
	}

	if GAME_STATE_LOAD_SAVE != expectedStates["GAME_STATE_LOAD_SAVE"] {
		t.Errorf("GAME_STATE_LOAD_SAVE should be %d, got %d", expectedStates["GAME_STATE_LOAD_SAVE"], GAME_STATE_LOAD_SAVE)
	}

	if GAME_STATE_SPECIAL_MENU != expectedStates["GAME_STATE_SPECIAL_MENU"] {
		t.Errorf("GAME_STATE_SPECIAL_MENU should be %d, got %d", expectedStates["GAME_STATE_SPECIAL_MENU"], GAME_STATE_SPECIAL_MENU)
	}
}

// Benchmark tests
func BenchmarkNewGameStateManager(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewGameStateManager()
	}
}

func BenchmarkUpdateGameState(b *testing.B) {
	manager := NewGameStateManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.UpdateGameState(GAME_STATE_MAIN_GAME)
	}
}
