package render

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNewViewSystem(t *testing.T) {
	windowWidth, windowHeight := 800, 600
	vs := NewViewSystemForTesting(windowWidth, windowHeight)

	// Test basic initialization
	if vs == nil {
		t.Fatal("NewViewSystem returned nil")
	}

	// Test window dimensions
	if vs.WindowWidth != windowWidth {
		t.Errorf("Expected WindowWidth %d, got %d", windowWidth, vs.WindowWidth)
	}
	if vs.WindowHeight != windowHeight {
		t.Errorf("Expected WindowHeight %d, got %d", windowHeight, vs.WindowHeight)
	}

	// Test camera initialization
	if vs.Camera == nil {
		t.Fatal("Camera should not be nil")
	}

	// Test matrices are initialized
	if vs.ProjectionMatrix == (mgl32.Mat4{}) {
		t.Error("ProjectionMatrix should be initialized")
	}
	if vs.ViewMatrix == (mgl32.Mat4{}) {
		t.Error("ViewMatrix should be initialized")
	}
}

func TestViewSystem_UpdateMatrices(t *testing.T) {
	vs := NewViewSystemForTesting(800, 600)

	// Update matrices
	vs.UpdateMatrices()

	// Matrices should be different (unless camera is at origin looking forward)
	// For this test, we'll just verify the method doesn't crash
	if vs.ProjectionMatrix == (mgl32.Mat4{}) {
		t.Error("ProjectionMatrix should not be zero matrix after update")
	}
	if vs.ViewMatrix == (mgl32.Mat4{}) {
		t.Error("ViewMatrix should not be zero matrix after update")
	}

	// Test that matrices are valid (not all zeros)
	projectionValid := false
	viewValid := false

	for i := 0; i < 16; i++ {
		if vs.ProjectionMatrix[i] != 0 {
			projectionValid = true
		}
		if vs.ViewMatrix[i] != 0 {
			viewValid = true
		}
	}

	if !projectionValid {
		t.Error("ProjectionMatrix should have non-zero values")
	}
	if !viewValid {
		t.Error("ViewMatrix should have non-zero values")
	}
}

func TestViewSystem_GetViewMatrix(t *testing.T) {
	vs := NewViewSystemForTesting(800, 600)

	viewMatrix := vs.GetViewMatrix()

	// Should return the same matrix as stored
	if viewMatrix != vs.ViewMatrix {
		t.Error("GetViewMatrix should return the stored ViewMatrix")
	}
}

func TestViewSystem_GetProjectionMatrix(t *testing.T) {
	vs := NewViewSystemForTesting(800, 600)

	projectionMatrix := vs.GetProjectionMatrix()

	// Should return the same matrix as stored
	if projectionMatrix != vs.ProjectionMatrix {
		t.Error("GetProjectionMatrix should return the stored ProjectionMatrix")
	}
}

func TestViewSystem_GetPerspectiveMatrix(t *testing.T) {
	tests := []struct {
		name         string
		fovDegrees   float32
		windowWidth  int
		windowHeight int
	}{
		{"Standard_FOV", 60.0, 800, 600},
		{"Wide_FOV", 90.0, 800, 600},
		{"Narrow_FOV", 30.0, 800, 600},
		{"Square_aspect_ratio", 60.0, 600, 600},
		{"Wide_aspect_ratio", 60.0, 1200, 600},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vs := NewViewSystemForTesting(tt.windowWidth, tt.windowHeight)
			matrix := vs.GetPerspectiveMatrix(tt.fovDegrees)

			// Test that the matrix is valid (not all zeros)
			if matrix == (mgl32.Mat4{}) {
				t.Error("Perspective matrix should not be zero matrix")
			}

			// Test that the matrix has expected properties
			// The matrix should be a valid perspective projection
			// We can test some basic properties without getting into OpenGL specifics
			hasNonZeroValues := false
			for i := 0; i < 16; i++ {
				if matrix[i] != 0 {
					hasNonZeroValues = true
					break
				}
			}
			if !hasNonZeroValues {
				t.Error("Perspective matrix should have non-zero values")
			}
		})
	}
}

func TestViewSystem_GetPerspectiveMatrix_DifferentAspectRatios(t *testing.T) {
	// Test different aspect ratios produce different matrices
	vs1 := NewViewSystemForTesting(800, 600)  // 4:3
	vs2 := NewViewSystemForTesting(1200, 600) // 2:1
	vs3 := NewViewSystemForTesting(600, 600)  // 1:1

	matrix1 := vs1.GetPerspectiveMatrix(60.0)
	matrix2 := vs2.GetPerspectiveMatrix(60.0)
	matrix3 := vs3.GetPerspectiveMatrix(60.0)

	// Different aspect ratios should produce different matrices
	if matrix1 == matrix2 {
		t.Error("Different aspect ratios should produce different matrices")
	}
	if matrix1 == matrix3 {
		t.Error("Different aspect ratios should produce different matrices")
	}
	if matrix2 == matrix3 {
		t.Error("Different aspect ratios should produce different matrices")
	}
}

func TestViewSystem_Isolation(t *testing.T) {
	// Test that multiple ViewSystem instances are independent
	vs1 := NewViewSystemForTesting(800, 600)  // 4:3 aspect ratio
	vs2 := NewViewSystemForTesting(1200, 600) // 2:1 aspect ratio

	// Different aspect ratios should produce different matrices
	matrix1 := vs1.GetPerspectiveMatrix(60.0)
	matrix2 := vs2.GetPerspectiveMatrix(60.0)

	if matrix1 == matrix2 {
		t.Error("Different aspect ratios should produce different matrices")
	}

	// Camera positions should be independent
	vs1.Camera.CameraFrom = mgl32.Vec3{1, 2, 3}
	vs2.Camera.CameraFrom = mgl32.Vec3{4, 5, 6}

	if vs1.Camera.CameraFrom == vs2.Camera.CameraFrom {
		t.Error("Camera positions should be independent between ViewSystem instances")
	}
}

func BenchmarkViewSystem_GetPerspectiveMatrix(b *testing.B) {
	vs := NewViewSystemForTesting(800, 600)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = vs.GetPerspectiveMatrix(60.0)
	}
}

func BenchmarkViewSystem_UpdateMatrices(b *testing.B) {
	vs := NewViewSystemForTesting(800, 600)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vs.UpdateMatrices()
	}
}
