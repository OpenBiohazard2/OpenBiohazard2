package render

import (
	"github.com/OpenBiohazard2/OpenBiohazard2/shader"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// OpenGLRenderer encapsulates common OpenGL rendering operations
type OpenGLRenderer struct {
	uniformLocations *shader.UniformLocations
}

// NewOpenGLRenderer creates a new OpenGL renderer
func NewOpenGLRenderer(uniformLocations *shader.UniformLocations) *OpenGLRenderer {
	return &OpenGLRenderer{
		uniformLocations: uniformLocations,
	}
}

// VertexAttribute defines a single vertex attribute
type VertexAttribute struct {
	Index      uint32
	Size       int32
	Type       uint32
	Offset     int32
	Normalized bool
}

// RenderConfig contains all the data needed to render an entity
type RenderConfig struct {
	// Core OpenGL objects
	VAO          uint32
	VBO          uint32
	VertexBuffer []float32
	Stride       int32

	// Vertex attributes
	Attributes []VertexAttribute

	// Uniforms
	ModelMatrix *mgl32.Mat4
	RenderType  int32
	TextureID   uint32

	// Rendering parameters
	DrawMode    uint32
	VertexCount int32
}

// RenderEntity performs the common OpenGL rendering pattern
func (r *OpenGLRenderer) RenderEntity(config RenderConfig) {
	// Early return if no vertex data
	if len(config.VertexBuffer) == 0 {
		return
	}

	// Set uniforms
	r.setUniforms(config)

	// Setup vertex data
	r.setupVertexData(config)

	// Setup vertex attributes
	r.setupVertexAttributes(config)

	// Setup texture
	r.setupTexture(config)

	// Draw
	r.draw(config)

	// Cleanup
	r.cleanup(config)
}

// setUniforms sets all the required uniforms
func (r *OpenGLRenderer) setUniforms(config RenderConfig) {
	// Set render type
	gl.Uniform1i(r.uniformLocations.RenderType, config.RenderType)

	// Set model matrix if provided
	if config.ModelMatrix != nil {
		gl.UniformMatrix4fv(r.uniformLocations.Model, 1, false, &config.ModelMatrix[0])
	}
}

// setupVertexData binds VAO, VBO and uploads vertex data
func (r *OpenGLRenderer) setupVertexData(config RenderConfig) {
	// Bind vertex array object
	gl.BindVertexArray(config.VAO)

	// Bind vertex buffer object
	gl.BindBuffer(gl.ARRAY_BUFFER, config.VBO)

	// Upload vertex data
	const floatSize = 4
	gl.BufferData(gl.ARRAY_BUFFER, len(config.VertexBuffer)*floatSize, gl.Ptr(config.VertexBuffer), gl.STATIC_DRAW)
}

// setupVertexAttributes configures all vertex attribute pointers
func (r *OpenGLRenderer) setupVertexAttributes(config RenderConfig) {
	for _, attr := range config.Attributes {
		gl.VertexAttribPointer(attr.Index, attr.Size, attr.Type, attr.Normalized, config.Stride, gl.PtrOffset(int(attr.Offset)))
		gl.EnableVertexAttribArray(attr.Index)
	}
}

// setupTexture binds the texture
func (r *OpenGLRenderer) setupTexture(config RenderConfig) {
	// Set diffuse texture uniform
	gl.Uniform1i(r.uniformLocations.Diffuse, 0)

	// Bind texture
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, config.TextureID)
}

// draw performs the actual drawing
func (r *OpenGLRenderer) draw(config RenderConfig) {
	gl.DrawArrays(config.DrawMode, 0, config.VertexCount)
}

// cleanup disables all vertex attribute arrays
func (r *OpenGLRenderer) cleanup(config RenderConfig) {
	for _, attr := range config.Attributes {
		gl.DisableVertexAttribArray(attr.Index)
	}
}

// Helper function to create a standard 3D entity render config
func (r *OpenGLRenderer) Create3DEntityConfig(
	vao, vbo uint32,
	vertexBuffer []float32,
	textureID uint32,
	modelMatrix *mgl32.Mat4,
	renderType int32,
) RenderConfig {
	const floatSize = 4
	stride := int32(8 * floatSize) // 3 pos + 2 tex + 3 normal

	attributes := []VertexAttribute{
		{Index: 0, Size: 3, Type: gl.FLOAT, Offset: 0, Normalized: false},             // Position
		{Index: 1, Size: 2, Type: gl.FLOAT, Offset: 3 * floatSize, Normalized: false}, // Texture UV
		{Index: 2, Size: 3, Type: gl.FLOAT, Offset: 5 * floatSize, Normalized: false}, // Normal
	}

	return RenderConfig{
		VAO:          vao,
		VBO:          vbo,
		VertexBuffer: vertexBuffer,
		Stride:       stride,
		Attributes:   attributes,
		ModelMatrix:  modelMatrix,
		RenderType:   renderType,
		TextureID:    textureID,
		DrawMode:     gl.TRIANGLES,
		VertexCount:  int32(len(vertexBuffer) / 8), // 8 floats per vertex
	}
}

