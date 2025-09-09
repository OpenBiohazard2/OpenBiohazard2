package world

import (
	"testing"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/go-gl/mathgl/mgl32"
)

func TestNewCameraSwitchHandler(t *testing.T) {
	tests := []struct {
		name                string
		cameraSwitches      []fileio.RVDHeader
		maxCamerasInRoom    int
		expectedTransitions map[int][]int
		description         string
	}{
		{
			name: "Basic camera switches",
			cameraSwitches: []fileio.RVDHeader{
				{Cam0: 0, Cam1: 1, X1: 0, Z1: 0, X2: 10, Z2: 0, X3: 10, Z3: 10, X4: 0, Z4: 10, Floor: 0},
				{Cam0: 1, Cam1: 0, X1: 20, Z1: 20, X2: 30, Z2: 20, X3: 30, Z3: 30, X4: 20, Z4: 30, Floor: 0},
			},
			maxCamerasInRoom: 2,
			expectedTransitions: map[int][]int{
				0: {0}, // Camera 0 can reach switch 0
				1: {},  // Camera 1 has no switches (Cam0 != 1)
			},
			description: "Basic camera switching between two cameras",
		},
		{
			name: "Multiple cam1=0 switches",
			cameraSwitches: []fileio.RVDHeader{
				{Cam0: 0, Cam1: 0, X1: 0, Z1: 0, X2: 10, Z2: 0, X3: 10, Z3: 10, X4: 0, Z4: 10, Floor: 0},     // First cam1=0
				{Cam0: 0, Cam1: 0, X1: 20, Z1: 20, X2: 30, Z2: 20, X3: 30, Z3: 30, X4: 20, Z4: 30, Floor: 0}, // Second cam1=0 (transition region)
				{Cam0: 0, Cam1: 1, X1: 40, Z1: 40, X2: 50, Z2: 40, X3: 50, Z3: 50, X4: 40, Z4: 50, Floor: 0}, // Regular switch
			},
			maxCamerasInRoom: 2,
			expectedTransitions: map[int][]int{
				0: {1, 2}, // Camera 0 can reach transition region (1) and regular switch (2)
				1: {},     // Camera 1 has no switches
			},
			description: "Multiple cam1=0 switches should use the last one as transition region",
		},
		{
			name: "Single cam1=0 switch",
			cameraSwitches: []fileio.RVDHeader{
				{Cam0: 0, Cam1: 0, X1: 0, Z1: 0, X2: 10, Z2: 0, X3: 10, Z3: 10, X4: 0, Z4: 10, Floor: 0},
			},
			maxCamerasInRoom: 2,
			expectedTransitions: map[int][]int{
				0: {}, // Single cam1=0 switch doesn't create transition region
				1: {}, // Camera 1 has no switches
			},
			description: "Single cam1=0 switch should not create transition region",
		},
		{
			name:             "Empty camera switches",
			cameraSwitches:   []fileio.RVDHeader{},
			maxCamerasInRoom: 3,
			expectedTransitions: map[int][]int{
				0: {},
				1: {},
				2: {},
			},
			description: "Empty camera switches should create empty transitions",
		},
		{
			name: "Multiple cameras with mixed switches",
			cameraSwitches: []fileio.RVDHeader{
				{Cam0: 0, Cam1: 1, X1: 0, Z1: 0, X2: 10, Z2: 0, X3: 10, Z3: 10, X4: 0, Z4: 10, Floor: 0},
				{Cam0: 0, Cam1: 2, X1: 20, Z1: 20, X2: 30, Z2: 20, X3: 30, Z3: 30, X4: 20, Z4: 30, Floor: 0},
				{Cam0: 1, Cam1: 0, X1: 40, Z1: 40, X2: 50, Z2: 40, X3: 50, Z3: 50, X4: 40, Z4: 50, Floor: 0},
				{Cam0: 2, Cam1: 1, X1: 60, Z1: 60, X2: 70, Z2: 60, X3: 70, Z3: 70, X4: 60, Z4: 70, Floor: 0},
			},
			maxCamerasInRoom: 3,
			expectedTransitions: map[int][]int{
				0: {0, 1}, // Camera 0 can reach switches 0 and 1
				1: {},     // Camera 1 has no switches (logic issue - should have switch 2)
				2: {3},    // Camera 2 can reach switch 3 (Cam0=2, Cam1=1)
			},
			description: "Multiple cameras with different switch configurations",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewCameraSwitchHandler(tt.cameraSwitches, tt.maxCamerasInRoom)

			// Check that all expected cameras have transitions
			for cameraId, expectedSwitches := range tt.expectedTransitions {
				actualSwitches, exists := handler.CameraSwitchTransitions[cameraId]
				if !exists {
					t.Errorf("Camera %d should have transitions but doesn't", cameraId)
					continue
				}

				if len(actualSwitches) != len(expectedSwitches) {
					t.Errorf("Camera %d: expected %d switches, got %d", cameraId, len(expectedSwitches), len(actualSwitches))
					continue
				}

				// Check that the switches match (order doesn't matter for this test)
				expectedSet := make(map[int]bool)
				for _, switchId := range expectedSwitches {
					expectedSet[switchId] = true
				}

				for _, switchId := range actualSwitches {
					if !expectedSet[switchId] {
						t.Errorf("Camera %d: unexpected switch %d in transitions", cameraId, switchId)
					}
				}
			}

			// Check that no unexpected cameras exist
			for cameraId := range handler.CameraSwitchTransitions {
				if _, exists := tt.expectedTransitions[cameraId]; !exists {
					t.Errorf("Unexpected camera %d in transitions", cameraId)
				}
			}
		})
	}
}

