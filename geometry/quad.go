package geometry

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Quad struct {
	VertexBuffer []float32
}

func NewQuad(corners [4]mgl32.Vec3) *Quad {
	vertices := make([][]float32, 4)
	for i := 0; i < 4; i++ {
		vertices[i] = []float32{corners[i].X(), corners[i].Y(), corners[i].Z()}
	}

	vertexBuffer := make([]float32, 0)
	// v0, v1, v2
	tri1Indices := [3]int{0, 1, 2}
	for _, index := range tri1Indices {
		vertexBuffer = append(vertexBuffer, vertices[index]...)
	}
	// v0, v3, v2
	tri2Indices := [3]int{0, 3, 2}
	for _, index := range tri2Indices {
		vertexBuffer = append(vertexBuffer, vertices[index]...)
	}
	return &Quad{
		VertexBuffer: vertexBuffer,
	}
}

func NewTexturedRectangle(vertices [4][]float32, uvs [4][]float32) *Quad {
	vertexBuffer := make([]float32, 0)

	// v0, v1, v2
	tri1Indices := [3]int{0, 1, 2}
	for _, index := range tri1Indices {
		vertexBuffer = append(vertexBuffer, vertices[index]...)
		vertexBuffer = append(vertexBuffer, uvs[index]...)
	}
	// v0, v3, v2
	tri2Indices := [3]int{0, 3, 2}
	for _, index := range tri2Indices {
		vertexBuffer = append(vertexBuffer, vertices[index]...)
		vertexBuffer = append(vertexBuffer, uvs[index]...)
	}

	return &Quad{
		VertexBuffer: vertexBuffer,
	}
}

func NewQuadMD1(vertices [4][]float32, uvs [4][]float32, normals [4][]float32) *Quad {
	vertexBuffer := make([]float32, 0)

	// MD1 vertex order is different (v0, v1, v3, v2)
	// Vertex order of other quads is (v0, v1, v2, v3)
	// v0, v1, v3
	tri1Indices := [3]int{0, 1, 3}
	for _, index := range tri1Indices {
		vertexBuffer = append(vertexBuffer, vertices[index]...)
		vertexBuffer = append(vertexBuffer, uvs[index]...)
		vertexBuffer = append(vertexBuffer, normals[index]...)
	}
	// v0, v2, v3
	tri2Indices := [3]int{0, 2, 3}
	for _, index := range tri2Indices {
		vertexBuffer = append(vertexBuffer, vertices[index]...)
		vertexBuffer = append(vertexBuffer, uvs[index]...)
		vertexBuffer = append(vertexBuffer, normals[index]...)
	}

	return &Quad{
		VertexBuffer: vertexBuffer,
	}
}
