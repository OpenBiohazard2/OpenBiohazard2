package render

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	CameraFrom mgl32.Vec3
	CameraTo   mgl32.Vec3
	CameraUp   mgl32.Vec3
	CameraFov  float32
}

func NewCamera(cameraFrom mgl32.Vec3, cameraTo mgl32.Vec3, cameraUp mgl32.Vec3, cameraFov float32) *Camera {
	return &Camera{
		CameraFrom: cameraFrom,
		CameraTo:   cameraTo,
		CameraUp:   cameraUp,
		CameraFov:  cameraFov,
	}
}

func (c *Camera) Update(cameraFrom mgl32.Vec3, cameraTo mgl32.Vec3, cameraFov float32) {
	c.CameraFrom = cameraFrom
	c.CameraTo = cameraTo
	c.CameraFov = cameraFov
}

func (c *Camera) BuildViewMatrix() mgl32.Mat4 {
	cameraFrom := c.CameraFrom
	cameraTo := c.CameraTo
	cameraUp := c.CameraUp
	return mgl32.LookAt(
		cameraFrom.X(), cameraFrom.Y(), cameraFrom.Z(),
		cameraTo.X(), cameraTo.Y(), cameraTo.Z(),
		cameraUp.X(), cameraUp.Y(), cameraUp.Z())
}

func (c *Camera) GetDirection() mgl32.Vec3 {
	return c.CameraTo.Sub(c.CameraFrom).Normalize()
}

// NormalizeMaskDepth normalizes the z coordinate to be between 0 and 1
// 0 is closer to the camera, 1 is farther from the camera
func (c *Camera) NormalizeMaskDepth(depth float32, projectionMatrix, viewMatrix mgl32.Mat4) float32 {
	cameraDir := c.GetDirection().Normalize()
	cameraFrom := c.CameraFrom
	transformMatrix := projectionMatrix.Mul4(viewMatrix)

	// Actual distance from camera is 32 * depth
	projectedPosition := cameraFrom.Add(cameraDir.Mul(depth * float32(32.0)))

	// Get its z coordinate on the screen
	renderPosition := transformMatrix.Mul4x1(mgl32.Vec4{projectedPosition.X(), projectedPosition.Y(), projectedPosition.Z(), 1})
	renderPosition = renderPosition.Mul(1 / renderPosition.W())
	return renderPosition.Z()
}
