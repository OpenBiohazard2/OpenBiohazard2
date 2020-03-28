package geometry

import (
  	"github.com/samuelyuan/openbiohazard2/fileio"
)

func NewMD1Geometry(meshData *fileio.MD1Output, textureData *fileio.TIMOutput) []float32 {
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

			tri := NewTriangleNormals([3][]float32{vertex0, vertex1, vertex2},
				[3][]float32{uv0, uv1, uv2},
				[3][]float32{normal0, normal1, normal2})
			vertexBuffer = append(vertexBuffer, tri.VertexBuffer...)
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

			quad := NewQuadMD1([4][]float32{vertex0, vertex1, vertex2, vertex3},
				[4][]float32{uv0, uv1, uv2, uv3},
				[4][]float32{normal0, normal1, normal2, normal3})
			vertexBuffer = append(vertexBuffer, quad.VertexBuffer...)
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
