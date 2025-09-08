package geometry

import (
	"testing"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
)

func TestBuildModelVertex(t *testing.T) {
	// Test basic vertex conversion
	vertex := fileio.MD1Vertex{
		X:    100,
		Y:    200,
		Z:    300,
		Zero: 0,
	}

	result := buildModelVertex(vertex)
	expected := []float32{100, 200, 300}

	if len(result) != 3 {
		t.Errorf("Expected 3 components, got %d", len(result))
	}

	for i, expectedVal := range expected {
		if result[i] != expectedVal {
			t.Errorf("Component %d: expected %f, got %f", i, expectedVal, result[i])
		}
	}
}

func TestBuildModelVertexNegativeValues(t *testing.T) {
	// Test with negative values
	vertex := fileio.MD1Vertex{
		X:    -50,
		Y:    -100,
		Z:    -150,
		Zero: 0,
	}

	result := buildModelVertex(vertex)
	expected := []float32{-50, -100, -150}

	for i, expectedVal := range expected {
		if result[i] != expectedVal {
			t.Errorf("Negative component %d: expected %f, got %f", i, expectedVal, result[i])
		}
	}
}

func TestBuildModelNormal(t *testing.T) {
	// Test normal conversion
	normal := fileio.MD1Vertex{
		X:    0,
		Y:    1000, // Normalized to 1.0
		Z:    0,
		Zero: 0,
	}

	result := buildModelNormal(normal)
	expected := []float32{0, 1000, 0}

	if len(result) != 3 {
		t.Errorf("Expected 3 components, got %d", len(result))
	}

	for i, expectedVal := range expected {
		if result[i] != expectedVal {
			t.Errorf("Normal component %d: expected %f, got %f", i, expectedVal, result[i])
		}
	}
}

func TestBuildTextureUV(t *testing.T) {
	// Test texture UV calculation
	textureData := &fileio.TIMOutput{
		ImageWidth:  256,
		ImageHeight: 256,
		NumPalettes: 4,
	}

	u := float32(128)
	v := float32(64)
	texturePage := uint16(2)

	result := buildTextureUV(u, v, texturePage, textureData)

	if len(result) != 2 {
		t.Errorf("Expected 2 components (U, V), got %d", len(result))
	}

	// Calculate expected values
	textureOffsetUnit := float32(textureData.ImageWidth) / float32(textureData.NumPalettes) // 256/4 = 64
	textureCoordOffset := textureOffsetUnit * float32(texturePage&3)                        // 64 * 2 = 128

	expectedU := (u + textureCoordOffset) / float32(textureData.ImageWidth) // (128 + 128) / 256 = 1.0
	expectedV := v / float32(textureData.ImageHeight)                       // 64 / 256 = 0.25

	expected := []float32{expectedU, expectedV}

	for i, expectedVal := range expected {
		if result[i] != expectedVal {
			t.Errorf("UV component %d: expected %f, got %f", i, expectedVal, result[i])
		}
	}
}

func TestBuildTextureUVEdgeCases(t *testing.T) {
	textureData := &fileio.TIMOutput{
		ImageWidth:  128,
		ImageHeight: 128,
		NumPalettes: 2,
	}

	// Test with zero values
	result := buildTextureUV(0, 0, 0, textureData)
	expected := []float32{0, 0}

	for i, expectedVal := range expected {
		if result[i] != expectedVal {
			t.Errorf("Zero UV component %d: expected %f, got %f", i, expectedVal, result[i])
		}
	}

	// Test with maximum values
	result = buildTextureUV(127, 127, 1, textureData)
	textureOffsetUnit := float32(128) / float32(2)         // 64
	textureCoordOffset := textureOffsetUnit * float32(1&3) // 64 * 1 = 64
	expectedU := (127 + textureCoordOffset) / float32(128) // (127 + 64) / 128 = 1.4921875
	expectedV := 127 / float32(128)                        // 0.9921875
	expected = []float32{expectedU, expectedV}

	for i, expectedVal := range expected {
		if result[i] != expectedVal {
			t.Errorf("Max UV component %d: expected %f, got %f", i, expectedVal, result[i])
		}
	}
}

func TestNewMD1GeometryEmpty(t *testing.T) {
	// Test with empty mesh data
	meshData := &fileio.MD1Output{
		Components: []fileio.MD1Object{},
		NumBytes:   0,
	}

	textureData := &fileio.TIMOutput{
		ImageWidth:  256,
		ImageHeight: 256,
		NumPalettes: 4,
	}

	result := NewMD1Geometry(meshData, textureData)

	if result == nil {
		t.Error("Expected empty slice, got nil")
	}

	if len(result) != 0 {
		t.Errorf("Expected empty result, got length %d", len(result))
	}
}

