package game

import (
	"testing"

	"github.com/OpenBiohazard2/OpenBiohazard2/world"
)

func TestNewGame(t *testing.T) {
	gameDef := NewGame(1, 2, 3)

	if gameDef.StageId != 1 {
		t.Errorf("Expected StageId 1, got %d", gameDef.StageId)
	}
	if gameDef.RoomId != 2 {
		t.Errorf("Expected RoomId 2, got %d", gameDef.RoomId)
	}
	if gameDef.CameraId != 3 {
		t.Errorf("Expected CameraId 3, got %d", gameDef.CameraId)
	}
	if gameDef.StateStatus != GAME_LOAD_ROOM {
		t.Errorf("Expected StateStatus %d, got %d", GAME_LOAD_ROOM, gameDef.StateStatus)
	}
	if gameDef.GameWorld == nil {
		t.Fatal("Expected GameWorld to be initialized, got nil")
	}
}

func TestChangeCamera_WithinBounds(t *testing.T) {
	gameDef := NewGame(1, 0, 0)
	gameDef.GameWorld.GameRoom = &world.Room{MaxCamerasInRoom: 5}

	gameDef.ChangeCamera(3)

	if gameDef.CameraId != 3 {
		t.Errorf("Expected CameraId 3, got %d", gameDef.CameraId)
	}
	if gameDef.StateStatus != GAME_LOAD_CAMERA {
		t.Errorf("Expected StateStatus %d, got %d", GAME_LOAD_CAMERA, gameDef.StateStatus)
	}
}

func TestChangeCamera_ClampsAboveMax(t *testing.T) {
	gameDef := NewGame(1, 0, 0)
	gameDef.GameWorld.GameRoom = &world.Room{MaxCamerasInRoom: 5}

	gameDef.ChangeCamera(10)

	expected := 4 // MaxCamerasInRoom - 1
	if gameDef.CameraId != expected {
		t.Errorf("Expected CameraId clamped to %d, got %d", expected, gameDef.CameraId)
	}
}

func TestChangeCamera_ClampsBelowZero(t *testing.T) {
	gameDef := NewGame(1, 0, 0)
	gameDef.GameWorld.GameRoom = &world.Room{MaxCamerasInRoom: 5}

	gameDef.ChangeCamera(-1)

	if gameDef.CameraId != 0 {
		t.Errorf("Expected CameraId clamped to 0, got %d", gameDef.CameraId)
	}
}
