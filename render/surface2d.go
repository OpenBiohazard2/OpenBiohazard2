package render

import (
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	IMAGE_SURFACE_WIDTH  = 320
	IMAGE_SURFACE_HEIGHT = 240
)

type Surface2D struct {
	ImagePixels        []uint16
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
		ImagePixels:        imagePixels,
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

	renderGameStateUniform := gl.GetUniformLocation(programShader, gl.Str("gameState\x00"))
	gl.Uniform1i(renderGameStateUniform, RENDER_GAME_STATE_BACKGROUND_SOLID)

	renderDef.RenderSurface2D(renderDef.VideoBuffer)
}

func (renderDef *RenderDef) RenderTransparentVideoBuffer() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	programShader := renderDef.ProgramShader

	// Activate shader
	gl.UseProgram(programShader)

	renderGameStateUniform := gl.GetUniformLocation(programShader, gl.Str("gameState\x00"))
	gl.Uniform1i(renderGameStateUniform, RENDER_GAME_STATE_BACKGROUND_TRANSPARENT)

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

	programShader := r.ProgramShader

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

	diffuseUniform := gl.GetUniformLocation(programShader, gl.Str("diffuse\x00"))
	gl.Uniform1i(diffuseUniform, 0)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, surface.TextureId)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertexBuffer)/5))

	// Cleanup
	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
}

func (surface *Surface2D) ClearSurface() {
	for i := 0; i < IMAGE_SURFACE_WIDTH*IMAGE_SURFACE_HEIGHT; i++ {
		surface.ImagePixels[i] = 0
	}
}

func (surface *Surface2D) UpdateSurface(newImagePixels []uint16) {
	UpdateTexture(surface.TextureId, newImagePixels, IMAGE_SURFACE_WIDTH, IMAGE_SURFACE_HEIGHT)
}

func copyPixels(sourcePixels [][]uint16, startX int, startY int, width int, height int,
	destPixels []uint16, destX int, destY int) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			destPixels[((destY+y)*IMAGE_SURFACE_WIDTH)+(destX+x)] = sourcePixels[startY+y][startX+x]
		}
	}
}

func copyPixelsTransparent(sourcePixels [][]uint16, startX int, startY int, width int, height int,
	destPixels []uint16, destX int, destY int) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := sourcePixels[startY+y][startX+x]
			newPixel := uint16(0)
			if pixel > 0 {
				newPixel = (1 << 15) | pixel
			}

			// Overwrite pixel if it's not transparent
			if newPixel > 0 {
				destPixels[((destY+y)*IMAGE_SURFACE_WIDTH)+(destX+x)] = newPixel
			}
		}
	}
}

// Multiply pixels by a brightness factor
// Less than 1.0 will darken it
func copyPixelsBrightness(sourcePixels [][]uint16, startX int, startY int, width int, height int,
	destPixels []uint16, destX int, destY int, factor float64) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := sourcePixels[startY+y][startX+x]
			newPixel := int(0)
			if pixel > 0 {
				pixelR := (pixel % 32)
				pixelG := ((pixel >> 5) % 32)
				pixelB := ((pixel >> 10) % 32)
				newR := int(math.Floor(float64(pixelR) * factor))
				newG := int(math.Floor(float64(pixelG) * factor))
				newB := int(math.Floor(float64(pixelB) * factor))
				newPixel = (1 << 15) | (newB << 10) | (newG << 5) | newR
			}

			// Overwrite pixel if it's not transparent
			if newPixel > 0 {
				destPixels[((destY+y)*IMAGE_SURFACE_WIDTH)+(destX+x)] = uint16(newPixel)
			}
		}
	}
}

func fillPixels(newImageColors []uint16, destX int, destY int, width int, height int,
	r int, g int, b int) {
	// Convert color to A1R5G5B5 format
	newR := uint16(math.Round(float64(r) / 8.0))
	newG := uint16(math.Round(float64(g) / 8.0))
	newB := uint16(math.Round(float64(b) / 8.0))
	if newR >= 32 {
		newR = 31
	}
	if newG >= 32 {
		newG = 31
	}
	if newB >= 32 {
		newB = 31
	}
	color := (1 << 15) | (newB << 10) | (newG << 5) | newR

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			newImageColors[((destY+y)*IMAGE_SURFACE_WIDTH)+(destX+x)] = color
		}
	}
}

func buildSurface2DVertexBuffer() []float32 {
	z := float32(0.999)
	return []float32{
		// (-1, 1, z)
		-1.0, 1.0, z, 0.0, 0.0,
		// (-1, -1, z)
		-1.0, -1.0, z, 0.0, 1.0,
		// (1, -1, z)
		1.0, -1.0, z, 1.0, 1.0,

		// (1, -1, z)
		1.0, -1.0, z, 1.0, 1.0,
		// (1, 1, z)
		1.0, 1.0, z, 1.0, 0.0,
		// (-1, 1, z)
		-1.0, 1.0, z, 0.0, 0.0,
	}
}