// Helper function to create a 2D entity render config
func (r *OpenGLRenderer) Create2DEntityConfig(
	vao, vbo uint32,
	vertexBuffer []float32,
	textureID uint32,
	modelMatrix *mgl32.Mat4,
	renderType int32,
) RenderConfig {
	const floatSize = 4
	stride := int32(5 * floatSize) // 3 pos + 2 tex

	attributes := []VertexAttribute{
		{Index: 0, Size: 3, Type: gl.FLOAT, Offset: 0, Normalized: false},             // Position
		{Index: 1, Size: 2, Type: gl.FLOAT, Offset: 3 * floatSize, Normalized: false}, // Texture UV
	}

	return RenderConfig{
		VAO:          vao,
		VBO:          vbo,
		VertexBuffer: vertexBuffer,
		Stride:       stride,
		Attributes:   attributes,
		ModelMatrix:  modelMatrix,
		RenderType:   renderType,
		TextureID:    textureID,
		DrawMode:     gl.TRIANGLES,
		VertexCount:  int32(len(vertexBuffer) / 5), // 5 floats per vertex
	}
}

// Helper function to create a debug entity render config
func (r *OpenGLRenderer) CreateDebugEntityConfig(
	vao, vbo uint32,
	vertexBuffer []float32,
	renderType int32,
) RenderConfig {
	const floatSize = 4
	stride := int32(3 * floatSize) // 3 pos only

	attributes := []VertexAttribute{
		{Index: 0, Size: 3, Type: gl.FLOAT, Offset: 0, Normalized: false}, // Position only
	}

	return RenderConfig{
		VAO:          vao,
		VBO:          vbo,
		VertexBuffer: vertexBuffer,
		Stride:       stride,
		Attributes:   attributes,
		ModelMatrix:  nil, // Debug entities don't use model matrix
		RenderType:   renderType,
		TextureID:    0, // Debug entities don't use textures
		DrawMode:     gl.TRIANGLES,
		VertexCount:  int32(len(vertexBuffer) / 3), // 3 floats per vertex
	}
}

// AnimatedEntityConfig contains data for rendering animated entities with bone transforms
type AnimatedEntityConfig struct {
	// Core OpenGL objects
	VAO          uint32
	VBO          uint32
	VertexBuffer []float32
	Stride       int32

	// Vertex attributes
	Attributes []VertexAttribute

	// Uniforms
	ModelMatrix *mgl32.Mat4
	RenderType  int32
	TextureID   uint32

	// Animation data
	BoneTransforms   []mgl32.Mat4
	ComponentOffsets []ComponentOffset
}

// ComponentOffset defines the vertex range for a mesh component
type ComponentOffset struct {
	StartIndex int32
	EndIndex   int32
}

// RenderAnimatedEntity performs rendering for animated entities with bone transforms
func (r *OpenGLRenderer) RenderAnimatedEntity(config AnimatedEntityConfig) {
	// Early return if no vertex data
	if len(config.VertexBuffer) == 0 {
		return
	}

	// Set uniforms
	r.setUniforms(RenderConfig{
		ModelMatrix: config.ModelMatrix,
		RenderType:  config.RenderType,
		TextureID:   config.TextureID,
	})

	// Setup vertex data
	r.setupVertexData(RenderConfig{
		VAO:          config.VAO,
		VBO:          config.VBO,
		VertexBuffer: config.VertexBuffer,
	})

	// Setup vertex attributes
	r.setupVertexAttributes(RenderConfig{
		Stride:     config.Stride,
		Attributes: config.Attributes,
	})

	// Setup texture
	r.setupTexture(RenderConfig{
		TextureID: config.TextureID,
	})

	// Render all components with bone transforms
	r.renderAnimatedComponents(config)

	// Cleanup
	r.cleanup(RenderConfig{
		Attributes: config.Attributes,
	})
}

// renderAnimatedComponents draws all mesh components with bone transforms
func (r *OpenGLRenderer) renderAnimatedComponents(config AnimatedEntityConfig) {
	for i, offset := range config.ComponentOffsets {
		// Set bone transform uniform
		if i < len(config.BoneTransforms) {
			gl.UniformMatrix4fv(r.uniformLocations.BoneOffset, 1, false, &config.BoneTransforms[i][0])
		}

		// Calculate vertex range
		vertOffset := int32(offset.StartIndex / 8) // 8 floats per vertex (3 pos + 2 tex + 3 normal)
		numVertices := int32((offset.EndIndex - offset.StartIndex) / 8)

		// Draw component
		gl.DrawArrays(gl.TRIANGLES, vertOffset, numVertices)
	}
}

// Helper function to create an animated entity render config
func (r *OpenGLRenderer) CreateAnimatedEntityConfig(
	vao, vbo uint32,
	vertexBuffer []float32,
	textureID uint32,
	modelMatrix *mgl32.Mat4,
	renderType int32,
	boneTransforms []mgl32.Mat4,
	componentOffsets []ComponentOffset,
) AnimatedEntityConfig {
	const floatSize = 4
	stride := int32(8 * floatSize) // 3 pos + 2 tex + 3 normal

	attributes := []VertexAttribute{
		{Index: 0, Size: 3, Type: gl.FLOAT, Offset: 0, Normalized: false},             // Position
		{Index: 1, Size: 2, Type: gl.FLOAT, Offset: 3 * floatSize, Normalized: false}, // Texture UV
		{Index: 2, Size: 3, Type: gl.FLOAT, Offset: 5 * floatSize, Normalized: false}, // Normal
	}

	return AnimatedEntityConfig{
		VAO:              vao,
		VBO:              vbo,
		VertexBuffer:     vertexBuffer,
		Stride:           stride,
		Attributes:       attributes,
		ModelMatrix:      modelMatrix,
		RenderType:       renderType,
		TextureID:        textureID,
		BoneTransforms:   boneTransforms,
		ComponentOffsets: componentOffsets,
	}
}
