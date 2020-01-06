package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Player struct {
	Position      mgl32.Vec3
	RotationAngle float32
	PoseNumber    int
}

// Position is in world space
// Rotation angle is in degrees
func NewPlayer(initialPosition mgl32.Vec3, initialRotationAngle float32) *Player {
	return &Player{
		Position:      initialPosition,
		RotationAngle: initialRotationAngle,
		PoseNumber:    -1,
	}
}

func (p *Player) GetModelMatrix() mgl32.Mat4 {
	modelMatrix := mgl32.Ident4()
	modelMatrix = modelMatrix.Mul4(mgl32.Translate3D(p.Position.X(), p.Position.Y(), p.Position.Z()))
	modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(float32(p.RotationAngle))))
	return modelMatrix
}
