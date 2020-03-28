package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type SceneMD1Entity struct {
	TextureId          uint32     // texture id in OpenGL
	VertexBuffer       []float32  // 3 elements for x,y,z, 2 elements for texture u,v, and 3 elements for normal x,y,z
	ModelPosition      mgl32.Vec3 // Position in world space
	VertexArrayObject  uint32
	VertexBufferObject uint32
}

func (r *RenderDef) RenderMD1Entity(entity SceneMD1Entity, renderType int32) {
	vertexBuffer := entity.VertexBuffer
	textureId := entity.TextureId
	modelPosition := entity.ModelPosition

	modelMatrix := mgl32.Ident4()
	modelMatrix = modelMatrix.Mul4(mgl32.Translate3D(modelPosition.X(), modelPosition.Y(), modelPosition.Z()))

	if len(vertexBuffer) == 0 {
		return
	}

	programShader := r.ProgramShader
	renderTypeUniform := gl.GetUniformLocation(programShader, gl.Str("renderType\x00"))
	gl.Uniform1i(renderTypeUniform, renderType)

	modelLoc := gl.GetUniformLocation(programShader, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelLoc, 1, false, &modelMatrix[0])

	floatSize := 4

	// 3 floats for vertex, 2 floats for texture UV, 3 float for normals
	stride := int32(8 * floatSize)

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

	// Normal
	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, stride, gl.PtrOffset(5*floatSize))
	gl.EnableVertexAttribArray(2)

	diffuseUniform := gl.GetUniformLocation(programShader, gl.Str("diffuse\x00"))
	gl.Uniform1i(diffuseUniform, 0)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, textureId)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertexBuffer)/8))

	// Cleanup
	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
	gl.DisableVertexAttribArray(2)
}
