package render

import (
	"github.com/OpenBiohazard2/OpenBiohazard2/geometry"
	"github.com/OpenBiohazard2/OpenBiohazard2/resource"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Entity2D struct {
	TextureId          uint32    // texture id in OpenGL
	VertexBuffer       []float32 // 3 elements for x,y,z and 2 elements for texture u,v
	VertexArrayObject  uint32
	VertexBufferObject uint32
}

func NewEntity2D() *Entity2D {
	var vao uint32
	gl.GenVertexArrays(1, &vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)

	var texId uint32
	gl.GenTextures(1, &texId)

	return &Entity2D{
		TextureId:          texId,
		VertexBuffer:       []float32{},
		VertexArrayObject:  vao,
		VertexBufferObject: vbo,
	}
}

func (r *RenderDef) RenderEntity2D(entity *Entity2D, renderType int32) {
	// Skip if entity is nil or has no vertex data
	if entity == nil || len(entity.VertexBuffer) == 0 {
		return
	}

	// Create render config for 2D entity (position + texture)
	config := r.Renderer.Create2DEntityConfig(
		entity.VertexArrayObject,
		entity.VertexBufferObject,
		entity.VertexBuffer,
		entity.TextureId,
		nil, // No model matrix for scene entities
		renderType,
	)

	// Render the entity
	r.Renderer.RenderEntity(config)
}

func (entity *Entity2D) DeleteEntity2D() {
	gl.DeleteVertexArrays(1, &entity.VertexArrayObject)
	gl.DeleteBuffers(1, &entity.VertexBufferObject)
}

func (entity *Entity2D) SetTexture(imagePixels []uint16, imageWidth int32, imageHeight int32) {
	UpdateTexture(entity.TextureId, imagePixels, imageWidth, imageHeight)
}

func (entity *Entity2D) SetMesh(vertexBuffer []float32) {
	entity.VertexBuffer = vertexBuffer
}

func (entity *Entity2D) UpdateTextureFromImage(newImage *resource.Image16Bit) {
	UpdateTexture(entity.TextureId, newImage.GetPixelsForRendering(), int32(newImage.GetWidth()), int32(newImage.GetHeight()))
}

func NewBackgroundImageEntity() *Entity2D {
	backgroundImageEntity := NewEntity2D()

	// The background image is a rectangle that covers the entire screen
	// It should be drawn in the back
	rect := geometry.NewFullScreenQuad(geometry.BACKGROUND_DEPTH)
	backgroundImageEntity.SetMesh(rect.VertexBuffer)
	return backgroundImageEntity
}

func (renderDef *RenderDef) RenderSolidVideoBuffer() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Activate shader
	renderDef.ShaderSystem.Use()

	// Use cached uniform location for better performance
	renderDef.ShaderSystem.SetGameState(RENDER_GAME_STATE_BACKGROUND_SOLID)

	renderDef.RenderEntity2D(renderDef.VideoBuffer, RENDER_GAME_STATE_BACKGROUND_SOLID)
}

func (renderDef *RenderDef) RenderTransparentVideoBuffer() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Activate shader
	renderDef.ShaderSystem.Use()

	// Use cached uniform location for better performance
	renderDef.ShaderSystem.SetGameState(RENDER_GAME_STATE_BACKGROUND_TRANSPARENT)

	renderDef.RenderEntity2D(renderDef.VideoBuffer, RENDER_GAME_STATE_BACKGROUND_TRANSPARENT)
}