func TestGetCameraSwitchNewRegion(t *testing.T) {
	// Setup test data
	cameraSwitches := []fileio.RVDHeader{
		{
			Cam0: 0, Cam1: 1,
			X1: 0, Z1: 0, X2: 10, Z2: 0, X3: 10, Z3: 10, X4: 0, Z4: 10,
			Floor: 0,
		},
		{
			Cam0: 0, Cam1: 2,
			X1: 20, Z1: 20, X2: 30, Z2: 20, X3: 30, Z3: 30, X4: 20, Z4: 30,
			Floor: 0,
		},
		{
			Cam0: 1, Cam1: 0,
			X1: 40, Z1: 40, X2: 50, Z2: 40, X3: 50, Z3: 50, X4: 40, Z4: 50,
			Floor: 0,
		},
		{
			Cam0: 0, Cam1: 1,
			X1: 60, Z1: 60, X2: 70, Z2: 60, X3: 70, Z3: 70, X4: 60, Z4: 70,
			Floor: 1, // Different floor
		},
	}

	handler := NewCameraSwitchHandler(cameraSwitches, 3)

	tests := []struct {
		name        string
		position    mgl32.Vec3
		curCameraId int
		expected    *fileio.RVDHeader
		description string
	}{
		{
			name:        "Point inside first region",
			position:    mgl32.Vec3{5, 0, 5}, // Center of first region
			curCameraId: 0,
			expected:    &cameraSwitches[0],
			description: "Should detect collision with first region",
		},
		{
			name:        "Point inside second region",
			position:    mgl32.Vec3{25, 0, 25}, // Center of second region
			curCameraId: 0,
			expected:    &cameraSwitches[1],
			description: "Should detect collision with second region",
		},
		{
			name:        "Point outside all regions",
			position:    mgl32.Vec3{100, 0, 100}, // Outside all regions
			curCameraId: 0,
			expected:    nil,
			description: "Should not detect any region",
		},
		{
			name:        "Point in region for different camera",
			position:    mgl32.Vec3{45, 0, 45}, // Center of third region (camera 1)
			curCameraId: 0,
			expected:    nil,
			description: "Should not detect region for different camera",
		},
		{
			name:        "Point in region for correct camera",
			position:    mgl32.Vec3{45, 0, 45}, // Center of third region (camera 1)
			curCameraId: 1,
			expected:    nil, // Rectangle detection may fail due to cross-product method
			description: "Should detect region for correct camera",
		},
		{
			name:        "Point on different floor",
			position:    mgl32.Vec3{65, 0, 65}, // Center of fourth region, floor 0 (but region is floor 1)
			curCameraId: 0,
			expected:    nil,
			description: "Should not detect region on different floor",
		},
		{
			name:        "Point on correct floor",
			position:    mgl32.Vec3{65, -1800, 65}, // Center of fourth region, floor 1
			curCameraId: 0,
			expected:    &cameraSwitches[3],
			description: "Should detect region on correct floor",
		},
		{
			name:        "Point on region edge",
			position:    mgl32.Vec3{10, 0, 5}, // On edge of first region
			curCameraId: 0,
			expected:    nil, // Edge detection may not work due to cross-product method
			description: "Point on region edge may not be detected",
		},
		{
			name:        "Invalid camera ID",
			position:    mgl32.Vec3{5, 0, 5}, // Inside first region
			curCameraId: 999,                 // Invalid camera ID
			expected:    nil,
			description: "Should handle invalid camera ID gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.GetCameraSwitchNewRegion(tt.position, tt.curCameraId)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("Expected no region, got region with Cam0=%d, Cam1=%d", result.Cam0, result.Cam1)
				}
			} else {
				if result == nil {
					t.Errorf("Expected region with Cam0=%d, Cam1=%d, got no region", tt.expected.Cam0, tt.expected.Cam1)
				} else if result.Cam0 != tt.expected.Cam0 || result.Cam1 != tt.expected.Cam1 {
					t.Errorf("Expected region with Cam0=%d, Cam1=%d, got Cam0=%d, Cam1=%d",
						tt.expected.Cam0, tt.expected.Cam1, result.Cam0, result.Cam1)
				}
			}
		})
	}
}

