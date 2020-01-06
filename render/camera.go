package render

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	CameraFrom mgl32.Vec3
	CameraTo   mgl32.Vec3
	CameraUp   mgl32.Vec3
}

func NewCamera(cameraFrom mgl32.Vec3, cameraTo mgl32.Vec3, cameraUp mgl32.Vec3) *Camera {
	return &Camera{
		CameraFrom: cameraFrom,
		CameraTo:   cameraTo,
		CameraUp:   cameraUp,
	}
}

func (c *Camera) GetViewMatrix() mgl32.Mat4 {
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
