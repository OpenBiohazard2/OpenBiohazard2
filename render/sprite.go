package render

import (
	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/OpenBiohazard2/OpenBiohazard2/geometry"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	RENDER_TYPE_SPRITE = 4
	SPRITE_FRAME_TIME  = 0.5 // in seconds
)

var (
	totalSpriteRuntime = float64(0)
	curSpriteFrame     = 0
)

type SpriteGroupEntity struct {
	SpriteTextureIndexMap map[int]int
	TextureIdPool         [][]uint32
	VertexBuffer          []float32
	VertexArrayObject     uint32
	VertexBufferObject    uint32
}

func NewSpriteGroupEntity(spriteData []fileio.SpriteData) *SpriteGroupEntity {
	spriteTextureIds := make([][]uint32, 0)
	for i := 0; i < len(spriteData); i++ {
		spriteFrames := BuildSpriteTexture(spriteData[i])
		spriteTextureIds = append(spriteTextureIds, spriteFrames)
	}

	spriteTextureIndexMap := make(map[int]int)
	for i := 0; i < len(spriteData); i++ {
		spriteTextureIndexMap[spriteData[i].Id] = i
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)

	return &SpriteGroupEntity{
		SpriteTextureIndexMap: spriteTextureIndexMap,
		TextureIdPool:         spriteTextureIds,
		VertexBuffer:          make([]float32, 0),
		VertexArrayObject:     vao,
		VertexBufferObject:    vbo,
	}
}

// Each sprite id has its own texture
// Build a texture for each frame
func BuildSpriteTexture(spriteData fileio.SpriteData) []uint32 {
	allFrameTextures := make([]uint32, 0)

	for _, frameData := range spriteData.FrameData {
		spriteId := frameData.SpriteId
		framePosition := spriteData.FramePositions[spriteId]

		frameHeight := int(frameData.SquareSide)
		frameWidth := int(frameData.SquareSide)

		if frameHeight == 0 || frameWidth == 0 {
			continue
		}

		startX := int(framePosition.ImageX)
		startY := int(framePosition.ImageY)

		frameImageColors := buildTexturePixels(spriteData.ImageData.PixelData, startX, startY, frameWidth, frameHeight)
		textureId := BuildTexture(frameImageColors, int32(frameWidth), int32(frameHeight))
		allFrameTextures = append(allFrameTextures, textureId)
	}

	return allFrameTextures
}

// buildTexturePixels extracts and processes pixel data for texture creation
func buildTexturePixels(pixelData [][]uint16, startX, startY, width, height int) []uint16 {
	texturePixels := make([]uint16, 0, width*height)
	
	for y := startY; y < startY+height; y++ {
		for x := startX; x < startX+width; x++ {
			curColor := pixelData[y][x]

			// Determine if pixel should be transparent
			// Set black to be transparent color
			newTextureColor := curColor
			if curColor > 0 {
				// Set alpha bit to 1
				newTextureColor = uint16(curColor) | (1 << 15)
			}
			texturePixels = append(texturePixels, newTextureColor)
		}
	}
	
	return texturePixels
}

func (renderDef *RenderDef) AddSprite(sprite fileio.ScriptInstrSceEsprOn) {
	spriteWidth := float32(1024 * 2)
	spriteCenter := mgl32.Vec3{float32(sprite.X), float32(sprite.Y), float32(sprite.Z)}
	viewMatrix := renderDef.ViewSystem.Camera.BuildViewMatrix()

	// Generate billboard sprite using geometry package
	rect := geometry.NewBillboardSprite(spriteCenter, spriteWidth, viewMatrix)
	renderDef.SceneSystem.SpriteGroupEntity.VertexBuffer = append(renderDef.SceneSystem.SpriteGroupEntity.VertexBuffer, rect.VertexBuffer...)
}

func RenderSprites(r *RenderDef, spriteGroupEntity *SpriteGroupEntity, timeElapsedSeconds float64) {
	vertexBuffer := spriteGroupEntity.VertexBuffer
	if len(vertexBuffer) == 0 {
		return
	}
	textureIds := spriteGroupEntity.TextureIdPool
	if len(textureIds) == 0 {
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

	// Create render config for 2D sprite (position + texture)
	config := r.Renderer.Create2DEntityConfig(
		spriteGroupEntity.VertexArrayObject,
		spriteGroupEntity.VertexBufferObject,
		vertexBuffer,
		textureIds[spriteIndex][curSpriteFrame],
		nil, // No model matrix for sprites
		RENDER_TYPE_SPRITE,
	)

	// Render the sprite
	r.Renderer.RenderEntity(config)
}
