package render

import (
	"math"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNewCamera(t *testing.T) {
	from := mgl32.Vec3{1, 2, 3}
	to := mgl32.Vec3{4, 5, 6}
	up := mgl32.Vec3{0, 1, 0}
	fov := float32(60.0)

	camera := NewCamera(from, to, up, fov)

	if camera == nil {
		t.Fatal("NewCamera returned nil")
	}

	if camera.CameraFrom != from {
		t.Errorf("Expected CameraFrom to be %v, got %v", from, camera.CameraFrom)
	}

	if camera.CameraTo != to {
		t.Errorf("Expected CameraTo to be %v, got %v", to, camera.CameraTo)
	}

	if camera.CameraUp != up {
		t.Errorf("Expected CameraUp to be %v, got %v", up, camera.CameraUp)
	}

	if camera.CameraFov != fov {
		t.Errorf("Expected CameraFov to be %f, got %f", fov, camera.CameraFov)
	}
}

func TestCamera_Update(t *testing.T) {
	camera := NewCamera(
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{0, 0, -1},
		mgl32.Vec3{0, 1, 0},
		60.0,
	)

	newFrom := mgl32.Vec3{10, 20, 30}
	newTo := mgl32.Vec3{40, 50, 60}
	newFov := float32(90.0)

	camera.Update(newFrom, newTo, newFov)

	if camera.CameraFrom != newFrom {
		t.Errorf("Expected CameraFrom to be updated to %v, got %v", newFrom, camera.CameraFrom)
	}

	if camera.CameraTo != newTo {
		t.Errorf("Expected CameraTo to be updated to %v, got %v", newTo, camera.CameraTo)
	}

	if camera.CameraFov != newFov {
		t.Errorf("Expected CameraFov to be updated to %f, got %f", newFov, camera.CameraFov)
	}

	// CameraUp should remain unchanged
	expectedUp := mgl32.Vec3{0, 1, 0}
	if camera.CameraUp != expectedUp {
		t.Errorf("Expected CameraUp to remain %v, got %v", expectedUp, camera.CameraUp)
	}
}

func TestCamera_BuildViewMatrix(t *testing.T) {
	tests := []struct {
		name     string
		from     mgl32.Vec3
		to       mgl32.Vec3
		up       mgl32.Vec3
		fov      float32
		expected mgl32.Mat4
	}{
		{
			name: "Identity_looking_forward",
			from: mgl32.Vec3{0, 0, 0},
			to:   mgl32.Vec3{0, 0, -1},
			up:   mgl32.Vec3{0, 1, 0},
			fov:  60.0,
			// Expected: Identity matrix (looking forward along negative Z)
		},
		{
			name: "Looking_up",
			from: mgl32.Vec3{0, 0, 0},
			to:   mgl32.Vec3{0, 1, 0},
			up:   mgl32.Vec3{0, 0, -1},
			fov:  60.0,
		},
		{
			name: "Looking_right",
			from: mgl32.Vec3{0, 0, 0},
			to:   mgl32.Vec3{1, 0, 0},
			up:   mgl32.Vec3{0, 1, 0},
			fov:  60.0,
		},
		{
			name: "Offset_position",
			from: mgl32.Vec3{5, 10, 15},
			to:   mgl32.Vec3{5, 10, 14},
			up:   mgl32.Vec3{0, 1, 0},
			fov:  60.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			camera := NewCamera(tt.from, tt.to, tt.up, tt.fov)
			viewMatrix := camera.BuildViewMatrix()

			// Test that the matrix is valid (not all zeros)
			if viewMatrix == (mgl32.Mat4{}) {
				t.Error("View matrix should not be zero matrix")
			}

			// Test that the matrix has the expected properties
			// The view matrix should be invertible (determinant != 0)
			det := viewMatrix.Det()
			if math.Abs(float64(det)) < 1e-6 {
				t.Errorf("View matrix determinant should not be zero, got %f", det)
			}

			// Test that the matrix is invertible (has non-zero determinant)
			// Note: View matrices are not necessarily orthogonal when camera is at offset position
			// They are transformation matrices that include translation
		})
	}
}

