package render

import (
	"github.com/go-gl/mathgl/mgl32"
)

type SceneMD1Entity struct {
	TextureId          uint32     // texture id in OpenGL
	VertexBuffer       []float32  // 3 elements for x,y,z, 2 elements for texture u,v, and 3 elements for normal x,y,z
	ModelPosition      mgl32.Vec3 // Position in world space
	RotationAngle      float32
	VertexArrayObject  uint32
	VertexBufferObject uint32
}

func (r *RenderDef) RenderStaticEntity(entity SceneMD1Entity, renderType int32) {
	// Calculate model matrix
	modelMatrix := mgl32.Ident4()
	modelMatrix = modelMatrix.Mul4(mgl32.Translate3D(entity.ModelPosition.X(), entity.ModelPosition.Y(), entity.ModelPosition.Z()))
	modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(float32(entity.RotationAngle))))

	// Create render config for 3D entity (position + texture + normal)
	config := r.Renderer.Create3DEntityConfig(
		entity.VertexArrayObject,
		entity.VertexBufferObject,
		entity.VertexBuffer,
		entity.TextureId,
		&modelMatrix,
		renderType,
	)

	// Render the entity
	r.Renderer.RenderEntity(config)
}
