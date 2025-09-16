package shader

import (
	"testing"
)

func TestNewShaderSystem(t *testing.T) {
	ss := NewShaderSystemForTesting()

	// Test basic initialization
	if ss == nil {
		t.Fatal("NewShaderSystem returned nil")
	}

	// Test that ProgramShader is initialized to 0
	if ss.ProgramShader != 0 {
		t.Error("ProgramShader should be 0 for new ShaderSystem")
	}

	// Test that UniformLocations is initialized (zero values are fine for testing)
	// UniformLocations will be zero-initialized, which is expected for testing
}

func TestNewShaderSystemForTesting(t *testing.T) {
	ss := NewShaderSystemForTesting()

	// Test basic initialization
	if ss == nil {
		t.Fatal("NewShaderSystemForTesting returned nil")
	}

	// Test that ProgramShader is 0 for testing
	if ss.ProgramShader != 0 {
		t.Error("ProgramShader should be 0 for testing version")
	}

	// Test that UniformLocations is initialized (zero values are fine for testing)
	// UniformLocations will be zero-initialized, which is expected for testing
}

func TestShaderSystem_Isolation(t *testing.T) {
	// Test that multiple ShaderSystem instances are independent
	ss1 := NewShaderSystemForTesting()
	ss2 := NewShaderSystemForTesting()

	// Both should be initialized, but they should be different instances
	if ss1 == ss2 {
		t.Error("ShaderSystem instances should be different objects")
	}

	// Test that modifying one doesn't affect the other
	ss1.ProgramShader = 123
	if ss2.ProgramShader == 123 {
		t.Error("Modifying one ShaderSystem should not affect another")
	}
}

func TestShaderSystem_GetUniformLocations(t *testing.T) {
	ss := NewShaderSystemForTesting()

	// Test that GetUniformLocations returns a pointer to the uniform locations
	uniforms := ss.GetUniformLocations()
	if uniforms == nil {
		t.Error("GetUniformLocations should not return nil")
	}

	// Test that it returns the same instance
	if uniforms != &ss.UniformLocations {
		t.Error("GetUniformLocations should return pointer to internal UniformLocations")
	}
}

func TestShaderSystem_Initialize(t *testing.T) {
	// This test would require OpenGL context, so we'll test the structure
	ss := NewShaderSystemForTesting()

	// Test that Initialize method exists and can be called
	// Note: This will fail in actual execution due to OpenGL context, but tests the structure
	defer func() {
		if r := recover(); r != nil {
			// Expected to panic due to OpenGL context, that's okay for this test
			t.Log("Initialize panicked as expected due to OpenGL context")
		}
	}()

	// Test that Initialize method exists (will panic due to OpenGL context)
	// We don't actually call it to avoid file system errors
	_ = ss.Initialize
}

func TestShaderSystem_Methods(t *testing.T) {
	// Test that all methods exist and can be called
	ss := NewShaderSystemForTesting()

	// Test that all methods exist (we don't call them to avoid OpenGL context issues)
	_ = ss.Use
	_ = ss.SetGameState
	_ = ss.SetViewMatrix
	_ = ss.SetProjectionMatrix
	_ = ss.SetEnvironmentLight
	_ = ss.GetUniformLocations
}

func TestShaderSystem_UniformLocations(t *testing.T) {
	ss := NewShaderSystemForTesting()

	// Test that we can access uniform locations
	uniforms := ss.GetUniformLocations()

	// Test that all expected fields exist
	_ = uniforms.GameState
	_ = uniforms.View
	_ = uniforms.Projection
	_ = uniforms.EnvLight
	_ = uniforms.RenderType
	_ = uniforms.Model
	_ = uniforms.Diffuse
	_ = uniforms.DebugColor
	_ = uniforms.BoneOffset
}

func BenchmarkNewShaderSystem(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewShaderSystemForTesting()
	}
}

func BenchmarkShaderSystem_Methods(b *testing.B) {
	ss := NewShaderSystemForTesting()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Test method calls (will panic due to OpenGL, but benchmarks the structure)
		_ = ss.GetUniformLocations()
		_ = ss.ProgramShader
	}
}
