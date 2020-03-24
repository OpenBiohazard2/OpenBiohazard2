package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/samuelyuan/openbiohazard2/fileio"
)

const (
	RENDER_TYPE_SPRITE = 4
	SPRITE_FRAME_TIME  = 0.5 // in seconds
)

var (
	totalSpriteRuntime = float64(0)
	curSpriteFrame     = 0
)

// Each sprite id has its own texture
// Build a texture for each frame
func BuildSpriteTexture(spriteData fileio.SpriteData) []uint32 {
	allFrameTextures := make([]uint32, 0)

	for i := 0; i < len(spriteData.Positions); i++ {
		position := spriteData.Positions[i]

		// offset is from center and a negative value
		frameHeight := (-2) * int(position.OffsetY)
		frameWidth := (-2) * int(position.OffsetX)

		if frameWidth == 0 || frameHeight == 0 {
			continue
		}

		startX := int(position.X)
		startY := int(position.Y)

		// If the sprite center is at the edge of the image, the dimensions should be halved
		if startX+frameWidth > spriteData.ImageData.ImageWidth {
			frameWidth = (-1) * int(position.OffsetX)
		}

		if startY+frameHeight > spriteData.ImageData.ImageHeight {
			frameHeight = (-1) * int(position.OffsetY)
		}

		frameImageColors := make([]uint16, 0)
		for y := startY; y < startY+frameHeight; y++ {
			for x := startX; x < startX+frameWidth; x++ {
				curColor := spriteData.ImageData.PixelData[y][x]

				// Determine if pixel should be transparent
				// Set black to be transparent color
				newTextureColor := curColor
				if curColor > 0 {
					// Set alpha bit to 1
					newTextureColor = uint16(curColor) | (1 << 15)
				}
				frameImageColors = append(frameImageColors, newTextureColor)
			}
		}
		textureId := BuildTexture(frameImageColors, int32(frameWidth), int32(frameHeight))
		allFrameTextures = append(allFrameTextures, textureId)
	}

	return allFrameTextures
}

func RenderSprites(programShader uint32, sprites []fileio.ScriptSprite, textureIds [][]uint32, timeElapsedSeconds float64) {
	renderTypeUniform := gl.GetUniformLocation(programShader, gl.Str("renderType\x00"))
	gl.Uniform1i(renderTypeUniform, RENDER_TYPE_SPRITE)

	vertexBuffer := make([]float32, 0)
	spriteWidth := float32(1024 * 2)
	for _, sprite := range sprites {
		vertex1 := mgl32.Vec3{float32(sprite.X), float32(sprite.Y) - spriteWidth, float32(sprite.Z)}
		vertex2 := mgl32.Vec3{float32(sprite.X), float32(sprite.Y) - spriteWidth, float32(sprite.Z) + float32(spriteWidth)}
		vertex3 := mgl32.Vec3{float32(sprite.X) + float32(spriteWidth), float32(sprite.Y), float32(sprite.Z) + float32(spriteWidth)}
		vertex4 := mgl32.Vec3{float32(sprite.X) + float32(spriteWidth), float32(sprite.Y), float32(sprite.Z)}
		rect := buildTexturedRectangle(vertex1, vertex2, vertex3, vertex4)
		vertexBuffer = append(vertexBuffer, rect...)
	}

	if len(vertexBuffer) == 0 {
		return
	}

	// TODO: Calculate index based on id
	spriteIndex := 0

	// Check when to move on to the next frame
	totalSpriteRuntime += timeElapsedSeconds
	if totalSpriteRuntime > SPRITE_FRAME_TIME {
		totalSpriteRuntime = 0
		curSpriteFrame++
		if curSpriteFrame >= len(textureIds[spriteIndex]) {
			curSpriteFrame = 0
		}
	}

	floatSize := 4

	// 3 floats for vertex, 2 floats for texture UV
	stride := int32(5 * floatSize)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
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
	gl.BindTexture(gl.TEXTURE_2D, textureIds[spriteIndex][curSpriteFrame])

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertexBuffer)/5))

	// Cleanup
	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
}

func buildTexturedRectangle(corner1 mgl32.Vec3, corner2 mgl32.Vec3, corner3 mgl32.Vec3, corner4 mgl32.Vec3) []float32 {
	rectBuffer := make([]float32, 0)
	vertex1 := []float32{corner1.X(), corner1.Y(), corner1.Z()}
	vertex2 := []float32{corner2.X(), corner2.Y(), corner2.Z()}
	vertex3 := []float32{corner3.X(), corner3.Y(), corner3.Z()}
	vertex4 := []float32{corner4.X(), corner4.Y(), corner4.Z()}

	rectBuffer = append(rectBuffer, vertex1...)
	rectBuffer = append(rectBuffer, 0.0, 0.0)
	rectBuffer = append(rectBuffer, vertex2...)
	rectBuffer = append(rectBuffer, 1.0, 0.0)
	rectBuffer = append(rectBuffer, vertex3...)
	rectBuffer = append(rectBuffer, 1.0, 1.0)

	rectBuffer = append(rectBuffer, vertex1...)
	rectBuffer = append(rectBuffer, 0.0, 0.0)
	rectBuffer = append(rectBuffer, vertex4...)
	rectBuffer = append(rectBuffer, 0.0, 1.0)
	rectBuffer = append(rectBuffer, vertex3...)
	rectBuffer = append(rectBuffer, 1.0, 1.0)
	return rectBuffer
}
