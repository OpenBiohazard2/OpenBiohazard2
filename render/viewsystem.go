package render

import (
	"github.com/go-gl/mathgl/mgl32"
)

// ViewSystem manages camera, view matrices, and view calculations
type ViewSystem struct {
	// Core camera
	Camera *Camera

	// Matrices (computed from camera + window)
	ProjectionMatrix mgl32.Mat4
	ViewMatrix       mgl32.Mat4

	// Window dimensions (needed for projection)
	WindowWidth  int
	WindowHeight int
}

// NewViewSystem creates a new view system with default camera
func NewViewSystem(windowWidth, windowHeight int) *ViewSystem {
	cameraUp := mgl32.Vec3{0, -1, 0}
	camera := NewCamera(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 0, 0}, cameraUp, DEFAULT_FOV_DEGREES)

	return &ViewSystem{
		Camera:           camera,
		ProjectionMatrix: mgl32.Perspective(mgl32.DegToRad(DEFAULT_FOV_DEGREES), float32(ASPECT_RATIO), NEAR_PLANE, FAR_PLANE),
		ViewMatrix:       mgl32.Ident4(),
		WindowWidth:      windowWidth,
		WindowHeight:     windowHeight,
	}
}

// NewViewSystemForTesting creates a ViewSystem without OpenGL dependencies for testing
func NewViewSystemForTesting(windowWidth, windowHeight int) *ViewSystem {
	cameraUp := mgl32.Vec3{0, -1, 0}
	camera := NewCamera(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 0, 0}, cameraUp, DEFAULT_FOV_DEGREES)

	return &ViewSystem{
		Camera:           camera,
		ProjectionMatrix: mgl32.Perspective(mgl32.DegToRad(DEFAULT_FOV_DEGREES), float32(ASPECT_RATIO), NEAR_PLANE, FAR_PLANE),
		ViewMatrix:       mgl32.Ident4(),
		WindowWidth:      windowWidth,
		WindowHeight:     windowHeight,
	}
}

// UpdateMatrices recalculates projection and view matrices
func (vs *ViewSystem) UpdateMatrices() {
	vs.ProjectionMatrix = vs.GetPerspectiveMatrix(vs.Camera.CameraFov)
	vs.ViewMatrix = vs.Camera.BuildViewMatrix()
}

// GetViewMatrix returns the current view matrix
func (vs *ViewSystem) GetViewMatrix() mgl32.Mat4 {
	return vs.ViewMatrix
}

// GetProjectionMatrix returns the current projection matrix
func (vs *ViewSystem) GetProjectionMatrix() mgl32.Mat4 {
	return vs.ProjectionMatrix
}

// GetPerspectiveMatrix calculates perspective matrix for given FOV
func (vs *ViewSystem) GetPerspectiveMatrix(fovDegrees float32) mgl32.Mat4 {
	ratio := float64(vs.WindowWidth) / float64(vs.WindowHeight)
	return mgl32.Perspective(mgl32.DegToRad(fovDegrees), float32(ratio), NEAR_PLANE, FAR_PLANE)
}