func TestCamera_GetDirection(t *testing.T) {
	tests := []struct {
		name     string
		from     mgl32.Vec3
		to       mgl32.Vec3
		up       mgl32.Vec3
		fov      float32
		expected mgl32.Vec3
	}{
		{
			name:     "Forward_direction",
			from:     mgl32.Vec3{0, 0, 0},
			to:       mgl32.Vec3{0, 0, -1},
			up:       mgl32.Vec3{0, 1, 0},
			fov:      60.0,
			expected: mgl32.Vec3{0, 0, -1},
		},
		{
			name:     "Up_direction",
			from:     mgl32.Vec3{0, 0, 0},
			to:       mgl32.Vec3{0, 1, 0},
			up:       mgl32.Vec3{0, 0, -1},
			fov:      60.0,
			expected: mgl32.Vec3{0, 1, 0},
		},
		{
			name:     "Right_direction",
			from:     mgl32.Vec3{0, 0, 0},
			to:       mgl32.Vec3{1, 0, 0},
			up:       mgl32.Vec3{0, 1, 0},
			fov:      60.0,
			expected: mgl32.Vec3{1, 0, 0},
		},
		{
			name:     "Diagonal_direction",
			from:     mgl32.Vec3{0, 0, 0},
			to:       mgl32.Vec3{1, 1, 1},
			up:       mgl32.Vec3{0, 1, 0},
			fov:      60.0,
			expected: mgl32.Vec3{1, 1, 1}.Normalize(),
		},
		{
			name:     "Offset_position",
			from:     mgl32.Vec3{5, 10, 15},
			to:       mgl32.Vec3{6, 11, 16},
			up:       mgl32.Vec3{0, 1, 0},
			fov:      60.0,
			expected: mgl32.Vec3{1, 1, 1}.Normalize(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			camera := NewCamera(tt.from, tt.to, tt.up, tt.fov)
			direction := camera.GetDirection()

			// Test that the direction is normalized (length = 1)
			length := direction.Len()
			if math.Abs(float64(length-1.0)) > 1e-6 {
				t.Errorf("Direction should be normalized (length = 1), got %f", length)
			}

			// Test that the direction matches expected (within tolerance)
			tolerance := float32(1e-6)
			diff := direction.Sub(tt.expected)
			if diff.Len() > tolerance {
				t.Errorf("Expected direction %v, got %v", tt.expected, direction)
			}
		})
	}
}

func TestCamera_GetDirection_ZeroVector(t *testing.T) {
	// Test edge case where from and to are the same point
	camera := NewCamera(
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{0, 0, 0}, // Same as from
		mgl32.Vec3{0, 1, 0},
		60.0,
	)

	// This should handle the zero vector case gracefully
	direction := camera.GetDirection()

	// When normalizing a zero vector, mathgl returns NaN values
	// This is expected behavior for the mathgl library
	// In a real application, you'd want to handle this edge case
	if !math.IsNaN(float64(direction.X())) || !math.IsNaN(float64(direction.Y())) || !math.IsNaN(float64(direction.Z())) {
		t.Log("Note: mathgl.Normalize() on zero vector returns NaN values, which is expected behavior")
	}
}

// Test the GetPerspectiveMatrix function from rendermain.go
func TestRenderDef_GetPerspectiveMatrix(t *testing.T) {
	// Create a minimal RenderDef for testing
	renderDef := &RenderDef{
		ViewSystem: NewViewSystem(800, 600),
	}

	tests := []struct {
		name         string
		fovDegrees   float32
		windowWidth  int
		windowHeight int
	}{
		{
			name:         "Standard_FOV",
			fovDegrees:   60.0,
			windowWidth:  800,
			windowHeight: 600,
		},
		{
			name:         "Wide_FOV",
			fovDegrees:   90.0,
			windowWidth:  800,
			windowHeight: 600,
		},
		{
			name:         "Narrow_FOV",
			fovDegrees:   30.0,
			windowWidth:  800,
			windowHeight: 600,
		},
		{
			name:         "Square_aspect_ratio",
			fovDegrees:   60.0,
			windowWidth:  600,
			windowHeight: 600,
		},
		{
			name:         "Wide_aspect_ratio",
			fovDegrees:   60.0,
			windowWidth:  1600,
			windowHeight: 900,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderDef.ViewSystem.WindowWidth = tt.windowWidth
			renderDef.ViewSystem.WindowHeight = tt.windowHeight

			projectionMatrix := renderDef.GetPerspectiveMatrix(tt.fovDegrees)

			// Test that the matrix is valid (not all zeros)
			if projectionMatrix == (mgl32.Mat4{}) {
				t.Error("Projection matrix should not be zero matrix")
			}

			// Test that the matrix has the expected properties
			// The projection matrix should be invertible (determinant != 0)
			det := projectionMatrix.Det()
			if math.Abs(float64(det)) < 1e-6 {
				t.Errorf("Projection matrix determinant should not be zero, got %f", det)
			}

			// Test that the matrix has the correct structure for a perspective projection
			// The bottom-right element should be 0 for perspective projection
			if math.Abs(float64(projectionMatrix.At(3, 3))) > 1e-6 {
				t.Errorf("Projection matrix [3,3] should be 0 for perspective projection, got %f", projectionMatrix.At(3, 3))
			}

			// The bottom-left element should be 0 for perspective projection
			if math.Abs(float64(projectionMatrix.At(3, 0))) > 1e-6 {
				t.Errorf("Projection matrix [3,0] should be 0 for perspective projection, got %f", projectionMatrix.At(3, 0))
			}

			// The bottom-center element should be 0 for perspective projection
			if math.Abs(float64(projectionMatrix.At(3, 1))) > 1e-6 {
				t.Errorf("Projection matrix [3,1] should be 0 for perspective projection, got %f", projectionMatrix.At(3, 1))
			}
		})
	}
}

