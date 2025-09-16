package render

import (
	"testing"

	"github.com/OpenBiohazard2/OpenBiohazard2/shader"
	"github.com/go-gl/mathgl/mgl32"
)

func TestOpenGLRenderer_Create3DEntityConfig(t *testing.T) {
	// Create mock uniform locations
	uniformLocations := &shader.UniformLocations{
		RenderType: 1,
		Model:      2,
		Diffuse:    3,
	}

	renderer := NewOpenGLRenderer(uniformLocations)

	// Test data
	vao := uint32(100)
	vbo := uint32(200)
	textureID := uint32(300)
	renderType := int32(5)

	// Create vertex buffer (8 floats per vertex: 3 pos + 2 tex + 3 normal)
	vertexBuffer := []float32{
		// Vertex 1
		1.0, 2.0, 3.0, // Position
		0.1, 0.2, // Texture UV
		0.0, 1.0, 0.0, // Normal
		// Vertex 2
		4.0, 5.0, 6.0, // Position
		0.3, 0.4, // Texture UV
		1.0, 0.0, 0.0, // Normal
	}

	modelMatrix := mgl32.Ident4()

	// Create config
	config := renderer.Create3DEntityConfig(vao, vbo, vertexBuffer, textureID, &modelMatrix, renderType)

	// Verify config
	if config.VAO != vao {
		t.Errorf("Expected VAO %d, got %d", vao, config.VAO)
	}

	if config.VBO != vbo {
		t.Errorf("Expected VBO %d, got %d", vbo, config.VBO)
	}

	if config.TextureID != textureID {
		t.Errorf("Expected TextureID %d, got %d", textureID, config.TextureID)
	}

	if config.RenderType != renderType {
		t.Errorf("Expected RenderType %d, got %d", renderType, config.RenderType)
	}

	if config.DrawMode != 4 { // gl.TRIANGLES
		t.Errorf("Expected DrawMode %d, got %d", 4, config.DrawMode)
	}

	// Verify stride (8 floats * 4 bytes = 32)
	expectedStride := int32(8 * 4)
	if config.Stride != expectedStride {
		t.Errorf("Expected Stride %d, got %d", expectedStride, config.Stride)
	}

	// Verify vertex count (16 floats / 8 floats per vertex = 2 vertices)
	expectedVertexCount := int32(2)
	if config.VertexCount != expectedVertexCount {
		t.Errorf("Expected VertexCount %d, got %d", expectedVertexCount, config.VertexCount)
	}

	// Verify attributes
	if len(config.Attributes) != 3 {
		t.Errorf("Expected 3 attributes, got %d", len(config.Attributes))
	}

	// Verify position attribute
	posAttr := config.Attributes[0]
	if posAttr.Index != 0 || posAttr.Size != 3 || posAttr.Offset != 0 {
		t.Errorf("Position attribute incorrect: Index=%d, Size=%d, Offset=%d", posAttr.Index, posAttr.Size, posAttr.Offset)
	}

	// Verify texture attribute
	texAttr := config.Attributes[1]
	expectedTexOffset := int32(3 * 4) // 3 floats * 4 bytes
	if texAttr.Index != 1 || texAttr.Size != 2 || texAttr.Offset != expectedTexOffset {
		t.Errorf("Texture attribute incorrect: Index=%d, Size=%d, Offset=%d", texAttr.Index, texAttr.Size, texAttr.Offset)
	}

	// Verify normal attribute
	normalAttr := config.Attributes[2]
	expectedNormalOffset := int32(5 * 4) // 5 floats * 4 bytes
	if normalAttr.Index != 2 || normalAttr.Size != 3 || normalAttr.Offset != expectedNormalOffset {
		t.Errorf("Normal attribute incorrect: Index=%d, Size=%d, Offset=%d", normalAttr.Index, normalAttr.Size, normalAttr.Offset)
	}
}

