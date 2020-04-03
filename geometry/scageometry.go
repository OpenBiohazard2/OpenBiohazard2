package geometry

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/samuelyuan/openbiohazard2/fileio"
)

func NewSlopedRectangle(entity fileio.CollisionEntity) *Quad {
	// Types 0 and 1 starts from x-axis
	// Types 2 and 3 starts from z-axis
	switch entity.SlopeType {
	case 0:
		vertices := [4]mgl32.Vec3{
			mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)},
			mgl32.Vec3{float32(entity.X), 0, float32(entity.Z + entity.Density)},
			mgl32.Vec3{float32(entity.X + entity.Width), float32(entity.SlopeHeight), float32(entity.Z + entity.Density)},
			mgl32.Vec3{float32(entity.X + entity.Width), float32(entity.SlopeHeight), float32(entity.Z)},
		}
		return NewQuad(vertices)
	case 1:
		vertices := [4]mgl32.Vec3{
			mgl32.Vec3{float32(entity.X), float32(entity.SlopeHeight), float32(entity.Z)},
			mgl32.Vec3{float32(entity.X), float32(entity.SlopeHeight), float32(entity.Z + entity.Density)},
			mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z + entity.Density)},
			mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z)},
		}
		return NewQuad(vertices)
	case 2:
		vertices := [4]mgl32.Vec3{
			mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)},
			mgl32.Vec3{float32(entity.X), float32(entity.SlopeHeight), float32(entity.Z + entity.Density)},
			mgl32.Vec3{float32(entity.X + entity.Width), float32(entity.SlopeHeight), float32(entity.Z + entity.Density)},
			mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z)},
		}
		return NewQuad(vertices)
	case 3:
		vertices := [4]mgl32.Vec3{
			mgl32.Vec3{float32(entity.X), float32(entity.SlopeHeight), float32(entity.Z)},
			mgl32.Vec3{float32(entity.X), 0, float32(entity.Z + entity.Density)},
			mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z + entity.Density)},
			mgl32.Vec3{float32(entity.X + entity.Width), float32(entity.SlopeHeight), float32(entity.Z)},
		}
		return NewQuad(vertices)
	}

	return nil
}
