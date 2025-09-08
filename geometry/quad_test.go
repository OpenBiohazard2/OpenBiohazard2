package geometry

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNewQuad(t *testing.T) {
	// Test basic quad creation
	corners := [4]mgl32.Vec3{
		{0, 0, 0},
		{1, 0, 0},
		{1, 0, 1},
		{0, 0, 1},
	}

	quad := NewQuad(corners)

	// Check vertices are stored correctly
	if quad.Vertices != corners {
		t.Errorf("Expected vertices %v, got %v", corners, quad.Vertices)
	}

	// Check vertex buffer has correct length (6 vertices * 3 components = 18 floats)
	expectedLength := 6 * 3 // 6 vertices, 3 components each
	if len(quad.VertexBuffer) != expectedLength {
		t.Errorf("Expected vertex buffer length %d, got %d", expectedLength, len(quad.VertexBuffer))
	}

	// Check that vertex buffer contains the expected vertices in correct order
	// First triangle: v0, v1, v2
	expectedFirstTriangle := []float32{0, 0, 0, 1, 0, 0, 1, 0, 1}
	for i, expected := range expectedFirstTriangle {
		if quad.VertexBuffer[i] != expected {
			t.Errorf("First triangle vertex %d: expected %f, got %f", i, expected, quad.VertexBuffer[i])
		}
	}

	// Second triangle: v0, v3, v2
	expectedSecondTriangle := []float32{0, 0, 0, 0, 0, 1, 1, 0, 1}
	startIdx := 9 // Start of second triangle
	for i, expected := range expectedSecondTriangle {
		if quad.VertexBuffer[startIdx+i] != expected {
			t.Errorf("Second triangle vertex %d: expected %f, got %f", i, expected, quad.VertexBuffer[startIdx+i])
		}
	}
}

func TestNewQuadFourPoints(t *testing.T) {
	// Test quad creation from xz pairs
	xzPairs := [4][]float32{
		{0, 0},
		{1, 0},
		{1, 1},
		{0, 1},
	}

	quad := NewQuadFourPoints(xzPairs)

	// Check that vertices have y=0 and correct x,z values
	expectedCorners := [4]mgl32.Vec3{
		{0, 0, 0},
		{1, 0, 0},
		{1, 0, 1},
		{0, 0, 1},
	}

	for i, expected := range expectedCorners {
		if quad.Vertices[i] != expected {
			t.Errorf("Vertex %d: expected %v, got %v", i, expected, quad.Vertices[i])
		}
	}
}

func TestNewRectangle(t *testing.T) {
	// Test rectangle creation
	rect := NewRectangle(2, 3, 4, 5)

	expectedCorners := [4]mgl32.Vec3{
		{2, 0, 3}, // x, y, z
		{2, 0, 8}, // x, y, z + depth
		{6, 0, 8}, // x + width, y, z + depth
		{6, 0, 3}, // x + width, y, z
	}

	for i, expected := range expectedCorners {
		if rect.Vertices[i] != expected {
			t.Errorf("Rectangle corner %d: expected %v, got %v", i, expected, rect.Vertices[i])
		}
	}
}

func TestNewTexturedRectangle(t *testing.T) {
	// Test textured rectangle creation
	vertices := [4][]float32{
		{0, 0, 0},
		{1, 0, 0},
		{1, 0, 1},
		{0, 0, 1},
	}
	uvs := [4][]float32{
		{0, 0},
		{1, 0},
		{1, 1},
		{0, 1},
	}

	quad := NewTexturedRectangle(vertices, uvs)

	// Check that vertex buffer has correct length
	// 6 vertices * (3 position + 2 UV) = 30 floats
	expectedLength := 6 * 5
	if len(quad.VertexBuffer) != expectedLength {
		t.Errorf("Expected vertex buffer length %d, got %d", expectedLength, len(quad.VertexBuffer))
	}

	// Check first triangle: v0 + uv0, v1 + uv1, v2 + uv2
	expectedFirstTriangle := []float32{
		0, 0, 0, 0, 0, // v0 + uv0
		1, 0, 0, 1, 0, // v1 + uv1
		1, 0, 1, 1, 1, // v2 + uv2
	}
	for i, expected := range expectedFirstTriangle {
		if quad.VertexBuffer[i] != expected {
			t.Errorf("First triangle vertex %d: expected %f, got %f", i, expected, quad.VertexBuffer[i])
		}
	}
}

func TestNewQuadMD1(t *testing.T) {
	// Test MD1 quad creation
	vertices := [4][]float32{
		{0, 0, 0},
		{1, 0, 0},
		{1, 0, 1},
		{0, 0, 1},
	}
	uvs := [4][]float32{
		{0, 0},
		{1, 0},
		{1, 1},
		{0, 1},
	}
	normals := [4][]float32{
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
	}

	quad := NewQuadMD1(vertices, uvs, normals)

	// Check that vertex buffer has correct length
	// 6 vertices * (3 position + 2 UV + 3 normal) = 48 floats
	expectedLength := 6 * 8
	if len(quad.VertexBuffer) != expectedLength {
		t.Errorf("Expected vertex buffer length %d, got %d", expectedLength, len(quad.VertexBuffer))
	}

	// Check first triangle: v0 + uv0 + normal0, v1 + uv1 + normal1, v3 + uv3 + normal3
	// MD1 uses different vertex order: v0, v1, v3
	expectedFirstTriangle := []float32{
		0, 0, 0, 0, 0, 0, 1, 0, // v0 + uv0 + normal0
		1, 0, 0, 1, 0, 0, 1, 0, // v1 + uv1 + normal1
		0, 0, 1, 0, 1, 0, 1, 0, // v3 + uv3 + normal3
	}
	for i, expected := range expectedFirstTriangle {
		if quad.VertexBuffer[i] != expected {
			t.Errorf("First triangle vertex %d: expected %f, got %f", i, expected, quad.VertexBuffer[i])
		}
	}
}

func TestNewQuadEdgeCases(t *testing.T) {
	// Test with negative coordinates
	corners := [4]mgl32.Vec3{
		{-1, -1, -1},
		{1, -1, -1},
		{1, -1, 1},
		{-1, -1, 1},
	}

	quad := NewQuad(corners)
	if quad.Vertices != corners {
		t.Errorf("Expected vertices %v, got %v", corners, quad.Vertices)
	}

	// Test with zero-size quad
	zeroCorners := [4]mgl32.Vec3{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}

	zeroQuad := NewQuad(zeroCorners)
	expectedLength := 6 * 3
	if len(zeroQuad.VertexBuffer) != expectedLength {
		t.Errorf("Expected vertex buffer length %d, got %d", expectedLength, len(zeroQuad.VertexBuffer))
	}
}
