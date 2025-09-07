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
	// Skip
	if surface == nil {
		return
	}

	vertexBuffer := surface.VertexBuffer
	if len(vertexBuffer) == 0 {
		return
	}

	floatSize := 4

	// 3 floats for vertex, 2 floats for texture UV
	stride := int32(5 * floatSize)

	vao := surface.VertexArrayObject
	gl.BindVertexArray(vao)

	vbo := surface.VertexBufferObject
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexBuffer)*floatSize, gl.Ptr(vertexBuffer), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Texture
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, stride, gl.PtrOffset(3*floatSize))
	gl.EnableVertexAttribArray(1)

	// Use cached uniform location for better performance
	gl.Uniform1i(r.UniformLocations.Diffuse, 0)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, surface.TextureId)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertexBuffer)/5))

	// Cleanup
	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
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
