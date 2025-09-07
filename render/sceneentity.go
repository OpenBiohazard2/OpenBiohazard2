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
	// Skip
	if entity == nil {
		return
	}

	vertexBuffer := entity.VertexBuffer
	if len(vertexBuffer) == 0 {
		return
	}

	// Use cached uniform location for better performance
	gl.Uniform1i(r.UniformLocations.RenderType, renderType)

	floatSize := 4

	// 3 floats for vertex, 2 floats for texture UV
	stride := int32(5 * floatSize)

	vao := entity.VertexArrayObject
	gl.BindVertexArray(vao)

	vbo := entity.VertexBufferObject
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexBuffer)*floatSize, gl.Ptr(vertexBuffer), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Texture
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, stride, gl.PtrOffset(3*floatSize))
	gl.EnableVertexAttribArray(1)

	// Use cached uniform location for better performance
	gl.Uniform1i(r.UniformLocations.Diffuse, 0)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, entity.TextureId)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertexBuffer)/5))

	// Cleanup
	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
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
