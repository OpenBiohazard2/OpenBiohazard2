package game

import "testing"

func TestGetRoomFilename(t *testing.T) {
	tests := []struct {
		name       string
		stageId    int
		roomId     int
		playerNum  int
		expectPath string
	}{
		{"stage1_room0_player0", 1, 0x00, 0, "data/Pl0/RDP/ROOM1000.RDT"},
		{"stage1_room1a_player0", 1, 0x1a, 0, "data/Pl0/RDP/ROOM11a0.RDT"},
		{"stage2_roomFF_player1", 2, 0xff, 1, "data/Pl0/RDP/ROOM2ff1.RDT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameDef := &GameDef{StageId: tt.stageId, RoomId: tt.roomId}
			got := gameDef.GetRoomFilename(tt.playerNum)
			if got != tt.expectPath {
				t.Errorf("Expected %q, got %q", tt.expectPath, got)
			}
		})
	}
}

func TestGetBackgroundImageNumber(t *testing.T) {
	tests := []struct {
		name     string
		stageId  int
		roomId   int
		cameraId int
		expected int
	}{
		{"stage1_room0_camera0", 1, 0, 0, 0},
		{"stage1_room0_camera1", 1, 0, 1, 1},
		{"stage1_room1_camera0", 1, 1, 0, 16},
		{"stage2_room0_camera0", 2, 0, 0, 512},
		{"stage2_room3_camera2", 2, 3, 2, 512 + 48 + 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameDef := &GameDef{StageId: tt.stageId, RoomId: tt.roomId, CameraId: tt.cameraId}
			got := gameDef.GetBackgroundImageNumber()
			if got != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, got)
			}
		})
	}
}

func TestNextRoom(t *testing.T) {
	gameDef := &GameDef{RoomId: 5, CameraId: 3}
	gameDef.NextRoom()

	if gameDef.RoomId != 6 {
		t.Errorf("Expected RoomId to advance to 6, got %d", gameDef.RoomId)
	}
	if gameDef.CameraId != 0 {
		t.Errorf("Expected CameraId to reset to 0, got %d", gameDef.CameraId)
	}
}

func TestNextRoom_ClampsAtUpperBound(t *testing.T) {
	gameDef := &GameDef{RoomId: 31}
	gameDef.NextRoom()

	if gameDef.RoomId != 31 {
		t.Errorf("Expected RoomId to clamp at 31, got %d", gameDef.RoomId)
	}
}

func TestPrevRoom(t *testing.T) {
	gameDef := &GameDef{RoomId: 5, CameraId: 3}
	gameDef.PrevRoom()

	if gameDef.RoomId != 4 {
		t.Errorf("Expected RoomId to decrease to 4, got %d", gameDef.RoomId)
	}
	if gameDef.CameraId != 0 {
		t.Errorf("Expected CameraId to reset to 0, got %d", gameDef.CameraId)
	}
}

func TestPrevRoom_ClampsAtLowerBound(t *testing.T) {
	gameDef := &GameDef{RoomId: 0}
	gameDef.PrevRoom()

	if gameDef.RoomId != 0 {
		t.Errorf("Expected RoomId to clamp at 0, got %d", gameDef.RoomId)
	}
}
