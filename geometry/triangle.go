package geometry

type Triangle struct {
	VertexBuffer []float32
}

func NewTriangleNormals(vertices [3][]float32, uvs [3][]float32, normals [3][]float32) *Triangle {
	vertexBuffer := make([]float32, 0)

	// v0, v1, v2
	tri1Indices := [3]int{0, 1, 2}
	for _, index := range tri1Indices {
		vertexBuffer = append(vertexBuffer, vertices[index]...)
		vertexBuffer = append(vertexBuffer, uvs[index]...)
		vertexBuffer = append(vertexBuffer, normals[index]...)
	}

	return &Triangle{
		VertexBuffer: vertexBuffer,
	}
}