func TestNewMD1GeometryWithTriangles(t *testing.T) {
	// Test with triangle data
	meshData := &fileio.MD1Output{
		Components: []fileio.MD1Object{
			{
				TriangleVertices: []fileio.MD1Vertex{
					{X: 0, Y: 0, Z: 0},
					{X: 100, Y: 0, Z: 0},
					{X: 50, Y: 100, Z: 0},
				},
				TriangleNormals: []fileio.MD1Vertex{
					{X: 0, Y: 0, Z: 1000},
					{X: 0, Y: 0, Z: 1000},
					{X: 0, Y: 0, Z: 1000},
				},
				TriangleIndices: []fileio.MD1TriangleIndex{
					{
						IndexVertex0: 0, IndexNormal0: 0,
						IndexVertex1: 1, IndexNormal1: 1,
						IndexVertex2: 2, IndexNormal2: 2,
					},
				},
				TriangleTextures: []fileio.MD1TriangleTexture{
					{
						U0: 0, V0: 0,
						U1: 100, V1: 0,
						U2: 50, V2: 100,
						Page: 0,
					},
				},
			},
		},
		NumBytes: 100,
	}

	textureData := &fileio.TIMOutput{
		ImageWidth:  256,
		ImageHeight: 256,
		NumPalettes: 4,
	}

	result := NewMD1Geometry(meshData, textureData)

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Each triangle should have 3 vertices * (3 position + 2 UV + 3 normal) = 24 floats
	expectedLength := 1 * 3 * 8 // 1 triangle * 3 vertices * 8 components
	if len(result) != expectedLength {
		t.Errorf("Expected length %d, got %d", expectedLength, len(result))
	}
}

func TestNewMD1GeometryWithQuads(t *testing.T) {
	// Test with quad data
	meshData := &fileio.MD1Output{
		Components: []fileio.MD1Object{
			{
				QuadVertices: []fileio.MD1Vertex{
					{X: 0, Y: 0, Z: 0},
					{X: 100, Y: 0, Z: 0},
					{X: 100, Y: 0, Z: 100},
					{X: 0, Y: 0, Z: 100},
				},
				QuadNormals: []fileio.MD1Vertex{
					{X: 0, Y: 1000, Z: 0},
					{X: 0, Y: 1000, Z: 0},
					{X: 0, Y: 1000, Z: 0},
					{X: 0, Y: 1000, Z: 0},
				},
				QuadIndices: []fileio.MD1QuadIndex{
					{
						IndexVertex0: 0, IndexNormal0: 0,
						IndexVertex1: 1, IndexNormal1: 1,
						IndexVertex2: 2, IndexNormal2: 2,
						IndexVertex3: 3, IndexNormal3: 3,
					},
				},
				QuadTextures: []fileio.MD1QuadTexture{
					{
						U0: 0, V0: 0,
						U1: 100, V1: 0,
						U2: 100, V2: 100,
						U3: 0, V3: 100,
						Page: 0,
					},
				},
			},
		},
		NumBytes: 100,
	}

	textureData := &fileio.TIMOutput{
		ImageWidth:  256,
		ImageHeight: 256,
		NumPalettes: 4,
	}

	result := NewMD1Geometry(meshData, textureData)

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Each quad should have 6 vertices * (3 position + 2 UV + 3 normal) = 48 floats
	expectedLength := 1 * 6 * 8 // 1 quad * 6 vertices * 8 components
	if len(result) != expectedLength {
		t.Errorf("Expected length %d, got %d", expectedLength, len(result))
	}
}

func TestNewMD1GeometryMixed(t *testing.T) {
	// Test with both triangles and quads
	meshData := &fileio.MD1Output{
		Components: []fileio.MD1Object{
			{
				// Triangle data
				TriangleVertices: []fileio.MD1Vertex{
					{X: 0, Y: 0, Z: 0},
					{X: 50, Y: 0, Z: 0},
					{X: 25, Y: 50, Z: 0},
				},
				TriangleNormals: []fileio.MD1Vertex{
					{X: 0, Y: 0, Z: 1000},
					{X: 0, Y: 0, Z: 1000},
					{X: 0, Y: 0, Z: 1000},
				},
				TriangleIndices: []fileio.MD1TriangleIndex{
					{
						IndexVertex0: 0, IndexNormal0: 0,
						IndexVertex1: 1, IndexNormal1: 1,
						IndexVertex2: 2, IndexNormal2: 2,
					},
				},
				TriangleTextures: []fileio.MD1TriangleTexture{
					{
						U0: 0, V0: 0,
						U1: 50, V1: 0,
						U2: 25, V2: 50,
						Page: 0,
					},
				},
				// Quad data
				QuadVertices: []fileio.MD1Vertex{
					{X: 100, Y: 0, Z: 0},
					{X: 200, Y: 0, Z: 0},
					{X: 200, Y: 0, Z: 100},
					{X: 100, Y: 0, Z: 100},
				},
				QuadNormals: []fileio.MD1Vertex{
					{X: 0, Y: 1000, Z: 0},
					{X: 0, Y: 1000, Z: 0},
					{X: 0, Y: 1000, Z: 0},
					{X: 0, Y: 1000, Z: 0},
				},
				QuadIndices: []fileio.MD1QuadIndex{
					{
						IndexVertex0: 0, IndexNormal0: 0,
						IndexVertex1: 1, IndexNormal1: 1,
						IndexVertex2: 2, IndexNormal2: 2,
						IndexVertex3: 3, IndexNormal3: 3,
					},
				},
				QuadTextures: []fileio.MD1QuadTexture{
					{
						U0: 0, V0: 0,
						U1: 100, V1: 0,
						U2: 100, V2: 100,
						U3: 0, V3: 100,
						Page: 1,
					},
				},
			},
		},
		NumBytes: 200,
	}

	textureData := &fileio.TIMOutput{
		ImageWidth:  256,
		ImageHeight: 256,
		NumPalettes: 4,
	}

	result := NewMD1Geometry(meshData, textureData)

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// 1 triangle (24 floats) + 1 quad (48 floats) = 72 floats
	expectedLength := (1 * 3 * 8) + (1 * 6 * 8) // 24 + 48 = 72
	if len(result) != expectedLength {
		t.Errorf("Expected length %d, got %d", expectedLength, len(result))
	}
}
