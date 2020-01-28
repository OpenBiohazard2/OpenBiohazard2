package render

import (
	"../fileio"
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

func BuildEntityComponentVertices(meshData *fileio.MD1Output, textureData *fileio.TIMOutput) []float32 {
	vertexBuffer := make([]float32, 0)

	for _, entityModel := range meshData.Components {
		// Triangles
		for j := 0; j < len(entityModel.TriangleIndices); j++ {
			triangleIndex := entityModel.TriangleIndices[j]
			textureInfo := entityModel.TriangleTextures[j]

			vertex0 := buildModelVertex(entityModel.TriangleVertices[triangleIndex.IndexVertex0])
			uv0 := buildTextureUV(float32(textureInfo.U0), float32(textureInfo.V0), textureInfo.Page, textureData)
			normal0 := buildModelNormal(entityModel.TriangleNormals[triangleIndex.IndexNormal0])

			vertex1 := buildModelVertex(entityModel.TriangleVertices[triangleIndex.IndexVertex1])
			uv1 := buildTextureUV(float32(textureInfo.U1), float32(textureInfo.V1), textureInfo.Page, textureData)
			normal1 := buildModelNormal(entityModel.TriangleNormals[triangleIndex.IndexNormal1])

			vertex2 := buildModelVertex(entityModel.TriangleVertices[triangleIndex.IndexVertex2])
			uv2 := buildTextureUV(float32(textureInfo.U2), float32(textureInfo.V2), textureInfo.Page, textureData)
			normal2 := buildModelNormal(entityModel.TriangleNormals[triangleIndex.IndexNormal2])

			// v0, v1, v2
			vertexBuffer = append(vertexBuffer, vertex0...)
			vertexBuffer = append(vertexBuffer, uv0...)
			vertexBuffer = append(vertexBuffer, normal0...)

			vertexBuffer = append(vertexBuffer, vertex1...)
			vertexBuffer = append(vertexBuffer, uv1...)
			vertexBuffer = append(vertexBuffer, normal1...)

			vertexBuffer = append(vertexBuffer, vertex2...)
			vertexBuffer = append(vertexBuffer, uv2...)
			vertexBuffer = append(vertexBuffer, normal2...)
		}

		// Quads
		for j := 0; j < len(entityModel.QuadIndices); j++ {
			quadIndex := entityModel.QuadIndices[j]
			textureInfo := entityModel.QuadTextures[j]

			vertex0 := buildModelVertex(entityModel.QuadVertices[quadIndex.IndexVertex0])
			uv0 := buildTextureUV(float32(textureInfo.U0), float32(textureInfo.V0), textureInfo.Page, textureData)
			normal0 := buildModelNormal(entityModel.QuadNormals[quadIndex.IndexNormal0])

			vertex1 := buildModelVertex(entityModel.QuadVertices[quadIndex.IndexVertex1])
			uv1 := buildTextureUV(float32(textureInfo.U1), float32(textureInfo.V1), textureInfo.Page, textureData)
			normal1 := buildModelNormal(entityModel.QuadNormals[quadIndex.IndexNormal1])

			vertex2 := buildModelVertex(entityModel.QuadVertices[quadIndex.IndexVertex2])
			uv2 := buildTextureUV(float32(textureInfo.U2), float32(textureInfo.V2), textureInfo.Page, textureData)
			normal2 := buildModelNormal(entityModel.QuadNormals[quadIndex.IndexNormal2])

			vertex3 := buildModelVertex(entityModel.QuadVertices[quadIndex.IndexVertex3])
			uv3 := buildTextureUV(float32(textureInfo.U3), float32(textureInfo.V3), textureInfo.Page, textureData)
			normal3 := buildModelNormal(entityModel.QuadNormals[quadIndex.IndexNormal3])

			// v0, v1, v3
			vertexBuffer = append(vertexBuffer, vertex0...)
			vertexBuffer = append(vertexBuffer, uv0...)
			vertexBuffer = append(vertexBuffer, normal0...)

			vertexBuffer = append(vertexBuffer, vertex1...)
			vertexBuffer = append(vertexBuffer, uv1...)
			vertexBuffer = append(vertexBuffer, normal1...)

			vertexBuffer = append(vertexBuffer, vertex3...)
			vertexBuffer = append(vertexBuffer, uv3...)
			vertexBuffer = append(vertexBuffer, normal3...)

			// v0, v2, v3
			vertexBuffer = append(vertexBuffer, vertex0...)
			vertexBuffer = append(vertexBuffer, uv0...)
			vertexBuffer = append(vertexBuffer, normal0...)

			vertexBuffer = append(vertexBuffer, vertex2...)
			vertexBuffer = append(vertexBuffer, uv2...)
			vertexBuffer = append(vertexBuffer, normal2...)

			vertexBuffer = append(vertexBuffer, vertex3...)
			vertexBuffer = append(vertexBuffer, uv3...)
			vertexBuffer = append(vertexBuffer, normal3...)
		}
	}
	return vertexBuffer
}

func buildModelVertex(vertex fileio.MD1Vertex) []float32 {
	return []float32{float32(vertex.X), float32(vertex.Y), float32(vertex.Z)}
}

func buildTextureUV(u float32, v float32, texturePage uint16, textureData *fileio.TIMOutput) []float32 {
	textureOffsetUnit := float32(textureData.ImageWidth) / float32(textureData.NumPalettes)
	textureCoordOffset := textureOffsetUnit * float32(texturePage&3)

	newU := (float32(u) + textureCoordOffset) / float32(textureData.ImageWidth)
	newV := float32(v) / float32(textureData.ImageHeight)
	return []float32{newU, newV}
}

func buildModelNormal(normal fileio.MD1Vertex) []float32 {
	return []float32{float32(normal.X), float32(normal.Y), float32(normal.Z)}
}
