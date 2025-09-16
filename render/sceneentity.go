package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type SceneEntity struct {
	TextureId          uint32    // texture id in OpenGL
	VertexBuffer       []float32 // 3 elements for x,y,z and 2 elements for texture u,v
	VertexArrayObject  uint32
	VertexBufferObject uint32
}

func NewSceneEntity() *SceneEntity {
	var vao uint32
	gl.GenVertexArrays(1, &vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)

	var texId uint32
	gl.GenTextures(1, &texId)

	return &SceneEntity{
		TextureId:          texId,
		VertexBuffer:       []float32{},
		VertexArrayObject:  vao,
		VertexBufferObject: vbo,
	}
}

func (r *RenderDef) RenderSceneEntity(entity *SceneEntity, renderType int32) {
	// Skip if entity is nil or has no vertex data
	if entity == nil || len(entity.VertexBuffer) == 0 {
		return
	}

	// Create renderer
	renderer := NewOpenGLRenderer(r.ShaderSystem.GetUniformLocations())

	// Create render config for 2D entity (position + texture)
	config := renderer.Create2DEntityConfig(
		entity.VertexArrayObject,
		entity.VertexBufferObject,
		entity.VertexBuffer,
		entity.TextureId,
		nil, // No model matrix for scene entities
		renderType,
	)

	// Render the entity
	renderer.RenderEntity(config)
}

func (entity *SceneEntity) DeleteSceneEntity() {
	gl.DeleteVertexArrays(1, &entity.VertexArrayObject)
	gl.DeleteBuffers(1, &entity.VertexBufferObject)
}

func (entity *SceneEntity) SetTexture(imagePixels []uint16, imageWidth int32, imageHeight int32) {
	UpdateTexture(entity.TextureId, imagePixels, imageWidth, imageHeight)
}

func (entity *SceneEntity) SetMesh(vertexBuffer []float32) {
	entity.VertexBuffer = vertexBuffer
}

func (entity *SceneEntity) UpdateSurface(newImagePixels []uint16, imageWidth int32, imageHeight int32) {
	UpdateTexture(entity.TextureId, newImagePixels, imageWidth, imageHeight)
}
