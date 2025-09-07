package render

import (
	"github.com/OpenBiohazard2/OpenBiohazard2/geometry"
)

func NewBackgroundImageEntity() *SceneEntity {
	backgroundImageEntity := NewSceneEntity()

	// The background image is a rectangle that covers the entire screen
	// It should be drawn in the back
	z := float32(0.999)

	vertices := [4][]float32{
		{-1.0, 1.0, z},
		{-1.0, -1.0, z},
		{1.0, -1.0, z},
		{1.0, 1.0, z},
	}
	uvs := [4][]float32{
		{0.0, 0.0},
		{0.0, 1.0},
		{1.0, 1.0},
		{1.0, 0.0},
	}
	rect := geometry.NewTexturedRectangle(vertices, uvs)
	backgroundImageEntity.SetMesh(rect.VertexBuffer)
	return backgroundImageEntity
}
