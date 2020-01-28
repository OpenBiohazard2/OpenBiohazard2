package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	RENDER_TYPE_DEBUG = -1
)

func RenderCameraSwitches(programShader uint32, cameraSwitchDebugEntity *DebugEntity) {
	renderTypeUniform := gl.GetUniformLocation(programShader, gl.Str("renderType\x00"))
	gl.Uniform1i(renderTypeUniform, RENDER_TYPE_DEBUG)

	RenderDebugEntities(programShader, []*DebugEntity{cameraSwitchDebugEntity})
}

func RenderDebugEntities(programShader uint32, debugEntities []*DebugEntity) {
	renderTypeUniform := gl.GetUniformLocation(programShader, gl.Str("renderType\x00"))
	gl.Uniform1i(renderTypeUniform, RENDER_TYPE_DEBUG)

	floatSize := 4

	for _, debugEntity := range debugEntities {
		entityVertexBuffer := debugEntity.VertexBuffer
		if len(entityVertexBuffer) == 0 {
			continue
		}

		// 3 floats for vertex
		stride := int32(3 * floatSize)

		vao := debugEntity.VertexArrayObject
		gl.BindVertexArray(vao)

		vbo := debugEntity.VertexBufferObject
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(entityVertexBuffer)*floatSize, gl.Ptr(entityVertexBuffer), gl.STATIC_DRAW)

		// Position attribute
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, gl.PtrOffset(0))
		gl.EnableVertexAttribArray(0)

		diffuseUniform := gl.GetUniformLocation(programShader, gl.Str("diffuse\x00"))
		gl.Uniform1i(diffuseUniform, 0)

		debugColorLoc := gl.GetUniformLocation(programShader, gl.Str("debugColor\x00"))
		color := debugEntity.Color
		gl.Uniform4f(debugColorLoc, color[0], color[1], color[2], color[3])

		// Draw triangles
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(entityVertexBuffer)/3))

		// Cleanup
		gl.DisableVertexAttribArray(0)
	}
}