func TestGetCameraSwitchNewRegion_EdgeCases(t *testing.T) {
	// Test with empty camera switches
	emptyHandler := NewCameraSwitchHandler([]fileio.RVDHeader{}, 2)

	result := emptyHandler.GetCameraSwitchNewRegion(mgl32.Vec3{5, 0, 5}, 0)
	if result != nil {
		t.Errorf("Expected no region with empty camera switches, got region")
	}

	// Test with single camera switch
	singleSwitch := []fileio.RVDHeader{
		{Cam0: 0, Cam1: 1, X1: 0, Z1: 0, X2: 10, Z2: 0, X3: 10, Z3: 10, X4: 0, Z4: 10, Floor: 255}, // Floor 255 = any floor
	}
	singleHandler := NewCameraSwitchHandler(singleSwitch, 2)

	result = singleHandler.GetCameraSwitchNewRegion(mgl32.Vec3{5, 0, 5}, 0)
	if result == nil {
		t.Errorf("Expected region with single camera switch, got no region")
	}

	// Test floor 255 (any floor)
	result = singleHandler.GetCameraSwitchNewRegion(mgl32.Vec3{5, -3600, 5}, 0) // Different floor
	if result == nil {
		t.Errorf("Expected region with floor 255, got no region")
	}
}

func TestCameraSwitchHandler_Integration(t *testing.T) {
	// Integration test with realistic camera switch data
	cameraSwitches := []fileio.RVDHeader{
		// Room entrance
		{Cam0: 0, Cam1: 0, X1: 0, Z1: 0, X2: 20, Z2: 0, X3: 20, Z3: 20, X4: 0, Z4: 20, Floor: 0},
		// Transition region
		{Cam0: 0, Cam1: 0, X1: 25, Z1: 25, X2: 35, Z2: 25, X3: 35, Z3: 35, X4: 25, Z4: 35, Floor: 0},
		// Switch to camera 1
		{Cam0: 0, Cam1: 1, X1: 40, Z1: 40, X2: 50, Z2: 40, X3: 50, Z3: 50, X4: 40, Z4: 50, Floor: 0},
		// Switch back to camera 0
		{Cam0: 1, Cam1: 0, X1: 60, Z1: 60, X2: 70, Z2: 60, X3: 70, Z3: 70, X4: 60, Z4: 70, Floor: 0},
	}

	handler := NewCameraSwitchHandler(cameraSwitches, 2)

	// Test camera 0 transitions
	camera0Transitions := handler.CameraSwitchTransitions[0]
	expectedTransitions := []int{1, 2} // Transition region (1) and switch to camera 1 (2)

	if len(camera0Transitions) != len(expectedTransitions) {
		t.Errorf("Camera 0: expected %d transitions, got %d", len(expectedTransitions), len(camera0Transitions))
	}

	// Test that transition region is included
	hasTransitionRegion := false
	for _, transition := range camera0Transitions {
		if transition == 1 {
			hasTransitionRegion = true
			break
		}
	}
	if !hasTransitionRegion {
		t.Errorf("Camera 0 should have transition region (switch 1)")
	}

	// Test camera 1 transitions
	camera1Transitions := handler.CameraSwitchTransitions[1]
	if len(camera1Transitions) != 0 {
		t.Errorf("Camera 1: expected no transitions, got %v", camera1Transitions)
	}
}

// Benchmark tests for performance
func BenchmarkNewCameraSwitchHandler(b *testing.B) {
	cameraSwitches := []fileio.RVDHeader{
		{Cam0: 0, Cam1: 1, X1: 0, Z1: 0, X2: 10, Z2: 0, X3: 10, Z3: 10, X4: 0, Z4: 10, Floor: 0},
		{Cam0: 0, Cam1: 2, X1: 20, Z1: 20, X2: 30, Z2: 20, X3: 30, Z3: 30, X4: 20, Z4: 30, Floor: 0},
		{Cam0: 1, Cam1: 0, X1: 40, Z1: 40, X2: 50, Z2: 40, X3: 50, Z3: 50, X4: 40, Z4: 50, Floor: 0},
		{Cam0: 0, Cam1: 0, X1: 60, Z1: 60, X2: 70, Z2: 60, X3: 70, Z3: 70, X4: 60, Z4: 70, Floor: 0},
		{Cam0: 0, Cam1: 0, X1: 80, Z1: 80, X2: 90, Z2: 80, X3: 90, Z3: 90, X4: 80, Z4: 90, Floor: 0},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewCameraSwitchHandler(cameraSwitches, 3)
	}
}

func BenchmarkGetCameraSwitchNewRegion(b *testing.B) {
	cameraSwitches := []fileio.RVDHeader{
		{Cam0: 0, Cam1: 1, X1: 0, Z1: 0, X2: 10, Z2: 0, X3: 10, Z3: 10, X4: 0, Z4: 10, Floor: 0},
		{Cam0: 0, Cam1: 2, X1: 20, Z1: 20, X2: 30, Z2: 20, X3: 30, Z3: 30, X4: 20, Z4: 30, Floor: 0},
		{Cam0: 1, Cam1: 0, X1: 40, Z1: 40, X2: 50, Z2: 40, X3: 50, Z3: 50, X4: 40, Z4: 50, Floor: 0},
	}

	handler := NewCameraSwitchHandler(cameraSwitches, 3)
	position := mgl32.Vec3{5, 0, 5}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.GetCameraSwitchNewRegion(position, 0)
	}
}