func TestOpenGLRenderer_Create2DEntityConfig(t *testing.T) {
	uniformLocations := &shader.UniformLocations{}
	renderer := NewOpenGLRenderer(uniformLocations)

	vao := uint32(100)
	vbo := uint32(200)
	textureID := uint32(300)
	renderType := int32(1)

	// Create vertex buffer (5 floats per vertex: 3 pos + 2 tex)
	vertexBuffer := []float32{
		// Vertex 1
		1.0, 2.0, 3.0, // Position
		0.1, 0.2, // Texture UV
		// Vertex 2
		4.0, 5.0, 6.0, // Position
		0.3, 0.4, // Texture UV
	}

	modelMatrix := mgl32.Ident4()

	config := renderer.Create2DEntityConfig(vao, vbo, vertexBuffer, textureID, &modelMatrix, renderType)

	// Verify basic properties
	if config.VAO != vao || config.VBO != vbo || config.TextureID != textureID {
		t.Error("Basic config properties incorrect")
	}

	// Verify stride (5 floats * 4 bytes = 20)
	expectedStride := int32(5 * 4)
	if config.Stride != expectedStride {
		t.Errorf("Expected Stride %d, got %d", expectedStride, config.Stride)
	}

	// Verify vertex count (10 floats / 5 floats per vertex = 2 vertices)
	expectedVertexCount := int32(2)
	if config.VertexCount != expectedVertexCount {
		t.Errorf("Expected VertexCount %d, got %d", expectedVertexCount, config.VertexCount)
	}

	// Verify attributes (should have 2: position + texture)
	if len(config.Attributes) != 2 {
		t.Errorf("Expected 2 attributes, got %d", len(config.Attributes))
	}
}

func TestOpenGLRenderer_CreateDebugEntityConfig(t *testing.T) {
	uniformLocations := &shader.UniformLocations{}
	renderer := NewOpenGLRenderer(uniformLocations)

	vao := uint32(100)
	vbo := uint32(200)
	renderType := int32(-1)

	// Create vertex buffer (3 floats per vertex: position only)
	vertexBuffer := []float32{
		1.0, 2.0, 3.0, // Position
		4.0, 5.0, 6.0, // Position
	}

	config := renderer.CreateDebugEntityConfig(vao, vbo, vertexBuffer, renderType)

	// Verify basic properties
	if config.VAO != vao || config.VBO != vbo {
		t.Error("Basic config properties incorrect")
	}

	// Verify no texture
	if config.TextureID != 0 {
		t.Errorf("Expected TextureID 0, got %d", config.TextureID)
	}

	// Verify no model matrix
	if config.ModelMatrix != nil {
		t.Error("Expected ModelMatrix to be nil for debug entities")
	}

	// Verify stride (3 floats * 4 bytes = 12)
	expectedStride := int32(3 * 4)
	if config.Stride != expectedStride {
		t.Errorf("Expected Stride %d, got %d", expectedStride, config.Stride)
	}

	// Verify vertex count (6 floats / 3 floats per vertex = 2 vertices)
	expectedVertexCount := int32(2)
	if config.VertexCount != expectedVertexCount {
		t.Errorf("Expected VertexCount %d, got %d", expectedVertexCount, config.VertexCount)
	}

	// Verify attributes (should have 1: position only)
	if len(config.Attributes) != 1 {
		t.Errorf("Expected 1 attribute, got %d", len(config.Attributes))
	}

	posAttr := config.Attributes[0]
	if posAttr.Index != 0 || posAttr.Size != 3 || posAttr.Offset != 0 {
		t.Errorf("Position attribute incorrect: Index=%d, Size=%d, Offset=%d", posAttr.Index, posAttr.Size, posAttr.Offset)
	}
}

func TestOpenGLRenderer_NewOpenGLRenderer(t *testing.T) {
	uniformLocations := &shader.UniformLocations{
		RenderType: 1,
		Model:      2,
		Diffuse:    3,
	}

	renderer := NewOpenGLRenderer(uniformLocations)

	if renderer == nil {
		t.Fatal("NewOpenGLRenderer returned nil")
	}

	if renderer.uniformLocations != uniformLocations {
		t.Error("UniformLocations not set correctly")
	}
}

// Test that the OpenGLRenderer can be created and configured correctly
func TestOpenGLRenderer_Configuration(t *testing.T) {
	uniformLocations := &shader.UniformLocations{
		RenderType: 1,
		Model:      2,
		Diffuse:    3,
	}

	renderer := NewOpenGLRenderer(uniformLocations)

	if renderer == nil {
		t.Fatal("NewOpenGLRenderer returned nil")
	}

	if renderer.uniformLocations != uniformLocations {
		t.Error("UniformLocations not set correctly")
	}
}