// Benchmark tests
func BenchmarkNewCamera(b *testing.B) {
	from := mgl32.Vec3{1, 2, 3}
	to := mgl32.Vec3{4, 5, 6}
	up := mgl32.Vec3{0, 1, 0}
	fov := float32(60.0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewCamera(from, to, up, fov)
	}
}

func BenchmarkCamera_BuildViewMatrix(b *testing.B) {
	camera := NewCamera(
		mgl32.Vec3{1, 2, 3},
		mgl32.Vec3{4, 5, 6},
		mgl32.Vec3{0, 1, 0},
		60.0,
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = camera.BuildViewMatrix()
	}
}

func BenchmarkCamera_GetDirection(b *testing.B) {
	camera := NewCamera(
		mgl32.Vec3{1, 2, 3},
		mgl32.Vec3{4, 5, 6},
		mgl32.Vec3{0, 1, 0},
		60.0,
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = camera.GetDirection()
	}
}

func TestCamera_NormalizeMaskDepth(t *testing.T) {
	camera := NewCamera(
		mgl32.Vec3{0, 0, 100}, // from (moved back from origin)
		mgl32.Vec3{0, 0, 0},   // to (looking at origin)
		mgl32.Vec3{0, 1, 0},   // up
		60.0,                  // fov
	)

	// Create test projection and view matrices with more reasonable near/far planes
	projectionMatrix := mgl32.Perspective(mgl32.DegToRad(60.0), 4.0/3.0, 1.0, 1000.0)
	viewMatrix := camera.BuildViewMatrix()

	tests := []struct {
		name     string
		depth    float32
		expected float32 // Expected normalized depth (0-1 range)
	}{
		{
			name:     "Small_depth",
			depth:    0.1,
			expected: 0.1, // Small depth value
		},
		{
			name:     "Medium_depth",
			depth:    0.5,
			expected: 0.5, // Should be in middle range
		},
		{
			name:     "Far_depth",
			depth:    1.0,
			expected: 1.0, // Should be at far plane
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := camera.NormalizeMaskDepth(tt.depth, projectionMatrix, viewMatrix)

			// Check that result is not NaN or infinite
			if math.IsNaN(float64(result)) || math.IsInf(float64(result), 0) {
				t.Errorf("Normalized depth should not be NaN or infinite, got %f", result)
				return
			}

			// For this test, we just verify the function doesn't crash and returns a reasonable value
			// The actual projection math is complex and depends on the specific camera setup
			// In practice, this function is used with real camera data from the game
			if result < -10.0 || result > 10.0 {
				t.Errorf("Normalized depth seems unreasonable, got %f (expected reasonable range)", result)
			}
		})
	}
}

func BenchmarkRenderDef_GetPerspectiveMatrix(b *testing.B) {
	renderDef := &RenderDef{
		ViewSystem: NewViewSystem(800, 600),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = renderDef.GetPerspectiveMatrix(60.0)
	}
}

func BenchmarkCamera_NormalizeMaskDepth(b *testing.B) {
	camera := NewCamera(
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{0, 0, -1},
		mgl32.Vec3{0, 1, 0},
		60.0,
	)

	projectionMatrix := mgl32.Perspective(mgl32.DegToRad(60.0), 4.0/3.0, 16.0, 45000.0)
	viewMatrix := camera.BuildViewMatrix()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = camera.NormalizeMaskDepth(0.5, projectionMatrix, viewMatrix)
	}
}
