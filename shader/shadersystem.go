package shader

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// UniformLocations caches OpenGL uniform locations for performance
type UniformLocations struct {
	// Main rendering uniforms
	GameState  int32
	View       int32
	Projection int32
	EnvLight   int32

	// Entity rendering uniforms
	RenderType int32
	Model      int32
	Diffuse    int32

	// Debug rendering uniforms
	DebugColor int32
	BoneOffset int32
}

// ShaderSystem manages shaders and uniform locations
type ShaderSystem struct {
	ProgramShader    uint32
	UniformLocations UniformLocations
}

// NewShaderSystem creates a new shader system
func NewShaderSystem() *ShaderSystem {
	return &ShaderSystem{
		ProgramShader:    0,
		UniformLocations: UniformLocations{},
	}
}

// Initialize loads and compiles shaders, then caches uniform locations
func (ss *ShaderSystem) Initialize(vertexPath, fragPath string) error {
	// Load and compile shaders
	shader, err := NewShader(vertexPath, fragPath)
	if err != nil {
		return fmt.Errorf("failed to create shader: %w", err)
	}

	ss.ProgramShader = shader.ProgramShader

	// Cache all uniform locations for performance
	ss.cacheUniformLocations()

	return nil
}

// Use activates the shader program
func (ss *ShaderSystem) Use() {
	gl.UseProgram(ss.ProgramShader)
}

// SetGameState sets the game state uniform
func (ss *ShaderSystem) SetGameState(gameState int32) {
	gl.Uniform1i(ss.UniformLocations.GameState, gameState)
}

// SetViewMatrix sets the view matrix uniform
func (ss *ShaderSystem) SetViewMatrix(viewMatrix mgl32.Mat4) {
	gl.UniformMatrix4fv(ss.UniformLocations.View, 1, false, &viewMatrix[0])
}

// SetProjectionMatrix sets the projection matrix uniform
func (ss *ShaderSystem) SetProjectionMatrix(projectionMatrix mgl32.Mat4) {
	gl.UniformMatrix4fv(ss.UniformLocations.Projection, 1, false, &projectionMatrix[0])
}

// SetEnvironmentLight sets the environment light uniform
func (ss *ShaderSystem) SetEnvironmentLight(envLight [3]float32) {
	gl.Uniform3fv(ss.UniformLocations.EnvLight, 1, &envLight[0])
}

// SetRenderType sets the render type uniform
func (ss *ShaderSystem) SetRenderType(renderType int32) {
	gl.Uniform1i(ss.UniformLocations.RenderType, renderType)
}

// SetModelMatrix sets the model matrix uniform
func (ss *ShaderSystem) SetModelMatrix(modelMatrix mgl32.Mat4) {
	gl.UniformMatrix4fv(ss.UniformLocations.Model, 1, false, &modelMatrix[0])
}

// SetDiffuse sets the diffuse texture uniform
func (ss *ShaderSystem) SetDiffuse(textureUnit int32) {
	gl.Uniform1i(ss.UniformLocations.Diffuse, textureUnit)
}

// SetDebugColor sets the debug color uniform
func (ss *ShaderSystem) SetDebugColor(color [4]float32) {
	gl.Uniform4f(ss.UniformLocations.DebugColor, color[0], color[1], color[2], color[3])
}

// SetBoneOffset sets the bone offset uniform
func (ss *ShaderSystem) SetBoneOffset(boneOffset mgl32.Mat4) {
	gl.UniformMatrix4fv(ss.UniformLocations.BoneOffset, 1, false, &boneOffset[0])
}

// GetUniformLocations returns the cached uniform locations
func (ss *ShaderSystem) GetUniformLocations() *UniformLocations {
	return &ss.UniformLocations
}

// cacheUniformLocations caches all uniform locations to avoid expensive gl.GetUniformLocation calls every frame
func (ss *ShaderSystem) cacheUniformLocations() {
	programShader := ss.ProgramShader

	// Main rendering uniforms
	ss.UniformLocations.GameState = gl.GetUniformLocation(programShader, gl.Str("gameState\x00"))
	ss.UniformLocations.View = gl.GetUniformLocation(programShader, gl.Str("view\x00"))
	ss.UniformLocations.Projection = gl.GetUniformLocation(programShader, gl.Str("projection\x00"))
	ss.UniformLocations.EnvLight = gl.GetUniformLocation(programShader, gl.Str("envLight\x00"))

	// Entity rendering uniforms
	ss.UniformLocations.RenderType = gl.GetUniformLocation(programShader, gl.Str("renderType\x00"))
	ss.UniformLocations.Model = gl.GetUniformLocation(programShader, gl.Str("model\x00"))
	ss.UniformLocations.Diffuse = gl.GetUniformLocation(programShader, gl.Str("diffuse\x00"))

	// Debug rendering uniforms
	ss.UniformLocations.DebugColor = gl.GetUniformLocation(programShader, gl.Str("debugColor\x00"))
	ss.UniformLocations.BoneOffset = gl.GetUniformLocation(programShader, gl.Str("boneOffset\x00"))
}

// NewShaderSystemForTesting creates a ShaderSystem without OpenGL dependencies for testing
func NewShaderSystemForTesting() *ShaderSystem {
	return &ShaderSystem{
		ProgramShader:    0, // No OpenGL context in tests
		UniformLocations: UniformLocations{},
	}
}
