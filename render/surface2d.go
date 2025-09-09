package render

import (
	"github.com/OpenBiohazard2/OpenBiohazard2/geometry"
	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	IMAGE_SURFACE_WIDTH  = 320
	IMAGE_SURFACE_HEIGHT = 240
)

var (
	screenImage = NewImage16Bit(0, 0, IMAGE_SURFACE_WIDTH, IMAGE_SURFACE_HEIGHT)
)

type Surface2D struct {
	TextureId          uint32    // texture id in OpenGL
	VertexBuffer       []float32 // 3 elements for x,y,z and 2 elements for texture u,v
	VertexArrayObject  uint32
	VertexBufferObject uint32
}

func NewSurface2D() *Surface2D {
	var vao uint32
	gl.GenVertexArrays(1, &vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)

	imagePixels := make([]uint16, IMAGE_SURFACE_WIDTH*IMAGE_SURFACE_HEIGHT)
	return &Surface2D{
		TextureId:          BuildTexture(imagePixels, IMAGE_SURFACE_WIDTH, IMAGE_SURFACE_HEIGHT),
		VertexBuffer:       buildSurface2DVertexBuffer(),
		VertexArrayObject:  vao,
		VertexBufferObject: vbo,
	}
}

func (renderDef *RenderDef) RenderSolidVideoBuffer() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	programShader := renderDef.ProgramShader

	// Activate shader
	gl.UseProgram(programShader)

	// Use cached uniform location for better performance
	gl.Uniform1i(renderDef.UniformLocations.GameState, RENDER_GAME_STATE_BACKGROUND_SOLID)

	renderDef.RenderSurface2D(renderDef.VideoBuffer)
}

func (renderDef *RenderDef) RenderTransparentVideoBuffer() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	programShader := renderDef.ProgramShader

	// Activate shader
	gl.UseProgram(programShader)

	// Use cached uniform location for better performance
	gl.Uniform1i(renderDef.UniformLocations.GameState, RENDER_GAME_STATE_BACKGROUND_TRANSPARENT)

	renderDef.RenderSurface2D(renderDef.VideoBuffer)
}

func (r *RenderDef) RenderSurface2D(surface *Surface2D) {
	// Skip if surface is nil or has no vertex data
	if surface == nil || len(surface.VertexBuffer) == 0 {
		return
	}

	// Create renderer
	renderer := NewOpenGLRenderer(&r.UniformLocations)

	// Create render config for 2D surface (position + texture)
	config := renderer.Create2DEntityConfig(
		surface.VertexArrayObject,
		surface.VertexBufferObject,
		surface.VertexBuffer,
		surface.TextureId,
		nil, // No model matrix for surfaces
		RENDER_GAME_STATE_BACKGROUND_TRANSPARENT,
	)

	// Render the surface
	renderer.RenderEntity(config)
}

func (surface *Surface2D) UpdateSurface(newImage *Image16Bit) {
	UpdateTexture(surface.TextureId, newImage.GetPixelsForRendering(), int32(newImage.GetWidth()), int32(newImage.GetHeight()))
}

func buildSurface2DVertexBuffer() []float32 {
	z := float32(0.999)
	vertices := [4][]float32{
		{-1.0, 1.0, z},
		{-1.0, -1.0, z},
		{1.0, -1.0, z},
		{1.0, 1.0, z},
	}
	uvs := [4][]float32{
		{0.0, 0.0},
		{0.0, 1.0},
		{1.0, 1.0},
		{1.0, 0.0},
	}
	return geometry.NewTexturedRectangle(vertices, uvs).VertexBuffer
}
