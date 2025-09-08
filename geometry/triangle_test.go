package geometry

import (
	"testing"
)

func TestNewTriangleNormals(t *testing.T) {
	vertices := [3][]float32{
		{0, 0, 0}, // v0
		{1, 0, 0}, // v1
		{0, 1, 0}, // v2
	}
	uvs := [3][]float32{
		{0, 0}, // uv0
		{1, 0}, // uv1
		{0, 1}, // uv2
	}
	normals := [3][]float32{
		{0, 0, 1}, // n0
		{0, 0, 1}, // n1
		{0, 0, 1}, // n2
	}

	triangle := NewTriangleNormals(vertices, uvs, normals)

	// Test vertex buffer length: 3 vertices * (3 pos + 2 uv + 3 normal) = 24 floats
	expectedLength := 3 * 8
	if len(triangle.VertexBuffer) != expectedLength {
		t.Errorf("Expected vertex buffer length %d, got %d", expectedLength, len(triangle.VertexBuffer))
	}

	// Test first vertex (v0): position + uv + normal
	expectedFirstVertex := []float32{0, 0, 0, 0, 0, 0, 0, 1} // pos + uv + normal
	for i, expected := range expectedFirstVertex {
		if triangle.VertexBuffer[i] != expected {
			t.Errorf("First vertex[%d]: expected %f, got %f", i, expected, triangle.VertexBuffer[i])
		}
	}

	// Test second vertex (v1): position + uv + normal
	expectedSecondVertex := []float32{1, 0, 0, 1, 0, 0, 0, 1} // pos + uv + normal
	startIndex := 8
	for i, expected := range expectedSecondVertex {
		actual := triangle.VertexBuffer[startIndex+i]
		if actual != expected {
			t.Errorf("Second vertex[%d]: expected %f, got %f", i, expected, actual)
		}
	}

	// Test third vertex (v2): position + uv + normal
	expectedThirdVertex := []float32{0, 1, 0, 0, 1, 0, 0, 1} // pos + uv + normal
	startIndex = 16
	for i, expected := range expectedThirdVertex {
		actual := triangle.VertexBuffer[startIndex+i]
		if actual != expected {
			t.Errorf("Third vertex[%d]: expected %f, got %f", i, expected, actual)
		}
	}
}

func TestTriangleWithDifferentValues(t *testing.T) {
	vertices := [3][]float32{
		{10, 20, 30}, // v0
		{11, 21, 31}, // v1
		{12, 22, 32}, // v2
	}
	uvs := [3][]float32{
		{0.1, 0.2}, // uv0
		{0.3, 0.4}, // uv1
		{0.5, 0.6}, // uv2
	}
	normals := [3][]float32{
		{1, 0, 0}, // n0
		{0, 1, 0}, // n1
		{0, 0, 1}, // n2
	}

	triangle := NewTriangleNormals(vertices, uvs, normals)

	// Test that all values are correctly interleaved
	expectedValues := []float32{
		// Vertex 0: pos + uv + normal
		10, 20, 30, 0.1, 0.2, 1, 0, 0,
		// Vertex 1: pos + uv + normal
		11, 21, 31, 0.3, 0.4, 0, 1, 0,
		// Vertex 2: pos + uv + normal
		12, 22, 32, 0.5, 0.6, 0, 0, 1,
	}

	for i, expected := range expectedValues {
		actual := triangle.VertexBuffer[i]
		if actual != expected {
			t.Errorf("VertexBuffer[%d]: expected %f, got %f", i, expected, actual)
		}
	}
}

func TestTriangleEmptyInput(t *testing.T) {
	// Test with zero values
	vertices := [3][]float32{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}
	uvs := [3][]float32{
		{0, 0},
		{0, 0},
		{0, 0},
	}
	normals := [3][]float32{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}

	triangle := NewTriangleNormals(vertices, uvs, normals)

	// Should still create a valid triangle with all zeros
	expectedLength := 3 * 8
	if len(triangle.VertexBuffer) != expectedLength {
		t.Errorf("Expected vertex buffer length %d, got %d", expectedLength, len(triangle.VertexBuffer))
	}

	// All values should be zero
	for i, value := range triangle.VertexBuffer {
		if value != 0 {
			t.Errorf("Expected all zeros, but VertexBuffer[%d] = %f", i, value)
		}
	}
}
